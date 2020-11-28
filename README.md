# Hypertrace Go Agent Example

This repository shows a example on how to use [Go Agent](https://github.com/hypertrace/goagent).

![Screenshot](screenshot.png)

## Running it locally

**Run backend:**

```bash
go run backend/main.go
```

**Run frontend:**

```bash
go run frontend/main.go
```

**Run Hypertrace and MySQL:**

```bash
docker-compose -f docker-compose-mysql.yml -f docker-compose-hypertrace.yml up --renew-anon-volumes
```

Once everything is up and running you can curl the frontend:

```bash
curl -i http://localhost:8081
```

## Running it docker

```bash
docker-compose -f docker-compose.yml -f docker-compose-hypertrace.yml up --renew-anon-volumes
```
