# DevOps Info Service (Go)

A lightweight Go web service providing system and runtime information. This application demonstrates multi-stage Docker builds with compiled languages.

## Overview

This service exposes HTTP endpoints returning JSON data about the host system and Go runtime. It serves as the bonus task for Lab 02, showcasing the dramatic size reduction achievable with multi-stage builds for compiled languages.

## Prerequisites

- Go 1.22+ (for local development)
- Docker (for containerized deployment)

## Local Development

```bash
cd app_go
go run main.go
```

The server starts on `http://0.0.0.0:8080` by default.

## API Endpoints

### `GET /`

Returns comprehensive JSON with service metadata and system information.

```bash
curl http://localhost:8080/
```

### `GET /health`

Returns health status with timestamp and uptime.

```bash
curl http://localhost:8080/health
```

## Configuration

| Variable | Description  | Default |
|----------|--------------|---------|
| `PORT`   | Port number  | `8080`  |

## Docker

### Building the Image

```bash
docker build -t devops-info-service-go .
```

### Running the Container

```bash
docker run -p 8080:8080 devops-info-service-go
```

### Image Size

The multi-stage build with scratch base image results in an extremely small image:

| Image | Size |
|-------|------|
| devops-info-service-go | 6.72 MB |
| devops-info-service (Python) | 225 MB |

**Size reduction: 97%** (33x smaller)
