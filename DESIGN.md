# Design URLWatch

Document d'architecture du projet. Sera complete au fur et a mesure de l'implementation.

## Vue d'ensemble

URLWatch recoit un lot d'URLs, les verifie en parallele (avec limite de concurrence et timeout via `context`), puis stocke et expose les resultats via une API REST.

## Packages prevus

| Package | Role |
|---------|------|
| `cmd/urlwatch` | Point d'entree, assemblage des dependances |
| `internal/domain` | Types metier, erreurs, interfaces (`Checker`, `Store`) |
| `internal/checker` | Verification HTTP concrete |
| `internal/pool` | Worker pool concurrent (fan-out / fan-in) |
| `internal/store` | Persistance en memoire (SQLite en bonus) |
| `internal/api` | Handlers HTTP, DTO JSON, middleware |

## Decisions (partie 0)

- Module Go : `github.com/Cengokill/Examen-final-go`
- Inversion de dependance : les packages techniques dependront de `domain`, pas l'inverse
- `main.go` reste mince : uniquement le cablage
