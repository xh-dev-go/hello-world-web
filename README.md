# Hello World Web Server

A simple Go web server that returns request information in YAML or JSON format.
The purpose of this application is to easily get request information, which is useful for web development and debugging connectivity.

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
The app is simple, just start up a web server and receive for request. I receive the request and return the request information, it cloud be used for web dev for checking connectivity and the request informat (including url, ip, headers, etc)

To start the server, run the following command from the project root
```shell
# add env PORT variable if you want the server the use custom server port 
./hello-world-web
```

To verify the request data:
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

## Run over docker container
The application inside the container listens on port `8080` by default. You can map any host port to the container's port `8080`.

### Running on the Default Port

To run the container and map your local port `8888` to the container's port `8080`:

```shell
# 8888 is the port expected to be access from host machine
docker run -p 8888:8080 {name:version}
```
