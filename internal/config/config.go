package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Username       string `yaml:"username"`
	Password       string `yaml:"password,omitempty"`
	PasswordHash   string `yaml:"password_hash,omitempty"`
	URL            string `yaml:"URL"`
	UseEncryption  bool   `yaml:"use_encryption,omitempty"`
}

var AppConfig Config
var decryptedPassword string

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "subcli", "config.yaml"), nil
}

func LoadConfig() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	file, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("could not open config file: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&AppConfig); err != nil {
		return fmt.Errorf("could not decode config: %v", err)
	}

	// Handle encrypted password
	if AppConfig.PasswordHash != "" && AppConfig.UseEncryption {
		decryptedPassword, err = decryptPassword(AppConfig.PasswordHash, AppConfig.Username)
		if err != nil {
			return fmt.Errorf("could not decrypt password: %v", err)
		}
	} else if AppConfig.Password != "" {
		decryptedPassword = AppConfig.Password
	}

	return nil
}

func GetPassword() string {
	if decryptedPassword != "" {
		return decryptedPassword
	}
	return AppConfig.Password
}

func SaveConfig() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a copy to avoid modifying the original
	configToSave := AppConfig
	
	// If using encryption, clear the plain text password
	if configToSave.UseEncryption && configToSave.PasswordHash != "" {
		configToSave.Password = ""
	}
	
	// If not using encryption, clear the hash
	if !configToSave.UseEncryption && configToSave.Password != "" {
		configToSave.PasswordHash = ""
	}

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	return encoder.Encode(&configToSave)
}

// InteractiveSetup runs the first-time setup wizard
func InteractiveSetup() error {
	fmt.Println("╔═══════════════════════════════════════════╗")
	fmt.Println("║     subCli - First Time Configuration     ║")
	fmt.Println("╚═══════════════════════════════════════════╝")
	fmt.Println()

	// Check if config already exists
	configPath, _ := getConfigPath()
	if _, err := os.Stat(configPath); err == nil {
		fmt.Print("Configuration file already exists. Overwrite? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			fmt.Println("Setup cancelled.")
			return nil
		}
	}

	// Get server URL
	fmt.Print("Subsonic Server URL (e.g., https://music.example.com): ")
	var url string
	fmt.Scanln(&url)
	url = strings.TrimSpace(url)
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	// Get username
	fmt.Print("Username: ")
	var username string
	fmt.Scanln(&username)
	username = strings.TrimSpace(username)
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	// Get password securely
	fmt.Print("Password (hidden): ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("could not read password: %v", err)
	}
	fmt.Println() // New line after password input
	password := string(passwordBytes)

	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	// Ask about encryption
	fmt.Println()
	fmt.Println("Password Storage Options:")
	fmt.Println("  1. Encrypted (recommended)")
	fmt.Println("  2. Plain text (less secure)")
	fmt.Print("Choose option (1-2) [1]: ")
	var encChoice string
	fmt.Scanln(&encChoice)
	if encChoice == "" {
		encChoice = "1"
	}

	useEncryption := encChoice == "1"

	// Build config
	AppConfig = Config{
		Username:      username,
		URL:           url,
		UseEncryption: useEncryption,
	}

	if useEncryption {
		encrypted, err := encryptPassword(password, username)
		if err != nil {
			return fmt.Errorf("could not encrypt password: %v", err)
		}
		AppConfig.PasswordHash = encrypted
		decryptedPassword = password
	} else {
		AppConfig.Password = password
		decryptedPassword = password
	}

	// Save config
	if err := SaveConfig(); err != nil {
		return fmt.Errorf("could not save config: %v", err)
	}

	fmt.Println()
	fmt.Println("✓ Configuration saved successfully!")
	fmt.Printf("  Config location: %s\n", configPath)
	fmt.Println()

	// Test connection
	fmt.Print("Testing connection... ")
	// This will be called from main after InitPlayer
	
	return nil
}

// encryptPassword encrypts the password using AES-256
func encryptPassword(password, username string) (string, error) {
	// Derive a key from the username
	key := sha256.Sum256([]byte(username + "subcli-salt-v1"))

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt
	ciphertext := gcm.Seal(nonce, nonce, []byte(password), nil)
	
	// Encode to base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decryptPassword decrypts the password using AES-256
func decryptPassword(encryptedPassword, username string) (string, error) {
	// Decode from base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedPassword)
	if err != nil {
		return "", err
	}

	// Derive the same key
	key := sha256.Sum256([]byte(username + "subcli-salt-v1"))

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}


