# ws

ws is a simple command line websocket client designed for exploring and debugging websocket servers. ws includes readline-style keyboard shortcuts, persistent history, and colorization.

## Installation

```
go install github.com/weiwenchen2022/ws@latest
```

## Usage

Simply run ws with the destination URL. For security some sites check the origin header. ws will automatically send the destination URL as the origin. If this doesn't work you can specify it directly with the `--origin` parameter.

```
$ ws ws://localhost:8000/ws
> {"type": "echo", "payload": "Hello, world"}
< {"type":"echo","payload":"Hello, world"}
> {"type": "broadcast", "payload": "Hello, world"}
< {"type":"broadcast","payload":"Hello, world"}
< {"type":"broadcastResult","payload":"Hello, world","listenerCount":1}
> ^D
EOF
```
