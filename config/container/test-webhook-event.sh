#!/bin/bash

# Usage: ./send_webhook.sh <base_url> <external_ref> [status]
# Example: ./send_webhook.sh https://example.com 123456789 completed

# Validate required inputs
BASE_URL="$1"
EXTERNAL_REF="$2"
STATUS="${3:-completed}" # Default status is 'completed'

if [ -z "$BASE_URL" ] || [ -z "$EXTERNAL_REF" ]; then
  echo -e "\033[31mError: base_url and external_ref are required.\033[0m"
  echo "Usage: $0 <base_url> <external_ref> [status]"
  exit 1
fi

# Compose full URL
URL="${BASE_URL}/prod/webhook_events" 

# Construct JSON payload
DATA=$(jq -n \
  --arg ref "$EXTERNAL_REF" \
  --arg status "$STATUS" \
  '{external_reference: $ref, status: $status}')

# Send POST request and capture response and status code
response=$(curl -s -w "\n%{http_code}" -X POST "$URL" \
  -H "Content-Type: application/json" \
  -d "$DATA")

# Split body and status code
body=$(echo "$response" | sed '$d')
code=$(echo "$response" | tail -n1)

# Print result with color based on status code
if [[ "$code" -ge 200 && "$code" -lt 300 ]]; then
  echo -e "\033[32mStatus: $code\033[0m"
else
  echo -e "\033[31mStatus: $code\033[0m"
fi

echo "Response body:"
echo "$body"
