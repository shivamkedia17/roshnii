# Roshnii - Backend Microservices

This repository contains the backend microservices for Roshnii, a cloud-native photo storage and management service. The project utilizes a microservices architecture implemented in Go, designed for scalability and cloud deployment.

## Overview

Roshnii aims to provide a secure and efficient platform for users to store, organize, and share their photos. This backend system handles user authentication, image metadata management, and provides the foundation for future features like albums, tagging, sharing, and ML-based image analysis.

The project emphasizes cloud computing principles, leveraging containerization with Docker and demonstrating patterns suitable for deployment on platforms like Kubernetes (GKE).

## Features

### Implemented
*   **Microservices Architecture:** Initial decomposition into `server` (API Gateway/BFF) and placeholder `faces` services.
*   **Authentication:**
    *   Google OAuth 2.0 for primary user login.
    *   JWT (JSON Web Tokens) for stateless session management via secure cookies.
    *   Development-only login endpoint (`/api/auth/dev/login`) for easier testing.
*   **Image Metadata API:**
    *   Upload image metadata (POST `/api/upload`).
    *   List user's image metadata (GET `/api/images`).
    *   Get specific image metadata (GET `/api/image/:id`).
    *   *(Note: Actual file storage is not yet implemented; currently handles metadata only).*
*   **Database:** PostgreSQL schema for users and image metadata.
*   **Containerization:** Dockerfile for building service images and Docker Compose configurations for local development (`docker-compose.yml`) and testing (`docker-compose.test.yml`).
*   **API Specification:** OpenAPI 3.0 definition (`api/openapi.yaml`).
*   **Testing:**
    *   Integration tests for the `server` service interacting with the database.
    *   Basic API endpoint testing script (`test_api.sh`).

### Planned / Future Work
*   Actual image file storage implementation (e.g., Google Cloud Storage, AWS S3).
*   Implementation of Album, Tag, Share, and Search features/APIs.
*   Development of the `faces` microservice logic (e.g., background processing, ML tasks).
*   Asynchronous communication between services (e.g., using a message queue like Pub/Sub or RabbitMQ).
*   Robust CI/CD pipeline for automated testing and deployment.
*   Comprehensive logging, monitoring, and alerting setup.
*   Enhanced security measures (rate limiting, advanced validation).

## Architecture Overview

The backend follows a microservices pattern:

*   **`server` service:** Acts as the main entry point, handling API requests, authentication, and core metadata management by interacting with the database.
*   **`faces` service:** (Placeholder) Intended for future background/ML tasks, decoupled from the main request flow.
*   **PostgreSQL:** Relational database for storing user information and image metadata.
*   **Docker/Docker Compose:** Used for containerization and local environment orchestration.

*(Refer to the detailed project report for architectural diagrams and in-depth discussion).*

## Technology Stack

*   **Language:** Go (Golang) 1.23+
*   **Framework:** Gin (Web Framework)
*   **Database:** PostgreSQL 15+
*   **Containerization:** Docker & Docker Compose
*   **Configuration:** Viper, `.env` files
*   **Authentication:** Google OAuth 2.0, JWT (golang-jwt/jwt/v5)
*   **API Specification:** OpenAPI 3.0
*   **Testing:** Go standard testing library, `stretchr/testify`, `net/http/httptest`

## Prerequisites

Before you begin, ensure you have the following installed:

*   **Docker:** Latest version (\url{https://docs.docker.com/get-docker/})
*   **Docker Compose:** Usually included with Docker Desktop, otherwise install separately (\url{https://docs.docker.com/compose/install/})
*   **Go:** Version 1.23 or higher (only needed if running tests outside Docker or modifying code) (\url{https://go.dev/doc/install})
*   **Git:** For cloning the repository (\url{https://git-scm.com/downloads})
*   **curl & jq:** Command-line tools used by the API test script (`test_api.sh`).
*   **Google OAuth Credentials:** You need to create credentials in the Google Cloud Console (\url{https://console.cloud.google.com/apis/credentials}) for a Web Application:
    *   Client ID
    *   Client Secret
    *   Set the **Authorized redirect URIs** to `http://localhost:8080/api/auth/google/callback` (or adjust host/port if your setup differs).

## Getting Started & Local Setup

1.  **Clone the Repository:**
    ```bash
    git clone <your-repo-url>
    cd roshnii/backend # Navigate to the backend directory
    ```

2.  **Create Configuration File (`.env`):**
    Copy the example environment file provided for the `server` service to the project root (`backend/`) and name it `.env`. This file will be used by Docker Compose for both development and testing setups.
    ```bash
    cp services/server/cmd/app.env .env
    ```

3.  **Edit `.env` File:**
    Open the newly created `.env` file in the `backend/` directory and **update the variables** according to your environment. **This step is crucial.**

    *   **Development Database:** (Used by `docker-compose.yml`)
        *   `POSTGRES_USER`: Username for the development database (e.g., `roshnii`).
        *   `POSTGRES_PASSWORD`: Password for the development database (e.g., `abcd1234`).
        *   `POSTGRES_DB`: Name of the development database (e.g., `roshnii_db`).
        *   `POSTGRES_URL`: The connection string the Go application will use to connect to the *development* database container. **Important:** Use the service name defined in `docker-compose.yml` as the host (default is `db`).
            ```
            # Example for docker-compose.yml
            POSTGRES_URL=postgresql://roshnii:abcd1234@db:5432/roshnii_db?sslmode=disable
            ```

    *   **Test Database:** (Used by `docker-compose.test.yml`)
        *   `TEST_POSTGRES_USER`: Username for the *test* database (e.g., `roshnii_test`).
        *   `TEST_POSTGRES_PASSWORD`: Password for the *test* database (e.g., `testpassword`).
        *   `TEST_POSTGRES_DB`: Name of the *test* database (e.g., `roshnii_test_db`).
        *   `TEST_POSTGRES_URL`: The connection string used during integration tests to connect to the *test* database container. Use the service name `test-db` from `docker-compose.test.yml` as the host.
            ```
            # Example for docker-compose.test.yml
            TEST_POSTGRES_URL=postgresql://roshnii_test:testpassword@test-db:5432/roshnii_test_db?sslmode=disable
            ```

    *   **Authentication:**
        *   `GOOGLE_CLIENT_ID`: Your Google OAuth Client ID.
        *   `GOOGLE_CLIENT_SECRET`: Your Google OAuth Client Secret.
        *   `JWT_SECRET`: A strong, random secret key for signing JWTs. Generate one using:
            ```bash
            openssl rand -base64 32
            ```
        *   `TOKEN_DURATION`: How long JWT sessions are valid (e.g., `24h`).

    *   **Application URLs:**
        *   `FRONTEND_URL`: The URL of your frontend application (used for redirect after OAuth login, e.g., `http://localhost:3000`).
        *   `SERVER_HOST`: Host the Go server listens on *inside* the container (usually `0.0.0.0`).
        *   `SERVER_PORT`: Port the Go server listens on (e.g., `8080`).

    *   **Other:**
        *   `ENVIRONMENT`: Set to `development` for local running/testing to enable the dev login endpoint.

## Running the Application (Local Development)

1.  **Ensure Docker is running.**
2.  **Ensure you have created and correctly configured the `.env` file in the `backend/` directory.**
3.  **Start the services using Docker Compose:**
    ```bash
    docker-compose up --build
    ```
    *   `--build`: Forces Docker Compose to build the images (e.g., `backend`) if they don't exist or if the `Dockerfile` or context has changed.
    *   This command will:
        *   Start the PostgreSQL database container (`db`).
        *   Initialize the database schema using `db/schema.sql` on the first run.
        *   Build the `backend` service image.
        *   Start the `backend` service container, connecting it to the `db` container.
        *   Show logs from both containers in your terminal.

4.  **Access the Service:**
    The backend service should now be running and accessible at `http://localhost:8080`.
    You can check the health endpoint:
    ```bash
    curl http://localhost:8080/health
    ```
    You should see: `{"status":"UP"}`

5.  **Stopping the Services:**
    Press `Ctrl+C` in the terminal where `docker-compose up` is running. To remove the containers and network (but preserve the database volume):
    ```bash
    docker-compose down
    ```
    To remove the database volume as well (lose all dev data):
    ```bash
    docker-compose down -v
    ```

## Running Tests

There are multiple ways to run the tests:

### 1. Integration Tests (Recommended - via Docker Compose)

This is the most reliable method as it runs tests in an isolated containerized environment with a dedicated test database.

1.  **Ensure Docker is running.**
2.  **Ensure the `.env` file in the `backend/` directory is configured correctly, especially the `TEST_POSTGRES_*` variables.**
3.  **Run the test suite using the test-specific Docker Compose file:**
    ```bash
    docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
    ```
    *   `-f docker-compose.test.yml`: Specifies the test environment configuration.
    *   `--build`: Builds necessary images (like the test runner if needed, though often uses base Go).
    *   `--abort-on-container-exit`: Stops all containers when the test runner container (`backend-test-runner`) finishes, making it easy to see the final test results.

4.  **Cleanup Test Environment:** After the tests complete, tear down the test containers and network:
    ```bash
    docker-compose -f docker-compose.test.yml down
    ```

### 2. Integration Tests (Alternative - Local Go + Docker DB)

This method runs the Go tests directly on your machine but still uses Docker for the test database.

1.  **Ensure Docker is running.**
2.  **Ensure the `.env` file is configured, particularly `TEST_POSTGRES_URL`.** Make sure the hostname/port in `TEST_POSTGRES_URL` is accessible from your local machine (e.g., `localhost:5433` if using the default port mapping in `docker-compose.test.yml`).
3.  **Start *only* the test database container:**
    ```bash
    docker-compose -f docker-compose.test.yml up -d test-db
    ```
    *   `-d`: Runs the container in detached mode.
4.  **Ensure Go (1.23+) is installed locally.**
5.  **Navigate to the server command directory and run tests:**
    ```bash
    cd services/server/cmd
    go test -v ./...
    cd ../../.. # Go back to backend root
    ```
6.  **Stop the test database container when finished:**
    ```bash
    docker-compose -f docker-compose.test.yml down
    ```

### 3. API Tests (Shell Script)

This script performs black-box checks against a *running* development instance of the API.

1.  **Ensure the application is running locally using `docker-compose up`.**
2.  **Ensure `curl` and `jq` are installed.**
3.  **Run the script from the `backend/` directory:**
    ```bash
    ./test_api.sh
    ```
    The script will attempt to use the development login, upload/list image metadata, etc.

## Configuration

Application configuration is managed via environment variables loaded using Viper.

*   A `.env` file should be placed in the `backend/` root directory for local development and testing (used by Docker Compose).
*   See `services/server/cmd/app.env` for an example and list of all available variables.
*   See `shared/pkg/config/config.go` for how configuration is loaded and defaults are set.
*   **Key variables to set in `.env`:** `POSTGRES_URL`, `TEST_POSTGRES_URL`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `TEST_POSTGRES_USER`, `TEST_POSTGRES_PASSWORD`, `TEST_POSTGRES_DB`, `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `JWT_SECRET`, `FRONTEND_URL`, `ENVIRONMENT`.

## API Documentation

The API is documented using the OpenAPI 3.0 standard.

*   The specification file is located at: `api/openapi.yaml`.
*   You can use tools like the [Swagger Editor](https://editor.swagger.io/) or [Swagger UI](https://swagger.io/tools/swagger-ui/) to view and interact with the API definition. Load the `openapi.yaml` file into these tools.

## License

*(Choose and add a license, e.g., MIT or Apache 2.0)*
