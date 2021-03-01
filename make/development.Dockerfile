# build the go binary
FROM golang:latest as build

# private repo pkg's
RUN go env -w GOPRIVATE=git.iconmobile.com/shared

ARG service_name

WORKDIR /go/src/github.com/iconmobile-dev/backend-coding-challenge/

# Copy everything
COPY . .

# Build the Go service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/${service_name} cmd/${service_name}/main.go

# build final alpine image
FROM alpine:3.12.0

# allow user and source root to be passed as args at default to sensibles
ARG service_dir=/service/
ARG service_port
ARG service_name

# install required libs
RUN apk update && apk --no-cache --update add ca-certificates

# set local directory
WORKDIR ${service_dir}

# copy final go binary from the build stage
COPY --from=build /build/${service_name} ${service_dir}/${service_name}

# set SERVICE_NAME env variable
ENV SERVICE_NAME ${service_name}

# finally expose the port and run the process
EXPOSE ${service_port}
CMD ./${SERVICE_NAME}