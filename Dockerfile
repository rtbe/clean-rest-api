FROM golang:latest AS build

WORKDIR /go/src
COPY . . 
RUN go get -d -v 
RUN go build -o /go/bin/app

FROM gcr.io/distroless/base-debian10
 
COPY --from=build /go/bin/app /
COPY .env /
COPY swagger.yaml /
COPY /internal/database/migrate/migrations /migrations
 
ENTRYPOINT ["/app"]