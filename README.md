# Gin Ping Example

## Voraussetzungen

* Go >= 1.20
* Git (optional, um das Repository zu klonen)

## Installation

1. Repository klonen (optional):

```bash
git clone <REPO_URL>
cd <REPO_NAME>
```

2. Abhängigkeiten installieren:

```bash
go mod tidy
```

## Projektstruktur

```
.
├── main.go       # Hauptdatei mit dem Gin-Server
└── go.mod        # Go Moduldatei
```

## Server starten

```bash
go run main.go
```

Der Server läuft dann standardmäßig auf `http://localhost:8080`.

## Testen

```bash
curl http://localhost:8080/ping
```

Erwartete Antwort:

```json
{
  "message": "pong"
}
```

## Hinweise

* Um den Port zu ändern, passe `router.Run(":PORT")` in `main.go` an.
* Für den produktiven Einsatz empfiehlt sich ein Reverse Proxy wie Nginx.
