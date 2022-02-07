# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.17-buster AS build

WORKDIR /app

# Create dir which can be copied and chowned to app image
Run mkdir db/

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /ipwatcher

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build --chown=nonroot:nonroot /app/db/ /db/

COPY --from=build /ipwatcher /ipwatcher

USER nonroot:nonroot

VOLUME /db

ENTRYPOINT ["/ipwatcher"]