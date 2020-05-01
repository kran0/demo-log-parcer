Demo log parcer

# What for?

Implements:
 - Read config from ENV;
 - Read [nginx][link_nginx_home] access.log from file or stdin;
 - Select "columns" from log file based on file structure (using [gonx][link_gonx_home]);
 - Group, Sum, Limit entries;
 - Display the result in human readable format.

# How to run?

## Environment variables

| Name | type | Description |
|:----:|:-----|:------------|
| PARCER_FILENAME      | string | **Required**. Using stdin if set to "```-```" |
| PARCER_LIMIT         | int    | Default:```10```                                      |
| PARCER_LOGFORMAT     | string | Default:```$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent $request_time "$http_referer" "$http_user_agent" [upstream: $upstream_addr $upstream_status] request_id=$upstream_http_x_request_id``` |
| PARCER_HUMANREADABLE | bool   | Default:```false```                                   |

## Build winth Golang (1.8+) and run

```bash
  cd demo-parcer
  go get -d -v
  go build -o ./demo-parcer
```

Run:

```bash
  export PARCER_FILENAME=../examples/access.log # stdin: PARCER_FILENAME=-
  ./demo-parcer
```

## Build and run in container

```bash
   docker build -t demo-parcer:latest .
```

Run:

```bash
   cat ./examples/*.log |
    docker run -i --rm\
     -e PARCER_FILENAME=-\
     demo-parcer:latest
```

---
[link_nginx_home]:https://nginx.org/
[link_gonx_home]:https://github.com/satyrius/gonx
