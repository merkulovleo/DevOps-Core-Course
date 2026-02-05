# Lab 01: DevOps Info Service

## Framework Choice

**Framework:** Flask

**Rationale:** Flask is a lightweight, well-documented micro-framework that provides the minimal scaffolding needed for a REST API without imposing unnecessary structure. It has a large ecosystem, extensive community support, and is straightforward to set up for a service with a small number of endpoints. Compared to Django (which is heavier and designed for larger applications), Flask allows fine-grained control with minimal boilerplate. Compared to FastAPI, Flask has broader adoption and does not require an ASGI server, making deployment simpler for this use case.

## Best Practices Applied

### 1. Clean Code Organization (PEP 8)

The application follows PEP 8 style guidelines:
- Consistent 4-space indentation
- Descriptive variable and function names
- Module-level docstring and function docstrings
- Imports grouped by standard library, then third-party packages

### 2. Proper Error Handling

Custom error handlers are registered for HTTP 404 and 500 errors, returning structured JSON responses instead of default HTML error pages:

```python
@app.errorhandler(404)
def not_found(error):
    return jsonify({"error": "Not Found", "path": request.path}), 404
```

### 3. Logging

The application uses Python's built-in `logging` module with a configured format including timestamps, log level, and logger name. Log level is controlled via the `DEBUG` environment variable.

### 4. Pinned Dependencies

All dependencies in `requirements.txt` are pinned to exact versions (`flask==3.1.0`) to ensure reproducible builds.

### 5. Configuration via Environment Variables

The application reads `HOST`, `PORT`, and `DEBUG` from environment variables with sensible defaults, following the twelve-factor app methodology.

## API Documentation

### GET /

Returns service metadata, system information, runtime details, and request data.

```bash
curl http://localhost:5000/
```

Example response:

```json
{
  "service": {
    "name": "DevOps Info Service",
    "version": "1.0.0",
    "description": "A web service providing system and runtime information"
  },
  "system": {
    "hostname": "my-machine",
    "platform": "Darwin",
    "platform_version": "...",
    "architecture": "arm64",
    "cpu_count": 8
  },
  "runtime": {
    "python_version": "3.12.0",
    "uptime_seconds": 42.5,
    "current_time": "2026-01-28T12:00:00+00:00"
  },
  "request": {
    "client_ip": "127.0.0.1",
    "user_agent": "curl/8.0",
    "method": "GET",
    "path": "/"
  },
  "endpoints": [
    {"path": "/", "method": "GET", "description": "Service info and metadata"},
    {"path": "/health", "method": "GET", "description": "Health check"}
  ]
}
```

### GET /health

Returns health check status.

```bash
curl http://localhost:5000/health
```

Example response:

```json
{
  "status": "healthy",
  "timestamp": "2026-01-28T12:00:00+00:00",
  "uptime_seconds": 42.5
}
```

## Screenshots

Screenshots demonstrating working endpoints are located in the `screenshots/` directory.

## Challenges and Solutions

1. **Timezone-aware timestamps:** Used `datetime.now(timezone.utc)` instead of `datetime.utcnow()` (deprecated in Python 3.12+) to produce timezone-aware ISO 8601 timestamps.

2. **Configuration flexibility:** Implemented environment variable configuration with sensible defaults so the application works out of the box while remaining configurable for different environments.

## GitHub Community Engagement

- Starred the course repository and the `simple-container-com/api` project.
- Followed the course instructors and classmates on GitHub.
- **Why starring and following matter in open source:** Starring repositories signals interest and helps maintainers gauge community engagement. It also helps other developers discover useful projects through trending lists and recommendations. Following contributors keeps you informed about their new projects and activity, fostering collaboration and knowledge sharing within the community.
