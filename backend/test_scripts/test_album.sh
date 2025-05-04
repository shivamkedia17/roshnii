# Setup: Login & get token
TOKEN=$(curl -s -X POST -H "Content-Type: application/json" -d '{"email":"test@example.com","name":"Test User"}' http://localhost:8080/api/auth/dev/login | jq -r '.token')

# Create an album
curl -s -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"My Vacation Photos","description":"Photos from my recent trip"}' \
  http://localhost:8080/api/albums | jq

# List albums
curl -s -X GET \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/albums | jq

# Get a specific album (replace 1 with actual album ID)
curl -s -X GET \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/albums/1 | jq

# Update an album
curl -s -X PUT \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"My Amazing Vacation","description":"Updated description"}' \
  http://localhost:8080/api/albums/1 | jq

# Add an image to an album (replace with actual image ID)
curl -s -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"image_id":"abc123-uuid-here"}' \
  http://localhost:8080/api/albums/1/images | jq

# List images in an album
curl -s -X GET \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/albums/1/images | jq

# Remove an image from an album
curl -s -X DELETE \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/albums/1/images/abc123-uuid-here | jq

# Delete an album
curl -s -X DELETE
