# pausarr

A lightweight automation tool for qBittorrent, designed to automatically pause torrents based on configurable seeding ratio, seeding time, and tracker privacy settings. This project is written in Go and is easily deployable using Docker.

## Features
- Automatically pauses torrents that exceed a specified ratio or seeding time
- Optionally skips private trackers
- Skips torrents tagged with `no-pause`
- Configurable check intervals
- Designed for easy deployment with Docker

## How It Works
`pausarr` connects to your qBittorrent instance and periodically checks all uploading torrents. If a torrent exceeds the configured maximum ratio or maximum seeding duration, and does not have the `no-pause` tag, it will be paused. Torrents from private trackers can be excluded from pausing.

## Tag Exclusion
Torrents tagged with `no-pause` will be ignored by the automation and will not be paused, regardless of their ratio or seed time.

## Environment Variables
The application is configured via environment variables. These can be set in your `.env` file or directly in your Docker configuration.

| Variable                | Description                                                      | Example Value                                  | Required |
|-------------------------|------------------------------------------------------------------|------------------------------------------------|----------|
| `QBITTORRENT_URL`       | URL to your qBittorrent Web UI                                   | `https://your-qbittorrent.example.com`         | Yes      |
| `QBITTORRENT_USERNAME`  | Username for qBittorrent                                         | `admin`                                        | Yes      |
| `QBITTORRENT_PASSWORD`  | Password for qBittorrent                                         | `yourpassword`                                 | Yes      |
| `MAX_RATIO`             | Maximum allowed seeding ratio before pausing                     | `2.0`                                          | No (default: 2.0) |
| `MAX_SEED_HOURS`        | Maximum allowed seeding time in hours before pausing             | `4320` (6 months)                              | No (default: 4320) |
| `SKIP_PRIVATE_TRACKERS` | Skip pausing torrents from private trackers (`true` or `false`)  | `true`                                         | No (default: false) |
| `CHECK_INTERVAL_MINUTES`| How often to check torrents, in minutes                          | `10`                                           | No (default: 60) |

### Example `.env`
```
QBITTORRENT_URL=https://your-qbittorrent.example.com
QBITTORRENT_USERNAME=admin
QBITTORRENT_PASSWORD=yourpassword
MAX_RATIO=2.0
MAX_SEED_HOURS=4320
SKIP_PRIVATE_TRACKERS=true
CHECK_INTERVAL_MINUTES=10
```

## Docker Images from GHCR

Official images are published to GitHub Container Registry (GHCR):

```
docker pull ghcr.io/felipemarinho97/pausarr:latest
```

- `latest` is always the most recent build from the default branch.
- Versioned tags (e.g., `v1.0.0`) are also available if a release is published.

You can use these images directly in your deployments without building locally.

### Example: Running from GHCR

```sh
docker run --env-file .env --restart unless-stopped ghcr.io/felipemarinho97/pausarr:latest
```

Or with environment variables inline:

```sh
docker run -e QBITTORRENT_URL=https://your-qbittorrent.example.com \
           -e QBITTORRENT_USERNAME=admin \
           -e QBITTORRENT_PASSWORD=yourpassword \
           -e MAX_RATIO=2.0 \
           -e MAX_SEED_HOURS=4320 \
           -e SKIP_PRIVATE_TRACKERS=true \
           -e CHECK_INTERVAL_MINUTES=10 \
           --restart unless-stopped \
           ghcr.io/felipemarinho97/pausarr:latest
```

## Docker Deployment

### 1. Using Docker Compose (Recommended)

Create a `docker-compose.yaml` file in your project directory:

```yaml
version: '3.8'
services:
  pausarr:
    image: ghcr.io/felipemarinho97/pausarr:latest
    environment:
      - QBITTORRENT_URL=https://your-qbittorrent.example.com
      - QBITTORRENT_USERNAME=admin
      - QBITTORRENT_PASSWORD=yourpassword
      - MAX_RATIO=2.0
      - MAX_SEED_HOURS=4320
      - SKIP_PRIVATE_TRACKERS=true
      - CHECK_INTERVAL_MINUTES=10
    restart: unless-stopped
```

Then run the container:

```sh
docker-compose up -d
```


## Logging
The application will output logs to stdout, including information about configuration, torrent checks, and actions taken.

## Contributing
Pull requests and issues are welcome! Please open an issue if you encounter bugs or have feature requests.

## License
MIT License
