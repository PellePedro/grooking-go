FROM golang as builder

# Build Arguments
ARG build
ARG version
ARG program

WORKDIR /src

COPY . /src

RUN   go mod download \
      && CGO_ENABLED=0 go build \
      -ldflags="-s -w -X main.Version=${version} -X main.Build=${build}" -o /app ${program}

FROM scratch
COPY --from=builder /app /app

EXPOSE 86

ENTRYPOINT ["./app"]
