# Chirpy

Chirpy is a Go HTTP API backed by PostgreSQL.

## API Docs

The API reference lives in [`docs/api_docs.md`](docs/api_docs.md).

## Run Locally

1. Create a `.env` file with:
   - `DB_URL`
   - `PLATFORM`
   - `JWT_SECRET`
   - `POLKA_KEY`
2. Start the server:

```bash
go run ./cmd/chirpy
```

The server listens on `http://localhost:8080`.
