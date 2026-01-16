# netboot

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.25.5-blue.svg)

**netboot** is a modern, lightweight PXE Boot Manager written in Go. It allows you to easily manage and boot ISO images over the network using iPXE. It features a great web interface built with the [IBM Carbon Design System](https://carbondesignsystem.com/).

## Features

-   🚀 **Lightweight Backend**: Built with Go, efficient and fast.
-   ✨ **UI**: Professional interface using IBM Carbon Design System (Gray 100 Dark Theme).
-   📂 **ISO Management**: Upload, list, and delete ISO images directly from the browser.
-   🔒 **Concurrency Safe**: Robust handling of concurrent requests.
-   🐳 **Docker Ready**: Easy deployment with Docker.
-   🛠 **iPXE Generation**: Automatically generates compatible iPXE boot scripts.

## Getting Started

### Prerequisites

-   Go 1.25+ (for building from source)
-   Docker (optional, for containerized deployment)

### Installation

#### From Source

1.  Clone the repository:
    ```bash
    git clone https://github.com/yourusername/netboot.git
    cd netboot
    ```

2.  Build the application:
    ```bash
    go build -o netboot ./cmd/netboot
    ```

3.  Run the server:
    ```bash
    ./netboot
    ```

4.  Access the web interface at `http://localhost:8080`.

#### Using Docker

```bash
docker build -t netboot .
docker run -p 8080:8080 -v $(pwd)/isos:/data/isos netboot
```

## Configuration

 The application is configured via Environment Variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server listening port | `:8080` |
| `ISO_DIR` | Directory to store ISO files | `./isos` |
| `UPLOAD_DIR` | Directory to temporary upload files | `./uploads` |

## Usage

1.  Open the web interface.
2.  Drag and drop an ISO file (e.g., `ubuntu.iso`, `proxmox.iso`) into the uploader.
3.  Configure your client machine to boot via network (PXE) pointing to this server.
4.  The client will load the iPXE script from `http://<server-ip>:8080/boot.ipxe`.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.