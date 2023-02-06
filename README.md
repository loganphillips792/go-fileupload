Simple app to practice uploading and downloading files in React and GO

# To Do

- Convert to Vite
- React Drag and Drop file upload
- Manitine components

# Running the Application

## Frontend

- In a terminal:
    - cd frontend/fileupload
    - npm install
    - npm start
    - Access the UI at http://localhost:3000/app

## Backend

Create ```.env``` file at root of backend

```
GORILLA_SESSIONS_HASH_KEY=<some_hash_key_here>
GORILLA_SESSIONS_BLOCK_KEY=<some_block_key_here>
```

- In a terminal:
    - cd backend
    - go run cmd/main.go

### Using Docker

1. ```cd backend```
2. ```docker build -t getting-started-go --file=./Dockerfile .```
3. ```docker run --publish 8000:8000 getting-started-go```
4. ```curl 127.0.0.1:8000/api/hello```

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
