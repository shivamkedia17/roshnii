# Roshnii - Cloud Photo Management Application

Roshnii is a cloud-native photo storage and management service allowing users to store, organize, and share their photos. This application aims to provide a feature-rich personal photo gallery with a modern, responsive interface.

## Prerequisites

- Docker and Docker Compose
- Go 1.23 or higher (for development)
- Bun (for frontend development)
- Google OAuth credentials (see below on how to setup)

## Getting Started (Takes 5-7 mins)

### Creating OAuth Credentials

1. Create a Google OAuth Client ID and Secret from Google's Cloud Console:
   - APIs & Services > Credentials > Create Credentials > OAuth Client ID
   - Application Type > Web Application
2. Configure Allowed Domains:
   - Authorized JavaScript Origins: "http://127.0.0.1:8080"
   - Authorized Redirect URIs: "http://127.0.0.1:8080/api/auth/google/callback"

### Configuration

1. Clone the repository. Ensure you are in the `backend/` directory.
2. Create a `.env` file in the `backend/` directory with the following configuration:

```env
ENVIRONMENT=development

# Database settings
POSTGRES_USER=roshnii
POSTGRES_PASSWORD=your_password   # change as needed
POSTGRES_DB=roshnii_db            
DB_HOST=127.0.0.1
DB_PORT=5432
POSTGRES_URL=postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${DB_HOST}:${DB_PORT}/${POSTGRES_DB}?sslmode=disable

# Google OAuth Credentials
GOOGLE_CLIENT_ID=your_google_client_id.apps.googleusercontent.com      # change as needed
GOOGLE_CLIENT_SECRET=your_google_client_secret                         # change as needed 
TOKEN_DURATION=24h
JWT_SECRET=your_jwt_secret                                             # change as needed, see instructions in backend/README.md to generate 

# Frontend URL (for redirects after OAuth)
FRONTEND_URL=http://127.0.0.1:8080

# Server settings
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

3. Build the frontend app:
```sh
bun run build
```

4. Build the backend server service:
```sh
go build ./services/server/cmd/main.go
```

5. Run the database using docker compose
```sh
docker compose up --build
```

6. Start the backend server microservice:
```sh
./main
``` 

---

## Features

- **Modern User Interface**: Clean, responsive design for managing your photo collection
- **Secure Authentication**: Google OAuth 2.0 integration for seamless and secure login
- **Photo Management**:
  - Upload and store photos in various formats (JPEG, PNG, GIF, WebP)
  - View photo details including metadata
  - Organize photos in custom albums
  - Delete unwanted photos
- **Album Organization**: Create, manage, and delete photo albums
- **Multi-user Support**: Each user has their own private photo collection
- **Responsive Design**: Works on desktop and mobile devices

## Technology Stack

### Backend
- **Language**: Go (Golang) 1.23+
- **Web Framework**: Gin
- **Database**: PostgreSQL
- **Authentication**: Google OAuth 2.0, JWT
- **Storage**: Local file storage (with support for cloud storage planned)
- **Containerization**: Docker & Docker Compose

### Frontend
- **Framework**: React + Vite
- **State Management**: React Query
- **Styling**: Custom CSS with responsive design

## Architecture

Roshnii follows a microservices architecture:

- **Server Service**: Handles API requests, authentication, and core data management
- **PostgreSQL**: Stores user information, image metadata, and album relationships
- **Frontend**: React application for user interaction

