#!/bin/bash

# Check if base_url was passed as an argument
base_url="$1"
if [ -z "$base_url" ]; then
    echo "Usage: $0 <base_url>"
    exit 1
fi

# Function to check if required dependencies are installed
check_dependencies() {
    if ! command -v jq &> /dev/null; then
        echo -e "\033[31mError: jq is not installed. Please install jq to proceed.\033[0m"
        exit 1
    fi
}

# Function to print status and body with color-coded output
print_response() {
    local status_code=$1
    local body=$2

    if [[ "$status_code" -ge 200 && "$status_code" -lt 300 ]]; then
        echo -e "\033[32mStatus: $status_code\033[0m"
    elif [[ "$status_code" -ge 400 && "$status_code" -lt 500 ]]; then
        echo -e "\033[34mStatus: $status_code\033[0m"
    else
        echo -e "\033[31mStatus: $status_code\033[0m"
    fi

    if [[ -z "$body" ]]; then
        echo -e "\033[33mNo Body\033[0m"
    else
        echo -e "$body" | jq .
    fi
    echo "---------------------------------------------"
}

# Function to generate random number
generate_random_number() {
    echo $((RANDOM % (999999999999 - 1000 + 1) + 1000))
}

# Function to check health
check_health() {
    echo "Checking healthcheck..."
    health_response=$(curl --silent --write-out "HTTPSTATUS:%{http_code}" --location "$base_url/healthcheck")
    health_body=$(echo "$health_response" | sed -e 's/HTTPSTATUS\:.*//g')
    health_status=$(echo "$health_response" | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')

    print_response "$health_status" "$health_body"

    if [[ "$health_status" -ne 200 ]]; then
        echo "Healthcheck failed, exiting..."
        exit 1
    fi
}

# --- Main Script ---

check_dependencies
check_health

for i in {1..5}; do
    random_amount=$((RANDOM % 1000000 + 1))
    rand_int=$(generate_random_number)

    echo "Creating payment #$i"

    response=$(curl --silent --write-out "HTTPSTATUS:%{http_code}" --location "$base_url/api/v1/payments/authorize" \
        --header 'Content-Type: application/json' \
        --data '{
            "amount": '"$random_amount"',
            "external_reference": "'"$rand_int"'",
            "payment_method": "pix"
        }')

    body=$(echo "$response" | sed -e 's/HTTPSTATUS\:.*//g')
    status_code=$(echo "$response" | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')

    print_response "$status_code" "$body"
done
