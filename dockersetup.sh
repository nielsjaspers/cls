#!/bin/bash

# Define the Docker image and container name
IMAGE_NAME="cls"
CONTAINER_NAME="cls-server-docker"

# Check if the Docker image exists
if docker images | grep -q "^${IMAGE_NAME}\s"; then
  echo "Docker image '${IMAGE_NAME}' exists. Deleting..."
  docker rmi -f "${IMAGE_NAME}"
else
  echo "Docker image '${IMAGE_NAME}' does not exist. Proceeding..."
fi

# Build the Docker image
echo "Building Docker image '${IMAGE_NAME}' from Dockerfile..."
docker build -t "${IMAGE_NAME}" .

# Run the Docker container
echo "Running Docker container '${CONTAINER_NAME}'..."
docker run -d -p 443:443 --name "${CONTAINER_NAME}" "${IMAGE_NAME}"

# Confirm the container is running
if [ $? -eq 0 ]; then
  echo "Docker container '${CONTAINER_NAME}' is running successfully."
else
  echo "Failed to run Docker container '${CONTAINER_NAME}'. Check logs for errors."
fi

