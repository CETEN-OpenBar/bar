<div align="center">

<img src="docs/assets/logo.png" alt="Le Bar logo" width="360">

# Le Bar

Software for running a student bar: self-service kiosk, ordering from your phone, the counter, stock and the cash register.

[![Release](https://img.shields.io/github/v/release/CETEN-OpenBar/bar?style=flat-square&color=e8a33d)](https://github.com/CETEN-OpenBar/bar/releases/latest)
[![CI](https://img.shields.io/github/actions/workflow/status/CETEN-OpenBar/bar/ci.yml?branch=main&style=flat-square&label=ci)](https://github.com/CETEN-OpenBar/bar/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/CETEN-OpenBar/bar?style=flat-square&color=brightgreen)](LICENSE)
[![Last commit](https://img.shields.io/github/last-commit/CETEN-OpenBar/bar?style=flat-square)](https://github.com/CETEN-OpenBar/bar/commits/main)
[![Contributors](https://img.shields.io/github/contributors/CETEN-OpenBar/bar?style=flat-square)](https://github.com/CETEN-OpenBar/bar/graphs/contributors)

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Svelte](https://img.shields.io/badge/SvelteKit-FF3E00?style=for-the-badge&logo=svelte&logoColor=white)
![TypeScript](https://img.shields.io/badge/TypeScript-3178C6?style=for-the-badge&logo=typescript&logoColor=white)
![MongoDB](https://img.shields.io/badge/MongoDB-47A248?style=for-the-badge&logo=mongodb&logoColor=white)
![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-06B6D4?style=for-the-badge&logo=tailwindcss&logoColor=white)
![OpenAPI](https://img.shields.io/badge/OpenAPI-6BA539?style=for-the-badge&logo=openapiinitiative&logoColor=white)

</div>

<p align="center">
  <img src="docs/assets/screenshot.png" alt="Ordering screen on the kiosk" width="80%">
</p>

## What it does

- Members tap their NFC card at the kiosk, type their PIN and order. They can also scan a QR code to log in from their phone.
- From their phone, members can order, check their balance and history, and top up online with HelloAsso.
- Behind the counter, staff see orders arrive live, add money to accounts and close the day.
- Admins manage products, menus, stock, restocks, the shopping list, accounts and the cash register.
- Two small daemons on the bar machine drive the NFC reader, the door and the cash drawer.

## Repository

| Folder | What's inside |
|---|---|
| [`backend/`](backend/) | Go API (Echo + MongoDB) |
| [`frontend/`](frontend/) | SvelteKit app: kiosk, mobile, counter and admin |
| [`nfc/`](nfc/) | NFC card reader daemon |
| [`porte/`](porte/) | Door, cash drawer and fan daemon |
| [`bar.openapi.yml`](bar.openapi.yml) | API spec, the source of both the Go server and the TS client |

## Quick start

You need Docker and [mise](https://mise.jdx.dev/getting-started.html) (or Nix, see the [development guide](docs/DEVELOPMENT.md#prerequisites)).

```sh
mise install    # Go, Node and Task at the right versions
task setup      # config files and dependencies
task dev        # MongoDB, API on :8080 (reloads on save), app on :5173
task seed       # sample products and test accounts
```

Open <http://localhost:5173/borne>, type `1`, press <kbd>Enter</kbd>, then enter the PIN `1234`. That logs you in as the default admin without a card reader.

## Documentation

- [Development guide](docs/DEVELOPMENT.md): setup, configuration, regenerating the API, checks, releases
- [Architecture](docs/ARCHITECTURE.md): how the pieces fit together and where things live

## Contributing

Open an issue or a pull request. There are [issue templates](https://github.com/CETEN-OpenBar/bar/issues/new/choose) and a PR template.

<a href="https://github.com/CETEN-OpenBar/bar/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=CETEN-OpenBar/bar" alt="Contributors">
</a>

## License

[MIT](LICENSE)
