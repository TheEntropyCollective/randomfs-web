# RandomFS Web Server

A modern, responsive web server for the Owner Free File System, providing an intuitive browser-based experience for storing and retrieving files using randomized blocks on IPFS.

## Overview

RandomFS Web Server is a standalone Go executable that provides a beautiful, user-friendly web interface for interacting with RandomFS. It combines the power of RandomFS with a modern web interface, making it easy to upload, store, and share files through the browser.

## Features

- **Web Interface**: Beautiful, responsive web UI for file management
- **Drag & Drop Upload**: Intuitive file upload with drag-and-drop support
- **File Management**: View and manage stored files
- **rfs:// URL Generation**: Automatic URL creation for file sharing
- **Direct File Access**: Access files via `/rfs/{encodedURL}` endpoints
- **Real-time Statistics**: Live system metrics display
- **RESTful API**: Full API for programmatic access
- **IPFS Integration**: Seamless integration with IPFS for decentralized storage
- **Local Caching**: Fast access with intelligent caching
- **Cross-platform**: Runs on Windows, macOS, and Linux

## Quick Start

### Prerequisites
- Go 1.21 or later
- IPFS daemon running (optional, can run without IPFS for testing)

### Building and Running

```bash
# Clone the repository
git clone https://github.com/TheEntropyCollective/randomfs
cd randomfs/randomfs-web

# Install dependencies
go mod tidy

# Build the executable
go build -o randomfs-web main.go

# Run the server
./randomfs-web
```

Then visit `http://localhost:8080` to access the web interface.

### Using Make

```bash
# Build and run
make run

# Run on custom port
make run-port PORT=3000

# Run without IPFS (for testing)
make run-no-ipfs

# Show all available commands
make help
```

## Configuration

The server can be configured using environment variables:

```bash
# Server port (default: 8080)
export RANDOMFS_PORT=3000

# Data directory (default: ./data)
export RANDOMFS_DATA_DIR=/path/to/data

# IPFS API endpoint (default: http://localhost:5001)
export RANDOMFS_IPFS_API=http://localhost:5001

# Run the server
./randomfs-web
```

## API Endpoints

### File Operations
- `POST /api/v1/store` - Upload a file
- `GET /api/v1/retrieve/{hash}` - Download a file by hash
- `GET /rfs/{encodedURL}` - Access file via encoded rfs:// URL

### System Information
- `GET /api/v1/stats` - Get system statistics

### Web Interface
- `GET /` - Main web interface
- `GET /index.html` - Web interface (same as root)

## Usage Examples

### Upload a File via API
```bash
curl -X POST -F "file=@/path/to/file.txt" http://localhost:8080/api/v1/store
```

Response:
```json
{
  "url": "rfs://randomfs/v4/1024/file.txt/1640995200/QmHash...",
  "hash": "QmHash..."
}
```

### Get System Statistics
```bash
curl http://localhost:8080/api/v1/stats
```

Response:
```json
{
  "files_stored": 42,
  "blocks_generated": 156,
  "total_size": 1048576,
  "cache_hits": 89,
  "cache_misses": 12
}
```

### Access File via Web URL
```bash
# The rfs:// URL is base64 encoded in the web path
curl http://localhost:8080/rfs/cmQ6Ly9yYW5kb21mcy92NC8xMDI0L2ZpbGUudHh0LzE2NDA5OTUyMDAvUW1IYXNoLi4u
```

## File Structure

```
randomfs-web/
├── main.go              # Main server executable
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
├── Makefile             # Build and run commands
├── web/
│   └── index.html       # Web interface
├── data/                # Local data storage (created automatically)
└── README.md            # This file
```

## Development

### Local Development
```bash
# Install dependencies
go mod tidy

# Run with hot reload (requires air or similar)
air

# Run tests
go test ./...

# Build for different platforms
GOOS=linux GOARCH=amd64 go build -o randomfs-web-linux main.go
GOOS=darwin GOARCH=amd64 go build -o randomfs-web-macos main.go
GOOS=windows GOARCH=amd64 go build -o randomfs-web-windows.exe main.go
```

### Testing
```bash
# Run all tests
make test

# Run without IPFS for testing
make run-no-ipfs
```

## Deployment

### Docker
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o randomfs-web main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/randomfs-web .
COPY --from=builder /app/web ./web
EXPOSE 8080
CMD ["./randomfs-web"]
```

### Systemd Service
```ini
[Unit]
Description=RandomFS Web Server
After=network.target

[Service]
Type=simple
User=randomfs
WorkingDirectory=/opt/randomfs-web
ExecStart=/opt/randomfs-web/randomfs-web
Environment=RANDOMFS_PORT=8080
Environment=RANDOMFS_DATA_DIR=/var/lib/randomfs
Restart=always

[Install]
WantedBy=multi-user.target
```

### Nginx Reverse Proxy
```nginx
server {
    listen 80;
    server_name randomfs.example.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Troubleshooting

### Common Issues

1. **IPFS Connection Failed**
   - Ensure IPFS daemon is running: `ipfs daemon`
   - Check IPFS API endpoint: `curl http://localhost:5001/api/v0/version`

2. **Port Already in Use**
   - Change port: `RANDOMFS_PORT=3000 ./randomfs-web`
   - Check what's using the port: `lsof -i :8080`

3. **Permission Denied**
   - Ensure write permissions to data directory
   - Run with appropriate user permissions

4. **File Upload Fails**
   - Check file size limits (default: 32MB)
   - Verify IPFS is accessible
   - Check server logs for detailed error messages

### Logs
The server provides detailed logging. Common log messages:
- `RandomFS Web Server starting on port 8080`
- `Data directory: ./data`
- `Web interface available at: http://localhost:8080`

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 