# Development guide

How to run, change and release Le Bar. For how the code is organised, see [ARCHITECTURE.md](ARCHITECTURE.md).

## Prerequisites

Works on Linux, macOS and WSL. You need [Docker](https://docs.docker.com/get-docker/) with Compose, plus Go 1.27, Node.js 24 and [Task](https://taskfile.dev). Pick one way to get the last three:

- **mise** (simplest): [install mise](https://mise.jdx.dev/getting-started.html), then run `mise install` in the repo. It reads [mise.toml](../mise.toml).
- **Nix**: `nix develop`, or [direnv](https://direnv.net/) with the provided `.envrc`. The flake also brings gopls, gofumpt, the Svelte language server and `openapi-generator-cli`.
- **By hand**: install Go, Node and Task yourself at those versions.

## Running locally

```sh
task setup    # first time: checks your tools, creates backend/.env and frontend/static/config.json, installs dependencies
task dev      # MongoDB + backend + frontend in one terminal, Ctrl-C to stop
task seed     # sample products and test accounts, needed to sign in on /client
```

`task setup` never overwrites a `.env` or `config.json` you already have, so you can run it again at any time.

`task dev` starts MongoDB, Mongo Express and a dev login server in Docker, then runs the backend on <http://localhost:8080> and the frontend on <http://localhost:5173>. The backend runs under [air](https://github.com/air-verse/air) and rebuilds when you save a `.go` file. If the build fails, the error shows up in the log and the previous version keeps running. The very first start takes a little longer while Go compiles everything.

Ctrl-C stops the backend and the frontend. MongoDB keeps running, with your data, until you run `task down`.

Mongo Express is at <http://localhost:8081>, with no login, and the dev login server is at <http://localhost:8082>. Like MongoDB, they only accept connections from your own machine.

MongoDB runs as a replica set, because the backend uses multi-document transactions and change streams. The compose file sets that up. If you bring your own Mongo, initialise a replica set first.

`task --list` shows every task. If you use [zellij](https://zellij.dev/), `zellij --layout layout.kdl` opens the frontend, backend, lazygit and an editor in tabs.

## Configuration

### Backend: `backend/.env`

`task setup` creates it from [`.env.example`](../backend/.env.example). Those values work for local development as they are: sign-in on `/client` goes through the dev login server, and HelloAsso stays off until you fill it in.

| Variable | Default | Notes |
|---|---|---|
| `BAR_MONGO_URI` | `mongodb://127.0.0.1:27017` | |
| `BAR_MONGO_DB` | `bar` | |
| `BAR_MONGO_TIMEOUT` | `3000` | Query timeout, in milliseconds |
| `BAR_API_PORT` | `:8080` | |
| `BAR_API_BASE_PATH` | `http://localhost:8080` | Public URL of the API, used for the OAuth redirect |
| `BAR_API_FRONTEND_BASE_PATH` | `http://localhost:5173` | Where to send users back after login |
| `BAR_API_SESSION_SECRET` / `BAR_API_ADMIN_SESSION_SECRET` | | Cookie signing keys. Use long random values outside development. |
| `BAR_API_LOCAL_TOKEN` | | Shared secret for the bar's own machines (see below) |
| `BAR_OAUTH_PROVIDER` | `google` | `google`, or `oidc` for any OpenID Connect server. `.env.example` sets `oidc` for the dev login server. |
| `BAR_OAUTH_GOOGLE_CLIENT_ID` / `_SECRET` | | Google sign-in |
| `BAR_OAUTH_OIDC_ISSUER` | | Issuer URL when the provider is `oidc`, `http://localhost:8082/default` for the dev server |
| `BAR_OAUTH_OIDC_CLIENT_ID` / `_SECRET` | | Client credentials for that issuer. The dev server accepts any. |
| `HELLOASSO_CLIENT_ID` / `_SECRET` / `_SLUG` / `_URL` | | Online top-ups. Optional. |
| `HELLOASSO_CHECKOUT_PROCESSING_INTERVAL` | `10m` | How often pending HelloAsso checkouts are reconciled |
| `BAR_STORAGE_STORAGE_PATH` | `./storage` | Uploaded pictures |
| `BAR_DISCORD_WEBHOOK_URL` | | Restock notifications. Optional. |
| `BAR_LOG_LEVEL` | `info` | `trace` … `panic` |

#### Local token

Some routes (card login, the kiosk QR code, the door and cash drawer) are only meant for the bar's own machines. Those requests carry an `X-Local-Token` header that must equal `BAR_API_LOCAL_TOKEN`. The frontend sends the `local_token` from its `config.json`, so only deploy that token to the machines in the bar.

#### Google sign-in (optional)

1. In the [Google Cloud Console](https://console.cloud.google.com/), create a project, then go to *APIs & Services → Credentials*.
2. Create an *OAuth client ID* of type *Web application*.
3. Add `<BAR_API_BASE_PATH>/auth/google/callback` (for example `http://localhost:8080/auth/google/callback`) as an authorised redirect URI.
4. Put the client ID and secret in `BAR_OAUTH_GOOGLE_CLIENT_ID` and `BAR_OAUTH_GOOGLE_CLIENT_SECRET`, and set `BAR_OAUTH_PROVIDER=google`.

#### HelloAsso (optional)

Leave the `HELLOASSO_*` credentials empty and online top-ups are disabled; the backend logs that once at startup. With credentials set, the backend connects at startup and retries every 2 minutes if HelloAsso can't be reached.

### Frontend: `frontend/static/config.json`

Loaded by the browser at runtime, so one build can be pointed at any backend. `task setup` creates it from [`config.example.json`](../frontend/static/config.example.json):

| Key | Meaning |
|---|---|
| `api` | Backend base URL |
| `apiws` | Backend WebSocket URL |
| `local_api` | The `porte` daemon on the bar machine (door, cash drawer) |
| `local_token` | Must match `BAR_API_LOCAL_TOKEN` |

## Signing in during development

On startup the backend creates a super-admin with card `1` and PIN `1234` if no account has that card yet.

You don't need an NFC reader: open <http://localhost:5173/borne>, type `1`, press <kbd>Enter</kbd>, then enter `1234`. Typing on `/borne` does the same thing as tapping a card.

After `task seed` there are three more test cards, all with PIN `1234`:

| Card | Role | Balance |
|---|---|---|
| `2` | member | 0 € |
| `3` | student | 20 € |
| `4` | admin | 0 € |

Any other card ID starts the onboarding flow with the PIN `1234`.

### On `/client`

`/client` signs in through the dev login server instead of Google. On <http://localhost:5173/auth>, click the sign-in button and type a username on the page that opens. It logs you in as `<username>@telecomnancy.net`.

As with Google, that only works for an account that already exists. After `task seed`, use `membre`, `etudiant` or `admin` (cards 2, 3 and 4). To use any other name, link it to a card first: sign in on `/borne` with a new card number and go through onboarding, which shows a QR code to link an account. A phone can't reach `localhost`, but the backend logs the QR code's URL (`QR code generated for account …`) at the `debug` level that `.env.example` sets. Open that URL in your browser and type the name there.

The main interfaces:

| URL | Interface |
|---|---|
| `/borne` | Kiosk on the bar's touchscreen |
| `/client` | Members on their phone (sign-in via `/auth`) |
| `/comptoir` | Counter staff |
| `/admin`, `/panel` | Administration |

## Changing the API

[`bar.openapi.yml`](../bar.openapi.yml) is the contract. Both sides are generated from it, so never edit the generated code by hand.

The root `task setup` doesn't install the generators. Install them once with `task backend:setup` (oapi-codegen, gomodifytags) and `task frontend:setup` (openapi-generator-cli, which needs Java). Nix users already have openapi-generator-cli. Be aware that `task backend:setup` installs the latest oapi-codegen, which currently produces code that doesn't build against the committed runtime (see #160).

1. Edit `bar.openapi.yml`.
2. Regenerate both sides:
   ```sh
   (cd backend && task regen)     # → backend/autogen/bar.gen.go
   (cd frontend && task regen)    # → frontend/src/lib/api/
   ```
3. Implement the new handler on `*api.Server` in `backend/api/`. The build fails until every operation in the spec has an implementation.
4. Commit the spec and both generated outputs together.

The HelloAsso client is generated the same way from their public spec: `cd backend && task update-helloasso`.

## Checks before a pull request

```sh
# backend (repeat go build / go vet in nfc/ and porte/ if you touched them)
cd backend && go build ./... && go vet ./... && go test ./...

# frontend
cd frontend && npm run check && npm run lint && npm run build
```

CI only runs the Go build and tests and the frontend build ([build.yml](../.github/workflows/build.yml)), so run the other checks yourself.

Format Go with `gofumpt` and the frontend with `npm run format`.

## Docker

```sh
docker compose --profile dev up --build
```

builds and runs the backend on <http://localhost:9999> (reading `backend/.env`) and the frontend on <http://localhost:9998> (mounting `frontend/static/config.dev.json` as its config), next to MongoDB.

Image definitions: [backend/Dockerfile](../backend/Dockerfile), [frontend/Dockerfile](../frontend/Dockerfile).

## Hardware daemons

Both run on the bar's own machine and are separate Go modules.

[`nfc/`](../nfc/) reads an ACR122U NFC reader over PC/SC and broadcasts each card UID as `{"uid": "..."}` on `ws://localhost:3737`. The kiosk listens there. Building it needs `pcscd`, `libpcsclite-dev`, `pkg-config` and CGO. It ships a systemd unit, [nfc.service](../nfc/nfc.service).

[`porte/`](../porte/) exposes `POST /porte`, `/caisse` and `/ventilo` on `:8080`. It checks the card and PIN against the backend, then opens the door or the cash drawer, or switches the fan. Configure it with [`.env.example`](../porte/.env.example) (`BAR_API_URL`, `BAR_LOCAL_TOKEN`). Systemd unit: [porte.service](../porte/porte.service).

Both daemons update themselves from the latest GitHub release.

## Releases and deployment

Publishing a GitHub release:

- builds and pushes `ghcr.io/ceten-openbar/bar-backend` and `bar-frontend`, then calls the deployment hook ([docker.yml](../.github/workflows/docker.yml)).
- builds the `nfc` and `porte` binaries and attaches them to the release ([nfc.yml](../.github/workflows/nfc.yml), [porte.yml](../.github/workflows/porte.yml)), where the daemons pick them up.

Releases are tagged `vX.Y.Z`.
