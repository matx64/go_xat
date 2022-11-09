# 💬 go_xat

A fast, minimal, private & open source Chat Room App written in Go.


### 🔨 Tech

- Go 1.19
- Html, Javascript, Bootstrap5
- Redis
- Docker
- Fly.io


### 🏠 Run it locally

```bash
# build and run a redis container
docker run -dp 6379:6379 --rm --name goxat-redis redis

# run the app
go mod tidy && go run .
```