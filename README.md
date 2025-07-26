# 📱💻 Twarit - Simple File Sharing

> **Twarit** is a lightweight, secure file sharing application that enables seamless file transfers between your PC and mobile devices over your local network.

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Cross--Platform-brightgreen)](https://github.com/satvikg7/twarit)

## ✨ Features

### 🚀 **Modern Upload Experience**

-   **Real-time Progress Tracking**: Watch your uploads with live progress bars, speed monitoring, and time estimates
-   **Drag & Drop Interface**: Intuitive file selection with modern UI/UX design
-   **Multiple File Support**: Upload multiple files simultaneously with batch processing
-   **Mobile Optimized**: Responsive design that works perfectly on phones and tablets

### 🔒 **Security & Privacy**

-   **Local Network Only**: Files never leave your network - no cloud storage involved
-   **Direct Transfer**: Peer-to-peer file sharing between your devices
-   **Secure Connection**: Built-in security measures for safe file transfers

### 🌐 **Cross-Platform Compatibility**

-   **Universal Access**: Works on any device with a web browser
-   **QR Code Connection**: Quick setup via QR code scanning
-   **No App Installation**: Pure web-based solution

### 📁 **File Management**

-   **File Explorer**: Browse and download files from your PC remotely
-   **File Type Detection**: Intelligent file type recognition with icons
-   **Large File Support**: Handle files of any size with optimized transfer

## 🛠️ Quick Start

### Prerequisites

-   [Go 1.21+](https://golang.org/dl/) installed on your system
-   [Task](https://taskfile.dev/) for build automation

### Installation

1. **Clone the repository**

    ```bash
    git clone https://github.com/satvikg7/twarit.git
    cd twarit
    ```

2. **Install dependencies**

    ```bash
    go mod download
    ```

3. **Build and run**

    ```bash
    task && ./twarit.exe
    ```

    _Or manually:_

    ```bash
    go run cmd/twarit/twarit.go
    ```

4. **Connect your devices**
    - Scan the QR code displayed in the terminal with your mobile device
    - Or visit the displayed URL on any device on the same network

## 📱 Usage

### Sending Files

1. Navigate to the **Send Files** section
2. Drag and drop files or click to browse
3. Monitor real-time upload progress
4. Files are automatically saved to the `Received` folder

### Browsing Files

1. Use the **Browse Files** section
2. Navigate through your PC's file system
3. Download files directly to your device
4. Files open in new tabs for immediate viewing

## 🏗️ Project Structure

```
twarit/
├── cmd/twarit/           # Main application entry point
├── internal/
│   ├── handler/          # HTTP request handlers
│   ├── router/           # Route definitions
│   └── server/           # Server configuration
├── pkg/qr/               # QR code generation
├── src/                  # Frontend assets
│   ├── index.html        # Homepage
│   ├── send.html         # Upload interface
│   ├── explore.html      # File browser
│   ├── main.css          # Homepage styles
│   ├── send.css          # Upload page styles
│   └── upload.js         # Upload functionality
└── Received/             # Uploaded files destination
```

## 🔧 Technical Details

### Backend (Go)

-   **Framework**: Native Go HTTP server with Gorilla WebSocket
-   **File Handling**: Multipart form processing with progress tracking
-   **APIs**: RESTful endpoints with real-time WebSocket updates

### Frontend (Vanilla JS)

-   **No Dependencies**: Pure HTML, CSS, and JavaScript
-   **Modern ES6+**: Classes, async/await, and modern web APIs
-   **Responsive Design**: CSS Grid and Flexbox layouts
-   **Progressive Enhancement**: Works without JavaScript for basic functionality

### Key Technologies

-   **WebSocket**: Real-time upload progress updates
-   **XMLHttpRequest**: File upload with progress monitoring
-   **FormData API**: Modern file handling
-   **CSS Grid/Flexbox**: Responsive layouts
-   **QR Code Generation**: Easy device connection

## 🚀 Development

### Adding New Features

1. **Backend Changes**: Add handlers in `internal/handler/`
2. **Frontend Updates**: Modify files in `src/`
3. **Routing**: Update `internal/router/router.go`

### Building for Production

```bash
task build
```

This creates optimized binaries for multiple platforms.

## 🤝 Contributing

Contributions are welcome! Here's how to get started:

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

### Development Guidelines

-   Follow Go best practices and conventions
-   Write clean, commented code
-   Test thoroughly on multiple devices
-   Ensure responsive design compatibility

## 📝 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

-   Built with ❤️ by [SatvikG7](https://github.com/satvikg7)
-   Inspired by the need for simple, secure file sharing
-   Thanks to the Go and web development communities

---

**Star ⭐ this repository if you find it helpful!**
