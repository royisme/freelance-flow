# macOS Signing and Notarization for Wails

This directory contains scripts for signing and notarizing the Tally Wails application for distribution outside the Mac App Store.

## Prerequisites

### 1. Apple Developer Certificate

You need a **Developer ID Application** certificate to sign apps for distribution outside the App Store.

**Certificate Setup:**
1. Request a Developer ID Application certificate in your Apple Developer account
2. Export the certificate and private key as a `.p12` file
3. Import into your keychain:
   ```bash
   security import Certificate.p12 -k ~/Library/Keychains/login.keychain-db
   ```

### 2. App-Specific Password

To notarize, you need an App-Specific Password for your Apple ID.

**Setup:**
1. Go to https://appleid.apple.com/account/manage
2. Create an App-Specific Password
3. Store it in your keychain with the label `AC_PASSWORD`:
   ```bash
   security add-generic-password -s "AC_PASSWORD" -a "your@email.com" -w "xxxx-xxxx-xxxx-xxxx" -k ~/Library/Keychains/login.keychain-db
   ```

### 3. Environment Variables

Add these to your shell profile (`.zshrc` or `.bashrc`):

```bash
# Apple Developer credentials
export APPLE_ID="your@email.com"
export AC_PASSWORD="@keychain:AC_PASSWORD"
```

The `@keychain:` prefix tells the notarytool to fetch the password from your keychain instead of using the literal string.

## Build Targets

The Wails app can be built for multiple macOS architectures:

```bash
# ARM64 (Apple Silicon)
make build-darwin-arm64

# AMD64 (Intel)
make build-darwin-amd64

# Both architectures
make build-darwin-all
```

Output is placed in `build/bin/Tally.app`.

## Signing and Notarization

### The Workflow

1. **Build**: Compile the Wails application
2. **Sign**: Apply code signature with hardened runtime
3. **Notarize**: Submit to Apple for automated security checks
4. **Staple**: Attach the notarization ticket to the executable

### Running the Script

```bash
# Make executable (first time only)
chmod +x scripts/sign-and-notarize.sh

# Run with default version from wails.json
./scripts/sign-and-notarize.sh

# Run with custom version
./scripts/sign-and-notarize.sh 1.2.3
```

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `APPLE_ID` | Your Apple ID email | Yes (for notarization) |
| `AC_PASSWORD` | App-Specific Password or `@keychain:AC_PASSWORD` | Yes (for notarization) |

Without notarization credentials, the script will sign the app but skip notarization.

## Manual Commands

If you need to run steps individually:

```bash
# 1. Build
wails build -platform darwin/arm64

# 2. Sign
codesign --sign --verbose \
    --keychain ~/Library/Keychains/login.keychain-db \
    "Developer ID Application: shaoqing zhu (UL5T9WXYJT)" \
    --options runtime \
    build/bin/Tally.app

# 3. Notarize
xcrun notarytool submit build/bin/Tally.app \
    --apple-id "your@email.com" \
    --password "@keychain:AC_PASSWORD" \
    --team-id UL5T9WXYJT \
    --wait

# 4. Staple
xcrun stapler staple build/bin/Tally.app

# 5. Verify
codesign --verify --verbose --deep --strict build/bin/Tally.app
```

## Troubleshooting

### "Certificate not found in keychain"

Run this to verify your certificate is imported:
```bash
security find-identity -v -p codesigning ~/Library/Keychains/login.keychain-db
```

### "Notarization failed"

- Verify your Apple ID and App-Specific Password are correct
- Check your internet connection
- Apple services status: https://developer.apple.com/system-status/

### "Staple failed"

Try again - Apple's servers can be unreliable:
```bash
xcrun stapler staple -v build/bin/Tally.app
```

### "Signature verification failed"

Re-sign and notarize:
```bash
codesign --remove-signatures build/bin/Tally.app
./scripts/sign-and-notarize.sh
```

## References

- [Wails macOS Signing](https://wails.io/docs/guides/signing/)
- [Apple Developer - Notarizing macOS Software](https://developer.apple.com/documentation/xcode/notarizing_macos_software_before_distribution)
- [Apple Developer - Creating App-Specific Passwords](https://support.apple.com/en-us/102414)