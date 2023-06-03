Simple app to practice uploading and downloading files in React and GO

# To Do

- Convert to Vite
- React Drag and Drop file upload
- Manitine components
- Convert to Keycloak for Authentication: https://mikebolshakov.medium.com/keycloak-with-go-web-services-why-not-f806c0bc820a
    - https://gruchalski.com/posts/2020-09-03-keycloak-with-docker-compose/
    - https://www.keycloak.org/getting-started/getting-started-docker
    - https://subscription.packtpub.com/book/cloud-and-networking/9781800562493/9/ch09lvl1sec59/integrating-with-golang-applications
- Make image cards look like: https://mantine.dev/core/card/

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
DB_HOST=db
DB_NAME=fileupload
DB_PASSWORD=12345
DB_USERNAME=logan
DB_PORT=5432
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

## To Run Postgres DB

1. `docker-compose up -d db`
    - Make sure that the DB port is set to 5433 in the env file (if you are running the Go app from outside a docker container)

# Uploading Images
## Through UI

1. Choose image to upload
2. If a name is provided, it will save the image with that name + original file extension. If there is no name provided, it will set the name to be the current UNIX timestamp.

## Through Curl Command

- `curl -X POST -F "file_name=example_name" -F "file=@48yI8S4.jpeg" http://localhost:8000/uploadfile/`

# Downloading files

## Through UI

1. Go to http://localhost:8000/download_csv/ in another tab

# Running unit tests

1. go test -coverprofile=coverage.out
2. go tool cover -html=coverage.out

# Run in Docker

1. ```docker build -t getting-started-go --file=Dockerfile .```
2 ```docker run --publish 8080:8080 getting-started-go```

# Run in Docker-Compose

1. Go to root of directory
2. docker-compose up

If you make any code changes, then run the following command to rebuild the code through docker-compose: docker-compose up backend --build

Entering psql

- docker exec -it go-fileupload-db-1 bash
- psql -U logan -d fileupload