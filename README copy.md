This repository contains a Golang project that is designed to run in both development and production environments using Docker and Docker Compose.

# Prerequisites
Docker installed on your machine.
Docker Compose installed on your machine.

# Getting Started
## Development Mode
To start the application in development mode with live reloading, use the following command:

```bash
make run
```
This will build the Docker image and start the container with live reloading enabled using **air** library. Any changes you make to the source code will trigger a rebuild and restart of the service.

## Production Mode
To deploy the application in production mode, use the following command:

```bash
make build TAG=myTAG #to build the image

make run-prod #to build and run the application
```

This command will build the optimized Docker image for production and start the container. The application will be accessible on the **3000** port.
