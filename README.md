# GitHub Actions Aggregator Service

A Go-based service to aggregate and analyze data from GitHub Actions workflows across multiple repositories. This application provides insights into workflow runs, success rates, failure rates, and other statistics over customizable time ranges.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Authentication](#authentication)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## Features

- **OAuth 2.0 Authentication**: Securely authenticate users via GitHub OAuth.
- **Data Collection**:
  - **Webhooks**: Receive real-time updates on workflow events.
  - **Polling**: Periodically poll GitHub API to ensure data completeness.
- **Data Aggregation**: Compute statistics like success rates, failure rates, and more.
- **API Endpoints**: Expose RESTful APIs for accessing aggregated data.
- **Background Processing**: Use worker pools to handle asynchronous tasks.
- **Configurable**: Easily adjust settings like polling intervals and webhook secrets.
- **Secure**: Validate webhook payloads and protect routes with authentication middleware.

## Prerequisites

- **Go**: Version 1.18 or higher.
- **GitHub Account**: For OAuth authentication and API access.
- **PostgreSQL**: For storing data.
- **Redis** (optional): For caching (if implemented).
- **Docker** (optional): For containerization and deployment.

## Installation

1. **Clone the Repository**

   ```bash
   git clone https://github.com/yourusername/github-actions-aggregator.git
   cd github-actions-aggregator
   ```

2. **Install Dependencies**

   ```bash
   go mod download
   ```

3. **Set Up Environment Variables**

   Create a `.env` file or export the required environment variables:

   ```bash
   export GITHUB_CLIENT_ID="your_github_client_id"
   export GITHUB_CLIENT_SECRET="your_github_client_secret"
   export GITHUB_ACCESS_TOKEN="your_github_access_token"
   export GITHUB_WEBHOOK_SECRET="your_webhook_secret"
   export DATABASE_URL="postgres://username:password@localhost:5432/yourdbname?sslmode=disable"
   export SERVER_PORT="8080"
   ```

## Configuration

Configuration can be managed via a `config.yaml` file in the `configs/` directory or through environment variables.

**Example `config.yaml`:**

```yaml
server:
  port: "8080"

log:
  level: "info"

github:
  client_id: "your_github_client_id"
  client_secret: "your_github_client_secret"
  access_token: "your_github_access_token"
  webhook_secret: "your_webhook_secret"
```

**Note:** Environment variables override values in the configuration file.

## Usage

### Running the Application

1. **Run Database Migrations**

   ```bash
   ./scripts/migrate.sh
   ```

2. **Start the Application**

   ```bash
   go run cmd/server/main.go
   ```

### Accessing the Application

- **Login with GitHub**: Navigate to `http://localhost:8080/login` to authenticate via GitHub.
- **API Requests**: Use tools like `curl` or Postman to interact with the API endpoints.

## API Endpoints

### Authentication

- `GET /login`: Redirects the user to GitHub for OAuth authentication.
- `GET /callback`: Handles the OAuth callback from GitHub.

### Workflow Statistics

- `GET /workflows/:id/stats`: Retrieves statistics for a specific workflow.

  **Query Parameters:**

  - `start_time` (optional): Start of the time range (ISO 8601 format).
  - `end_time` (optional): End of the time range (ISO 8601 format).

  **Example Request:**

  ```http
  GET /workflows/123/stats?start_time=2023-09-01T00:00:00Z&end_time=2023-09-30T23:59:59Z
  ```

  **Example Response:**

  ```json
  {
    "workflow_id": 123,
    "workflow_name": "CI Build and Test",
    "total_runs": 200,
    "success_count": 150,
    "failure_count": 30,
    "cancelled_count": 10,
    "timed_out_count": 5,
    "action_required_count": 5,
    "success_rate": 75.0,
    "failure_rate": 15.0,
    "cancelled_rate": 5.0,
    "timed_out_rate": 2.5,
    "action_required_rate": 2.5,
    "start_time": "2023-09-01T00:00:00Z",
    "end_time": "2023-09-30T23:59:59Z"
  }
  ```

## Authentication

### Setting Up OAuth with GitHub

1. **Register a New OAuth Application**

   - Go to [GitHub Developer Settings](https://github.com/settings/developers).
   - Click on **"New OAuth App"**.
   - Fill in the application details:
     - **Application Name**
     - **Homepage URL**: `http://localhost:8080`
     - **Authorization Callback URL**: `http://localhost:8080/callback`
   - Obtain your **Client ID** and **Client Secret**.

2. **Configure Application Credentials**

   Set your `GITHUB_CLIENT_ID` and `GITHUB_CLIENT_SECRET` in your environment variables or `config.yaml`.

### Permissions and Scopes

Ensure that your GitHub OAuth application has the necessary scopes:

- `read:user`
- `repo`
- `workflow`

## Testing

### Running Unit Tests

```bash
go test ./tests/unit/...
```

### Running Integration Tests

```bash
go test ./tests/integration/...
```

### Test Coverage

You can generate a test coverage report using:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Contributing

Contributions are welcome! Please follow these steps:

1. **Fork the Repository**

   Click on the "Fork" button at the top right of the repository page.

2. **Clone Your Fork**

   ```bash
   git clone https://github.com/yourusername/github-actions-aggregator.git
   ```

3. **Create a Feature Branch**

   ```bash
   git checkout -b feature/your-feature-name
   ```

4. **Commit Your Changes**

   ```bash
   git commit -am "Add new feature"
   ```

5. **Push to Your Fork**

   ```bash
   git push origin feature/your-feature-name
   ```

6. **Create a Pull Request**

   Go to the original repository and open a pull request.

## License

This project is licensed under the [MIT License](LICENSE).

---

**Disclaimer:** This project is not affiliated with GitHub. Ensure compliance with GitHub's [Terms of Service](https://docs.github.com/en/github/site-policy/github-terms-of-service) when using their APIs.