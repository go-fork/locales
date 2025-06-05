# Scripts Directory

This directory contains automation scripts for release management and project maintenance.

## Available Scripts

### Release Management

- **`create_release_templates.sh`** - Creates release documentation templates for new versions
- **`archive_release.sh`** - Archives completed releases and prepares for next version

## Usage

### Creating Release Templates

```bash
# Create templates for next version (auto-increment)
./scripts/create_release_templates.sh

# Create templates for specific version
./scripts/create_release_templates.sh v0.2.0
```

### Archiving Completed Release

```bash
# Archive current release and prepare for next
./scripts/archive_release.sh v0.1.3

# Archive with custom next version
./scripts/archive_release.sh v0.1.3 v0.2.0
```

## Script Requirements

- **Git**: All scripts require git to be available in PATH
- **Bash**: Scripts are written in bash and require bash 4.0+
- **Permission**: Make sure scripts are executable (`chmod +x scripts/*.sh`)

## Workflow Integration

These scripts are designed to work with:
- GitHub Actions workflows in `.github/workflows/`
- Release documentation structure in `releases/`
- Automated version tagging and release creation

## Development

When modifying scripts:
1. Test locally before committing
2. Update this README if adding new scripts
3. Follow existing naming conventions
4. Add proper error handling and validation
