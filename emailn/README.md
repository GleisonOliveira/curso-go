# Emailn

## Pré-requisitos

- Go 1.26+
- Docker

## Subindo o banco de dados

```bash
docker compose up -d
```

## Baixando as dependências

```bash
go mod download
```

## Executando as migrations

```bash
# Rodar todas as migrations pendentes (DSN resolvida do .env)
go run ./cmd/migrate

# Rollback de 1 step
go run ./cmd/migrate -down

# Rollback de N steps
go run ./cmd/migrate -steps -3

# Subir N steps
go run ./cmd/migrate -steps 2
```

## Subindo a API com hot reload (Air)

```bash
go tool air
```


## Keycloak

### Machine to machine
```
curl --location 'http://localhost:8080/realms/emailn/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'client_id=emailn' \
--data-urlencode 'client_secret=secret' \
--data-urlencode 'grant_type=client_credentials'
```

### Username and password
```
curl --location 'http://localhost:8080/realms/emailn/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'client_id=emailn_users' \
--data-urlencode 'password=admin' \
--data-urlencode 'grant_type=password' \
--data-urlencode 'username=gleison'
```

### Username and password OAuth
```
http://localhost:8080/realms/emailn/protocol/openid-connect/auth?client_id=emailn&redirect_uri=http://localhost:8082/auth/callback&response_type=code&scope=openid&state=168e4224ae9060856521f14a36b565df
```

```
curl --location 'http://localhost:8080/realms/emailn/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'client_id=emailn' \
--data-urlencode 'client_secret=6FYfuOLWqhOnQA7zb5gbCjBf1JAKPvSG' \
--data-urlencode 'grant_type=authorization_code' \
--data-urlencode 'code=code' \
--data-urlencode 'redirect_uri=http://localhost:8082/auth/callback'
```