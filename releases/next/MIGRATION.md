# Migration Guide - v0.1.1

## Overview
This guide helps you migrate from the previous version to v0.1.1.

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
```go
// Old way (previous version)
oldFunction(param1, param2)

// New way (v0.1.1)
newFunction(param1, param2, newParam)
```

#### Removed Functions
- `removedFunction()` - Use `newAlternativeFunction()` instead

#### Changed Types
```go
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
```

### Configuration Changes
If you're using configuration files:

```yaml
# Old configuration format
old_setting: value
deprecated_option: true

# New configuration format
new_setting: value
# deprecated_option removed
new_option: false
```

## Step-by-Step Migration

### Step 1: Update Dependencies
```bash
go get go.fork.vn/locales@v0.1.1
go mod tidy
```

### Step 2: Update Import Statements
```go
// If import paths changed
import (
    "go.fork.vn/locales" // Updated import
)
```

### Step 3: Update Code
Replace deprecated function calls:

```go
// Before
result := locales.OldFunction(param)

// After
result := locales.NewFunction(param, defaultValue)
```

### Step 4: Update Configuration
Update your configuration files according to the new schema.

### Step 5: Run Tests
```bash
go test ./...
```

## Common Issues and Solutions

### Issue 1: Function Not Found
**Problem**: `undefined: locales.OldFunction`  
**Solution**: Replace with `locales.NewFunction`

### Issue 2: Type Mismatch
**Problem**: `cannot use int as int64`  
**Solution**: Cast the value or update variable type

## Getting Help
- Check the [documentation](https://pkg.go.dev/go.fork.vn/locales@v0.1.1)
- Search [existing issues](https://github.com/go-fork/locales/issues)
- Create a [new issue](https://github.com/go-fork/locales/issues/new) if needed

## Rollback Instructions
If you need to rollback:

```bash
go get go.fork.vn/locales@previous-version
go mod tidy
```

Replace `previous-version` with your previous version tag.

---
**Need Help?** Feel free to open an issue or discussion on GitHub.
