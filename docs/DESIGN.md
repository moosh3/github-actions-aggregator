# Design

Instead of querying the GitHub API for each request, we will query our database. This will be faster and cheaper. At the organization level, a webhook is created that sends events to our server.

## API

### Endpoints

#### GET /repositories/:id/workflows

#### GET /workflows/:id/stats

#### GET /workflows/:id/runs

#### GET /runs/:id

#### GET /jobs/:id
