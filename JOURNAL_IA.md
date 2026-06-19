# Journal d'utilisation de l'IA

Notes sur l'usage de l'IA pour l'examen URLWatch (l'entièreté du document a été reformulé par l'IA au format markdown).

## Partie 0 — Mise en place

### Ce que j'ai demande a l'IA

- Generer l'arborescence suggeree par le sujet (`cmd/`, `internal/domain`, `checker`, `pool`, `store`, `api`)
- Proposer un `.gitignore` standard pour Go
- Ecrire un `main.go` vierge

### Accepte

- La structure de dossiers du sujet, sans modification
- Le module `github.com/Cengokill/Examen-final-go` (aligne sur l'URL du depot Git)
- Les fichiers `README.md` et `DESIGN.md` en squelette pour les remplir plus tard

### Modifie / rejete

- Pas de code metier dans la partie 0 : seulement des `package` vides dans `internal/`
- Le `main.go` n'ecoute pas encore de port HTTP (prevu pour les parties suivantes)

### Pourquoi

La partie 0 ne demande que l'init du module et du depot Git. Ajouter de la logique maintenant compliquerait les commits incrementaux prevus par le sujet.
