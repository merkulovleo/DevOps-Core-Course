# Lab 02: Docker Containerization

## Overview

This document describes the Dockerization process for the DevOps Info Service application.

## Docker Best Practices Applied

### 1. Multi-Stage Build (Bonus)

The Dockerfile uses a multi-stage build approach:

- **Stage 1 (builder)**: Installs dependencies in a virtual environment
- **Stage 2 (runtime)**: Copies only the virtual environment and application code

**Benefits**:
- Smaller final image size (no build tools or cache)
- Cleaner separation of build and runtime environments
- Improved security (fewer packages in final image)

### 2. Non-Root User

```dockerfile
RUN groupadd --gid 1000 appgroup && \
    useradd --uid 1000 --gid 1000 --shell /bin/bash --create-home appuser
USER appuser
```

**Why**: Running as non-root prevents potential container escape attacks and limits damage if the application is compromised.

### 3. Specific Base Image Version

```dockerfile
FROM python:3.13-slim
```

**Why**:
- `python:3.13-slim` is a minimal image (~120 MB vs ~1 GB for full image)
- Specific version ensures reproducible builds
- `slim` variant excludes unnecessary packages, reducing attack surface

### 4. Layer Ordering and Caching

```dockerfile
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY app.py .
```

**Why**: Copying `requirements.txt` before application code allows Docker to cache the dependency installation layer. When only application code changes, Docker reuses the cached dependencies layer.

### 5. Minimal File Copying

Only necessary files are copied to the image:
- `requirements.txt` (for dependencies)
- `app.py` (application code)

Excluded via `.dockerignore`:
- Virtual environment (`venv/`)
- Tests (`tests/`)
- Documentation (`docs/`, `*.M@`)
- Git files (`.git/`)
- IDE files (`.idea/`, `.vscode/`)
- Python cache (`__pycache__/`)

### 6. Health Check

```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD python -c "import urllib.request; urllib.request.urlopen('http://localhost:5000/health')" || exit 1
```

**Why**: Allows orchestrators (Docker Swarm, Kubernetes) to monitor container health and restart unhealthy containers automatically.

## Image Information

| Property | Value |
|----------|------|
| Base Image | `python:3.13-slim` |
| Final Size | 225 MB |
| Exposed Port | 5000 |
| User | `appuser` (non-root) |
| Working Directory | `/app` |

## Build Process

### Build the image

```bash
docker build -t devops-info-service .
```

Terminal output:
```
[+] Building 23.9s (15/15) FINISHED                                            docker:desktop-linux
 => [internal] load build definition from Dockerfile                                           0.0s
 => => transferring dockerfile: 1.24kB                                                         0.0s
 => [internal] load metadata for docker.io/library/python:3.13-slim                            5.2s
 => [internal] load .dockerignore                                                              0.0s
 => => transferring context: 543B                                                              0.0s
 => [builder 1/5] FROM docker.io/library/python:3.13-slim@sha256:49b618b8afc2742b94fa8419d8f  12.4s
 => [internal] load build context                                                              0.0s
 => => transferring context: 3.17kB                                                            0.0s
 => [runtime 2/6] RUN groupadd --gid 1000 appgroup &&     useradd --uid 1000 --gid 1000 --she  0.3s
 => [builder 2/5] WORKDIR /build                                                               0.2s
 => [builder 3/5] RUN python -m venv /opt/venv                                                 1.5s
 => [runtime 3/6] WORKDIR /app                                                                 0.0s
 => [builder 4/5] COPY requirements.txt .                                                      0.0s
 => [builder 5/5] RUN pip install --no-cache-dir -r requirements.txt                           3.8s
 => [runtime 4/6] COPY --from=builder /opt/venv /opt/venv                                      0.1s
 => [runtime 5/6] COPY app.py .                                                                0.0s
 => [runtime 6/6] RUN chown -R appuser:appgroup /app                                           0.1s
 => exporting to image                                                                         0.5s
 => => exporting layers                                                                        0.4s
 => => naming to docker.io/library/devops-info-service:latest                                  0.0s
```

### Run the container

```bash
docker run -p 5000:5000 devops-info-service
```

Terminal output:
```
2026-02-05 12:06:14,599 [INFO] __main__: Starting DevOps Info Service on 0.0.0.0:5000 (debug=False)
 * Serving Flask app 'app'
 * Debug mode: off
 * Running on all addresses (0.0.0.0)
 * Running on http://127.0.0.1:5000
 * Running on http://172.17.0.2:5000
```

### Verify the application

```bash
curl http://localhost:5000/
```

Response:
```json
{
  "endpoints": [
    {"description": "Service info and metadata", "method": "GET", "path": "/"},
    {"description": "Health check", "method": "GET", "path": "/health"}
  ],
  "request": {
    "client_ip": "192.168.65.1",
    "method": "GET",
    "path": "/",
    "user_agent": "curl/8.7.1"
  },
  "runtime": {
    "current_time": "2026-02-05T12:06:19.962139+00:00",
    "python_version": "3.13.12",
    "timezone": "UTC",
    "uptime_seconds": 5.36
  },
  "service": {
    "description": "A web service providing system and runtime information",
    "name": "DevOps Info Service",
    "version": "1.0.0"
  },
  "system": {
    "architecture": "aarch64",
    "cpu_count": 12,
    "hostname": "9fa1c9cfdd6a",
    "platform": "Linux",
    "platform_version": "#1 SMP Thu Mar 20 16:32:56 UTC 2025"
  }
}
```

### Verify non-root user

```bash
docker exec <container_id> whoami
```

Output:
```
appuser
```

### Image size

```bash
docker images devops-info-service
```

Output:
```
REPOSITORY            TAG       IMAGE ID       CREATED          SIZE
devops-info-service   latest    e1219da81235   19 seconds ago   225MB
```

## Docker Hub

### Repository URL

```
https://hub.docker.com/r/merkulovleo/devops-info-service
```

### Push Commands

```bash
docker tag devops-info-service merkulovleo/devops-info-service:latest
docker push merkulovleo/devops-info-service:latest
```

## Bonus: Multi-Stage Build with Compiled Language (Go)

To demonstrate the full power of multi-stage builds, a Go version of the service was created in `app_go/`.

### Go Dockerfile Strategy

```dockerfile
# Stage 1: Builder - compile the Go application
FROM golang:1.22-alpine AS builder

WORKDIR /build
COPY go.mod ./
RUN go mod download
COPY main.go .

# Build static binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server main.go


# Stage 2: Runtime - minimal scratch image
FROM scratch AS runtime

COPY --from=builder /app/server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
```

### Key Optimizations

1. **Static Binary**: `CGO_ENABLED=0` creates a fully static binary with no external dependencies
2. **Strip Debug Info**: `-ldflags="-s -w"` removes debug symbols, reducing binary size
3. **Scratch Base**: Using `scratch` (empty image) as the final base - only the binary is included

### Size Comparison

| Image | Base | Size |
|-------|------|------|
| devops-info-service (Python) | python:3.13-slim | 225 MB |
| devops-info-service-go | scratch | **6.72 MB** |

**Size reduction: 97%** (33x smaller than Python version)

### Build Output

```
[+] Building 8.8s (12/12) FINISHED                          docker:desktop-linux
 => [builder 1/6] FROM docker.io/library/golang:1.22-alpine                 2.6s
 => [builder 2/6] WORKDIR /build                                            0.2s
 => [builder 3/6] COPY go.mod ./                                            0.0s
 => [builder 4/6] RUN go mod download                                       0.1s
 => [builder 5/6] COPY main.go .                                            0.0s
 => [builder 6/6] RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w"    2.9s
 => [runtime 1/1] COPY --from=builder /app/server /server                   0.0s
 => exporting to image                                                      0.2s
```

### Verification

```bash
$ curl http://localhost:8080/
{
  "service": {
    "name": "DevOps Info Service (Go)",
    "version": "1.0.0",
    "description": "A web service providing system and runtime information"
  },
  "system": {
    "hostname": "9107f38b59ef",
    "platform": "linux",
    "architecture": "arm64",
    "num_cpu": 12,
    "go_version": "go1.22.12"
  },
  "runtime": {
    "uptime_seconds": 4.37,
    "current_time": "2026-02-05T12:23:32Z",
    "timezone": "UTC"
  }
}
```

### Security Benefits of Smaller Images

1. **Minimal Attack Surface**: The scratch image contains nothing but the binary - no shell, no utilities, no package manager
2. **No CVEs from Base OS**: Since there's no OS layer, there are no OS-level vulnerabilities
3. **Cannot Shell Into Container**: Attackers cannot get a shell even if they exploit the application
4. **Faster Scanning**: Vulnerability scanners complete almost instantly
5. **Immutable Runtime**: The container cannot be modified at runtime

## Multi-Stage Build Analysis (Python)

### Size Comparison

| Build Type | Size |
|------------|------|
| Single-stage (python:3.13) | ~1 GB |
| Single-stage (python:3.13-slim) | ~180 MB |
| Multi-stage (python:3.13-slim) | 225 MB |

### Security Benefits

1. **Reduced Attack Surface**: Build tools and development dependencies are not included in the final image
2. **Fewer Vulnerabilities**: Less packages mean fewer potential CVEs
3. **No Build Artifacts**: Source code and build cache are not exposed

## Challenges and Solutions

### Challenge 1: Virtual Environment Path

**Problem**: Initially, dependencies installed globally were not copied correctly between stages.

**Solution**: Used a virtual environment at a fixed path (`/opt/venv`) that can be copied as a single directory.

### Challenge 2: Health Check without curl

**Problem**: `python:3.13-slim` doesn't include `curl` or `wget`.

**Solution**: Used Python's built-in `urllib` module for health checks.

### Challenge 3: File Permissions

**Problem**: Files copied to the container are owned by root by default.

**Solution**: Added `chown -R appuser:appgroup /app` before switching to the non-root user.