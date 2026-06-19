# Design URLWatch

Document d'architecture du projet.
L'entièreté de ce document a été formaté et reformulé pour qu'il soit bien écrit, par l'IA.

## Vue d'ensemble

URLWatch recoit un lot d'URLs, les verifie en parallele (avec limite de concurrence et timeout via `context`), puis stocke et expose les resultats via une API REST.

## Packages prévus

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

## Décisions (partie 1 domaine)

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

### Erreurs (partie 5.3)

- **Sentinelle** : `ErrBatchNotFound` pour un id inconnu dans `Store.Get`
- **Erreur personnalisee** : `ValidationError` avec le champ fautif (`Field`) et un message
- **Wrapping** : `ValidateBatchInput` utilise `fmt.Errorf("...: %w", err)` pour enrichir le contexte
- **Detection** : `errors.Is` pour la sentinelle, `errors.As` pour extraire `*ValidationError` meme wrappée
- **API** : `api.StatusFromError` traduit en 404 / 400 / 500 selon le type d'erreur

## Decisions (partie 2 — pool concurrent)

### Architecture

- `pool.Runner` prend un `domain.Checker` en dependance (injection, testable)
- `pool.Options` : `Concurrency`, `BatchTimeout`, `URLTimeout`

### Fan-out / fan-in (channels)

| Canal | Direction | Buffer | Justification |
|-------|-----------|--------|---------------|
| `jobs` | `chan string` | `len(urls)` | Le fan-out envoie toutes les URLs sans bloquer en attendant un worker |
| `results` | `chan CheckResult` | `len(urls)` | Evite le deadlock quand plusieurs workers terminent en meme temps (TP 4a) |

- **Fan-out** : goroutine `distribuerURLs` envoie les URLs sur `jobs`, puis `close(jobs)`
- **Workers** : nombre fixe = `Concurrency` (jamais une goroutine par URL)
- **Fan-in** : goroutine qui `wg.Wait()` puis `close(results)` ; le caller fait `range results`

### Context

- `context.WithTimeout(parent, BatchTimeout)` pour le lot entier
- `context.WithTimeout(batchCtx, URLTimeout)` par URL dans chaque worker
- `select` sur `ctx.Done()` dans le fan-out, les workers et l'envoi des resultats

### Synchronisation

- `sync.WaitGroup` : attendre la fin des N workers avant de fermer `results`
- `sync.Mutex` dans `MockChecker` uniquement (compteur de goroutines actives pour les tests)
- Pas de `sync.Once` : une seule goroutine ferme `results` apres `wg.Wait()`

### Checker

- `checker.HTTPChecker` : vrai GET avec `http.NewRequestWithContext`
- `checker.MockChecker` : delai simule + suivi du pic de concurrence pour les tests

## Decisions (partie 3 — API REST)

### Choix net/http (pas Gin)

- Déjà vu en TP exercice-5a : `ServeMux`, handlers, extraction d'id dans le path
- Pas de dépendance externe pour un microservice simple
- Middlewares maison (logging, recovery) proches du pattern Gin vu en 5b

### Endpoints

- `POST /v1/checks` : valide, exécute le pool, persiste, renvoie `201`
- `GET /v1/checks/{id}` : lit le store ou `404 batch_not_found`
- `GET /healthz` : sonde simple, exclue du middleware slog

### JSON

- DTOs API separes du domaine (`batch_id`, `ok`, `up`/`down` en snake_case)
- Erreurs uniformes `{ "error": { "code", "message" } }`

### Logging

- `slog` + handler JSON sur stdout
- Niveau via `LOG_LEVEL`
- Middleware : `method`, `path`, `status`, `duration_ms`, `batch_id` si connu

### Bonus recovery

- `recoveryMiddleware` : panic -> log + `500 internal`

## Decisions (partie 4 les tests)

- Tests table-driven : validation de l'API, validation du domaine, ComputeSummary, pool, store, handlers httptest
- Mocking : `checker.MockChecker` (Checker), `mockStore` + `MemoryStore` (Store)
- httptest : POST 201, GET 200, GET 404, POST 400
- Context : annulation manuelle + timeout global lot (`TestRunContextAnnulationOuTimeout`)
- Suite propre sous `go test -race ./...`
