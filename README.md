# Hypertrace Go Agent Example

This repository shows a example on how to use [Go Agent](https://github.com/hypertrace/goagent).

![Screenshot](screenshot.png)

**Run server:**

```bash
go run server/main.go
```

**Run client:**

```bash
go run client/main.go
```

**Run hypertrace and mysql:**

```bash
docker-compose -f docker-compose.yml -f docker-compose-hypertrace.yml up --renew-anon-volumes
```
