FROM golang:1.16-alpine as base
WORKDIR /src

COPY go.mod go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -a -o /app

FROM alpine
RUN apk --no-cache add ca-certificates

VOLUME /assets
COPY ./assets /assets
COPY --from=base /app /app 

ENV PORT=8005
EXPOSE 8005

ENTRYPOINT "/app"
CMD ""