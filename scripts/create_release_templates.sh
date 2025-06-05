#!/bin/bash

# create_release_templates.sh
# Creates release documentation templates for new versions

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

# Function to get latest version from git tags
get_latest_version() {
    local latest_tag
    latest_tag=$(git tag --sort=-version:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | head -n1 2>/dev/null || echo "v0.0.0")
    echo "$latest_tag"
}

# Function to create release templates
create_templates() {
    local target_version="$1"
    local releases_dir="releases"
    local next_dir="${releases_dir}/next"
    
    print_info "Creating release templates for version: $target_version"
    
    # Ensure releases directory structure exists
    mkdir -p "$next_dir"
    
    # Check if templates already exist
    if [[ -f "$next_dir/RELEASE_NOTES.md" ]]; then
        print_warning "Release templates already exist in $next_dir"
        read -p "Do you want to overwrite them? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_info "Keeping existing templates"
            return 0
        fi
    fi
    
    # Create RELEASE_NOTES.md template
    cat > "$next_dir/RELEASE_NOTES.md" << EOF
# Release Notes - $target_version

## Overview
Brief description of this release and its main purpose.

## What's New
### 🚀 Features
- New feature 1
- New feature 2

### 🐛 Bug Fixes
- Fix for issue #X
- Fix for issue #Y

### 🔧 Improvements
- Performance improvement 1
- Code quality improvement 2

### 📚 Documentation
- Updated documentation for feature X
- Added examples for use case Y

## Breaking Changes
### ⚠️ Important Notes
- Breaking change 1 (if any)
- Breaking change 2 (if any)

## Migration Guide
See [MIGRATION.md](./MIGRATION.md) for detailed migration instructions.

## Dependencies
### Updated
- dependency-name: vX.Y.Z → vA.B.C

### Added
- new-dependency: vX.Y.Z

### Removed
- removed-dependency: vX.Y.Z

## Performance
- Benchmark improvement: X% faster in scenario Y
- Memory usage: X% reduction in scenario Z

## Security
- Security fix for vulnerability X
- Updated dependencies with security patches

## Testing
- Added X new test cases
- Improved test coverage to X%

## Contributors
Thanks to all contributors who made this release possible:
- @contributor1
- @contributor2

## Download
- Source code: [go.fork.vn/locales@$target_version]
- Documentation: [pkg.go.dev/go.fork.vn/locales@$target_version]

---
Release Date: $(date '+%Y-%m-%d')
EOF

    # Create RELEASE_SUMMARY.md template
    cat > "$next_dir/RELEASE_SUMMARY.md" << EOF
# $target_version Release Summary

## Quick Overview
One-line summary of what this release brings.

## Key Highlights
- 🎉 **Major Feature**: Description of the most important feature
- 🚀 **Performance**: Key performance improvements
- 🔧 **Developer Experience**: Improvements for developers using this library

## Stats
- **Issues Closed**: X
- **Pull Requests Merged**: Y
- **New Contributors**: Z
- **Files Changed**: A
- **Lines Added**: B
- **Lines Removed**: C

## Impact
Brief description of how this release affects users and the ecosystem.

## Next Steps
What to expect in the next release cycle.

---
**Full Release Notes**: [RELEASE_NOTES.md](./RELEASE_NOTES.md)  
**Migration Guide**: [MIGRATION.md](./MIGRATION.md)  
**Release Date**: $(date '+%Y-%m-%d')
EOF

    # Create MIGRATION.md template
    cat > "$next_dir/MIGRATION.md" << EOF
# Migration Guide - $target_version

## Overview
This guide helps you migrate from the previous version to $target_version.

## Prerequisites
- Go 1.23 or later
- Previous version installed

## Quick Migration Checklist
- [ ] Update import statements (if changed)
- [ ] Update function calls (if signatures changed)
- [ ] Update configuration (if format changed)
- [ ] Run tests to ensure compatibility
- [ ] Update documentation references

## Breaking Changes

### API Changes
#### Changed Functions
\`\`\`go
// Old way (previous version)
oldFunction(param1, param2)

// New way ($target_version)
newFunction(param1, param2, newParam)
\`\`\`

#### Removed Functions
- \`removedFunction()\` - Use \`newAlternativeFunction()\` instead

#### Changed Types
\`\`\`go
// Old type definition
type OldConfig struct {
    Field1 string
    Field2 int
}

// New type definition
type NewConfig struct {
    Field1 string
    Field2 int64 // Changed from int
    Field3 bool  // New field
}
\`\`\`

### Configuration Changes
If you're using configuration files:

\`\`\`yaml
# Old configuration format
old_setting: value
deprecated_option: true

# New configuration format
new_setting: value
# deprecated_option removed
new_option: false
\`\`\`

## Step-by-Step Migration

### Step 1: Update Dependencies
\`\`\`bash
go get go.fork.vn/locales@$target_version
go mod tidy
\`\`\`

### Step 2: Update Import Statements
\`\`\`go
// If import paths changed
import (
    "go.fork.vn/locales" // Updated import
)
\`\`\`

### Step 3: Update Code
Replace deprecated function calls:

\`\`\`go
// Before
result := locales.OldFunction(param)

// After
result := locales.NewFunction(param, defaultValue)
\`\`\`

### Step 4: Update Configuration
Update your configuration files according to the new schema.

### Step 5: Run Tests
\`\`\`bash
go test ./...
\`\`\`

## Common Issues and Solutions

### Issue 1: Function Not Found
**Problem**: \`undefined: locales.OldFunction\`  
**Solution**: Replace with \`locales.NewFunction\`

### Issue 2: Type Mismatch
**Problem**: \`cannot use int as int64\`  
**Solution**: Cast the value or update variable type

## Getting Help
- Check the [documentation](https://pkg.go.dev/go.fork.vn/locales@$target_version)
- Search [existing issues](https://github.com/go-fork/locales/issues)
- Create a [new issue](https://github.com/go-fork/locales/issues/new) if needed

## Rollback Instructions
If you need to rollback:

\`\`\`bash
go get go.fork.vn/locales@previous-version
go mod tidy
\`\`\`

Replace \`previous-version\` with your previous version tag.

---
**Need Help?** Feel free to open an issue or discussion on GitHub.
EOF

    print_success "Created release templates in $next_dir/"
    print_info "Next steps:"
    echo "  1. Edit the templates with actual release information"
    echo "  2. Update version-specific details"
    echo "  3. Review and commit the changes"
    echo "  4. Use archive_release.sh when ready to release"
}

# Main script logic
main() {
    local target_version=""
    local version_type="patch"
    
    # Check if we're in a git repository
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        print_error "This script must be run from within a git repository"
        exit 1
    fi
    
    # Parse command line arguments
    case "${1:-}" in
        -h|--help)
            echo "Usage: $0 [VERSION] [VERSION_TYPE]"
            echo ""
            echo "Creates release documentation templates for a new version."
            echo ""
            echo "Arguments:"
            echo "  VERSION      Target version (e.g., v0.2.0). If not provided, auto-increments latest."
            echo "  VERSION_TYPE Type of version increment: major, minor, patch (default: patch)"
            echo ""
            echo "Examples:"
            echo "  $0                    # Auto-increment patch version"
            echo "  $0 v0.2.0            # Create templates for specific version"
            echo "  $0 \"\" minor          # Auto-increment minor version"
            exit 0
            ;;
        "")
            # Auto-increment version
            local latest_version
            latest_version=$(get_latest_version)
            target_version=$(get_next_version "$latest_version" "$version_type")
            print_info "Auto-incrementing from $latest_version to $target_version"
            ;;
        *)
            target_version="$1"
            version_type="${2:-patch}"
            
            # Validate version format
            if [[ ! "$target_version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
                print_error "Invalid version format: $target_version"
                print_error "Expected format: vX.Y.Z (e.g., v0.2.0)"
                exit 1
            fi
            ;;
    esac
    
    create_templates "$target_version"
}

# Run main function with all arguments
main "$@"
