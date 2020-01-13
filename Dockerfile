FROM golang

ARG releaseVersion="v0.3.2020"
ENV VERSION=$releaseVersion

## Go supports 'cross-complication' for many platforms - change goos=windows,linux,
##   (i.e. set to 'windows' and get an .exe)

## Setting these finals controls binary that is outputed.
ARG goos="linux"
ENV GOOS=$goos
ARG goarch="amd64"
ENV GOARCH=$goarch

## Dockerfile will build a binary release of tio.go
##
## Here's how to get a build going using this docker file:
##
## Using docker:
##    $ docker build --tag tiogo:v0.3.2020 .
##     ... (build output)
##
##    $ docker run --tty --interactive --rm tiogo:v0.3.2020
##     ... (docker drops you into working folder with a binary already built. :-)
##
##    root@4f51ab2342123:/tiogo# ./tio help
##
##    Get SecretKey and AccessKey from your Tenable.io instance:
##
##           https://docs.tenable.com/cloud/Content/Settings/GenerateAPIKey.htm
##
##
## NOTE: GOFLAGS won't be needed in go1.12 and beyond
##       This is what allows the hermetic builds - the whole 'vendor' folder holds the packages for this release
##
ARG goflags="-mod=vendor"
ENV GOFLAGS=$goflags

## 1. Make a directory and copy the code into it.
RUN mkdir /tiogo

ADD . /tiogo/

## 2. Move into the directory and start the build.
WORKDIR /tiogo

## 3A. Generate embeds the binaries of "config/" folder into our final build.
RUN go generate -tags release ./...
## 3B. Execute the release using
RUN go test -tags release  -v ./...

## Update debian pacakges and grab jq :-)
RUN DEBIAN_FRONTEND=noninteractive apt update
RUN DEBIAN_FRONTEND=noninteractive apt-get -y install jq

## -1 == -ONE!!!111
RUN GIT_HASH=$(git rev-list -1 master | cut -b1-8) && go build \
    -tags release \
    -ldflags \
    "-X github.com/whereiskurt/tiogo/internal/app/cmd/vm.ReleaseVersion=$VERSION \
    -X github.com/whereiskurt/tiogo/internal/app/cmd/vm.GitHash=$GIT_HASH" \
    -o ./tio \
    cmd/tio.go

## 4. Invoke tio vm helps for demonstration :-)
RUN ./tio help