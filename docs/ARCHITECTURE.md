# Architecture

How the pieces of Le Bar fit together. To run or configure them, see [DEVELOPMENT.md](DEVELOPMENT.md).

```
                    bar.openapi.yml
               generates │         │ generates
                         ▼         ▼
   MongoDB ◄──── backend (Go) ◄──── frontend (SvelteKit)
                    │                  │            │
           HelloAsso, Discord     nfc daemon    porte daemon
                                 (card reader)  (door, cash drawer)
```

| Part | Role |
|---|---|
| [`bar.openapi.yml`](../bar.openapi.yml) | The API contract. The Go server code and the TypeScript client are both generated from it. |
| [`backend/`](../backend/) | Go API on Echo, storing everything in MongoDB |
| [`frontend/`](../frontend/) | One SvelteKit app that holds every interface |
| [`nfc/`](../nfc/) | Runs on the bar machine. Reads NFC cards and sends their IDs to the kiosk over a WebSocket. |
| [`porte/`](../porte/) | Runs on the bar machine. Checks a card and PIN with the backend, then opens the door or the cash drawer. |
| [`migration/`](../migration/) | Old one-off import from the previous MariaDB system |

## Backend

```
backend/
├── cmd/bar/          entrypoint
├── api/              HTTP handlers
├── autogen/          generated from the spec, don't edit
└── internal/
    ├── db/           DBackend interface + its MongoDB implementation
    ├── models/       generated types plus a Mongo _id
    ├── config/       environment variables
    ├── hash/         argon2 for PINs and passwords
    ├── storage/      uploaded pictures on disk
    └── webhook/      Discord notifications
```

`api.Server` implements the interface generated from the spec, so each operation in `bar.openapi.yml` is one method in `backend/api/`. A few file names are French: `passerCommande.go` places an order, `validerCommande.go` is the counter validating it, `ajoutProduit.go` adds a product, `course.go` is the shopping list, `deleted.go` is the trash.

### Auth

There are three cookie sessions. `BAR_SESS` is the logged-in member, `BAR_ADMIN_SESS` is an admin, and `BAR_ONBOARD_SESS` is a new card that hasn't been registered yet. `AuthMiddleware` loads them on every request, and each handler checks the one it needs. Requests from the bar's own machines also send an `X-Local-Token` header.

### Database

All access goes through the `DBackend` interface in `internal/db/database.go`. Deletes are soft: they set `deleted_at`, and reads skip those documents until someone restores or purges them from the trash. MongoDB has to run as a replica set, because anything that moves money (orders, refills, cash movements) uses transactions and the live balance uses a change stream.

### Background jobs

A ticker reconciles pending HelloAsso payments. If HelloAsso is unreachable at startup, the backend retries every 2 minutes and online top-ups stay off until then.

## Frontend

The browser loads `/config.json` at runtime to find the backend, so one build works for any deployment. The generated API client lives in `src/lib/api/` and `src/lib/requests/requests.ts` wraps it.

| Route | Who uses it |
|---|---|
| `/borne` | The kiosk at the bar: card login, ordering, carousel |
| `/client` | Members on their phone: ordering, history, online top-up |
| `/comptoir` | Counter staff: order queue, refills, end of day |
| `/admin` | Admins: accounts, products, categories, carousel, trash |
| `/panel` | Admins: menus, restocks, treasury, stock checks and other tools |
| `/auth` | Sign-in (Google, or the dev login server locally) |
