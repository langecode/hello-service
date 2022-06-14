##
## Build
##
FROM golang:1.18 AS build

WORKDIR /go/src/app
ADD . /go/src/app

RUN go get -d -v ./...
RUN go vet -v
RUN CGO_ENABLED=0 go build -o /go/bin/app

##
## Distribution image
##
FROM gcr.io/distroless/static

EXPOSE 8080/tcp

COPY --from=build /go/bin/app /
CMD ["/app"]
