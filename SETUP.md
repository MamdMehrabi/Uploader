# Setup Guide - Fixing Telegram Bot Token Error

## Common Error: "Telegram Bot Token not configured"

This error occurs when the bot token is missing or not properly configured. Follow these steps:

### Step 1: Create .env file

```bash
cp .env.example .env
```

Or create it manually:

```bash
touch .env
```

### Step 2: Get Your Telegram Bot Token

1. Open Telegram and search for **@BotFather**
2. Send the command: `/newbot`
3. Follow the instructions:
   - Choose a name for your bot
   - Choose a username (must end with 'bot')
4. Copy the token that BotFather gives you (looks like: `123456789:ABCdefGHIjklMNOpqrsTUVwxyz`)

### Step 3: Add Token to .env file

Open `.env` file and add:

```
TELEGRAM_BOT_TOKEN=your_actual_token_here
```

**Important:** Replace `your_actual_token_here` with the actual token from BotFather.

Example:
```
TELEGRAM_BOT_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz
```

### Step 4: (Optional) Set Default Chat ID

You can also set a default chat ID so you don't need to enter it every time:

```
DEFAULT_CHAT_ID=123456789
```

Or use your username:
```
DEFAULT_CHAT_ID=@yourusername
```

**To find your Chat ID:**
- Message [@userinfobot](https://t.me/userinfobot) on Telegram
- It will reply with your user ID

### Step 5: Restart the Server

After updating `.env`, restart your Go server:

```bash
# Stop the server (Ctrl+C) and restart
go run main.go
```

### Step 6: Verify Configuration

Visit `http://localhost:3000/api/health` in your browser. You should see:

```json
{
  "status": "ok",
  "botToken": "configured"
}
```

If you see `"botToken": "missing"`, check that:
- The `.env` file exists in the project root
- The token is correctly formatted (no extra spaces)
- You restarted the server after creating/editing `.env`

## Common Issues

### Issue: "Unauthorized: Invalid Telegram Bot Token"

**Solution:** 
- Double-check your token in `.env` file
- Make sure there are no extra spaces or quotes
- Get a fresh token from @BotFather if needed

### Issue: "chat not found"

**Solution:**
- Make sure you've started a conversation with your bot first
- Send `/start` to your bot in Telegram
- Then try uploading again

### Issue: Server doesn't read .env file

**Solution:**
- Make sure `.env` is in the same directory as `main.go`
- Check that the file is named exactly `.env` (not `.env.txt` or `env`)
- Restart the server after creating/editing `.env`

## Quick Test

1. Make sure `.env` exists with valid `TELEGRAM_BOT_TOKEN`
2. Start server: `go run main.go`
3. Open: `http://localhost:3000`
4. Enter your Chat ID
5. Select a file
6. Click Upload

If you still get errors, check the server console for detailed error messages.


