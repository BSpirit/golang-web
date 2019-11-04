# Create Database
```
sqlite3 test.db < schema.sql
```

# Compile for linux using docker
```
docker run --rm -v "$PWD":/usr/src/web -w /usr/src/web golang:1.13 go build -v
```