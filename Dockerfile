FROM golang:latest AS build
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o out app/main.go
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build /go/src/app /myapp
EXPOSE 8080
WORKDIR myapp
RUN ls -al
ENTRYPOINT [ "./out" ]