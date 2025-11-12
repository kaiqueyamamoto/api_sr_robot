#!/bin/bash
set -e

echo "🔧 Installing dependencies..."
go mod download

echo "📦 Installing Swag..."
go install github.com/swaggo/swag/cmd/swag@latest

echo "📖 Generating Swagger documentation..."
# Try multiple paths for swag
if command -v swag &> /dev/null; then
    swag init -g main.go --output ./docs
elif [ -f ~/go/bin/swag ]; then
    ~/go/bin/swag init -g main.go --output ./docs
elif [ -f $GOPATH/bin/swag ]; then
    $GOPATH/bin/swag init -g main.go --output ./docs
else
    echo "⚠️  Warning: swag not found, skipping documentation generation"
    mkdir -p ./docs
fi

echo "🔨 Building application..."
go build -o chatserver main.go

echo "✅ Build completed successfully!"

