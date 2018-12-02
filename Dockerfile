FROM golang:1.11-alpine as golang

# FROM will allow us to run a docker container interactive
#docker run -it --rm golang:1.11-alpine

## Could add this, but not sure we need it.
## RUN apk add git
## RUN apk --no-cache add ca-certificates

