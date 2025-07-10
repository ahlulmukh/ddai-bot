# DDAI Depin Bot

Automated tool for DDAI Depin

## Requirements

- **DDAI Account**: [Register on DDAI Network](https://app.ddai.space/register?ref=r9zs0fHw)
- **Proxy**: Optional, but recommended for multi-account management
- **Captcha Service**: Configure in `config.json` (Private/2Captcha/AntiCaptcha)

## Features

- **Auto Registration**: Create DDAI accounts with random emails and usernames
- **Auto Task Claiming**: Automatically claims eligible mission rewards
- **Random Email Domain Generator**: Uses various email domains for registration
- **Proxy Support**: Rotate through proxies for multi-account creation
- **Detailed Logging**: Comprehensive activity logs with status indicators
- **Account Storage**: Saves successful accounts for future use
- **Captcha Solving Integration**: Multiple captcha solving services supported

## Checking System Architecture

To determine your system architecture, run the following command:

```bash
uname -m
```

- `x86_64`: Use the `linux-amd64` version.
- `aarch64`: Use the `linux-arm64` version.

## Configuration

The bot is configured through the `config.json` file. Here's a sample configuration:

```json
{
  "captchaServices": {
    "captchaUsing": "private",
    "urlPrivate": "https://your-captcha-service.com",
    "antiCaptchaApikey": ["your-anti-captcha-apikey"],
    "captcha2Apikey": ["your-2captcha-apikey"]
  }
}
```

### Configuration Options:

- **captchaServices**: Captcha solving service configuration

  - **captchaUsing**: Choose between "private", "antiCaptcha", or "2captcha"
  - **urlPrivate**: URL for your private captcha service (if using "private")
  - **antiCaptchaApikey**: API key for Anti-Captcha service
  - **captcha2Apikey**: API key for 2Captcha service

### Proxy Configuration:

Create a `proxy.txt` file with one proxy per line in the following format:

```
http://username:password@ip:port
http://ip:port
socks5://username:password@ip:port
socks5://ip:port
```

## Usage Instructions

### For Windows

1. Download the latest release from the [GitHub Releases page](https://github.com/ahlulmukh/ddai-bot/releases/latest)
2. Create/edit `config.json` and optionally `proxy.txt` files
3. Run File

### For Linux (AMD64)

1. Download the latest release from the [GitHub Releases page](https://github.com/ahlulmukh/ddai-bot/releases/latest) or use wget:

   ```bash
   wget https://github.com/ahlulmukh/ddai-bot/releases/latest/download/ddai-bot-amd64
   ```

2. Set execution permissions:

   ```bash
   chmod +x ddai-bot-amd64
   ```

3. Run the application:
   ```bash
   ./ddai-bot-amd64
   ```

### For Linux (ARM64)

1. Download the latest release from the [GitHub Releases page](https://github.com/ahlulmukh/ddai-bot/releases/latest) or use wget:

   ```bash
   wget https://github.com/ahlulmukh/ddai-bot/releases/latest/download/ddai-bot-arm64
   ```

2. Set execution permissions:

   ```bash
   chmod +x ddai-bot-arm64
   ```

3. Run the application:
   ```bash
   ./ddai-bot-arm64
   ```

## Account Management

### Saved Accounts

Successfully registered accounts are saved in `accounts.txt` with the format:

```
email:password
```

These accounts can be used later for additional tasks or verification purposes.

### Running Saved Accounts

The bot can also run tasks on previously created accounts. Place your accounts in `runaccounts.txt` with the format:

```
email:password
```

Then select "Run Accounts" from the menu to perform tasks with these existing accounts.

## Troubleshooting

### Captcha Issues

- **Private Captcha Service**: Ensure your private captcha service URL is correct and the service is running
- **API Key Errors**: Verify your Anti-Captcha or 2Captcha API keys are valid and have sufficient balance
- **Timeout Errors**: Try increasing the timeout settings if captchas take too long to solve

### Proxy Problems

- **Connection Errors**: Test your proxies independently to verify they're working
- **Rate Limiting**: If experiencing rate limits, try using more proxies or increasing the delay between operations
- **Authentication Failures**: Double-check proxy username/password format

### Task Claiming Issues

- **Mission Not Found**: Make sure your account has eligible missions available
- **Failed Claims**: Some social media tasks (like "Invite Friends") cannot be automatically claimed
- **API Changes**: If the DDAI API changes, update to the latest version of the bot

### Common Errors

- **HTTP 401**: Check your login credentials
- **HTTP 403**: Your IP might be banned or restricted; use a different proxy
- **HTTP 429**: Too many requests; increase delay between operations
- **HTTP 500**: Server-side issue, usually with the "Invite Friends" tasks; these can be safely ignored

## Stay Connected

Stay updated and connect with the community through the following channels:

- **Telegram Channel**: [Join on Telegram](https://t.me/elpuqus)
- **WhatsApp Channel**: [Join on WhatsApp](https://whatsapp.com/channel/0029VavBRhGBqbrEF9vxal1R)
- **Discord Server**: [Join on Discord](https://discord.com/invite/uKM4UCAccY)

## Support the Project

If you find this project helpful and would like to support its development, consider making a donation:

- **Solana**: `5jQMndHzWVH8MCitXdUEYmshJZXVKCzUG12mxJku4WdX`
- **EVM**: `0x72120c3c9cf3fee3ad57a34d5fcdbffe45f4ff28`
- **Bitcoin (BTC)**: `bc1ppfl3w3l4spnda7lawlhlycwuq2ekz74c936c8emwfprfu9hyun6sq5k6xl`

## Disclaimer

This tool is provided for educational and research purposes only. Please consider the following:

- **Terms of Service**: Using automated tools might violate the terms of service of DDAI Network
- **Legal Responsibility**: The developers bear no responsibility for any misuse of this software
- **Risk**: Use this tool at your own risk; account bans or restrictions may occur
- **Ethics**: Respect rate limits and don't overload services with excessive requests
- **Updates**: The DDAI API may change at any time, potentially breaking functionality

By using this software, you acknowledge that you understand and accept these risks and responsibilities.
