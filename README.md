# rigger

The following is Go project layout scaffold generated:

```
├── Dockerfile
├── Makefile
├── README.md
├── docker-compose.yml
├── main.go
├── api
├── config
│   └── config.yml
├── fomrs
├── global
│   ├── response
│   └── global.go
├── initialize
│   ├── config.go
│   ├── logger.go
│   ├── router.go
│   ├── sentinel.go
│   ├── srv_conn.go
│   └── validator.go
├── middlewares
│   ├── admin.go
│   ├── cors.go
│   ├── jwt.go
│   └── tracing.go
├── models
│   └── request.go
├── proto
│   └── demo.proto
├── router
│   └── base.go
├── utils
└── validator
    └── validators.go
    
```


## Installation

Download rigger by using:
```sh
$ go install github.com/yangjishen/rigger
```

## Create a new project

1. Going to your new project folder:
```sh
# change to project directory
$ cd $GOPATH/src/path/to/project
```

2. Run `rigger init` in the new project folder:

```sh
$ rigger init
```

## Run service by using:
```sh
$ make run
```# rigger
