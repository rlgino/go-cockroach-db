# Go with cockroach DB

## Diagram
![Diagram](./docs/diagram.svg)

## Requests
```shell
# List users
curl --location '{{server}}/user'
# Create users
curl --location '{{server}}/user' \
--header 'Content-Type: application/json' \
--data '{
    "user": "Usuario1",
    "password": "pwd1"
}'
```

## Build a docker image for m1 microchip
```shell
 docker buildx build --platform linux/amd64 -t <api-tag> .
```

## Steps
1. Create a minimum app
2. Connect with Cockroach
3. Build a Docker image
4. Create a GCP Account
5. Enable GCR (google container registry) and upload image in it
6. Deploy in Cloud run
   1. Add Environment variables
   2. Generate Service account (What's this?)
   3. Check request permission
7. Enable load balancing
### Sources
* [Build a Simple CRUD Go App with CockroachDB and the Go pgx Driver](https://www.cockroachlabs.com/docs/stable/build-a-go-app-with-cockroachdb)