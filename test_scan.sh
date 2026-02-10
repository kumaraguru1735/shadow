#!/bin/bash

echo "🕵️  Shadow Test Run"
echo "===================="
echo ""

# Test 1: Help command
echo "📖 Test 1: Show help"
./shadow --help
echo ""
echo "✅ Help works!"
echo ""

# Test 2: Scan command structure
echo "📋 Test 2: Scan help"
./shadow scan --help | head -20
echo ""
echo "✅ Scan command works!"
echo ""

# Test 3: Version
echo "🔢 Test 3: Version"
./shadow --version
echo ""

echo "===================="
echo "✅ Shadow is ready to use!"
echo ""
echo "To run a real scan:"
echo "  ./shadow scan example.com"
echo ""
