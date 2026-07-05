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
