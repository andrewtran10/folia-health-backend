# Folia Health Technical Challenge (Go)

I attempt to mirror the set up for the Laravel version as closely as I can.

# Deploying

## Locally building with Go

1. Install go with

```bash
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.25.1.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

2. Install the go-migrate binary with the sqlite driver and add GOBIN to path

```bash
go install -tags 'sqlite' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
export PATH=$PATH:$(go env GOPATH)/bin
```

3. Create a database file

```bash
touch database.sqlite
```

4. Run migrations

```bash
migrate -path ./migrations -database "sqlite://database.sqlite" up
```

5. Build

```bash
make build
```

6. Copy .env.example into .env

```bash
cp .env.example .env
```

7. Run

```bash
make run
```

## Docker

If you don't want to install go yourself, you can build the application using docker.

1. `make build-docker`
2. `touch database.sqlite`
3. `cp .env.example .env`
4. `make migrate-refresh`

NOTE: `migrate-fresh` will reset the database to a clean state, i.e no data will be persisted
