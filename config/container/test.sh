#!/usr/bin/env sh

set -e

echo "DATABASE_URI: $DATABASE_URI"

curl -fsSL https://pgp.mongodb.com/server-7.0.asc | gpg --dearmor -o /usr/share/keyrings/mongodb-server-7.0.gpg
    echo "deb [ signed-by=/usr/share/keyrings/mongodb-server-7.0.gpg ] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/7.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-7.0.list
    sudo apt-get update
    sudo apt-get install -y mongodb-mongosh

# Espera o MongoDB estar disponível
until mongosh "$DATABASE_URI" --eval 'db.runCommand({ ping: 1 })'; do
    echo "Waiting for MongoDB..."
    sleep 1
done

echo "MongoDB is available!"

# Run tests
go test ./... -coverprofile=coverage.out -cover -p 1

# Run coverage check
go tool go-test-coverage --config=./.testcoverage.yml
