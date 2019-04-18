FROM golang:1.8 as build

WORKDIR /app/src/sample-app
COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -ldflags="-w -s" -a -installsuffix cgo -o /tmp/api

FROM alpine:latest

COPY --from=build /tmp/api /app/api

CMD /app/api

