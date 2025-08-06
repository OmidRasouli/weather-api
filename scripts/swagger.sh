#!/bin/bash
set -e

# Check if swag is installed
if ! command -v swag &> /dev/null; then
    echo "Error: 'swag' command not found."
    echo "Please install it with the following command:"
    echo "go install github.com/swaggo/swag/cmd/swag@latest"
    echo "Then make sure your GOPATH/bin directory is in your PATH"
    echo "Typical locations:"
    echo "  - Linux/Mac: export PATH=\$PATH:\$HOME/go/bin"
    echo "  - Windows: add %USERPROFILE%\\go\\bin to your PATH"
    exit 1
fi

# Generate swagger documentation
swag init -g cmd/main.go --output docs

echo "Swagger documentation generated successfully!"
