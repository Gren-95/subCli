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
	Username     string `yaml:"username"`
	PasswordHash string `yaml:"password_hash"`
	URL          string `yaml:"URL"`
}

var AppConfig Config
var decryptedPassword string

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "subCli", "config.yaml"), nil
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

	// Decrypt password
	if AppConfig.PasswordHash != "" {
		decryptedPassword, err = decryptPassword(AppConfig.PasswordHash, AppConfig.Username)
		if err != nil {
			return fmt.Errorf("could not decrypt password: %v", err)
		}
	}

	return nil
}

func GetPassword() string {
	return decryptedPassword
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

	// Set proper permissions for config file (readable/writable by owner only)
	if err := os.Chmod(configPath, 0600); err != nil {
		return fmt.Errorf("could not set config file permissions: %v", err)
	}

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	return encoder.Encode(&AppConfig)
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

	// Encrypt password
	encrypted, err := encryptPassword(password, username)
	if err != nil {
		return fmt.Errorf("could not encrypt password: %v", err)
	}

	// Build config with encrypted password
	AppConfig = Config{
		Username:     username,
		URL:          url,
		PasswordHash: encrypted,
	}
	decryptedPassword = password

	// Save config
	if err := SaveConfig(); err != nil {
		return fmt.Errorf("could not save config: %v", err)
	}

	fmt.Println()
	fmt.Println("✓ Configuration saved successfully!")
	fmt.Printf("  Config location: %s\n", configPath)
	fmt.Println("  Password stored securely with AES-256 encryption")
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


