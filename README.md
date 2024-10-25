# GitHub Actions Aggregator

GitHub Actions Aggregator is a Go-based service designed to collect, aggregate, and analyze data from GitHub Actions workflows across multiple repositories. This application provides valuable insights into workflow runs, success rates, failure rates, and other statistics over customizable time ranges.

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

- **OAuth 2.0 Authentication**: Secure user authentication via GitHub OAuth.
- **Data Collection**:
  - **Webhooks**: Receive real-time updates on workflow events.
  - **Polling**: Periodically poll GitHub API to ensure data completeness.
- **Data Aggregation**: Compute statistics like success rates, failure rates, and more.
- **API Endpoints**: Expose RESTful APIs for accessing aggregated data.
- **Background Processing**: Use worker pools to handle asynchronous tasks.
- **Configurable**: Easily adjust settings like polling intervals and webhook secrets.
- **Secure**: Validate webhook payloads and protect routes with authentication middleware.

## Prerequisites

- Go (version 1.18 or higher)
- PostgreSQL
- GitHub Account (for OAuth authentication and API access)
- Docker (optional, for containerization)

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/github-actions-aggregator.git
   cd github-actions-aggregator
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Set up environment variables:
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

Example `config.yaml`:

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

Note: Environment variables override values in the configuration file.

## Usage

1. Run database migrations:
   ```
   ./scripts/migrate.sh up
   ```

2. Start the application:
   ```
   go run cmd/server/main.go
   ```

3. Access the application:
   - Login with GitHub: Navigate to `http://localhost:8080/login`
   - API Requests: Use tools like `curl` or Postman to interact with the API endpoints

## API Endpoints

- `GET /login`: Redirects the user to GitHub for OAuth authentication.
- `GET /callback`: Handles the OAuth callback from GitHub.
- `GET /workflows/:id/stats`: Retrieves statistics for a specific workflow.
- `GET /repositories/:id/workflows`: Get all workflows for a repository
- `GET /workflows/:id/runs`: Get all runs for a workflow
- `GET /runs/:id`: Get a specific run
- `GET /jobs/:id`: Get a specific job
- `GET /jobs/:id/steps`: Get all steps for a job
- `GET /jobs/:id/stats`: Get stats for a job

For detailed information on request parameters and response formats, please refer to the API documentation.

## Authentication

1. Register a new OAuth application on GitHub:
   - Go to GitHub Developer Settings
   - Click on "New OAuth App"
   - Fill in the application details:
     - Application Name
     - Homepage URL: `http://localhost:8080`
     - Authorization Callback URL: `http://localhost:8080/callback`
   - Obtain your Client ID and Client Secret

2. Configure application credentials:
   Set your `GITHUB_CLIENT_ID` and `GITHUB_CLIENT_SECRET` in your environment variables or `config.yaml`.

Ensure that your GitHub OAuth application has the necessary scopes: `read:user`, `repo`, and `workflow`.

## Testing

Run unit tests:
```
go test ./tests/unit/...
```

Run integration tests:
```
go test ./tests/integration/...
```

Generate a test coverage report:
```
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to your fork
5. Create a pull request

For more details, please refer to the [CONTRIBUTING.md](CONTRIBUTING.md) file.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

**Disclaimer:** This project is not affiliated with GitHub. Ensure compliance with GitHub's [Terms of Service](https://docs.github.com/en/github/site-policy/github-terms-of-service) when using their APIs.