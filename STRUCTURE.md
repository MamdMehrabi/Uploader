# Project Structure

This document describes the modular structure of the Telegram Uploader application.

## Directory Structure

```
telegram-uploader/
├── main.go                 # Application entry point
├── config/                 # Configuration management
│   └── config.go          # Loads environment variables and config
├── handlers/               # HTTP request handlers
│   ├── home.go            # Home page handler
│   ├── health.go          # Health check and max file size handlers
│   └── upload.go          # File upload handler
├── services/               # Business logic services
│   └── telegram.go        # Telegram API service
├── models/                 # Data models/types
│   └── response.go        # Response structures
├── middleware/             # HTTP middleware
│   └── filesize.go        # File size limit middleware
├── utils/                  # Utility functions
│   └── chatid.go          # Chat ID normalization utilities
└── public/                 # Static frontend files
    ├── index.html
    ├── style.css
    └── script.js
```

## Package Descriptions

### `main.go`
- Application entry point
- Initializes configuration
- Sets up routes and middleware
- Starts the HTTP server

### `config/`
- **config.go**: Loads environment variables from `.env` file
- Returns a `Config` struct with all application settings
- Handles defaults for missing configuration values

### `handlers/`
- **home.go**: Serves the main HTML page
- **health.go**: 
  - Health check endpoint (`/api/health`)
  - Max file size endpoint (`/api/max-file-size`)
- **upload.go**: Handles file upload requests (`/api/upload`)
  - Validates file size
  - Normalizes chat ID
  - Calls Telegram service

### `services/`
- **telegram.go**: Telegram Bot API integration
  - Sends files to Telegram
  - Handles proxy configuration
  - Error handling and response parsing

### `models/`
- **response.go**: Response structures
  - `UploadResponse`: Upload operation response
  - `HealthResponse`: Health check response
  - `MaxFileSizeResponse`: File size limit response

### `middleware/`
- **filesize.go**: Adds file size limits to request context
  - Sets `maxFileSizeBytes` and `maxFileSizeMB` in context

### `utils/`
- **chatid.go**: Chat ID normalization utilities
  - Handles username format (`@username`)
  - Validates numeric IDs
  - Auto-adds `@` prefix when needed

## Benefits of This Structure

1. **Separation of Concerns**: Each package has a single responsibility
2. **Maintainability**: Easy to find and modify specific functionality
3. **Testability**: Each component can be tested independently
4. **Scalability**: Easy to add new handlers, services, or utilities
5. **Code Reusability**: Services and utilities can be reused across handlers

## Adding New Features

### Adding a New Handler
1. Create a new file in `handlers/` (e.g., `handlers/newfeature.go`)
2. Implement handler function(s)
3. Register route in `main.go`

### Adding a New Service
1. Create a new file in `services/` (e.g., `services/newservice.go`)
2. Implement service methods
3. Use the service in handlers

### Adding a New Model
1. Add struct to `models/response.go` or create new file in `models/`
2. Use the model in handlers/services

### Adding Middleware
1. Create a new file in `middleware/` (e.g., `middleware/auth.go`)
2. Implement middleware function
3. Register in `main.go` using `r.Use()`

## Route Organization

All routes are registered in `main.go`:

```go
// Static files
r.Static("/static", "./public")

// Pages
r.GET("/", handlers.HomeHandler)

// API endpoints
r.GET("/api/health", handlers.HealthHandler(cfg.BotToken))
r.GET("/api/max-file-size", handlers.MaxFileSizeHandler)
r.POST("/api/upload", uploadHandler.HandleUpload)
```

## Configuration

Configuration is loaded once at startup in `main.go`:

```go
cfg := config.Load()
```

All configuration values are available through the `Config` struct:
- `BotToken`: Telegram bot token
- `DefaultChatID`: Default chat ID for uploads
- `ProxyURL`: HTTP/HTTPS proxy URL
- `Port`: Server port
- `MaxFileSizeMB`: Maximum file size in MB

