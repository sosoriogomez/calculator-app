# Full-Stack Calculator

A small calculator application with a Go REST API and a React + TypeScript frontend. The API owns calculation rules and validation; the browser adds immediate input validation and presents API errors clearly.

## Project Layout

```text
backend/
  cmd/server/                 Composition root and HTTP server entry point
  internal/domain/calculation Domain operations, value object, aggregate, and tests
  internal/application/       Calculate use case and application tests
  internal/interfaces/httpapi HTTP transport adapter and request/response tests
frontend/
  src/api/                    Typed API client
  src/components/             Calculator UI and component tests
  src/test/                    Test setup
docker-compose.yml            One-command production-style setup
```

## Prerequisites

- Go 1.22 or newer
- Node.js 22 or newer and npm
- Docker and Docker Compose (optional)

## Run Locally

Start the backend:

```bash
cd backend
go run ./cmd/server
```

In another terminal, start the frontend:

```bash
cd frontend
npm install
npm run dev
```

Open <http://localhost:5173>. Vite proxies `/api` requests to the Go server at `localhost:8080`.

## Run With Docker

From the repository root:

```bash
docker compose up --build
```

Open <http://localhost:3000>. The frontend container serves the built React app and proxies API requests to the backend container.

## Swagger Documentation

The OpenAPI specification is available at `GET /openapi.yaml`. When running locally, open <http://localhost:8080/docs> for Swagger UI. With Docker Compose, open <http://localhost:3000/docs>.

The Swagger UI loads its presentation assets from `unpkg.com`; the raw specification remains available locally even without that CDN connection at `/openapi.yaml`.

## API

The API uses one resource-oriented endpoint because every operation has the same lifecycle: submit operands, validate, calculate, and receive a result. The operation is explicit in the JSON resource rather than encoded into seven nearly identical URL paths.

`POST /api/v1/calculations`

All requests use `Content-Type: application/json` and return JSON. `a` is always required. `b` is required except for square root. Percentage means "what percentage is `a` of `b`?", calculated as `(a / b) * 100`.

### Examples

```bash
# Addition
curl -X POST http://localhost:8080/api/v1/calculations \
  -H 'Content-Type: application/json' \
  -d '{"operation":"add","a":12,"b":8}'
# {"operation":"add","result":20}

# Subtraction
curl -X POST http://localhost:8080/api/v1/calculations -H 'Content-Type: application/json' \
  -d '{"operation":"subtract","a":12,"b":8}'
# {"operation":"subtract","result":4}

# Multiplication
curl -X POST http://localhost:8080/api/v1/calculations -H 'Content-Type: application/json' \
  -d '{"operation":"multiply","a":12,"b":8}'
# {"operation":"multiply","result":96}

# Division
curl -X POST http://localhost:8080/api/v1/calculations -H 'Content-Type: application/json' \
  -d '{"operation":"divide","a":12,"b":8}'
# {"operation":"divide","result":1.5}

# Exponentiation
curl -X POST http://localhost:8080/api/v1/calculations -H 'Content-Type: application/json' \
  -d '{"operation":"power","a":2,"b":8}'
# {"operation":"power","result":256}

# Square root
curl -X POST http://localhost:8080/api/v1/calculations -H 'Content-Type: application/json' \
  -d '{"operation":"sqrt","a":81}'
# {"operation":"sqrt","result":9}

# Percentage
curl -X POST http://localhost:8080/api/v1/calculations -H 'Content-Type: application/json' \
  -d '{"operation":"percentage","a":25,"b":200}'
# {"operation":"percentage","result":12.5}
```

Health check: `GET /healthz` returns `{"status":"ok"}`.

Invalid requests use a consistent shape:

```json
{
  "error": {
    "code": "calculation_error",
    "message": "division by zero"
  }
}
```

Malformed JSON, missing fields, wrong types, and unknown operations return `400 Bad Request`. Valid JSON with an undefined or non-finite calculation returns `422 Unprocessable Entity`. The API rejects division by zero, square roots of negative values, `0^0`, zero raised to a negative exponent, percentage calculations with a zero base, and non-finite results from exponentiation.

## Tests and Coverage

Backend tests cover pure arithmetic and HTTP behavior:

```bash
cd backend
gofmt -w .
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

Frontend tests use Vitest and React Testing Library:

```bash
cd frontend
npm install
npm run test:coverage
```

The Vitest configuration enforces 90% line, function, branch, and statement coverage. The backend suite is designed for the same target; the exact percentage is emitted by `go tool cover`.

Coverage artifacts are committed for review:

- Backend profile: [`backend/coverage.out`](backend/coverage.out)
- Backend HTML report: [`backend/coverage.html`](backend/coverage.html)
- Frontend HTML report: [`frontend/coverage/index.html`](frontend/coverage/index.html)
- Frontend detailed report: [`frontend/coverage/lcov-report/index.html`](frontend/coverage/lcov-report/index.html)

The latest generated reports show 91.8% backend statement coverage and 100% frontend statements, lines, and functions with 91.48% frontend branch coverage. Run the commands above to regenerate them after code changes.

## Design Decisions

- **Single endpoint:** The operation parameter avoids duplicating transport code while still keeping the API versioned and explicit. The application service invokes one domain aggregate for every calculation.
- **DDD boundaries:** The calculation domain owns operations, the finite-number value object, invariants, and execution. The application layer translates commands into domain objects and classifies invalid commands. The HTTP package only handles transport concerns and is wired to the use case by the composition root.
- **Validation ownership:** JSON shape and numeric parsing happen at the HTTP boundary. Numeric value creation, required operands, and operation-specific rules live in the domain package so they apply consistently to every caller.
- **Numeric policy:** JSON numbers must decode as finite Go `float64` values, and non-finite results are rejected instead of returning `Infinity` or `NaN`.
- **Error contract:** Client mistakes are `400`; mathematically undefined or unrepresentable calculations are `422`. Both use a stable `{ error: { code, message } }` body.
- **Scope trade-off:** There is no persistence or authentication because calculations are stateless. CORS is limited to the local Vite origin; Docker uses an Nginx same-origin proxy.
