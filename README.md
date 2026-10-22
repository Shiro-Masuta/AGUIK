**Une bibliothèque de composants UI néobrutalistes en Go, construite sur [Gio](https://gioui.org/).**

AGUIK est une bibliothèque de composants d'interface graphique conçue pour créer des applications desktop en Go avec un style néobrutaliste, des composants réutilisables et une API de composition inspirée des primitives de layout de Gio.

Elle fournit un ensemble de composants visuels et interactifs, un système de thème centralisé et des helpers de layout pour simplifier la création d'interfaces modernes, cohérentes et expressives.

[**Dépôt GitHub**](https://github.com/Shiro-Masuta/AGUIK) · [**Documentation Gio**](https://gioui.org/)

---

## Aperçu

AGUIK s'appuie sur les capacités de rendu et de gestion des événements de Gio pour proposer une couche de composants réutilisables.

- **Style néobrutaliste** — bordures marquées, ombres décalées et interactions visuelles dynamiques.
- **Composants composables** — boutons, champs de saisie, cartes, contrôles, overlays et bien plus.
- **Thème centralisé** — couleurs, typographie, formes et espacements réunis dans une structure commune.
- **Layouts simplifiés** — helpers pour organiser les composants horizontalement, verticalement, en grille ou en couches.
- **Interactions intégrées** — clics, focus, états de sélection, tooltips, dialogues et notifications.
- **Natif Go** — une bibliothèque construite sur Gio, sans remplacer son modèle de rendu.

## Fonctionnalités

### Composants UI

| Catégorie   | Composants                                     |
| ----------- | ---------------------------------------------- |
| Texte       | `Text`, `Title`, `Subtitle`, `Body`, `Caption` |
| Icônes      | `Icon`, `IconText`, `IconButton`               |
| Actions     | `Button`, `Badge`, `Card`                      |
| Formulaires | `Input`, `TextArea`, `Checkbox`, `Switch`      |
| Sélection   | `RadioGroup`, `Select`, `Tabs`, `Sidebar`      |
| Feedback    | `Alert`, `Dialog`, `Toaster`, `Tooltip`        |
| Indicateurs | `ProgressBar`, `Loader`, `Divider`             |

Les composants spécialisés sont séparés dans `components/table`, `components/codeeditor`, `components/imageview`, `components/select` et `components/overlay`.

### Outils de composition

Le package `glayout` propose des helpers pour construire des interfaces sans répéter la logique de layout Gio :

- `Row` et `Column` : disposition horizontale et verticale.
- `Rigid` et `Flexed` : gestion de l'espace fixe et flexible.
- `Inset` : ajout de padding.
- `Grid` : disposition en grille.
- `Stack` : superposition de widgets et gestion des overlays.

### Système de thème

Le package `theme` centralise les tokens visuels :

- Couleurs sémantiques et surfaces.
- Formes, bordures et ombres.
- Espacements en dp.
- Typographie et configuration du shaper Gio.

Le thème intégré `theme.Ao()` propose une palette sombre avec des tons bleus profonds et des accents cyan, verts, jaunes et rouges.

## Prérequis

- Go **1.26.1**, conformément à la version déclarée par le module.
- Les dépendances système éventuellement nécessaires à Gio pour compiler et exécuter une application graphique sur ton système d'exploitation.

AGUIK utilise principalement [Gio v0.10.2](https://gioui.org/).

## Installation

Ajoute AGUIK à ton projet Go :

```bash
go get github.com/Shiro-Masuta/AGUIK
```

Importe ensuite les packages nécessaires :

```go
import (
    "github.com/Shiro-Masuta/AGUIK/components"
    "github.com/Shiro-Masuta/AGUIK/glayout"
    "github.com/Shiro-Masuta/AGUIK/theme"
    "github.com/Shiro-Masuta/AGUIK/interaction"
)
```

Importe uniquement les packages dont ton application a besoin.

## Démarrage rapide

Le dépôt contient une application showcase permettant de découvrir les composants et leur fonctionnement.

Depuis la racine du dépôt :

```bash
go run ./cmd/showcase
```

Pour vérifier la compilation du module :

```bash
go test ./...
```

Les tests couvrent les layouts, les utilitaires, le tri et plusieurs interactions Gio. Lance `go test ./... -cover` pour consulter la couverture par package.

## Exemple d'utilisation

Voici un exemple illustrant la configuration d'un thème et la création de composants persistants.

```go
package main

import (
    "gioui.org/font/gofont"
    "gioui.org/text"

    "github.com/Shiro-Masuta/AGUIK/components"
    "github.com/Shiro-Masuta/AGUIK/theme"
)

type AppState struct {
    SaveButton components.Button
    NameInput  components.Input
}

func NewAppState() *AppState {
    return &AppState{
        SaveButton: components.Button{
            Text:    "Enregistrer",
            Variant: components.ButtonPrimary,
            Size:    components.SizeMD,
        },
        NameInput: components.Input{
            Placeholder: "Votre nom",
            Width:       240,
        },
    }
}

func NewTheme() *theme.Theme {
    th := theme.Ao()
    th.Shaper = text.NewShaper(
        text.WithCollection(gofont.Collection()),
    )
    return &th
}
```

Cet exemple présente uniquement l'initialisation du thème et de l'état des composants. Pour les afficher, il faut les intégrer au cycle de rendu Gio et appeler leur méthode `Layout` avec le contexte courant.

### Rendre un composant

Les composants interactifs conservent leur état entre les frames.

```go
dims := state.SaveButton.Layout(gtx, th)

if state.SaveButton.Clicked() {
    save()
}

return dims
```

`Clicked()` représente le clic consommé pendant la frame courante : il ne s'agit pas d'un état permanent.

**Important :** n'initialise pas les composants stateful dans une fonction de rendu appelée à chaque frame. Conserve-les dans une structure d'état persistante afin de préserver leurs événements, leur focus et leurs animations.

## Architecture

Le projet est organisé en quatre packages principaux, accompagnés d'une application de démonstration.

```text
AGUIK/
├── components/       # Composants fondamentaux
│   ├── codeeditor/   # Éditeur de code avec gouttière
│   ├── imageview/    # Images avec surlignages
│   ├── overlay/      # Dialogues et notifications
│   ├── select/       # Dropdown
│   └── table/        # Tableau triable et réorganisable
├── glayout/          # Helpers de composition Gio
├── interaction/      # Hitboxes, focus et invalidation
├── theme/            # Tokens visuels et thème Ao
└── cmd/
    └── showcase/     # Application de démonstration
```

### `components`

Contient les composants UI, leur rendu, leurs états et leurs interactions. Le cadre commun `NeoFrame` fournit notamment les bordures, les ombres et les variations visuelles au survol et à l'appui.

### `glayout`

Simplifie l'utilisation de `layout.Flex`, `layout.Stack` et des espacements Gio.

### `interaction`

Centralise plusieurs mécanismes d'interaction :

- `ClickArea` : délimitation et gestion d'une zone interactive.
- `Invalidate` : demande d'une nouvelle frame après une mutation d'état.
- `ClearFocus` : suppression du focus courant.

### `theme`

Définit les couleurs, les formes, les espacements, la typographie et le thème `Ao()`.

### `cmd/showcase`

Présente les composants dans une application fonctionnelle et illustre leur composition, la gestion des états et l'utilisation des overlays.

## Exemple de composition

AGUIK permet de combiner ses composants avec les helpers de layout :

```go
return glayout.Column(th.Spacing.MD, layout.Start,
    glayout.Rigid(func(gtx layout.Context) layout.Dimensions {
        return components.Title(th, "Mon application")(gtx)
    }),
    glayout.Rigid(func(gtx layout.Context) layout.Dimensions {
        return state.NameInput.Layout(gtx, th)
    }),
    glayout.Rigid(func(gtx layout.Context) layout.Dimensions {
        return state.SaveButton.Layout(gtx, th)
    }),
)(gtx)
```

Cet extrait illustre la composition attendue ; les imports et la structure d'état doivent être adaptés au contexte de l'application. Pour des overlays tels que les dialogues et les notifications, privilégie `glayout.Stack`.

## Bonnes pratiques

Pour utiliser AGUIK efficacement :

1. **Conserve l'état.** Les composants interactifs doivent vivre entre les frames.
2. **Configure le thème.** Fournis un `Theme.Shaper` avant d'utiliser les composants textuels.
3. **Respecte le cycle Gio.** Appelle `Layout` avec le contexte courant à chaque frame.
4. **Gère l'invalidation.** Utilise `interaction.Invalidate(gtx)` lorsqu'une mutation d'état nécessite un nouveau rendu.
5. **Utilise les bonnes unités.** Les tailles et espacements utilisent des dp ; la typographie utilise des sp.
6. **Superpose correctement les éléments.** Utilise `glayout.Stack` pour les dialogues et les toasts.
7. **Évite les doubles layouts interactifs.** Un composant qui consomme des événements ne doit normalement être layouté qu'une fois par frame.

## Limites actuelles

AGUIK se concentre sur la couche UI et la composition visuelle. Certaines responsabilités restent à la charge de l'application qui l'utilise :

- Gestion des données et logique métier.
- Navigation entre les écrans.
- Validation des formulaires.
- Gestion de plusieurs thèmes.
- Chargement et fourniture des icônes.

La couverture des composants n'est pas encore exhaustive ; consulte `docs/TODO.md` pour les tests et améliorations prévus.

## Documentation

La documentation technique détaillée couvre notamment :

- Le cycle de vie des frames Gio.
- Le contrat de rendu et les dimensions des composants.
- Le système de thème et le cadre `NeoFrame`.
- Les packages `glayout` et `interaction`.
- Les API des composants.
- La composition des overlays.
- L'architecture du showcase.
- Les bonnes pratiques pour ajouter de nouveaux composants.

Consulte le [dépôt GitHub d'AGUIK](https://github.com/Shiro-Masuta/AGUIK) pour accéder au code source et à la documentation du projet.

## Contribution

Les contributions sont les bienvenues.

Pour proposer une amélioration :

1. Fork le dépôt.
2. Crée une branche dédiée à ta fonctionnalité.
3. Implémente les changements en respectant les conventions du projet.
4. Vérifie la compilation avec `go test ./...`.
5. Ouvre une pull request en décrivant les changements apportés.

Lors de l'ajout d'un composant, pense à réutiliser `NeoFrame`, les helpers de layout et les mécanismes d'interaction existants, puis à présenter son utilisation dans le showcase.

## Licence

La licence du projet est à préciser selon le fichier de licence présent dans le dépôt.

---

**AGUIK — Des interfaces Go expressives, composables et néobrutalistes.**

[GitHub — Shiro-Masuta/AGUIK](https://github.com/Shiro-Masuta/AGUIK)
