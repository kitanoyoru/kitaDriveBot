# kitaDriveBot

Personal Telegram bot that uploads PDFs to your Google Drive. Send a PDF with an optional caption as the destination folder path (for example `Work/Invoices/2026`).

## Features

- Accepts PDF documents from your Telegram account only
- Creates nested Google Drive folders on demand
- Uses long polling (no public webhook URL required)
- One-time Google OAuth setup with refresh token persistence

## Prerequisites

- Go 1.21+
- A Telegram bot token from [@BotFather](https://t.me/BotFather)
- Your Telegram user ID
- A Google Cloud project with OAuth credentials (Desktop app)

## Google Cloud setup

1. Create a project in [Google Cloud Console](https://console.cloud.google.com/).
2. Enable the **Google Drive API**.
3. Configure the OAuth consent screen.
4. Create OAuth client credentials of type **Desktop app**.
5. Add `http://127.0.0.1:8080/oauth/callback` as an authorized redirect URI.

Copy the client ID and client secret into your environment.

## Configuration

Copy the example env file and fill in values:

```bash
cp .env.example .env
```

| Variable | Required | Description |
|----------|----------|-------------|
| `TELEGRAM_BOT_TOKEN` | yes | Bot token from BotFather |
| `OWNER_TELEGRAM_ID` | yes | Your Telegram user ID |
| `GOOGLE_CLIENT_ID` | yes | Google OAuth client ID |
| `GOOGLE_CLIENT_SECRET` | yes | Google OAuth client secret |
| `GOOGLE_TOKEN_PATH` | no | Token file path (default `./data/token.json`) |
| `DRIVE_ROOT_FOLDER_ID` | no | Optional base folder ID in Drive |
| `DEFAULT_FOLDER_PATH` | no | Used when caption is empty (default `Telegram`) |
| `LOG_LEVEL` | no | `debug`, `info`, `warn`, or `error` |

The bot automatically loads `.env` from the project root when present.

## One-time Google authorization

Run the OAuth flow locally and complete it in your browser:

```bash
go run ./cmd/bot auth
```

This stores a refresh token at `GOOGLE_TOKEN_PATH`. Re-run only if the token is revoked or lost.

## Run locally

```bash
go run ./cmd/bot run
```

Or with Task:

```bash
task run
```

## Usage

1. Start the bot.
2. Send `/start` to see a short help message.
3. Send a PDF to the bot.
4. Optionally set the caption to a folder path such as `Receipts/2026/June`.
5. The bot uploads the file and replies with a Google Drive link.

If the caption is empty, the bot uses `DEFAULT_FOLDER_PATH`.

## Docker

Build and run with Docker Compose:

```bash
docker compose up --build
```

Mount `./data` so the Google token persists across restarts. Run `bot auth` on your host first, or exec into the container for initial auth if port `8080` is exposed.

## Development

```bash
task generate
task tidy
task test
task lint
```

## Limits

- Telegram Bot API file downloads are limited to 20 MB.
- The bot only accepts PDF files.
- Only the configured `OWNER_TELEGRAM_ID` can use the bot.

## License

See [LICENSE](LICENSE).
