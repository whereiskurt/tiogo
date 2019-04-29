FROM golang

ARG releaseVersion="v0.1.0"
ENV VERSION=$releaseVersion

ARG hash="0xABCD1234"
ENV HASH=$hash

ARG goos="linux"
ENV GOOS=$goos

## NOTE: GOFLAGS won't be needed in go1.12
ARG goflags="-mod=vendor"
ENV GOFLAGS=$goflags

RUN mkdir /tiogo

ADD . /tiogo/

WORKDIR /tiogo

RUN go test -v ./...

RUN go build \
    -tags release \
    --ldflags \
    "-X internal/app/cmd/vm.ReleaseVersion=$VERSION \
     -X internal/app/cmd/vm.GitHash=$HASH" \
    -o ./tio \
    cmd/tio.go