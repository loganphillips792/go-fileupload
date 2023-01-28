Simple app to practice uploading files in React and GO

# Running the Application

## Frontend

- In a terminal:
    - cd frontend/fileupload
    - npm install
    - npm start
    - Access the UI at http://localhost:3000/app

## Backend

- In a terminal:
    - cd backend
    - go run main.go

## To view SQL Lite Database

1. Connect to a new SQL Lite database
2. JDBC URL: ```jdbc:sqlite:/Users/logan/go-fileupload/backend/data.db```
3. Path: ```/Users/logan/go-fileupload/backend/data.db```
4. Connect

# Uploading Images
## Through UI

1. Choose image to upload
2. If a name is provided, it will save the image with that name + original file extension. If there is no name provided, it will set the name to be the current UNIX timestamp.

## Through Curl Command

# Downloading files

# Running unit tests

1. go test -coverprofile=coverage.out
2. go tool cover -html=coverage.out

# Run in Docker

1. ```docker build -t getting-started-go --file=build/Dockerfile .```
2 ```docker run --publish 8080:8080 getting-started-go```
