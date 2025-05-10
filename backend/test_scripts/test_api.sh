#!/bin/bash

# --- Configuration ---
BASE_URL="http://127.0.0.1:8080" # Adjust if your server runs elsewhere
API_URL="$BASE_URL/api"
# Use a test email for the dev login
TEST_EMAIL="testuser-$(date +%s)@example.com" # Add timestamp for uniqueness
echo "Using test email: $TEST_EMAIL"

# Store the token here
AUTH_TOKEN="KMpyTBhNN3fQRwIRquQIMeGiT/8NN0D/aKjVnJMvg2g="

# --- Helper Function to Make Authenticated Requests ---
# Usage: authenticated_request METHOD ENDPOINT [JSON_DATA]
authenticated_request() {
    local method="$1"
    local endpoint="$2"
    local data="$3"
    local args=()

    if [[ -z "$AUTH_TOKEN" ]]; then
        echo "Error: No AUTH_TOKEN set. Please login first."
        return 1
    fi

    args+=("-H" "Authorization: Bearer $AUTH_TOKEN")

    if [[ "$method" == "POST" && -n "$data" ]]; then
        args+=("-H" "Content-Type: application/json")
        args+=("-d" "$data")
    fi

    echo ">>> $method $API_URL$endpoint"
    curl -s -X "$method" "${args[@]}" "$API_URL$endpoint" | jq '.' # Pretty print JSON output
    echo "<<<"
    echo
}

# --- 1. Health Check (No Auth Needed) ---
echo "--- Testing Health Check ---"
curl -s "$BASE_URL/health" | jq '.'
echo

# --- 2. Development Login (Get Token) ---
echo "--- Testing Development Login ---"
login_response=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d "{\"email\": \"$TEST_EMAIL\", \"name\": \"Test User\"}" \
    "$API_URL/auth/dev/login")

echo "Login Response:"
echo "$login_response" | jq '.'

# Extract the token using jq
AUTH_TOKEN=$(echo "$login_response" | jq -r '.token')

if [[ -z "$AUTH_TOKEN" || "$AUTH_TOKEN" == "null" ]]; then
    echo "Error: Failed to get authentication token!"
    exit 1
else
    echo "Successfully obtained auth token."
fi
echo

# --- 3. Get Current User Info (/me) ---
echo "--- Testing GET /api/me ---"
authenticated_request GET "/me"
echo

# --- 4. List Initial Images (Should be Empty) ---
echo "--- Testing GET /api/images (Initial) ---"
authenticated_request GET "/images"
echo

# --- 5. Upload an Image ---
# Create a dummy image file for testing
echo "Creating dummy image file: test_image.jpg"
echo "This is a dummy JPEG file content." > test_image.jpg

echo "--- Testing POST /api/upload ---"
# Note: multipart/form-data upload with curl
# We still need the Authorization header
echo ">>> POST $API_URL/upload (with test_image.jpg)"
upload_response=$(curl -s -X POST \
    -H "Authorization: Bearer $AUTH_TOKEN" \
    -F "file=@test_image.jpg;type=image/jpeg" \
    "$API_URL/upload")

echo "Upload Response:"
echo "$upload_response" | jq '.'

# Extract the ID of the uploaded image
UPLOADED_IMAGE_ID=$(echo "$upload_response" | jq -r '.id')

if [[ -z "$UPLOADED_IMAGE_ID" || "$UPLOADED_IMAGE_ID" == "null" ]]; then
    echo "Error: Failed to get ID from image upload response!"
    # Clean up dummy file even on error
    rm test_image.jpg
    exit 1
else
    echo "Successfully uploaded image. ID: $UPLOADED_IMAGE_ID"
fi
echo

# Clean up dummy file
rm test_image.jpg

# --- 6. List Images Again (Should Contain the Uploaded One) ---
echo "--- Testing GET /api/images (After Upload) ---"
authenticated_request GET "/images"
echo

# --- 7. Get Specific Image Metadata ---
echo "--- Testing GET /api/image/$UPLOADED_IMAGE_ID ---"
if [[ -n "$UPLOADED_IMAGE_ID" ]]; then
    authenticated_request GET "/image/$UPLOADED_IMAGE_ID"
else
    echo "Skipping GET /image/:id because UPLOADED_IMAGE_ID is not set."
fi
echo

# --- 8. Test Get Non-Existent Image ---
echo "--- Testing GET /api/image/non-existent-uuid ---"
authenticated_request GET "/image/non-existent-uuid"
echo

echo "--- Testing Album Endpoints ---"
# Create an album
echo ">>> POST $API_URL/albums"
ALBUM_RESPONSE=$(curl -s -X POST \
    -H "Authorization: Bearer $AUTH_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"Test Album","description":"Created by test script"}' \
    "$API_URL/albums")

echo "Album Creation Response:"
echo "$ALBUM_RESPONSE" | jq '.'

# Extract the album ID
ALBUM_ID=$(echo "$ALBUM_RESPONSE" | jq -r '.id')

if [[ -z "$ALBUM_ID" || "$ALBUM_ID" == "null" ]]; then
    echo "Error: Failed to get album ID from creation response!"
    exit 1
else
    echo "Successfully created album with ID: $ALBUM_ID"
fi
echo

# List all albums
echo ">>> GET $API_URL/albums"
authenticated_request GET "/albums"
echo

# Get specific album
echo ">>> GET $API_URL/albums/$ALBUM_ID"
authenticated_request GET "/albums/$ALBUM_ID"
echo

# Add image to album (assuming we have an uploaded image)
if [[ -n "$UPLOADED_IMAGE_ID" ]]; then
    echo ">>> POST $API_URL/albums/$ALBUM_ID/images"
    authenticated_request POST "/albums/$ALBUM_ID/images" "{\"image_id\":\"$UPLOADED_IMAGE_ID\"}"
    echo

    # List images in album
    echo ">>> GET $API_URL/albums/$ALBUM_ID/images"
    authenticated_request GET "/albums/$ALBUM_ID/images"
    echo

    # Remove image from album
    echo ">>> DELETE $API_URL/albums/$ALBUM_ID/images/$UPLOADED_IMAGE_ID"
    authenticated_request DELETE "/albums/$ALBUM_ID/images/$UPLOADED_IMAGE_ID"
    echo
fi

# Update album
echo ">>> PUT $API_URL/albums/$ALBUM_ID"
authenticated_request PUT "/albums/$ALBUM_ID" "{\"name\":\"Updated Test Album\",\"description\":\"Updated by test script\"}"
echo

# Delete album
echo ">>> DELETE $API_URL/albums/$ALBUM_ID"
authenticated_request DELETE "/albums/$ALBUM_ID"
echo


# --- 9. Test Logout (Placeholder - Dev login doesn't use cookies easily) ---
# The dev login gives a token directly. The /logout endpoint primarily clears
# the cookie. We can call it, but it won't affect our $AUTH_TOKEN variable.
# In a real scenario using cookies, this would invalidate the session.
echo "--- Testing POST /api/auth/google/logout (Informational) ---"
authenticated_request POST "/auth/google/logout"
echo "Logout called. Note: This doesn't clear the token stored in this script."
echo

# --- End of Tests ---
echo "--- API Tests Completed ---"
