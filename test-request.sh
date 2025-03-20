#!/bin/bash

# Replace with your actual ngrok URL
NGROK_URL="https://1974-62-163-212-117.ngrok-free.app"

# Test creating a user with Treblle headers
echo "Creating a user with Treblle tracking headers..."
curl -X POST "${NGROK_URL}/api/v1/users" \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 12345" \
  -H "X-Trace-ID: test-trace-123" \
  -d '{"name": "John Doe", "email": "john@example.com"}'

echo -e "\n\nGetting all users with Treblle tracking headers..."
curl -X GET "${NGROK_URL}/api/v1/users" \
  -H "X-User-ID: 12345" \
  -H "X-Trace-ID: test-trace-456"
