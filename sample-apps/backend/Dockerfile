FROM golang:1.8 as build

WORKDIR /app/src/backend
COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -ldflags="-w -s" -a -installsuffix cgo -o /tmp/backend

FROM alpine:latest

COPY --from=build /tmp/backend /app/backend

CMD /app/backend
