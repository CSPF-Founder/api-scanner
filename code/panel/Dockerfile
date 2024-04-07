FROM golang:1.22.1-bookworm AS build-stage

RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid 1000 \
  app


WORKDIR /app/frontend

RUN apt install -y curl

RUN curl -fsSL https://deb.nodesource.com/setup_21.x | bash -
RUN apt-get install -y nodejs

COPY frontend/ ./

RUN npm install

RUN npm run build

# create a working directory inside the image
WORKDIR /app/panel

# copy Go modules and dependencies to image
COPY go.mod go.sum ./

# download Go modules and dependencies
RUN go mod download

# Install SQLite development package
RUN apt-get update && apt-get install -y libsqlite3-dev

# copy the source code into the image
COPY . .

RUN apt-get update && \
  apt-get install --no-install-recommends -y jq libcap2-bin ca-certificates && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN update-ca-certificates

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags=release -ldflags="-s -w" -o /app/panel/server ./cmd

# # Run the tests in the container
# FROM build-stage AS run-test-stage
# RUN go test -v ./...


# FROM debian:stable-slim
FROM alpine:latest

COPY --from=build-stage /etc/passwd /etc/passwd
COPY --from=build-stage /etc/group /etc/group
COPY --from=build-stage /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# RUN apt-get update && \
#   apt-get install --no-install-recommends -y jq libcap2-bin ca-certificates && \
#   apt-get clean && \
#   rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

WORKDIR /app/panel
COPY --from=build-stage /app/panel/server /app/panel/server

# Copy static files (todo!!, distribute only the dist folder)
# COPY ./static /app/panel/static
# DB migration folder
COPY ./db/migrations /app/panel/db/migrations

# Templates folder
# COPY ./templates /app/panel/templates

RUN chown -R app:app /app/

USER app

EXPOSE 443

ENTRYPOINT ["/app/panel/server"]
