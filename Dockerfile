FROM golang:alpine3.18 as build
WORKDIR /app

COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
EXPOSE 8080

RUN go build -o /app/application

FROM golang:alpine3.18
COPY --from=build /app/application /app/executable
ENTRYPOINT ["/app/executable"]

