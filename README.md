# Telegram File Uploader Web

A modern web application built with Go that allows you to upload files directly to Telegram through a beautiful web interface.

## Features

- 🚀 Fast and lightweight Go backend
- 🎨 Modern, responsive web UI
- 📤 Upload files up to 50MB to Telegram
- 💬 Support for captions
- 🔒 Secure file handling
- 📱 Mobile-friendly interface

## Prerequisites

- Go 1.21 or higher
- A Telegram Bot Token (get it from [@BotFather](https://t.me/BotFather))

## Setup

1. **Clone or navigate to the project directory:**
   ```bash
   cd Uplaoder-1
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Create a `.env` file:**
   ```bash
   cp .env.example .env
   ```

4. **Edit `.env` and add your Telegram Bot Token:**
   ```
   TELEGRAM_BOT_TOKEN=your_bot_token_here
   DEFAULT_CHAT_ID=your_chat_id_here  # Optional
   PORT=3000  # Optional, defaults to 3000
   ```

5. **Get your Telegram Bot Token:**
   - Open Telegram and search for [@BotFather](https://t.me/BotFather)
   - Send `/newbot` and follow the instructions
   - Copy the bot token and paste it in your `.env` file

6. **Get your Chat ID (optional):**
   - You can use your user ID (numeric) or username (e.g., @username)
   - To find your user ID, message [@userinfobot](https://t.me/userinfobot) on Telegram
   - Or leave `DEFAULT_CHAT_ID` empty and enter it in the web form

## Running the Application

1. **Start the server:**
   ```bash
   go run main.go
   ```

2. **Open your browser:**
   Navigate to `http://localhost:3000`

3. **Upload files:**
   - Enter your Chat ID (or use the default if configured)
   - Optionally add a caption
   - Select a file
   - Click "Upload to Telegram"

## Building for Production

```bash
go build -o telegram-uploader main.go
./telegram-uploader
```

## API Endpoints

### `POST /api/upload`
Upload a file to Telegram.

**Form Data:**
- `file` (required): The file to upload
- `chatId` (required): Telegram chat ID or username
- `caption` (optional): Caption for the file

**Response:**
```json
{
  "success": true,
  "message": "File uploaded successfully",
  "fileId": "file_id_from_telegram",
  "messageId": 123
}
```

### `GET /api/health`
Check server health and configuration status.

**Response:**
```json
{
  "status": "ok",
  "botToken": "configured"
}
```

## File Size Limits

The default file size limit is **20MB**, but you can configure it in your `.env` file:

```
MAX_FILE_SIZE_MB=20
```

Telegram Bot API supports files up to 50MB maximum. You can set any value between 1 and 50 MB.

## Proxy Support (For Blocked Countries)

If you're in a country where Telegram is blocked (like Iran), you need to configure a proxy:

1. Add to your `.env` file:
   ```
   HTTP_PROXY=socks5://127.0.0.1:1080
   ```
   Or for HTTP proxy:
   ```
   HTTP_PROXY=http://127.0.0.1:1080
   ```

2. Restart the server

See [PROXY_SETUP.md](PROXY_SETUP.md) for detailed proxy configuration instructions.

## Troubleshooting

### Error: "Telegram Bot Token not configured"

1. Make sure you have a `.env` file in the project root
2. Add your bot token: `TELEGRAM_BOT_TOKEN=your_token_here`
3. Restart the server after creating/editing `.env`
4. Check the health endpoint: `http://localhost:3000/api/health`

See [SETUP.md](SETUP.md) for detailed setup instructions.

### Error: "Unauthorized: Invalid Telegram Bot Token"

- Verify your token is correct (no extra spaces)
- Get a fresh token from [@BotFather](https://t.me/BotFather)
- Make sure the token format is correct: `123456789:ABCdefGHIjklMNOpqrsTUVwxyz`

### Error: "chat not found"

- Start a conversation with your bot first (send `/start`)
- Verify your Chat ID is correct
- Use numeric Chat ID or username (e.g., `@username`)

## Security Notes

- Never commit your `.env` file to version control
- Keep your bot token secure
- Consider adding authentication for production use

## License

MIT License - see LICENSE file for details

## Author

Mamd Mehrabi Rad

