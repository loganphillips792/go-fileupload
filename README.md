Simple app to practice uploading and downloading files in React and GO

# To Do

- Manitine components
- Convert to Keycloak for Authentication: https://mikebolshakov.medium.com/keycloak-with-go-web-services-why-not-f806c0bc820a
    - https://gruchalski.com/posts/2020-09-03-keycloak-with-docker-compose/
    - https://www.keycloak.org/getting-started/getting-started-docker
    - https://subscription.packtpub.com/book/cloud-and-networking/9781800562493/9/ch09lvl1sec59/integrating-with-golang-applications
    - https://levelup.gitconnected.com/building-micro-services-in-go-using-keycloak-for-authorisation-e00a29b80a43
- Use mantine cards: https://mantine.dev/core/card/

# Running the Application

## Frontend

- In a terminal:
    - cd frontend/fileupload
    - npm install
    - npm run dev
    - Access the UI at http://localhost:3000/app

## Backend

Create ```.env``` file at root of backend. Look at .env.example to see the contents

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

## Through Postman

# Downloading files

## Through UI

CSV: Go to http://localhost:8000/download_csv/ in another tab

## Through Curl Command

CSV: ```curl http://localhost:8000/download_csv/ --output l.download```

## Through Postman

CSV: If you click 'Send', you will see the CSV contents, but the file will not download. Click 'Send and Download' to save the file locally.

# Running unit tests

1. go test -coverprofile=coverage.out
2. go tool cover -html=coverage.out

# Run in Docker

1. ```docker build -t getting-started-go --file=Dockerfile .```
2 ```docker run --publish 8080:8080 getting-started-go```