# Social Media App Backend (Go)

## Features

- Microsoft OAuth login
- MySQL database
- RESTful API for posts (CRUD, filter, comments, likes, image upload)

## Structure

- `configs/` - Configuration files (MySQL, OAuth)
- `controllers/` - HTTP handlers
- `middleware/` - Auth, logging, etc.
- `models/` - DB models
- `migrations/` - DB migrations & seeds
- `services/` - Business logic
- `utils/` - Helpers
- `public/` - Uploaded images

## Setup

1. Configure MySQL and Microsoft OAuth in `configs/`
2. Run migrations
3. Start the server: `go run main.go`
