FROM golang:1.9

RUN mkdir -p /go/src/github.com/openfaas-incubator/faas-o6s/

WORKDIR /go/src/github.com/openfaas-incubator/faas-o6s

COPY . .

RUN gofmt -l -d $(find . -type f -name '*.go' -not -path "./vendor/*") && \
  VERSION=$(git describe --all --exact-match `git rev-parse HEAD` | grep tags | sed 's/tags\///') && \
  GIT_COMMIT=$(git rev-list -1 HEAD) && \
  CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w \
  -X github.com/openfaas-incubator/faas-o6s/pkg/version.Release=${VERSION} \
  -X github.com/openfaas-incubator/faas-o6s/pkg/version.SHA=${GIT_COMMIT}" \
  -a -installsuffix cgo -o faas-o6s .

FROM alpine:3.7

RUN addgroup -S app \
    && adduser -S -g app app \
    && apk --no-cache add ca-certificates

WORKDIR /home/app

COPY --from=0 /go/src/github.com/openfaas-incubator/faas-o6s/faas-o6s .

RUN chown -R app:app ./

USER app

ENTRYPOINT ["./faas-o6s"]
CMD ["-logtostderr"]
