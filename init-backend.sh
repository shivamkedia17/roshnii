#!/bin/bash
# setup.sh - Initialize a multi-microservice monorepo folder structure
#
# Usage:
#   ./setup.sh                => Sets up the base folder structure with example microservices (orders and users)
#   ./setup.sh -n serviceName => Sets up the base structure (if not already created) and adds a new service "serviceName"
#
# This script creates the following structure:
#
# microservices-monorepo/
# ├── services/
# │   ├── orders/           (example service)
# │   └── users/            (example service)
# ├── shared/
# │   └── pkg/
# ├── deployments/
# │   └── k8s/
# ├── docs/
# └── README.md
#
# For each microservice, a standard structure is created:
#  ├── cmd/<service-service>/main.go
#  ├── internal/api/
#  ├── internal/domain/
#  ├── internal/service/
#  ├── internal/repository/
#  ├── pkg/
#  ├── Dockerfile
#  ├── go.mod
#  └── README.md

set -e

BASE_DIR="backend"

# Create base directories if they don't exist
function create_base_structure() {
  echo "Creating base folder structure..."

  mkdir -p "$BASE_DIR"/services
  mkdir -p "$BASE_DIR"/shared/pkg

  # Create top-level README.md if it doesn't exist
  if [ ! -f "$BASE_DIR/README.md" ]; then
    echo "# Microservices Monorepo" > "$BASE_DIR/README.md"
  fi

  echo "Base structure created under $BASE_DIR"
}

# Create a sample microservice structure within services/<serviceName>
function create_microservice() {
  local serviceName="$1"
  local serviceDir="$BASE_DIR/services/$serviceName"

  if [ -d "$serviceDir" ]; then
    echo "Microservice '$serviceName' already exists in $serviceDir"
    return
  fi

  echo "Creating microservice: $serviceName"
  mkdir -p "$serviceDir"/cmd/
  mkdir -p "$serviceDir"/internal/api
  mkdir -p "$serviceDir"/internal/domain
  mkdir -p "$serviceDir"/internal/service
  mkdir -p "$serviceDir"/pkg

  # Create sample main.go in cmd/
  cat <<EOF > "$serviceDir/cmd/main.go"
package main

import "fmt"

func main() {
    fmt.Println("Starting ${serviceName} microservice...")
    // Initialize your service here
}
EOF

  # Create a placeholder Dockerfile
  cat <<EOF > "$serviceDir/Dockerfile"
# Dockerfile for the ${serviceName} microservice
FROM golang:1.20-alpine
WORKDIR /app
COPY . .
RUN go build -o ${serviceName}-service ./cmd
CMD [ "./${serviceName}-service" ]
EOF

  # Create a basic go.mod file
  cat <<EOF > "$serviceDir/go.mod"
module github.com/your-org/${serviceName}

go 1.20
EOF

  # Create a README.md for this microservice
  cat <<EOF > "$serviceDir/README.md"
# ${serviceName} Microservice

This directory holds the code for the ${serviceName} microservice.
EOF

  echo "Microservice '$serviceName' created successfully."
}

# Parse command line options
SERVICE_TO_CREATE=""

usage() {
  echo "Usage: $0 [-n serviceName]"
  echo "  -n, --new serviceName   Create a new microservice with the given name"
  exit 1
}

# Process arguments
while [[ "$#" -gt 0 ]]; do
  case $1 in
    -n|--new)
      if [[ -n "$2" ]]; then
        SERVICE_TO_CREATE="$2"
        shift
      else
        echo "Error: Missing service name after $1"
        usage
      fi
      ;;
    -*)
      echo "Unknown option: $1"
      usage
      ;;
    *)
      echo "Unknown argument: $1"
      usage
      ;;
  esac
  shift
done

# Execute the functions to build the structure
create_base_structure

# Create example microservices (orders and users) if they don't exist
[[ -d "$BASE_DIR/services/server" ]] || create_microservice "server"
[[ -d "$BASE_DIR/services/faces" ]] || create_microservice "faces"

# Optionally create an additional microservice if provided by the user
if [ -n "$SERVICE_TO_CREATE" ]; then
  create_microservice "$SERVICE_TO_CREATE"
fi

echo "Setup completed successfully."
