# Image Pull-Through Cache (IPTC)
Welcome to the Image Pull-Through Cache project! 
This application is designed to act as a pull-through cache for Docker images, making your life easier by caching images locally. 
This project aims to heal the constant pain of those annoying 429 status codes from Docker Hub and speeds up access times for frequently used images.

## Features

- **Caching**: Uses an LRU (Least Recently Used) cache to store Docker images locally.
- **Super Simple Configuration**: Set up through environment variables and Kubernetes secrets.
- **Secure**: Handles sensitive information like Docker Hub tokens using Kubernetes secrets.
- **Detailed Logging**: Structured logging with `zap` because we build with observability in mind.
- **Health Check**: Simple endpoint to check if the service is running.

## Getting Started

### Prerequisites

Before you begin, make sure you have:

- **Go** (version 1.19 or later)
- **Docker** (for building and pushing Docker images)
- **Kubernetes** (you should have a running cluster to deploy this in, or just picture it in your head)
- **kubectl** (you should know what this is)


