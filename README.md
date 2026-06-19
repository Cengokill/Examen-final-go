# URLWatch

Microservice Go de verification d'URLs en masse (projet d'examen).

## Prerequis

- Go 1.23 ou plus recent

## Commandes

```bash
# compiler tout le module
go build ./...

# analyse statique
go vet ./...

# tests
go test ./...

# lancer le binaire (partie 0 : message de demarrage uniquement)
go run ./cmd/urlwatch
```

## Exemple curl

Les endpoints REST seront documentes ici a partir de la partie API.

```bash
# a venir : creation d'un lot
# curl -X POST http://localhost:8080/batches ...

# a venir : lecture d'un lot
# curl http://localhost:8080/batches/{id}
```
