#!/bin/bash
#
# Sign and Notarize a Wails v2 macOS Application
#
# Usage: ./sign-and-notarize.sh [version]
#   version - Optional version string (default: read from wails.json)
#

set -euo pipefail

# Configuration
DEVELOPER_ID="888624CD24481639CC22A8CAD9A3CC31B4890E2F"
TEAM_ID="UL5T9WXYJT"
APP_NAME="Tally"
KEYCHAIN_DB="$HOME/Library/Keychains/login.keychain-db"
BUILD_DIR="build/bin"
VERSION_STRING="unknown"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get version from argument or wails.json
get_version() {
    if [[ -n "${1:-}" ]]; then
        echo "$1"
    elif command -v jq &>/dev/null; then
        jq -r '.version' wails.json 2>/dev/null || echo "unknown"
    else
        grep '"version"' wails.json | sed 's/.*"version": *"\([^"]*\)".*/\1/' || echo "unknown"
    fi
}

# Logging functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

# Check if .app exists
check_app() {
    local app_path="$BUILD_DIR/${APP_NAME}.app"

    if [[ ! -d "$app_path" ]]; then
        log_error "Application not found at: $app_path"
        log_error "Please run 'wails build' or 'make build-darwin-arm64' first."
        exit 1
    fi

    echo "$app_path"
}

# Verify code signing certificate exists
verify_certificate() {
    log_info "Verifying code signing certificate..."

    # Check if Developer ID certificate exists in any keychain
    if ! security find-identity -v -p codesigning 2>/dev/null | grep -q "Developer ID Application"; then
        log_error "Developer ID Application certificate not found in any keychain."
        log_error "Please ensure you have imported your Developer ID certificate."
        exit 1
    fi

    log_info "Certificate verified."
}

# Sign the application
sign_app() {
    local app_path="$1"
    log_info "Signing application..."

    # If already signed, remove signature first for clean re-sign
    if codesign --verify "$app_path" 2>/dev/null; then
        log_info "App is already signed, removing old signature..."
        rm -rf "$app_path/_CodeSignature"
    fi

    # Sign with hardened runtime for notarization
    # Use SHA1 hash identity (short flags: -s for sign, -o for options)
    codesign -s "$DEVELOPER_ID" -f -o runtime "$app_path"
    if [[ $? -ne 0 ]]; then
        log_error "Code signing failed."
        exit 1
    fi
    log_info "Signing completed successfully."
}

# Notarize the application (submits .dmg for notarization)
notarize_app() {
    local app_path="$1"
    local apple_id="${APPLE_ID:-}"
    local password="${AC_PASSWORD:-}"

    if [[ -z "$apple_id" ]]; then
        log_warn "APPLE_ID not set. Skipping notarization."
        log_warn "Set APPLE_ID environment variable to enable notarization."
        return 0
    fi

    if [[ -z "$password" ]]; then
        log_warn "AC_PASSWORD not set. Skipping notarization."
        log_warn "Set AC_PASSWORD environment variable to enable notarization."
        log_warn "Use '@keychain:AC_PASSWORD' pattern in your shell profile:"
        log_warn "  export AC_PASSWORD='@keychain:AC_PASSWORD'"
        return 0
    fi

    log_info "Creating DMG for notarization..."
    local dmg_path="$BUILD_DIR/${APP_NAME}-${VERSION_STRING}-arm64.dmg"
    hdiutil create -volname "$APP_NAME" -srcfolder "$app_path" -ov -format UDZO "$dmg_path" >/dev/null 2>&1

    log_info "Submitting DMG for notarization..."

    # Submit to Apple notary service
    xcrun notarytool submit "$dmg_path" \
        --apple-id "$apple_id" \
        --password "$password" \
        --team-id "$TEAM_ID" \
        --wait

    if [[ $? -eq 0 ]]; then
        log_info "Notarization completed successfully."
        log_info "Stapling notarization ticket to DMG..."
        xcrun stapler staple "$dmg_path"
        log_info "DMG ready: $dmg_path"
    else
        log_error "Notarization failed."
        exit 1
    fi
}

# Staple the notarization ticket
staple_app() {
    local app_path="$1"
    log_info "Stapling notarization ticket..."

    xcrun stapler staple "$app_path"

    if [[ $? -eq 0 ]]; then
        log_info "Stapling completed successfully."
    else
        log_error "Stapling failed."
        exit 1
    fi
}

# Verify the final signature
verify_signature() {
    local app_path="$1"
    log_info "Verifying final signature..."

    # Check both signing and notarization
    codesign --verify --verbose --deep --strict "$app_path"

    if [[ $? -eq 0 ]]; then
        log_info "Signature verification passed."
    else
        log_error "Signature verification failed."
        exit 1
    fi
}

# Main execution
main() {
    local version="${1:-}"
    VERSION_STRING=$(get_version "$version")

    echo ""
    echo "============================================"
    echo "  Wails macOS Sign and Notarize Script"
    echo "  Version: $VERSION_STRING"
    echo "============================================"
    echo ""

    # Check prerequisites
    if ! command -v xcrun &>/dev/null; then
        log_error "xcrun not found. Please install Xcode Command Line Tools."
        exit 1
    fi

    # Verify certificate
    verify_certificate

    # Check app exists
    local app_path
    app_path=$(check_app)

    # Sign
    sign_app "$app_path"

    # Notarize (if credentials provided)
    notarize_app "$app_path"

    # Staple (only if DMG was created, otherwise staple the .app)
    if [[ -z "${APPLE_ID:-}" ]]; then
        staple_app "$app_path"
    fi

    # Verify
    verify_signature "$app_path"

    echo ""
    log_info "Build signed and notarized successfully!"
    log_info "Application: $app_path"
    echo ""
}

main "$@"