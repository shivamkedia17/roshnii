# Roshnii - Cloud Photo Management Application

Roshnii is a robust cloud-native photo storage and management service allowing users to securely store, organize, and share their photos. This application provides a feature-rich personal photo gallery with a modern, responsive interface.

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

## Prerequisites

- Docker and Docker Compose
- Go 1.23 or higher (for development)
- Bun (for frontend development)
- Google OAuth credentials

## Getting Started

### Configuration

1. Clone the repository
2. Create a `.env` file in the `backend/` directory with the following configuration:

```env
ENVIRONMENT=development

# Database settings
POSTGRES_USER=roshnii
POSTGRES_PASSWORD=your_password
POSTGRES_DB=roshnii_db
DB_HOST=127.0.0.1
DB_PORT=5432
POSTGRES_URL=postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${DB_HOST}:${DB_PORT}/${POSTGRES_DB}?sslmode=disable

# Google OAuth Credentials
GOOGLE_CLIENT_ID=your_google_client_id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your_google_client_secret
TOKEN_DURATION=24h
JWT_SECRET=your_jwt_secret

# Frontend URL (for redirects after OAuth)
FRONTEND_URL=http://127.0.0.1:8080

# Server settings
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

