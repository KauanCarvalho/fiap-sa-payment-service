#!/bin/bash

# Function to check if required dependencies are installed
check_dependencies() {
    # Check if jq is installed
    if ! command -v jq &> /dev/null
    then
        echo -e "\033[31mError: jq is not installed. Please install jq to proceed.\033[0m"
        exit 1
    fi
}

# Function to print status and body with color-coded output
print_response() {
    local status_code=$1
    local body=$2

    if [[ "$status_code" -ge 200 && "$status_code" -lt 300 ]]; then
        echo -e "\033[32mStatus: $status_code\033[0m"  # Green for 200-299
    elif [[ "$status_code" -ge 400 && "$status_code" -lt 500 ]]; then
        echo -e "\033[34mStatus: $status_code\033[0m"  # Blue for 400-499
    else
        echo -e "\033[31mStatus: $status_code\033[0m"  # Red for all others
    fi

    if [[ -z "$body" ]]; then
        echo -e "\033[33mNo Body\033[0m"
    else
        echo -e "$body" | jq .  # Pretty print the JSON body
    fi
    echo "---------------------------------------------"
}

# Function to generate random number between a given range
generate_random_number() {
    echo $((RANDOM % (999999999999 - 1000 + 1) + 1000))
}

# Function to check the health status of the service
check_health() {
    echo "Checking healthcheck..."
    health_response=$(curl --silent --write-out "HTTPSTATUS:%{http_code}" --location "{{base_url}}/healthcheck")
    health_status=$(echo "$health_response" | sed -e 's/HTTPSTATUS\:.*//g')
    health_body=$(echo "$health_response" | sed -e 's/HTTPSTATUS\:.*//g')

    print_response "$health_status" "$health_body"

    if [[ "$health_status" -ne 200 ]]; then
        echo "Healthcheck failed, exiting..."
        exit 1
    fi
}

# Main script starts here

# Check if the required dependencies are installed
check_dependencies

# Perform healthcheck
check_health

# Loop to create 5 payments
for i in {1..5}; do
    random_amount=$((RANDOM % 1000000 + 1))  # Random amount between 1 and 1000000
    rand_int=$(generate_random_number)

    echo "Creating payment #$i"
    
    # Make the payment authorization API call
    response=$(curl --silent --write-out "HTTPSTATUS:%{http_code}" --location "{{base_url}}/api/v1/payments/authorize" \
        --header 'Content-Type: application/json' \
        --data '{
            "amount": '"$random_amount"',
            "external_reference": "'"$rand_int"'",
            "payment_method": "pix"
        }')

    status_code=$(echo "$response" | sed -e 's/HTTPSTATUS\:.*//g')
    body=$(echo "$response" | sed -e 's/HTTPSTATUS\:.*//g')

    print_response "$status_code" "$body"
done
