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

L'entièreté de ce fichier markdown a été reformulée et indentée correctement par l'IA.

## Decisions (partie 0)

- Module Go : `github.com/Cengokill/Examen-final-go`
- `main.go` reste mince : uniquement le cablage

## Decisions (partie 1 domaine)

### Types

- `CheckResult` : une URL, son code HTTP, un bool `Available`, la latence en ms et un message d'erreur optionnel (`omitempty` en JSON)
- `BatchSummary` : type separe pour le résumé (total, disponibles, echecs, duree totale en ms)
- `Batch` : id, `CreatedAt` en UTC, slice de resultats + resume pre-calcule

### Agrégation

- `ComputeSummary` parcourt la slice et utilise une map pour compter disponibles / echecs (vu en TP)
- `TotalDurationMs` = somme des `LatencyMs` de chaque URL (duree cumulee des requetes)
- `NewBatch` copie la slice des resultats avant de calculer le resume (evite les mutations externes)

### Interfaces

- `Checker.Check(ctx, url)` : le `context` permettra timeout et annulation (partie pool/checker)
- `Store.Save` / `Store.Get` : contrat minimal ; `Get` renverra `ErrBatchNotFound` si l'id n'existe pas
