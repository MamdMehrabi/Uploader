# Proxy Setup Guide for Blocked Countries (Iran, etc.)

If you're in a country where Telegram is blocked (like Iran), you need to configure a proxy to access Telegram's API.

## Quick Setup

### Step 1: Get a Proxy Server

You need a proxy server that can access Telegram. Options:

1. **VPN Proxy** - If you have a VPN, use its proxy settings
2. **SOCKS5 Proxy** - Many VPNs provide SOCKS5 proxies
3. **HTTP Proxy** - Standard HTTP proxy server

### Step 2: Add Proxy to .env file

Open your `.env` file and add:

```bash
HTTP_PROXY=http://127.0.0.1:1080
```

Or for SOCKS5:
```bash
HTTP_PROXY=socks5://127.0.0.1:1080
```

### Step 3: Common Proxy Formats

**HTTP Proxy:**
```
HTTP_PROXY=http://127.0.0.1:1080
HTTP_PROXY=http://username:password@proxy.example.com:8080
```

**SOCKS5 Proxy:**
```
HTTP_PROXY=socks5://127.0.0.1:1080
HTTP_PROXY=socks5://username:password@proxy.example.com:1080
```

**HTTPS Proxy:**
```
HTTPS_PROXY=https://proxy.example.com:8080
```

### Step 4: Restart Server

After adding the proxy, restart your server:

```bash
go run main.go
```

You should see a log message:
```
Proxy configured: http://127.0.0.1:1080
```

## Example .env File for Iran

```bash
# Telegram Bot Token
TELEGRAM_BOT_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz

# Default Chat ID (optional)
DEFAULT_CHAT_ID=123456789

# Server Port
PORT=3000

# Proxy for accessing Telegram API (required in blocked countries)
HTTP_PROXY=socks5://127.0.0.1:1080
```

## Testing Your Proxy

1. Make sure your proxy is running and accessible
2. Test the connection:
   ```bash
   curl -x http://127.0.0.1:1080 https://api.telegram.org
   ```
3. If you get a response, your proxy works!

## Common Proxy Services

### Using a VPN's Proxy

If you're using a VPN:
- Check your VPN settings for proxy configuration
- Usually found in Advanced/Proxy settings
- Common ports: 1080 (SOCKS5), 8080 (HTTP)

### Using Shadowsocks/V2Ray

If you're using Shadowsocks or V2Ray:
- They usually provide a local SOCKS5 proxy
- Default: `socks5://127.0.0.1:1080`
- Check your Shadowsocks/V2Ray client settings

### Using Telegram's Built-in Proxy

Telegram itself uses proxies, but for API access you need:
- A separate HTTP/SOCKS5 proxy
- Configure it in your `.env` file

## Troubleshooting

### Error: "Cannot connect to Telegram API"

**Solutions:**
1. Verify your proxy is running: `curl -x YOUR_PROXY https://api.telegram.org`
2. Check proxy URL format in `.env`
3. Make sure proxy allows HTTPS connections
4. Try a different proxy server

### Error: "Invalid proxy URL"

**Solutions:**
1. Check the format: `http://` or `socks5://` prefix required
2. Verify port number is correct
3. Check for typos in `.env` file

### Proxy Works but Still Getting Errors

**Solutions:**
1. Some proxies block Telegram - try a different proxy
2. Check if your proxy requires authentication
3. Verify bot token is still valid
4. Check server logs for detailed error messages

## Security Notes

- Never commit your `.env` file with proxy credentials
- Use secure proxies (avoid public/free proxies)
- Consider using authentication: `http://user:pass@proxy:port`

## Need Help?

If you're still having issues:
1. Check server logs for detailed error messages
2. Test proxy connectivity separately
3. Verify your bot token is correct
4. Make sure you've started a conversation with your bot first


