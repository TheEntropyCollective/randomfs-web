# RandomFS Web Interface

A modern, responsive web interface for the Owner Free File System, providing an intuitive browser-based experience for storing and retrieving files using randomized blocks on IPFS.

## Overview

RandomFS Web Interface is a standalone web application that provides a beautiful, user-friendly interface for interacting with RandomFS. It's built with vanilla HTML, CSS, and JavaScript, making it easy to deploy and customize.

## Features

- **Drag & Drop Upload**: Intuitive file upload with drag-and-drop support
- **File Management**: View and manage stored files
- **rd:// URL Generation**: Automatic URL creation for file sharing
- **Download Links**: Direct download functionality
- **Real-time Statistics**: Live system metrics display
- **Responsive Design**: Works perfectly on desktop, tablet, and mobile
- **Modern UI**: Clean, professional interface with smooth animations
- **Cross-browser Compatible**: Works in all modern browsers
- **No Dependencies**: Pure HTML/CSS/JS - no frameworks required

## Quick Start

### Standalone Usage
```bash
# Clone the repository
git clone https://github.com/TheEntropyCollective/randomfs-web
cd randomfs-web

# Serve with any web server
python3 -m http.server 8000
# or
npx serve .
# or
php -S localhost:8000
```

### With RandomFS HTTP Server
The web interface is designed to work seamlessly with the RandomFS HTTP Server:

```bash
# Start the HTTP server with web interface
./randomfs-http -web ./randomfs-web
```

Then visit `http://localhost:8080` to access the interface.

## File Structure

```
randomfs-web/
├── index.html          # Main application page
├── css/
│   ├── style.css       # Main stylesheet
│   └── responsive.css  # Responsive design rules
├── js/
│   ├── app.js          # Main application logic
│   ├── upload.js       # File upload handling
│   ├── api.js          # API communication
│   └── utils.js        # Utility functions
├── assets/
│   ├── icons/          # Application icons
│   └── images/         # Background images
└── README.md           # This file
```

## API Integration

The web interface communicates with the RandomFS HTTP Server API:

### Endpoints Used
- `POST /api/v1/store` - Upload files
- `GET /api/v1/stats` - Get system statistics
- `GET /api/v1/retrieve/{hash}` - Download files
- `GET /api/v1/parse/{rd_url}` - Parse rd:// URLs

### Configuration
The interface automatically detects the API endpoint. You can customize it by setting:

```javascript
// In js/api.js
const API_BASE = 'http://localhost:8080/api/v1';
```

## Usage Guide

### Uploading Files
1. **Drag & Drop**: Simply drag files from your file manager to the upload area
2. **Click to Browse**: Click the upload area to open a file browser
3. **Multiple Files**: Select multiple files at once for batch upload
4. **Progress Tracking**: See real-time upload progress
5. **rd:// URL Generation**: Get shareable rd:// URLs immediately after upload

### Managing Files
- **View Files**: See all uploaded files with metadata
- **Download**: Click download links to retrieve files
- **Copy URLs**: Copy rd:// URLs for sharing
- **File Info**: View file size, type, and upload date

### System Statistics
- **Live Updates**: Real-time statistics display
- **Performance Metrics**: Cache hits, file counts, and storage usage
- **System Health**: Monitor RandomFS system status

## Customization

### Styling
The interface uses CSS custom properties for easy theming:

```css
:root {
  --primary-color: #2563eb;
  --secondary-color: #64748b;
  --background-color: #f8fafc;
  --text-color: #1e293b;
  --border-color: #e2e8f0;
}
```

### Configuration
Modify `js/config.js` to customize behavior:

```javascript
const CONFIG = {
  maxFileSize: 100 * 1024 * 1024, // 100MB
  allowedTypes: ['*/*'],
  autoRefresh: true,
  refreshInterval: 5000
};
```

## Browser Support

- **Chrome**: 90+
- **Firefox**: 88+
- **Safari**: 14+
- **Edge**: 90+
- **Mobile Browsers**: iOS Safari 14+, Chrome Mobile 90+

## Development

### Local Development
```bash
# Clone repository
git clone https://github.com/TheEntropyCollective/randomfs-web
cd randomfs-web

# Start development server
python3 -m http.server 8000

# Open in browser
open http://localhost:8000
```

### Building for Production
The web interface is already optimized for production. For additional optimization:

```bash
# Minify CSS (requires Node.js)
npx clean-css-cli css/style.css -o css/style.min.css

# Minify JavaScript (requires Node.js)
npx terser js/app.js -o js/app.min.js
```

## Deployment

### Static Hosting
Deploy to any static hosting service:

- **GitHub Pages**: Push to `gh-pages` branch
- **Netlify**: Drag and drop the folder
- **Vercel**: Connect GitHub repository
- **AWS S3**: Upload files to S3 bucket
- **Nginx**: Serve from web root

### Docker
```dockerfile
FROM nginx:alpine
COPY . /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

### CDN Integration
For production deployments, consider using a CDN:

```html
<!-- Example with Cloudflare -->
<script src="https://cdnjs.cloudflare.com/ajax/libs/..."></script>
```

## API Documentation

### File Upload
```javascript
const formData = new FormData();
formData.append('file', file);

fetch('/api/v1/store', {
  method: 'POST',
  body: formData
})
.then(response => response.json())
.then(data => {
  console.log('Uploaded:', data.rd_url);
});
```

### Get Statistics
```javascript
fetch('/api/v1/stats')
.then(response => response.json())
.then(data => {
  console.log('Files stored:', data.files_stored);
});
```

### Download File
```javascript
window.location.href = `/api/v1/retrieve/${hash}`;
```

## Troubleshooting

### Common Issues

**Upload Fails**
- Check that the RandomFS HTTP Server is running
- Verify API endpoint configuration
- Check file size limits
- Ensure CORS is properly configured

**Files Not Loading**
- Verify API connectivity
- Check browser console for errors
- Ensure proper file permissions

**Styling Issues**
- Clear browser cache
- Check CSS file loading
- Verify responsive breakpoints

### Debug Mode
Enable debug mode by setting:

```javascript
localStorage.setItem('debug', 'true');
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test across different browsers
5. Submit a pull request

### Development Guidelines
- Use semantic HTML
- Follow CSS best practices
- Write clean, documented JavaScript
- Test on multiple devices and browsers
- Maintain accessibility standards

## License

MIT License - see LICENSE file for details.

## Related Projects

- [randomfs-core](https://github.com/TheEntropyCollective/randomfs-core) - Core library
- [randomfs-cli](https://github.com/TheEntropyCollective/randomfs-cli) - Command-line interface
- [randomfs-http](https://github.com/TheEntropyCollective/randomfs-http) - HTTP server 