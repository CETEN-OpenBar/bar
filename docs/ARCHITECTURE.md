# Architecture

Map of the repository: what each part is, how the parts talk to each other, and where to look when changing something.

## Overview

```
                         bar.openapi.yml  (single source of truth for the HTTP API)
                          │            │
              oapi-codegen│            │openapi-generator (typescript-axios)
                          ▼            ▼
   ┌──────────────────────────┐   ┌───────────────────────────────┐
   │ backend/  (Go, Echo)     │◄──│ frontend/  (SvelteKit, node)  │
   │ :8080                    │   │ :3000 prod / :5173 dev        │
   └───────┬──────────┬───────┘   └──────┬───────────────┬────────┘
           │          │                  │ ws            │ HTTP (local_api)
           ▼          ▼                  ▼               ▼
       MongoDB    HelloAsso API     nfc/ (:3737)     porte/ (:8080 on the
     (replica set) Discord webhook  card reader      bar machine) door,
     ./storage (files)              → card UID       cash drawer, fan
                                                         │
                                                         └──► backend /auth/card
```

The parts:

| Path | What it is |
|---|---|
| [backend/](../backend/) | Go API server (`module bar`). Echo + MongoDB. Serves every route in the OpenAPI spec. |
| [frontend/](../frontend/) | SvelteKit 1 / Svelte 4 app built with `adapter-node`. One app holds every UI (kiosk, client, counter, admin). |
| [bar.openapi.yml](../bar.openapi.yml) | OpenAPI 3.0.3 spec (~6.8k lines, 72 paths). Both the Go server and the TS client are generated from it. |
| [nfc/](../nfc/) | Small Go daemon for the bar machine. Reads an ACR122U NFC reader over PC/SC and pushes `{"uid": ...}` to WebSocket clients on `localhost:3737`. Updates itself from GitHub releases. |
| [porte/](../porte/) | Small Go daemon for the bar machine. Exposes `POST /porte`, `/caisse`, `/ventilo`. It logs in to the backend with the card and PIN (`/auth/card` with `X-Local-Token`), checks `/account/admin`, then drives the door, cash drawer or fan. Updates itself from GitHub releases. |
| [migration/](../migration/) + [backend/cmd/migration](../backend/cmd/migration/main.go) | One-off importer from the old MariaDB system into Mongo. The compose file starts MariaDB and phpMyAdmin for it. Legacy. |
| [docs/](./) | This file, the README screenshot, the maintenance notes. |

Each Go program is its own module (`backend`, `nfc`, `porte`). There is no `go.work`, and `go.work` is gitignored.

## Backend

Entrypoint: [backend/cmd/bar/main.go](../backend/cmd/bar/main.go). It loads the config, connects to Mongo, seeds a default super-admin (card `1`, PIN `1234`, errors ignored), and calls `api.NewServer(db).Serve()`.

```
backend/
├── autogen/                 GENERATED, do not hand-edit
│   ├── bar.gen.go           oapi-codegen output: types + echo server + strict server + embedded spec
│   ├── utils.go             hand-written helpers (OptionalString, ...)
│   └── helloasso/           generated HelloAsso v5 client + hand-written auth.go
├── api/                     HTTP handlers, methods on *api.Server
├── internal/
│   ├── config/              env config (caarlos0/env + godotenv), loaded in init()
│   ├── db/database.go       DBackend interface: the only persistence contract
│   ├── db/mongo/            the only DBackend implementation
│   ├── models/              thin wrappers: autogen type inlined + Mongo _id (+ helpers)
│   ├── hash/                argon2 for PINs and passwords
│   ├── storage/             files on local disk under BAR_STORAGE_STORAGE_PATH (pictures)
│   └── webhook/             Discord webhook (restock notifications)
└── cmd/{bar,migration}/
```

### Request path

1. `api.Server` embeds `db.DBackend` and implements the generated `autogen.ServerInterface`. Routes are registered with `autogen.RegisterHandlers(e, s)` in [api/impl.go](../backend/api/impl.go).
2. Middlewares, in order: logrus logger, a 15 MB body limit, injection of the session stores, a hand-rolled CORS middleware that reflects `Origin`, then `AuthMiddleware` ([api/middlewares.go](../backend/api/middlewares.go)).
3. `AuthMiddleware` decodes three gorilla cookie sessions and puts the results on the echo context:
   - `BAR_SESS`: the logged-in user (`userLogged`, `userAccount`).
   - `BAR_ADMIN_SESS`: the admin. It counts only if the account's role is member, admin or superadmin (`adminLogged`, `adminAccountRole`).
   - `BAR_ONBOARD_SESS`: an unknown card that has not been onboarded yet. The account lives in an in-memory `go-cache`, not in Mongo.
   Each handler then checks the flags it needs. The OpenAPI `security` blocks are documentation only and nothing enforces them.
4. `X-Local-Token` (header) marks requests from trusted bar hardware: the kiosk, `porte`, card login, QR auth. Handlers compare it to `BAR_API_LOCAL_TOKEN`.

### Handlers ([backend/api/](../backend/api/))

Several files use French names:

| File | Does |
|---|---|
| `auth.go`, `connectCard.go`, `connectPassword.go`, `borneQr.go` | Login: card+PIN (an unknown card creates an onboarding account), password, Google OAuth (connect, link, impromptu), QR-code login from the kiosk over WebSocket |
| `account.go`, `accountPin.go`, `accountPassword.go` | Account CRUD, CSV/JSON import, PIN and password changes, "wants to staff" toggle |
| `accountWatcher.go` | `GET /account/watch`: long-poll. One shared Mongo change stream on `accounts` wakes the waiting requests. The `commande` layouts in borne/client loop on it to keep the balance current. |
| `websocket.go` | Small WebSocket "rooms" helper used by the QR login flows |
| `passerCommande.go` | "Place order": `POST /account/transactions` |
| `validerCommande.go` | "Validate order": counter staff patch a transaction or item state and debit the balance |
| `ajoutProduit.go` / `items.go` | Create items (incl. menus) / item CRUD, pictures, "incoherent" items (stock/price sanity view) |
| `categories.go`, `carousel.go` | Categories and kiosk carousel (texts + images) |
| `refills.go` | Counter-side balance refills |
| `remoteRefills.go` | Online refills through HelloAsso checkout. A background ticker (`HELLOASSO_CHECKOUT_PROCESSING_INTERVAL`) reconciles started checkouts. Turned off while the HelloAsso client is down (it retries every 2 min). |
| `restock.go`, `validerRestock.go` | Stock deliveries; validating one sends a Discord webhook |
| `course.go` | "Shopping list": items to buy, computed from stock |
| `cashMovements.go` | Cash drawer movements (treasury) |
| `stars.go` | "Starrings": stars granted to accounts, listed per account and globally (admin `treso/stars`) |
| `transactions.go` | Transaction listing and search |
| `deleted.go` | The admin "trash": list, restore or permanently delete every soft-deleted entity |
| `errors.go`, `utils.go` | Standard error responses, helpers |

### Persistence

- Everything goes through the `DBackend` interface in [internal/db/database.go](../backend/internal/db/database.go). Mongo is the only implementation.
- Each entity has a `<entity>_crud.go` (Create/Get/Update/MarkDelete/UnMarkDelete/Delete/Restore, plus the deleted listing) and a `<entity>_misc.go` (queries and counts). The CRUD files are near copies of each other.
- **Soft delete**: `MarkDelete*` sets `deleted_at`/`deleted_by`. Every read filters on `deleted_at` missing or null. `Delete*` is the hard delete, reachable from the trash.
- Collections and indexes are declared in [internal/db/mongo/init.go](../backend/internal/db/mongo/init.go): transactions, refills, remote_refills, starrings, items, categories, carousel_texts, carousel_images, accounts, restocks, cash_movements. `accounts` has partial unique indexes on `id`, `card_id`, `google_id`, `email_address`.
- Mongo must run as a **replica set**. Multi-document transactions (`WithTransaction`) and change streams both need it. The root `docker-compose.yml` starts `rs0`.
- Models embed the generated type with `bson:",inline"`. `regen` runs `gomodifytags` to add `bson` tags to the generated structs, which makes the API types the storage schema as well.

### Config ([internal/config/config.go](../backend/internal/config/config.go))

Read from env or `backend/.env` (see [.env.example](../backend/.env.example)):
`BAR_MONGO_{URI,DB,TIMEOUT}`, `BAR_API_{SESSION_SECRET,ADMIN_SESSION_SECRET,PORT,BASE_PATH,FRONTEND_BASE_PATH,LOCAL_TOKEN}`, `BAR_OAUTH_GOOGLE_{CLIENT_ID,CLIENT_SECRET}`, `HELLOASSO_{URL,CLIENT_ID,CLIENT_SECRET,SLUG,CHECKOUT_PROCESSING_INTERVAL}`, `BAR_STORAGE_STORAGE_PATH`, `BAR_DISCORD_WEBHOOK_URL`, `BAR_LOG_LEVEL`.

The config is parsed in a package `init()`, so importing `bar/internal/config` from anywhere, tests included, reads the environment.

## Frontend

SvelteKit app, client-rendered in practice: runtime config comes from `/config.json` through axios on the browser side. Styling uses Tailwind and Preline, icons use `iconify-icon`.

```
frontend/src/
├── lib/
│   ├── api/                 GENERATED typescript-axios client (+ docs/). Do not hand-edit.
│   ├── requests/requests.ts one factory per API group, all wired with basePath=api() and X-Local-Token
│   ├── config/config.ts     loads /config.json → api, apiws, local_api, local_token
│   ├── local/local.ts       calls to the porte daemon (open_door, open_caisse, open_ventilo)
│   ├── store/store.ts       svelte store: current account, comptoir search
│   ├── utils.ts
│   └── components/          shared + per-app components (see below)
└── routes/                  one top-level route per "app"
```

### Route groups

| Route | Audience | Notes |
|---|---|---|
| `/` | | Redirects to `/auth`. Still holds leftover demo code that logs in with a hard-coded card. |
| `/auth` | everyone | Login page |
| `/borne` | **kiosk** on the bar machine | Card reader (nfc WebSocket), carousel, ordering, QR onboarding, `mobile` (QR login from a phone), `onboarding` |
| `/client` | members on their own device | Ordering, history, profile, **HelloAsso refill** (`index/refill`, `refill/callback`) |
| `/comptoir` | counter staff | Live order queue (`c/transactions`), refills (`c/refills`), day summary (`c/resume`). Opens the cash drawer through `porte`. |
| `/admin` | admins | Entity CRUD (accounts, categories, products, carousel), the trash (`restore/*`) |
| `/panel` | admins/members | "Modules" grid ([components/panel/modules.svelte](../frontend/src/lib/components/panel/modules.svelte)): menu creation, restrictions, nicknames, restocks, PIN reset, category ordering, treasury (cash, refills, HelloAsso refills, stars, product history), incoherent products, shopping list |
| `/test` | | Dev scratch page |

`/borne` and `/client` grew out of the same code. Under `lib/components/`, the `borne/` and `client/` copies of `confirm`, `qr`, `refills` and `transactions` are byte-identical, while `pin`, `items` and `profile` differ a little. The route pages have diverged further. `/admin` and `/panel` have nearly identical `+layout.svelte` files.

Runtime config (`static/config.json`, gitignored, copy from [config.example.json](../frontend/static/config.example.json)):
`api` (backend URL), `apiws` (backend WS URL), `local_api` (porte daemon on the bar machine), `local_token`.

## Code generation

Edit `bar.openapi.yml` first, then regenerate both sides and commit the generated code together with the spec:

| Side | Command | Output |
|---|---|---|
| backend | `cd backend && task regen` | `autogen/bar.gen.go` (oapi-codegen `types,server,strict-server,spec`, then `gomodifytags` adds bson tags) |
| frontend | `cd frontend && task regen` | `src/lib/api/` (openapi-generator `typescript-axios`, generator version pinned to 7.16.0 in `openapitools.json`) |
| HelloAsso | `cd backend && task update-helloasso` | re-downloads the HelloAsso spec and regenerates `autogen/helloasso/` |

`task setup` in `backend` installs `oapi-codegen@latest`, so a regen can produce different output on different machines.

## Running locally

```sh
docker compose up -d                 # mongo (rs0) + mongo-express
cp backend/.env.example backend/.env
cp frontend/static/config.example.json frontend/static/config.json
(cd backend && go run cmd/bar/main.go)
(cd frontend && npm i && npm run dev)
```

`docker compose --profile dev up` also builds and runs the backend (:9999) and frontend (:9998) containers. [layout.kdl](../layout.kdl) is a zellij layout that starts both dev servers. [flake.nix](../flake.nix) together with `.envrc` (direnv) provides the toolchain.

## CI/CD ([.github/workflows/](../.github/workflows/))

| Workflow | Trigger | Does |
|---|---|---|
| `build.yml` | every push | backend `go build` + `go test` (Go 1.20), frontend `npm install` + `npm run build` (Node 20). No lint, no `svelte-check`. |
| `docker.yml` | release published | Builds `ghcr.io/<repo>-backend` and `-frontend`, then `curl`s `secrets.CD_URL` to deploy |
| `nfc.yml`, `porte.yml` | release created | Builds the daemon binaries (CGO, pcsclite) and attaches them to the release. The daemons update themselves from these. |

Code Climate config: [.codeclimate.yml](../.codeclimate.yml), which excludes the generated code.
