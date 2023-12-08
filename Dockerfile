FROM golang:1.20 AS build-stage
# Set destination for COPY
WORKDIR /app
# Download Go modules
COPY go.mod go.sum ./
RUN go mod download
# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . ./
# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-users-service-app ./cmd/ogen
# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /go-users-service-app /go-users-service-app

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/go-users-service-app"]