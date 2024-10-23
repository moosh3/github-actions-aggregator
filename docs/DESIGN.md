# Design

## API Design

### Endpoint Structure

#### Authentication Endpoints:

- `POST /login`: User login via GitHub OAuth.
- `POST /logout`: User logout.

#### Data Retrieval Endpoints:

- `GET /stats`: Retrieve aggregated statistics.
- `GET /workflows`: List workflows with statuses.
- `GET /workflows/:id`: Get details of a specific workflow.

#### Data Management Endpoints:

- `POST /repositories`: Add a repository to monitor.
- `DELETE /repositories/:id`: Remove a repository.

### API Features

- **Pagination**: Implement for endpoints returning lists.
- **Filtering**: Allow filtering by date ranges, status, etc.
- **Error Handling**: Provide meaningful error messages and HTTP status codes.
- **Throttling**: Implement rate limiting to prevent abuse.

## Data Collection Strategy

Option A: Polling the GitHub API

- **Scheduled Fetching**: Use a scheduler (like cron jobs) to periodically fetch data from the GitHub API.
- **Rate Limiting**: Implement logic to handle GitHub's API rate limits (e.g., using conditional requests with ETags).

Option B: GitHub Webhooks
- **Real-Time Updates**: Set up webhooks to receive events when workflows are triggered, succeed, or fail.
- **Event Handling**: Create endpoints to handle incoming webhook events securely.

Recommended Approach
- **Combine Both**: Use webhooks for real-time updates and periodic polling to handle missed events or initial data population.

## Authentication and Authorization

- **GitHub Authentication**: Use a GitHub App or Personal Access Token (PAT) to authenticate API requests.
- **User Authentication**: Implement OAuth 2.0 to allow users to log in with their GitHub accounts.
- **Permissions**: Ensure the service respects user permissions and data privacy.

## Tech Stack

- **Programming Language**: Golang for its performance and concurrency support.
- **Web Framework**: Use gin-gonic/gin for routing and middleware support.
- **Database**: PostgreSQL for relational data storage and complex queries.
- **ORM**: Use gorm for interacting with the database.
- **Caching**: Implement Redis for caching frequently accessed data.
- **HTTP Client**: Use net/http and oauth2 packages for making authenticated requests to GitHub.

## Architectural Design

### Backend Services

- **API Server**: Handles incoming HTTP requests from the frontend.
- **Worker Service**: Processes background tasks like polling the GitHub API or handling webhook events.

### Data Flow

- **Data Ingestion**: Collect data via webhooks or polling.
- **Data Storage**: Store raw data in the database.
- **Data Processing**: Calculate statistics and store them in a summarized form.
- **API Layer**: Expose endpoints for the frontend to retrieve data.

### Database Schema Design

- Tables:
    - users: Stores user information.
    - repositories: Tracks repositories being monitored.
    - workflows: Contains workflow definitions.
    - workflow_runs: Records individual workflow executions.
    - statistics: Stores precomputed stats for quick access.

- Relationships:
    - Establish foreign keys between tables for referential integrity.

- Indexing:
    - Index columns frequently used in queries (e.g., timestamps, status).

### Implementing Concurrency

- **Goroutines**: Use for handling multiple simultaneous tasks (e.g., processing webhook events).
- **Channels**: Communicate between goroutines safely.
- **Mutexes**: Ensure thread-safe operations when accessing shared resources.

### Security Measures

- **Secure Communication**: Use HTTPS for all client-server communication.
- **Input Validation**: Sanitize and validate all inputs to prevent injection attacks.
- **Secrets Management**: Store sensitive information (like API keys) securely, using environment variables or a secrets manager.
- **Webhook Security**: Validate GitHub webhook payloads using the shared secret.

### Testing Strategy

- **Unit Tests**: Test individual functions and methods.
- **Integration Tests**: Test interactions with the GitHub API and the database.
- **End-to-End Tests**: Simulate user interactions with the API.
- **Mocking**: Use interfaces and mocking libraries to simulate external services.

### Example Workflow

- **User Authentication**: A user logs in via GitHub OAuth and grants permissions.
- **Repository Selection**: The user selects repositories to monitor.
- **Data Collection Initiation**: The system sets up webhooks and begins polling if necessary.
- **Data Ingestion**: Workflow events are received and stored.
- **Data Processing**: Statistics are calculated and updated in the database.
- **Data Presentation**: The frontend requests data via the API and presents it to the user.
