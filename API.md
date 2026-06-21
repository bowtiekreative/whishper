# Whishper API & Authentication

Whishper exposes a REST API for managing transcriptions. By default it is open
(no authentication), which preserves the original behaviour. Setting
`WHISHPER_API_KEY` turns authentication on for both the web UI and the API.

## Configuration

Add these to your `.env` file:

```env
# Enables auth. When empty, the API is open and no login is required.
WHISHPER_API_KEY=a-long-random-secret

# Credentials for the web UI login page.
WHISHPER_USERNAME=you@example.com
WHISHPER_PASSWORD=your-password
```

- `WHISHPER_API_KEY` is used **two** ways: as the `X-API-Key` header for
  programmatic access, and as the secret that signs web-UI login sessions.
- The web UI redirects to `/login` until you sign in with
  `WHISHPER_USERNAME` / `WHISHPER_PASSWORD`.

Generate a strong key, e.g.:

```bash
openssl rand -hex 32
```

## Authentication

There are two ways to authenticate a request:

1. **API key** (for scripts/integrations): send the header
   `X-API-Key: <WHISHPER_API_KEY>`.
2. **Session cookie** (the web UI): obtained by logging in via
   `POST /api/auth/login`.

`/api/auth/login` and `/api/auth/logout` are always reachable. Every other
`/api/*` route requires authentication when `WHISHPER_API_KEY` is set.

## Querying the API

Assuming Whishper is reachable at `http://localhost:8082`:

### List all transcriptions

```bash
curl http://localhost:8082/api/transcriptions \
  -H "X-API-Key: $WHISHPER_API_KEY"
```

### Get a single transcription

```bash
curl http://localhost:8082/api/transcriptions/<id> \
  -H "X-API-Key: $WHISHPER_API_KEY"
```

### Create a transcription (upload a file)

```bash
curl -X POST http://localhost:8082/api/transcriptions \
  -H "X-API-Key: $WHISHPER_API_KEY" \
  -F "file=@/path/to/audio.mp3" \
  -F "language=auto" \
  -F "modelSize=small" \
  -F "device=cpu"
```

### Create a transcription from a URL

```bash
curl -X POST http://localhost:8082/api/transcriptions \
  -H "X-API-Key: $WHISHPER_API_KEY" \
  -F "sourceUrl=https://youtube.com/watch?v=..." \
  -F "language=auto" \
  -F "modelSize=small" \
  -F "device=cpu"
```

### Delete a transcription

```bash
curl -X DELETE http://localhost:8082/api/transcriptions/<id> \
  -H "X-API-Key: $WHISHPER_API_KEY"
```

## Endpoint reference

| Method | Path                              | Description                         |
| ------ | --------------------------------- | ----------------------------------- |
| POST   | `/api/auth/login`                 | Log in, returns a session cookie    |
| POST   | `/api/auth/logout`                | Clear the session cookie            |
| GET    | `/api/transcriptions`             | List all transcriptions             |
| GET    | `/api/transcriptions/:id`         | Get one transcription               |
| POST   | `/api/transcriptions`             | Create a transcription job          |
| PATCH  | `/api/transcriptions`             | Update a transcription              |
| DELETE | `/api/transcriptions/:id`         | Delete a transcription              |
| GET    | `/api/translate/:id/:target`      | Translate a transcription           |

## Notes & limitations

- Uploaded media files are served by nginx at `/api/video/...` and are **not**
  behind the API key. If you need those protected too, add auth at the nginx
  layer.
- The "new transcription" dialog in the web UI now accepts **multiple files at
  once** — each selected file is queued as its own transcription job.
