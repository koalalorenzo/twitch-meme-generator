FROM golang:1.17 as base

RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /src
COPY go.mod go.sum /src/
RUN go mod download

COPY . /src/
RUN go build -ldflags '-extldflags "-static"' -a -o /app

FROM scratch

# Ensuring that generate images are cached in a temporary volume
VOLUME /tmp
ENV KTMG_TEMP_PATH="/tmp"

VOLUME /assets

COPY ./assets /assets
COPY --from=base /etc/ssl/certs /etc/ssl/certs
COPY --from=base /app /app

ENV PORT=8005
EXPOSE 8005

ENTRYPOINT ["/app"]
CMD []
