#!/bin/bash
# update_structure.sh
#
# This script sets up a folder structure for a microservices-based
# API backend. It creates separate cmd/ folders for each microservice,
# creates domain-specific folders within internal/, and an optional pkg/ folder.
#
# Run this script from the backend/ directory.
#
# NOTE: This script creates directories if they do not exist,
# and moves a few sample files to their new locations (if found).
#
# Always back up your code before running a reorganizing script.
#
set -e

echo "Creating microservices folder structure..."

# Create new service executables under cmd/
mkdir -p cmd/auth-service
mkdir -p cmd/photo-service
mkdir -p cmd/upload-service

# (Optional) If you have an old monolithic server in cmd/server, you might want to move its files.
if [ -d "cmd/server" ]; then
  echo "Moving existing cmd/server files to internal/server/ ..."
  mkdir -p internal/server
  # Move all .go files from cmd/server to internal/server
  mv cmd/server/*.go internal/server/ 2>/dev/null || true
  # Remove the now-empty cmd/server folder
  rmdir cmd/server 2>/dev/null || true
fi

# Create internal directories for domain-specific logic
mkdir -p internal/auth       # For OAuth, JWT, and user identity code.
mkdir -p internal/handlers   # For HTTP handlers (auth_handler.go, photo_handler.go, upload_handler.go, etc.)
mkdir -p internal/services   # For business logic, e.g., JWT generation.
# internal/config, internal/database, internal/models, internal/server remain as-is.

# Create optional pkg folder for shared libraries/adapters
mkdir -p pkg/utils

echo "Folder structure has been updated for a microservices-based API architecture."

echo "Suggested manual actions:"
echo "  * Place the auth-service main entry point in cmd/auth-service/main.go"
echo "  * Place the photo-service main entry point in cmd/photo-service/main.go"
echo "  * Verify that internal/server/router.go contains shared routing or move per-service routers as desired."
echo "  * Update your build process to compile each service separately."
