# URLWatch

Microservice Go de verification d'URLs en masse (projet d'examen).

## Prerequis

- Go 1.23 ou plus recent

## Commandes

```bash
go build ./...
go vet ./...
go test ./...
go run ./cmd/urlwatch
```

Variables d'environnement optionnelles :
- `LOG_LEVEL` : `DEBUG`, `INFO` (défaut), `WARN`, `ERROR`
- `PORT` : port d'écoute (défaut `8080`)

## Exemples curl

```bash
# sonde pour voir si ça fonctionne
curl http://localhost:8080/healthz

# créer et exécuter un lot
curl -X POST http://localhost:8080/v1/checks \
  -H "Content-Type: application/json" \
  -d '{"urls":["https://go.dev","https://exemple.invalid"],"options":{"concurrency":4,"timeout_ms":2000}}'

# relire un lot
curl http://localhost:8080/v1/checks/b_4f3c1a
```
