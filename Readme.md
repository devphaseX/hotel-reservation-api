#### Project: Hotel Reservation System

## Overview

The Hotel Reservation System is a Go-based web application designed to manage hotel reservations. It connects to a MongoDB database and utilizes JWT for authentication.

Prerequisites
Before setting up the project, ensure you have the following installed on your system:

Docker: Install Docker
Docker Compose (optional): Install Docker Compose
Go (if running locally): Install Go
MongoDB (for local development): Install MongoDB
Setting Up Environment Variables
Create a .env file in the root directory of your project to store environment-specific configurations. These variables will be used by the application to connect to the MongoDB database and set other configurations.

env

# .env file

```
MONGO_DB_NAME=hotel-reservation
MONGO_DB_URL=mongodb://localhost:27017
MONGO_DB_TEST_NAME=hotel-reservation-test
PORT=8080
JWT_SECRET=this is fucking secure
```

## Building the Docker Image

To containerize the application using Docker, follow these steps:

Navigate to the project directory where your Dockerfile is located.

Build the Docker image using the following command:

```
bash
docker build -t hotel-reservation-app .
```

This command creates a Docker image tagged hotel-reservation-app based on the instructions in the Dockerfile.

Running the Docker Container
After building the Docker image, you can run the application in a Docker container:

Run the container with the environment variables specified in the .env file. Use the following command:

bash

```
docker run --env-file .env -p 8080:8080 hotel-reservation-app
```

This command does the following:

Uses the .env file to set the environment variables inside the container.
Maps port 8080 on your local machine to port 8080 in the container.
Runs the hotel-reservation-app Docker image.
Testing the Application
To test the application, you can use tools like curl, Postman, or any browser to make HTTP requests to your API.

Access the application by opening a browser and navigating to:

http://localhost:8080

Additional Information
Development: If you prefer to run the application locally without Docker, make sure you have Go installed. You can start the application by setting the environment variables in your terminal and running the main Go file:

bash

```
export MONGO_DB_NAME=hotel-reservation
export MONGO_DB_URL=mongodb://localhost:27017
export MONGO_DB_TEST_NAME=hotel-reservation-test
export PORT=8080
export JWT_SECRET=this is fucking secure
```

```
go run main.go
```

Database: Ensure that MongoDB is running locally on mongodb://localhost:27017 or change the MONGO_DB_URL in your .env file to point to your MongoDB instance.

Security: Replace JWT_SECRET with a more secure and random string for production environments.
