# Hello World Web

A versatile Go-based command-line tool for web development and debugging. It functions as both an HTTP server that inspects and echoes request details, and a client to test endpoints and parse responses.

# Installation

You can get the application in several ways.

## Using Go

*(Requires a local Go installation)*

### Go Install

Install directly using the `go install` command:
```shell
go install github.com/xh-dev-go/hello-world-web@latest
```

## Executable
Build through go command. (requires golang installed) 
```shell
# run under the root of the project
go build

# a executable `hello-world-web` or `hello-world-web.exe` will be generated (depends on os platform)
```

## Docker Hub
Pull from Docker Hub
```shell
docker pull xethhung/hello-world-web
```

## Build from Dockerfile
To build the Docker image, run the following command from the project root, replacing `{name:version}` with your desired image name (e.g., `helloworld-web:latest`).

```shell
docker build -t {name:version} .
```

# How to use it
The app is simple and can be worked as `Server` and `Client` (can access through `hello-world-web client` or restful api).

## Server
The app is simple, just start up a web server and receive for request. It receive the request and return the request information, it cloud be used for web dev for checking connectivity and the request informat (including url, ip, headers, etc)

To start the server, run the following command from the project root
```shell
# add env PORT variable if you want the server the use custom server port 
./hello-world-web server
```
### Run over docker container
The application inside the container listens on port `8080` by default. You can map any host port to the container's port `8080`.

To run the container and map your local port `8888` to the container's port `8080`:

```shell
# 8888 is the port expected to be access from host machine
docker run -p 8888:8080 xethhung/hello-world-web
```
## Client
The `client` can be used to test a simple API. It can either return the full response body using the `test` command, or it can return digested data using `get-ip`, `get-headers`, and `proxy-chain`.

All client commands require the `--url` (or `-u`) flag to specify the target server.

### `test`
Fetches the full response from the server and prints it.
```shell
./hello-world-web client -u http://localhost:8080 test
```

### `get-ip`
Parses the server response and prints only the client's IP address. It correctly identifies the IP even when behind a proxy by checking the `X-Forwarded-For` header.

```shell
./hello-world-web client -u http://localhost:8080 get-ip
```

`get-headers`
Parses the server response and prints only the request headers in a clean JSON format.

```shell
./hello-world-web client -u http://localhost:8080 get-headers
```

`proxy-chain`
Visualizes the network path of the request, showing the origin IP, any intermediate proxies, and the final destination.

```shell
# Example without a proxy
./hello-world-web client -u http://localhost:8080 proxy-chain
# Output: 127.0.0.1:54321 -> [ no proxy ] -> localhost:8080
```

## Api
The connection can be simply tested with calling RESTFUL API.

```shell
curl {the url of you server, e.g. http://localhost:8080?format=json}
```

Simple result should similar
```json
{
  "host": "localhost:8111",
  "url": "/?format=json",
  "ip": "[::1]:49936",
  "referer": "",
  "headers": {
    "Accept": [
      "*/*"
    ],
    "User-Agent": [
      "curl/8.5.0"
    ]
  }
}
```