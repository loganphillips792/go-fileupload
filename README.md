Simple app to practice uploading and downloading files in React and GO

# To Do

- Manitine components
- Convert to Keycloak for Authentication: https://mikebolshakov.medium.com/keycloak-with-go-web-services-why-not-f806c0bc820a
    - https://gruchalski.com/posts/2020-09-03-keycloak-with-docker-compose/
    - https://www.keycloak.org/getting-started/getting-started-docker
    - https://subscription.packtpub.com/book/cloud-and-networking/9781800562493/9/ch09lvl1sec59/integrating-with-golang-applications
    - https://levelup.gitconnected.com/building-micro-services-in-go-using-keycloak-for-authorisation-e00a29b80a43
- Use mantine cards: https://mantine.dev/core/card/
- Implement version tags and include in Docker image name: https://stackoverflow.com/questions/66017161/add-a-tag-to-a-docker-image-if-theres-a-git-tag-using-github-action
- Containerize frontend app so that it can run on K8s: https://www.knowledgehut.com/blog/web-development/how-to-dockerize-react-app
- Use React Query
- Convert from Styled Components to CSS Modules
- Allow users to download different size images 

# Running the Application

## Frontend

- In a terminal:
    - cd frontend/fileupload
    - npm install
    - npm run dev
    - Access the UI at http://localhost:3000/app

- To format the files, run ```npm run format```. This configuration is managed by the .prettierrc file
    - To see all options, please see https://prettier.io/docs/en/options.html

## Backend

Create ```.env``` file at root of backend. Look at .env.example to see the contents

- In a terminal:
    - cd backend
    - go run cmd/main.go

### Using Docker

1. `cd backend`
2. `docker build -t getting-started-go --file=./Dockerfile .`
3. `docker run --publish 8000:8000 getting-started-go`
4. Create account and login (see below)
5. `curl -H 'Cookie: user_session=<cookie-value-here>' 127.0.0.1:8000/api/hello`

Create an account:

If you don't create an account first and try to log in, you will get the following error:

`{"message":"missing key in cookies"}`

```
curl -X POST -H "Content-Type: application/json" \
     -d '{"username": "your_username", "password": "your_password"}' \
     http://127.0.0.1:8000/register/
```

Login:

```
curl -v -X POST -H "Content-Type: application/json" \
     -d '{"username": "your_username", "password": "your_password"}' \
     http://127.0.0.1:8000/login/
```

## To Run Postgres DB

1. `docker-compose up -d db`
    - Make sure that the DB port is set to 5433 in the env file (if you are running the Go app from outside a docker container)

# Uploading Images
## Through UI

1. Choose image to upload
2. If a name is provided, it will save the image with that name + original file extension. If there is no name provided, it will set the name to be the current UNIX timestamp.

## Through Curl Command

- `curl -X POST -F "file_name=example_name" -F "file=@/Users/logan/Downloads/profile.jpg" http://localhost:8000/uploadfile/`

## Through Postman

# Downloading files

## Through UI

CSV: Go to http://localhost:8000/download_csv/ in another tab

## Through Curl Command

CSV: `curl http://localhost:8000/download_csv/ --output output.csv`

## Through Postman

CSV: If you click 'Send', you will see the CSV contents, but the file will not download. Click 'Send and Download' to save the file locally.

# Running unit tests

1. go test -coverprofile=coverage.out
2. go tool cover -html=coverage.out

# Run in Docker

1. ```docker build -t getting-started-go --file=Dockerfile .```
2 ```docker run --publish 8080:8080 getting-started-go```

# Formatting project

`gofmt -s -w .`
