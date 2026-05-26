# Chirpy API Docs

Base URL: `http://localhost:8080`

## Auth

- JWT-protected endpoints expect `Authorization: Bearer <access-token>`.
- Refresh/revoke endpoints also read the refresh token from `Authorization: Bearer <refresh-token>`.
- Polka webhooks expect `Authorization: ApiKey <POLKA_KEY>`.

## Endpoints

### `GET /api/healthz`

Health check endpoint.

Response:

```text
OK
```

### `POST /api/users`

Create a user.

Request:

```json
{
  "email": "user@example.com",
  "password": "secret"
}
```

Response `201 Created`:

```json
{
  "id": "uuid",
  "created_at": "2026-01-01T00:00:00Z",
  "updated_at": "2026-01-01T00:00:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false
}
```

### `PUT /api/users`

Update the authenticated user's email and password.

Headers:

```text
Authorization: Bearer <access-token>
```

Request:

```json
{
  "email": "new@example.com",
  "password": "new-secret"
}
```

Response `200 OK`: same shape as `POST /api/users`.

### `POST /api/login`

Authenticate a user and issue tokens.

Request:

```json
{
  "email": "user@example.com",
  "password": "secret"
}
```

Response `200 OK`:

```json
{
  "id": "uuid",
  "created_at": "2026-01-01T00:00:00Z",
  "updated_at": "2026-01-01T00:00:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false,
  "token": "jwt-access-token",
  "refresh_token": "refresh-token"
}
```

### `POST /api/refresh`

Exchange a refresh token for a new access token.

Headers:

```text
Authorization: Bearer <refresh-token>
```

Response `200 OK`:

```json
{
  "token": "jwt-access-token"
}
```

### `POST /api/revoke`

Revoke a refresh token.

Headers:

```text
Authorization: Bearer <refresh-token>
```

Response: `204 No Content`

### `POST /api/chirps`

Create a chirp for the authenticated user.

Headers:

```text
Authorization: Bearer <access-token>
```

Request:

```json
{
  "body": "Hello, Chirpy!"
}
```

Response `201 Created`:

```json
{
  "id": "uuid",
  "created_at": "2026-01-01T00:00:00Z",
  "updated_at": "2026-01-01T00:00:00Z",
  "body": "Hello, Chirpy!",
  "user_id": "uuid"
}
```

### `GET /api/chirps`

List chirps.

Query params:

- `author_id=<uuid>` filters by author.
- `sort=desc` sorts newest-first. Default order is ascending by creation time from storage.

Response `200 OK`:

```json
[
  {
    "id": "uuid",
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-01-01T00:00:00Z",
    "body": "Hello, Chirpy!",
    "user_id": "uuid"
  }
]
```

### `GET /api/chirps/{id}`

Fetch a chirp by ID.

Response `200 OK`: same object shape as `POST /api/chirps`.

### `DELETE /api/chirps/{id}`

Delete a chirp owned by the authenticated user.

Headers:

```text
Authorization: Bearer <access-token>
```

Response: `204 No Content`

### `POST /api/polka/webhooks`

Webhook endpoint for Polka upgrade events.

Headers:

```text
Authorization: ApiKey <POLKA_KEY>
```

Request:

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "uuid"
  }
}
```

Behavior:

- For `user.upgraded`, the target user is marked as Chirpy Red.
- Other event types return `204 No Content` and are ignored.

### `GET /admin/metrics`

Returns a small HTML page with file server hit metrics.

### `POST /admin/reset`

Resets users when `PLATFORM=dev`.

Response:

- `200 OK` on success
- `403 Forbidden` outside dev mode
