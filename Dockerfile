FROM golang:1.23-alpine AS build-env

WORKDIR /go/src

ENV TEMPL_VERSION="v0.2.771"

# install dependencies
RUN go install github.com/a-h/templ/cmd/templ@${TEMPL_VERSION}

# install dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN --mount=type=cache,target=/root/.cache/templ_gen \
    TEMPL_EXPERIMENT=rawgo templ generate
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o main main.go

# Start from alpine
FROM alpine:3.20
WORKDIR /app

# Copy over the built binary from our builder stage
COPY --from=build-env /go/src/main main
# Copy over public files from environment
COPY ./public public/

EXPOSE 3000
ENTRYPOINT ["/app/main"]
