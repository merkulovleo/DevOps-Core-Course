# DevOps Info Service

A Python web application built with Flask that provides system and runtime information via REST API endpoints.

## Overview

This service exposes HTTP endpoints returning JSON data about the host system, Python runtime, and incoming request metadata. It is designed as part of Lab 01 for the DevOps Engineering course.

## Prerequisites

- Python 3.10+
- pip

## Installation

```bash
cd app_python
python -m venv venv
source venv/bin/activate   # On Windows: venv\Scripts\activate
pip install -r requirements.txt
```

## Running the Application

```bash
python app.py
```

The server starts on `http://0.0.0.0:5000` by default.

## API Endpoints

### `GET /`

Returns comprehensive JSON with service metadata, system information, runtime details, request data, and available endpoints.

```bash
curl http://localhost:5000/
```

### `GET /health`

Returns health status with timestamp and uptime.

```bash
curl http://localhost:5000/health
```

## Configuration

The application is configurable via environment variables:

| Variable | Description           | Default   |
|----------|-----------------------|-----------|
| `HOST`   | Bind address          | `0.0.0.0` |
| `PORT`   | Port number           | `5000`    |
| `DEBUG`  | Enable debug mode     | `false`   |

Example:

```bash
HOST=127.0.0.1 PORT=8080 DEBUG=true python app.py
```
