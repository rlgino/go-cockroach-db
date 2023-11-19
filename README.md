# Go with cockroach DB (Users API)

## Microservices architecture (Approach)

### Authentication and Authorization:
Microservice responsible for handling user authentication and authorization to access specific wallet functions.

### User Management:
Handles the creation, updating, and deletion of user accounts, as well as the management of profile information.

### Account and Balance Management:
Manages accounts associated with the wallet and performs operations related to the balance, such as balance inquiries, fund loading, and transfers between accounts.

### Transaction History:
Records and manages the history of all transactions, providing detailed information about each operation.

### Payment Processing:
Responsible for processing payments, integrating with external payment gateways if necessary, and updating the user's account balance.

### Notifications:
Microservice that manages the sending of notifications, such as transaction alerts, security reminders, and other important communications.

### Security:
Handles the implementation of security measures, such as data encryption, fraud prevention, and the management of security tokens.

## Technologies and design

### About structure

### Plain HTTP server
For this entrypoint I used `http/net` package for serve HTTP requests, it's a simple way to develop an application in Golang. The `main.go` file is in `cmd/plain`.

### OpenAPI Generator
One of the most often used open-source libraries for using an OAS file is OpenAPI Generator. The purpose of it is to produce documentation for OAS 2.0 and OAS 3.0 papers. These documents can be altered using options, unique templates, and unique generators on your classpath.
This API documentation tool can automatically produce API documentation from source code. Java, Node.js, Python, PHP, Ruby, and NET are just a few of the many programming languages and frameworks it supports.
For this integration with Go I'm using [Ogen](https://ogen.dev/) which I found simple and useful.

#### How to use?
Creating and editing from [Swagger Editor](https://editor.swagger.io/) and then copy and paste my Swagger into the project and execute `go gen ./...` in the folder `cmd/ogen`.

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