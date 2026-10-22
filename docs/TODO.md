# TODO

## Fait

- [x] `ImageView` : `FitFill` à échelles X/Y séparées, rendu découpé, cache sans panic sur les images non comparables
- [x] `Table` : tri non destructif, boutons de réorganisation persistants
- [x] `Input`, `TextArea`, `CodeEditor`, `Select` : dimensions bornées par le parent
- [x] Invalidation après clic cohérente (`IconButton`, `Sidebar`)
- [x] Utilitaires communs dans `glayout/utils.go`
- [x] Tests : `NeoFrame`, interactions (`Button`, `IconButton`, `Checkbox`, `Switch`, `RadioGroup`, `Tabs`, `Dialog`, `Sidebar`, `Select`, `Table`) via `internal/uitest`
- [x] Tests des layouts et utilitaires publics de `glayout`

## Tests

- [x] Interactions de `Dialog`, `Tabs`, `RadioGroup`, `Checkbox`, `Switch`
- [ ] `Toaster` : expiration et invalidation à l'échéance
- [x] `Table` : réorganisation par clic sur les boutons, en plus du test de `applyMove`
- [ ] `CodeEditor` et `Input` sous contraintes réduites
- [ ] Mesurer la couverture et fixer un seuil minimal

## Qualité du code

- [x] Séparer mesure/calcul et dessin dans `Select` et `ImageView`
- [x] Exposer `Clicked()` comme méthode de lecture des événements ponctuels ; garder `Selected` et `Active` comme champs d'état durables
- [x] Centraliser les marges de `NeoFrame` dans `neoFrameOverhead`
- [x] Ajouter une configuration `golangci-lint` au projet

## API et thème

- [ ] Fournir un `Shaper` par défaut avec `theme.Ao()` ou un constructeur dédié
- [ ] Ajouter un thème clair
- [ ] Mesurer la place réelle dans `Select` pour ouvrir au-dessus ou au-dessous

## Projet

- [ ] Séparer `cmd/AltesseTools` de la bibliothèque
- [ ] Ajouter une CI : `go vet`, `go test`, `gofmt -l`
- [ ] Ajouter une licence
