FROM golang:1.16-alpine AS build_base

# install git and openssh
RUN apk update && apk add git
RUN apk add openssh

# make sure your domain is accepted
RUN mkdir -p /root/.ssh/
RUN touch /root/.ssh/known_hosts
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts

# set go env
RUN export GIT_SSH_COMMAND="ssh -i /root/.ssh/id_rsa"

# Set the Current Working Directory inside the container
WORKDIR /tmp/backend

# copy files.
COPY . .
RUN go mod download

COPY . .

# Unit tests
# RUN CGO_ENABLED=0 go test -v

# Build the Go app
RUN go build -o ./out/backend .

# Start fresh from a smaller image
FROM alpine:3.9 

COPY --from=build_base /tmp/backend/out/backend /app/backend

# This container exposes port 8083 to the outside world
EXPOSE 8082

# Run the binary program produced by `go install`
CMD ["/app/backend"]