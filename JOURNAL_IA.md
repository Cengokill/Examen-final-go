# Journal d'utilisation de l'IA

Notes sur l'usage de l'IA pour l'examen URLWatch (l'entièreté du document a été reformulée par l'IA au format markdown).

## Partie 0 — Mise en place

### Ce que j'ai demande a l'IA

- Générer l'arborescence suggeree par le sujet (`cmd/`, `internal/domain`, `checker`, `pool`, `store`, `api`)
- Proposer un `.gitignore` standard pour Go
- Ecrire un `main.go` vierge

### Accepte

- La structure de dossiers du sujet, sans modification
- Le module `github.com/Cengokill/Examen-final-go` (aligne sur l'URL du depot Git)
- Les fichiers `README.md` et `DESIGN.md` en squelette pour les remplir plus tard

### Modifie / rejete

- Pas de code metier dans la partie 0 : seulement des `package` vides dans `internal/`
- Le `main.go` n'écoute pas encore de port HTTP (prévu pour les parties suivantes)

### Pourquoi

La partie 0 ne demande que l'init du module et du depot Git.

## Partie 1 — Modélisation du domaine

### Ce que j'ai demandé à l'IA

- Me proposer des struct que je pourrais utiiser
- Interfaces `Checker` et `Store` comme dans le sujet

### Accepté

- Séparation en fichiers (`types.go`, `interfaces.go`, `summary.go`, `errors.go`)
- Certaines struct pertinentes
- Correction de mes tests unitaires sur l'agrégation car 2 tests OK mais le dernier ne passait pas

### Modifié / rejeté

- J'ai ajouté `BatchSummary` comme struct à part (plus lisible que des champs plats dans `Batch`)
- `TotalDurationMs` = somme des latences individuelles (pas le temps mur du batch, qui viendra du pool)

### Pourquoi

Le domaine doit rester indépendant de HTTP et du stockage. Les interfaces permettront de mocker le checker dans les tests.
