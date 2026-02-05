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
| Final Size | ~150 MB |
| Exposed Port | 5000 |
| User | `appuser` (non-root) |
| Working Directory | `/app` |

## Build Process

### Build the image

```bash
docker build -t devops-info-service .
```

Expected output:
```
[...] => [builder 1/4] FROM docker.io/library/python:3.13-slim
[...] => [builder 2/4] WORKDIR /build
[...] => [builder 3/4] RUN python -m venv /opt/venv
[...] => [builder 4/4] COPY requirements.txt .
[...] => [builder 5/5] RUN pip install --no-cache-dir -r requirements.txt
[...] => [runtime 1/5] FROM docker.io/library/python:3.13-slim
[...] => [runtime 2/5] RUN groupadd --gid 1000 appgroup ...
[...] => [runtime 3/5] WORKDIR /app
[...] => [runtime 4/5] COPY --from=builder /opt/venv /opt/venv
[...] => [runtime 5/5] COPY app.py .
[...] => exporting to image
```

### Run the container

```bash
docker run -p 5000:5000 devops-info-service
```

Expected output:
```
2025-01-28 12:00:00,000 [INFO] __main__: Starting DevOps Info Service on 0.0.0.0:5000 (debug=False)
 * Serving Flask app 'app'
 * Running on http://0.0.0.0:5000
```

### Verify the application

```bash
curl http://localhost:5000/
```

Expected JSON response with service information.

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

## Multi-Stage Build Analysis (Bonus)

### Size Comparison

| Build Type | Estimated Size |
|------------|----------------|
| Single-stage (python:3.13) | ~1 GB |
| Single-stage (python:3.13-slim) | ~180 MB |
| Multi-stage (python:3.13-slim) | ~150 MB |

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