# Search Service

GraphQL gateway that exposes flexible search over TeamUp groups. It reuses the existing `app-service` group API and adds title / tag / date filtering via GraphQL.

## Configuration

| Variable | Description | Default |
| --- | --- | --- |
| `SEARCH_SERVICE_PORT` | HTTP port for the GraphQL server. | `4000` |
| `GROUP_SERVICE_URL` | Base URL of the group service REST API. | `http://app-service:3001` |
| `GROUP_SERVICE_TIMEOUT_SECONDS` | HTTP client timeout when calling the group service. | `5` |
| `ENABLE_GRAPHQL_PLAYGROUND` | Enables the built-in GraphQL Playground UI when `true`. | `true` |
| `DEFAULT_SEARCH_LIMIT` | Fallback size for search responses. | `20` |
| `MAX_SEARCH_LIMIT` | Upper bound for the `limit` argument. | `100` |

## Running locally

```bash
SEARCH_SERVICE_PORT=4000 \
GROUP_SERVICE_URL="http://localhost:3001" \
go run ./cmd/main.go
```

Then open [http://localhost:4000/graphql](http://localhost:4000/graphql) to access the Playground.

## Sample query

```graphql
query SearchGroups {
  searchGroups(title: "study", tags: [STUDY, PROJECT], includeClosed: false, limit: 5) {
    id
    title
    tags
    closed
    date
  }
}
```
