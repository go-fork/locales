#!/bin/bash

# archive_release.sh
# Archives completed releases and prepares for next version

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to get the next semantic version
get_next_version() {
    local current_version="$1"
    local version_type="${2:-patch}"
    
    # Remove 'v' prefix if present
    local version="${current_version#v}"
    
    # Split version into parts
    IFS='.' read -ra VERSION_PARTS <<< "$version"
    local major="${VERSION_PARTS[0]}"
    local minor="${VERSION_PARTS[1]}"
    local patch="${VERSION_PARTS[2]}"
    
    case "$version_type" in
        "major")
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        "minor")
            minor=$((minor + 1))
            patch=0
            ;;
        "patch"|*)
            patch=$((patch + 1))
            ;;
    esac
    
    echo "v${major}.${minor}.${patch}"
}

# Function to validate version format
validate_version() {
    local version="$1"
    if [[ ! "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        print_error "Invalid version format: $version"
        print_error "Expected format: vX.Y.Z (e.g., v0.1.3)"
        return 1
    fi
    return 0
}

# Function to check if version exists in git tags
version_exists() {
    local version="$1"
    git tag | grep -q "^${version}$"
}

# Function to archive release
archive_release() {
    local release_version="$1"
    local next_version="$2"
    local releases_dir="releases"
    local next_dir="${releases_dir}/next"
    local archive_dir="${releases_dir}/${release_version}"
    
    print_info "Archiving release $release_version"
    
    # Check if next directory exists and has content
    if [[ ! -d "$next_dir" ]]; then
        print_error "No $next_dir directory found. Run create_release_templates.sh first."
        return 1
    fi
    
    # Check if required files exist
    local required_files=("RELEASE_NOTES.md" "RELEASE_SUMMARY.md" "MIGRATION.md")
    for file in "${required_files[@]}"; do
        if [[ ! -f "$next_dir/$file" ]]; then
            print_warning "Missing $file in $next_dir"
        fi
    done
    
    # Create archive directory
    mkdir -p "$archive_dir"
    
    # Check if archive directory already exists and has content
    if [[ -n "$(ls -A "$archive_dir" 2>/dev/null)" ]]; then
        print_warning "Archive directory $archive_dir already exists and is not empty"
        read -p "Do you want to overwrite it? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_info "Keeping existing archive"
            return 1
        fi
        rm -rf "$archive_dir"
        mkdir -p "$archive_dir"
    fi
    
    # Copy files from next to archive
    if ls "$next_dir"/*.md > /dev/null 2>&1; then
        cp "$next_dir"/*.md "$archive_dir/"
        print_success "Copied release documents to $archive_dir/"
    else
        print_warning "No markdown files found in $next_dir"
    fi
    
    # Update version references in archived files
    if [[ -f "$archive_dir/RELEASE_NOTES.md" ]]; then
        # Update any placeholder versions in the archived files
        sed -i.bak "s/\$target_version/$release_version/g" "$archive_dir/RELEASE_NOTES.md" 2>/dev/null || true
        sed -i.bak "s/TBD - Version/$release_version/g" "$archive_dir/RELEASE_NOTES.md" 2>/dev/null || true
        rm -f "$archive_dir/RELEASE_NOTES.md.bak" 2>/dev/null || true
    fi
    
    if [[ -f "$archive_dir/RELEASE_SUMMARY.md" ]]; then
        sed -i.bak "s/\$target_version/$release_version/g" "$archive_dir/RELEASE_SUMMARY.md" 2>/dev/null || true
        sed -i.bak "s/TBD - Version/$release_version/g" "$archive_dir/RELEASE_SUMMARY.md" 2>/dev/null || true
        rm -f "$archive_dir/RELEASE_SUMMARY.md.bak" 2>/dev/null || true
    fi
    
    if [[ -f "$archive_dir/MIGRATION.md" ]]; then
        sed -i.bak "s/\$target_version/$release_version/g" "$archive_dir/MIGRATION.md" 2>/dev/null || true
        sed -i.bak "s/TBD - Version/$release_version/g" "$archive_dir/MIGRATION.md" 2>/dev/null || true
        rm -f "$archive_dir/MIGRATION.md.bak" 2>/dev/null || true
    fi
    
    # Create new templates for next version
    print_info "Creating templates for next version: $next_version"
    
    # Clear the next directory
    rm -rf "$next_dir"
    mkdir -p "$next_dir"
    
    # Create new templates using the create_release_templates.sh script
    if [[ -x "./scripts/create_release_templates.sh" ]]; then
        ./scripts/create_release_templates.sh "$next_version"
    else
        # Fallback: create basic templates inline
        create_basic_templates "$next_version" "$next_dir"
    fi
    
    print_success "Release $release_version archived successfully!"
    print_info "Next steps:"
    echo "  1. Review the archived release in $archive_dir/"
    echo "  2. Edit the new templates in $next_dir/ for $next_version"
    echo "  3. Commit the changes"
    echo "  4. Create and push the release tag: git tag $release_version && git push origin $release_version"
}

# Function to create basic templates (fallback)
create_basic_templates() {
    local version="$1"
    local target_dir="$2"
    
    cat > "$target_dir/RELEASE_NOTES.md" << EOF
# Release Notes - $version

## Overview
Brief description of this release and its main purpose.

## What's New
### 🚀 Features
- New feature 1

### 🐛 Bug Fixes
- Fix for issue #X

### 🔧 Improvements
- Performance improvement 1

## Breaking Changes
### ⚠️ Important Notes
None in this release.

## Migration Guide
See [MIGRATION.md](./MIGRATION.md) for detailed migration instructions.

---
Release Date: TBD
EOF

    cat > "$target_dir/RELEASE_SUMMARY.md" << EOF
# $version Release Summary

## Quick Overview
One-line summary of what this release brings.

## Key Highlights
- 🎉 **Major Feature**: Description of the most important feature

---
**Full Release Notes**: [RELEASE_NOTES.md](./RELEASE_NOTES.md)  
**Migration Guide**: [MIGRATION.md](./MIGRATION.md)  
**Release Date**: TBD
EOF

    cat > "$target_dir/MIGRATION.md" << EOF
# Migration Guide - $version

## Overview
This guide helps you migrate from the previous version to $version.

## Quick Migration Checklist
- [ ] Update dependencies
- [ ] Run tests to ensure compatibility

## Step-by-Step Migration

### Step 1: Update Dependencies
\`\`\`bash
go get go.fork.vn/locales@$version
go mod tidy
\`\`\`

### Step 2: Run Tests
\`\`\`bash
go test ./...
\`\`\`

---
**Need Help?** Feel free to open an issue or discussion on GitHub.
EOF
}

# Main script logic
main() {
    local release_version=""
    local next_version=""
    
    # Check if we're in a git repository
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        print_error "This script must be run from within a git repository"
        exit 1
    fi
    
    # Parse command line arguments
    case "${1:-}" in
        -h|--help)
            echo "Usage: $0 RELEASE_VERSION [NEXT_VERSION]"
            echo ""
            echo "Archives a completed release and prepares templates for the next version."
            echo ""
            echo "Arguments:"
            echo "  RELEASE_VERSION  Version to archive (e.g., v0.1.3)"
            echo "  NEXT_VERSION     Next version to prepare (optional, auto-increments patch if not provided)"
            echo ""
            echo "Examples:"
            echo "  $0 v0.1.3                    # Archive v0.1.3, prepare v0.1.4"
            echo "  $0 v0.1.3 v0.2.0            # Archive v0.1.3, prepare v0.2.0"
            exit 0
            ;;
        "")
            print_error "Missing required argument: RELEASE_VERSION"
            echo "Use $0 --help for usage information"
            exit 1
            ;;
        *)
            release_version="$1"
            next_version="${2:-}"
            ;;
    esac
    
    # Validate release version format
    if ! validate_version "$release_version"; then
        exit 1
    fi
    
    # Auto-generate next version if not provided
    if [[ -z "$next_version" ]]; then
        next_version=$(get_next_version "$release_version" "patch")
        print_info "Auto-generating next version: $next_version"
    else
        if ! validate_version "$next_version"; then
            exit 1
        fi
    fi
    
    # Check if release version already exists as a tag (warning, not error)
    if version_exists "$release_version"; then
        print_warning "Git tag $release_version already exists"
        read -p "Continue with archiving? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_info "Cancelled by user"
            exit 0
        fi
    fi
    
    # Confirmation before proceeding
    echo ""
    print_info "Archive Plan:"
    echo "  Release Version: $release_version"
    echo "  Next Version:    $next_version"
    echo "  Action:          Move releases/next/ → releases/$release_version/"
    echo "  Action:          Create new templates in releases/next/ for $next_version"
    echo ""
    read -p "Proceed with archiving? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Cancelled by user"
        exit 0
    fi
    
    archive_release "$release_version" "$next_version"
}

# Run main function with all arguments
main "$@"
