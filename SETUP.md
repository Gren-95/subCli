# First-Time Setup

subCli includes an interactive setup wizard for easy configuration.

## Running Setup

```bash
subcli setup
```

## What It Does

The setup wizard will:

1. **Prompt for Server URL**
   - Example: `https://music.example.com`

2. **Ask for Username**
   - Your Subsonic username

3. **Securely Collect Password**
   - Password input is hidden (not echoed to screen)

4. **Choose Password Storage**
   - **Encrypted (Recommended)**: Password is encrypted using AES-256
   - **Plain text**: Password stored as-is (less secure)

5. **Save Configuration**
   - Saves to `~/.config/subcli/config.yaml`

6. **Test Connection**
   - Automatically tests your credentials

## Example Session

```
╔═══════════════════════════════════════════╗
║     subCli - First Time Configuration     ║
╚═══════════════════════════════════════════╝

Subsonic Server URL (e.g., https://music.example.com): https://my-server.com
Username: myuser
Password (hidden): ********

Password Storage Options:
  1. Encrypted (recommended)
  2. Plain text (less secure)
Choose option (1-2) [1]: 1

✓ Configuration saved successfully!
  Config location: /home/user/.config/subcli/config.yaml

Testing connection... ✓ Connection test successful!

You're all set! Try running:
  subcli --shuffle | mpv --playlist=-
```

## Security Features

### Encrypted Password Storage

When you choose encrypted storage:
- Password is encrypted using **AES-256-GCM**
- Encryption key is derived from your username using SHA-256
- Unique nonce (initialization vector) for each encryption
- Password is never stored in plain text
- Decrypted only when needed to make API calls

### Config File Example

**With Encryption (Recommended):**
```yaml
username: myuser
password_hash: encrypted_base64_string_here
URL: https://my-server.com
use_encryption: true
```

**Without Encryption:**
```yaml
username: myuser
password: mypassword
URL: https://my-server.com
use_encryption: false
```

## Reconfiguration

To change your settings, simply run setup again:

```bash
subcli setup
```

You'll be asked if you want to overwrite the existing configuration.

## Manual Configuration

You can also manually create/edit `~/.config/subcli/config.yaml`:

```yaml
username: your_username
password: your_password
URL: https://your-server.com
```

Note: Manually created configs use plain text passwords unless you specify `password_hash` and `use_encryption: true`.

## Security Best Practices

1. ✅ **Use encrypted storage** when running setup
2. ✅ **Set proper file permissions**:
   ```bash
   chmod 600 ~/.config/subcli/config.yaml
   ```
3. ✅ **Never commit** config files to version control
4. ✅ **Use HTTPS** for your Subsonic server URL
5. ⚠️ Avoid plain text passwords when possible

## Troubleshooting

### Setup Fails to Connect

If connection test fails:
- Check your server URL (include http:// or https://)
- Verify username and password
- Ensure server is accessible from your network
- Check firewall settings

### Config File Issues

Check config file location and permissions:
```bash
ls -la ~/.config/subcli/config.yaml
cat ~/.config/subcli/config.yaml
```

### Re-run Setup

Delete config and start over:
```bash
rm ~/.config/subcli/config.yaml
subcli setup
```

## What Happens on First Run

If you try to use subcli without running setup first:

```bash
$ subcli --shuffle
Error loading config: could not open config file: ...

Run 'subcli setup' to configure your connection.
```

The application will guide you to run setup first!

