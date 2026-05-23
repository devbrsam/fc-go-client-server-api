# fc-go-client-server-api

Cliente e servidor HTTP para consultar cotação USD/BRL com limites de tempo usando `context`.

## Requisitos

- Go 1.24+
- CGO habilitado (driver SQLite)

**Linux** (Debian/Ubuntu):

```bash
sudo apt update
sudo apt install -y libsqlite3-dev build-essential
```

**macOS** (Homebrew):

```bash
brew install sqlite
```

## Configuração

```bash
go mod download
```

## Execução

Inicie o servidor (terminal 1):

```bash
go run server.go
```

Execute o cliente (terminal 2):

```bash
go run client.go
```

O cliente grava a cotação em `cotacao.txt` no formato `Dólar: {bid}`.

As cotações ficam persistidas em `cotacoes.db`.

## Timeouts

| Etapa              | Limite |
|--------------------|--------|
| API externa        | 200ms  |
| Gravação no banco  | 10ms   |
| Cliente → servidor | 300ms  |
