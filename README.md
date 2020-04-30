Demo log parcer

# What for?

Implements:
 - Read config from ENV;
 - Read [nginx][link_nginx_home] access.log from file or stdin;
 - Select "columns" from log file based on file structure (using [gonx][link_gonx_home]);
 - Group, Sum, Limit entries;
 - Display the result in human readable format.

# How to run?

## Build winth Golang (1.8+) and run

```bash
  cd demo-parcer
  go get -d -v
  go build -o ./demo-parcer
```

Run:

```bash
  export PARCER_LIMIT=10
  export PARCER_FILENAME=/path/to/access.log # stdin: PARCER_FILENAME=-
  ./demo-parcer
```

## Build and run in container

```bash
   docker build -t demo-parcer:latest .
```

Run:

```bash
   cat access.log |
    docker run -i --rm\
     -e PARCER_FILENAME=-\
     -e PARCER_LIMIT=10\
     demo-parcer:latest
```

## Build and run with compose

**todo**

---
[link_nginx_home]:https://nginx.org/
[link_gonx_home]:https://github.com/satyrius/gonx
