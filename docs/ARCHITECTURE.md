# AGUIK

AGUIK est une bibliothèque de composants d'interface graphique en Go construite sur [Gio](https://gioui.org/). Elle fournit un style visuel néobrutaliste, des composants composables, des helpers de layout et une petite couche commune pour les interactions.

Cette documentation décrit le code présent dans le dépôt et son fonctionnement réel.

## 1. Vue d'ensemble

Le projet est organisé en quatre packages principaux :

```text
AGUIK2/
├── components/   Composants visuels et interactifs
├── glayout/      Helpers de composition de layouts Gio
├── interaction/  Hitboxes, focus et invalidation de frame
├── theme/        Tokens visuels et thème Ao
└── cmd/
    ├── showcase/ Exemple fonctionnel de tous les composants
```

Le module Go est `github.com/Shiro-Masuta/AGUIK` et dépend principalement de `gioui.org v0.10.2`. Le module déclare Go `1.26.1`.

AGUIK ne remplace pas Gio :

- Gio crée la fenêtre, reçoit les événements et déclenche les frames.
- Gio fournit `layout.Context`, `layout.Widget`, `widget.Clickable`, `widget.Editor`, `widget.Bool` et les opérations de dessin.
- AGUIK ajoute des composants, un thème commun et des fonctions de composition.

## 2. Démarrer l'exemple

Depuis la racine du dépôt :

```bash
go run ./cmd/showcase
```

Pour vérifier la compilation de tout le module :

```bash
go test ./...
```

Il n'y a actuellement pas de fichiers de test dans le dépôt. `go test ./...` sert donc surtout de vérification de compilation.

## 3. Cycle de vie Gio

L'application de démonstration suit le cycle Gio standard :

1. `main` crée une `app.Window` et appelle `app.Main()`.
2. `run` lit les événements de la fenêtre.
3. À chaque `app.FrameEvent`, `app.NewContext` construit un `layout.Context`.
4. `draw` compose les widgets et ajoute leurs opérations dans `gtx.Ops`.
5. `e.Frame(gtx.Ops)` envoie les opérations au moteur Gio.
6. Un clic, un changement de focus ou une invalidation provoque une nouvelle frame.

Le point essentiel est que `draw` peut être rappelé très souvent. Un composant avec de l'état ne doit donc pas être créé dans `draw` ou dans une fonction de layout appelée à chaque frame.

### État persistant

Les structures suivantes doivent vivre dans un état conservé entre les frames :

- `components.Button`, `components.IconButton`
- `components.Input`, `components.TextArea`
- `components.Checkbox`, `components.Switch`
- `components.RadioGroup`, `components.Select`, `components.Tabs`
- `components.Dialog`, `components.Toaster`, `components.Alert`
- `components.Sidebar`
- toute structure Gio contenant `widget.Clickable`, `widget.Editor`, `widget.Bool` ou un timer logique

Exemple :

```go
type appState struct {
    saveButton components.Button
    name       components.Input
    dialog     components.Dialog
}

state := &appState{
    saveButton: components.Button{
        Text:    "Enregistrer",
        Variant: components.ButtonPrimary,
        Size:    components.SizeMD,
    },
    name: components.Input{
        Placeholder: "Nom",
        Width:       240,
    },
}
```

La configuration fixe peut être initialisée une fois. Pendant le rendu, on appelle seulement `Layout` et on lit les propriétés exposées (`Clicked`, `Value`, `Selected`, `Active`, etc.).

## 4. Contrat de rendu des composants

Un widget Gio a généralement la forme suivante :

```go
func(gtx layout.Context) layout.Dimensions
```

Les composants AGUIK stateful ajoutent le thème :

```go
func (c *Component) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions
```

Le `layout.Context` contient notamment :

- les contraintes de taille en pixels ;
- la densité d'écran utilisée par `gtx.Dp` et `gtx.Sp` ;
- la source d'événements ;
- les opérations de dessin et d'invalidation.

Les dimensions retournées incluent la taille visuelle calculée par `NeoFrame`, donc la bordure, l'espace réservé au lift et l'ombre. Pour cette raison, il ne faut pas supposer que la taille retournée est exactement la taille du texte ou de l'icône.

### Mesure puis dessin

Plusieurs composants mesurent leur contenu avant de le dessiner :

1. ils enregistrent les opérations avec `op.Record` ;
2. ils exécutent le widget avec des contraintes de mesure ;
3. ils récupèrent `layout.Dimensions` ;
4. ils construisent leur cadre ;
5. ils rejouent les opérations via `op.CallOp` ou redessinent le contenu avec un offset.

Cette technique est utilisée par `Card`, les boutons, `Input`, `TextArea`, `Dialog`, `Alert`, `Select`, `Toaster` et les tooltips.

## 5. Le thème

`theme.Theme` regroupe les tokens utilisés par les composants :

```go
type Theme struct {
    Colors     Colors
    Shape      Shape
    Spacing    Spacing
    Typography Typography
    Shaper     *text.Shaper
}
```

### Couleurs

`theme.Colors` contient :

- `Background` : fond global de la fenêtre ;
- `Surface`, `SurfaceAlt` : surfaces de composants et de surbrillance ;
- `Foreground`, `Muted` : texte principal et secondaire ;
- `Border` : bordures et contours ;
- `Primary`, `Secondary`, `Success`, `Warning`, `Danger` : variantes sémantiques.

### Formes et espacement

`Shape` définit le rayon, la largeur de bordure et l'ombre. `Spacing` expose les tokens `XS`, `SM`, `MD`, `LG` et `XL`, en dp.

Le thème fourni par `theme.Ao()` est sombre, avec des tons bleus profonds et des accents cyan, verts, jaunes et rouges. Chaque appel à `Ao()` retourne une nouvelle valeur indépendante.

### Typographie

`Typography.FontSize` est la taille de base en sp. `Typography.Font` permet de fixer la police et son poids. Le showcase installe un shaper Gio avec la collection Go fonts :

```go
th := theme.Ao()
th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
```

Un `Shaper` est nécessaire pour `Text` et les composants qui affichent du texte. Les icônes Gio n'en ont pas besoin.

## 6. Le cadre néobrutaliste

`components.NeoFrame` est le socle visuel commun de `Card`, `Button`, `IconButton`, des champs et des contrôles de sélection.

Il dessine, dans cet ordre logique :

1. l'ombre décalée ;
2. la bordure ;
3. le fond intérieur.

Le cadre réserve toujours une marge autour de son rectangle pour les effets de repos, de survol et d'appui. Les états visuels sont :

- repos : ombre normale ;
- survol : le contenu monte légèrement et l'ombre se décale ;
- appui : le contenu revient vers l'ombre et l'ombre est masquée.

La hitbox n'est pas posée sur toute la toile du cadre. `interaction.ClickArea` la limite au `RestRect`, afin que la marge de l'ombre ne devienne pas une zone de clic inattendue.

## 7. Package `glayout`

`glayout` simplifie les primitives `layout.Flex`, `layout.Stack` et les espacements.

### `Rigid` et `Flexed`

```go
glayout.Row(th.Spacing.MD, layout.Middle,
    glayout.Rigid(left),
    glayout.Flexed(1, right),
)
```

`Rigid` conserve la taille naturelle du widget. `Flexed(weight, widget)` donne au widget une part de l'espace restant selon son poids.

### `Row` et `Column`

Ces fonctions insèrent automatiquement un `Spacer` entre les enfants. Elles prennent un espacement en dp et un alignement Gio.

```go
glayout.Column(th.Spacing.SM, layout.Start,
    glayout.Rigid(title),
    glayout.Rigid(field),
    glayout.Rigid(button),
)
```

### `Inset`

`glayout.Inset(space, widget)` applique un padding uniforme.

### `Grid`

`glayout.Grid(spacing, columns, children...)` répartit les widgets ligne par ligne, avec des colonnes flexibles de même largeur. Une dernière ligne peut être incomplète.

### `Stack`

`glayout.Stack(children...)` transforme chaque enfant en couche `layout.Expanded`. Le dernier enfant est dessiné au-dessus et reçoit la priorité visuelle et d'interaction dans les zones qu'il recouvre.

C'est la bonne primitive pour un dialog ou une pile de toasts. Un overlay ne doit pas être ajouté comme un enfant ordinaire d'une `Column`, car il recevrait seulement l'espace résiduel de cette colonne.

## 8. Package `interaction`

### `ClickArea`

`ClickArea` associe une zone rectangulaire à une cible implémentant :

```go
type Interactive interface {
    Layout(gtx layout.Context, w layout.Widget) layout.Dimensions
}
```

`widget.Clickable` et `widget.Bool` satisfont ce contrat. La fonction :

- applique l'offset de la hitbox ;
- pose un clip pointeur ;
- délègue à la cible Gio ;
- active le curseur pointeur.

### `Invalidate`

Gio ne redessine pas automatiquement après une simple mutation d'une structure Go. Un composant qui change d'état après un événement doit donc exécuter `interaction.Invalidate(gtx)`, qui envoie `op.InvalidateCmd{}`.

Cas concernés : ouverture d'un dialog, sélection d'une option, fermeture d'une alerte, clic d'un bouton et changement de sélection radio.

### Focus

`interaction.ClearFocus(gtx)` envoie un `key.FocusCmd` sans tag. Le showcase l'utilise quand un bouton ou le fond est cliqué afin de retirer le focus d'un éditeur.

## 9. Référence des composants

### Texte

Fonctions sans état :

- `Text(th, value, options...)` : texte de base ;
- `Title` : base x 1,75 et gras ;
- `Subtitle` : base x 1,25 et poids medium ;
- `Body` : alias sémantique de `Text` ;
- `Caption` : base x 0,85 et couleur `Muted`.

Options disponibles : `WithSize`, `WithColor`, `WithFont`, `WithAlignment`, `WithMaxLines`, `WithWeight` et `WithEllipsis`.

`WithEllipsis` force au moins une ligne et utilise le caractère d'ellipse comme truncateur.

### Icônes et texte avec icône

- `Icon` dessine un `*widget.Icon`, avec une taille par défaut de 24 dp ;
- `IconText` compose une icône et un texte dans un `Row` ;
- `WithIconSize` et `WithIconColor` modifient l'icône.

Une icône doit être chargée ou construite par l'application, puis passée au composant. AGUIK ne possède pas de catalogue d'icônes global.

### Boutons

`Button` est un composant stateful. Champs principaux : `Text`, `Icon`, `Variant`, `Size`, `Disabled`, `Clicked`.

Les variantes sont `ButtonDefault`, `ButtonPrimary`, `ButtonSecondary`, `ButtonSuccess`, `ButtonWarning` et `ButtonError`. Les tailles sont `SizeSM`, `SizeMD` et `SizeLG`.

```go
dims := state.saveButton.Layout(gtx, th)
if state.saveButton.Clicked {
    save()
}
return dims
```

`Clicked` représente le clic consommé pendant la frame courante : il ne faut pas le traiter comme un booléen permanent.

`IconButton` suit le même modèle, mais son contenu est une icône. Il peut afficher un tooltip directement avec `TooltipText` et `TooltipPlacement`.

### `Badge`

`Badge` est une étiquette non interactive. Il utilise les variantes et le cadre commun, mais ne possède aucune hitbox et ne réagit pas au survol.

### `Card`

`Card(th, content, options...)` enveloppe un widget avec padding, surface, bordure et ombre. Options :

- `WithCardPadding(dp)` ;
- `WithCardBackground(color)`.

La carte mesure le contenu avec une largeur réduite du padding, puis ajoute son propre cadre.

### Champs

`Input` encapsule `widget.Editor` en mode une ligne. `TextArea` utilise le même mécanisme en mode multiligne.

Champs communs : `Editor`, `Placeholder`, `Width`, `Size`, `Disabled`. `TextArea` possède aussi `Lines`, par défaut à 4. `Input` possède `Password`, qui masque le texte avec un caractère de masquage.

La valeur saisie est lue depuis le champ Gio :

```go
value := state.nameInput.Editor.Text()
```

Les champs sont rendus via `material.Editor`, mais leur fond, bordure, padding et état de focus sont fournis par AGUIK. Le focus est visible par une variation de l'état du cadre.

### Cases et interrupteurs

`Checkbox` encapsule `widget.Bool` et expose `Bool.Value`. Il accepte un label, une variante, une taille et `Disabled`.

`Switch` encapsule également `widget.Bool`, mais dessine une piste et un bouton mobile. `Bool.Value` vaut `true` lorsqu'il est activé.

Les couleurs actives passent par `toggleCheckedColors`. Avec `ButtonDefault`, l'état actif retombe volontairement sur `Primary` pour rester visuellement différent de l'état inactif.

### `RadioGroup`

`RadioGroup` représente une sélection exclusive. `Options` contient les labels, `Selected` contient l'index choisi, et `-1` signifie aucune sélection.

Le composant conserve un `widget.Clickable` par option. Si la taille de `Options` change, il redimensionne son tableau interne tout en conservant les clickables existants qui peuvent l'être.

### `Tabs`

`Tabs` affiche une rangée d'onglets. `Active` est l'index actif. Il est normalisé à zéro si sa valeur sort de la liste.

Les onglets réutilisent le cadre de bouton, mais l'appelant décide du conteneur et du fond autour d'eux.

### `Select`

`Select` est un dropdown avec :

- `Options` ;
- `Selected`, avec `-1` pour aucune sélection ;
- `Placeholder` ;
- `Width`, par défaut 200 dp ;
- `Size` et `Disabled`.

Le trigger et chaque option ont leur propre `widget.Clickable`. Quand le menu est ouvert, les options sont dessinées dans un panneau unique différé avec `op.Defer`, afin de passer au-dessus du contenu voisin. Le menu s'ouvre sous le trigger lorsqu'il y a assez de place, sinon au-dessus selon l'espace Gio disponible.

### `ProgressBar`

`ProgressBar` est purement visuelle. `Value` est automatiquement limitée à l'intervalle `[0, 1]`. `Width` vaut 200 dp par défaut et `Height` 12 dp par défaut.

```go
progress.Value = uploaded / total
```

Le remplissage utilise la variante choisie ; `ButtonDefault` utilise `Primary`.

### `Loader`

`Loader` est un spinner indéterminé. Il conserve son instant de démarrage et demande une invalidation à chaque layout pour continuer à tourner même sans événement utilisateur.

`Size` est le diamètre en dp, 24 dp par défaut. Comme il s'agit d'une animation continue, il ne faut pas en créer une nouvelle instance à chaque frame.

### `Divider`

`Divider` dessine une ligne horizontale ou verticale. `Thickness` utilise par défaut la largeur de bordure du thème. `Length` vaut par défaut toute la longueur disponible sur l'axe principal.

### `Sidebar`

`Sidebar` est un conteneur vertical d'`IconButton` avec un seul item actif. `Active` est l'index courant et `Items` fournit l'icône et le tooltip de chaque entrée.

La sidebar mesure sa largeur sur son contenu, mais force sa hauteur à la hauteur disponible. C'est pourquoi le showcase la place dans une ligne avec `Rigid`.

### `Alert`, `Dialog` et `Toaster`

`Alert` affiche une surface avec accent coloré, titre, message, icône optionnelle et fermeture optionnelle. Il expose `Closed()` pour permettre à un conteneur de le retirer.

`Dialog` utilise `Overlay` pour :

- assombrir toute la fenêtre ;
- bloquer les clics vers l'arrière-plan ;
- centrer sa carte ;
- fermer éventuellement au clic sur le fond ;
- détecter l'action cliquée via `Clicked` ;
- fermer automatiquement si `AutoClose` est activé.

`Clicked` est remis à `-1` au début de chaque layout. Après le layout, le code applicatif peut traiter `Actions[Clicked]`.

`Toaster` gère une pile d'`Alert` en bas à droite. `Show(ToastConfig)` ajoute un toast sans exposer les clickables internes. Une durée positive programme une échéance ; `Toaster.Layout` purge les entrées expirées et programme une invalidation exactement à la prochaine échéance.

### `Tooltip`

`Tooltip` enveloppe un widget non interactif et suit son survol avec `gesture.Hover`. `Placement` vaut `TooltipTop`, `TooltipBottom` ou `TooltipRight`.

Pour un composant déjà interactif, notamment `IconButton`, il ne faut pas ajouter un wrapper `Tooltip` autour de sa propre hitbox. `IconButton` réutilise déjà son état de survol et appelle `LayoutTooltipBubble` directement.

## 10. Composition des overlays

La composition du showcase est représentative :

```go
return glayout.Stack(
    backdrop,
    mainContent,
    dialog,
    toaster,
)(gtx)
```

L'ordre est important :

1. le fond cliquable ;
2. le contenu principal ;
3. le dialog ;
4. les toasts.

Le dernier enfant est dessiné au-dessus. Le backdrop sert à retirer le focus quand l'utilisateur clique dans le vide. Le backdrop du dialog est distinct et bloque les interactions derrière lui.

## 11. Architecture du showcase

`cmd/showcase/state.go` centralise les instances persistantes et leur configuration initiale.

`cmd/showcase/draw.go` compose les couches principales :

- fond et retrait du focus ;
- sidebar ;
- contenu principal ;
- dialog et traitement de ses actions ;
- toaster.

`profile_panel.go` montre un formulaire composé de `Input`, `TextArea`, `Checkbox`, `RadioGroup`, boutons et onglets.

`component_panel.go` montre `Badge`, `Switch`, `ProgressBar`, `Loader`, `Select` et l'ouverture d'un dialog.

`header.go` montre `Icon`, `Title`, `Caption`, `Badge` et `IconButton`.

`footer.go` affiche un état métier simple (`clickCount`, `lastAction`) séparé du rendu des composants.

## 12. Exemple de composition complet

```go
func screen(gtx layout.Context, th *theme.Theme, state *appState) layout.Dimensions {
    return glayout.Stack(
        glayout.Inset(th.Spacing.LG,
            glayout.Column(th.Spacing.MD, layout.Start,
                glayout.Rigid(components.Title(th, "Profil")),
                glayout.Rigid(func(gtx layout.Context) layout.Dimensions {
                    return state.name.Layout(gtx, th)
                }),
                glayout.Rigid(func(gtx layout.Context) layout.Dimensions {
                    dims := state.saveButton.Layout(gtx, th)
                    if state.saveButton.Clicked {
                        state.dialog.Open = true
                    }
                    return dims
                }),
            ),
        ),
        func(gtx layout.Context) layout.Dimensions {
            return state.dialog.Layout(gtx, th)
        },
    )(gtx)
}
```

Le même composant ne doit être layouté qu'une seule fois par frame, en particulier pour `Dialog`, `Toaster` et les composants qui consomment des événements. Si un composant est mesuré puis redessiné, il faut utiliser le pattern interne d'enregistrement/rejeu plutôt que d'appeler plusieurs fois une logique interactive non prévue pour cela.

## 13. Règles pratiques et limites connues

- Conserver les composants stateful dans une structure persistante.
- Appeler `Layout` à chaque frame avec le même thème et le contexte Gio courant.
- Appeler `interaction.Invalidate` après les mutations d'état déclenchées par interaction lorsqu'une nouvelle frame n'est pas déjà garantie.
- Utiliser des valeurs en dp pour les tailles et en sp pour la typographie.
- Fournir un `Theme.Shaper` avant d'utiliser du texte.
- Utiliser `glayout.Stack` pour les overlays.
- Ne pas utiliser `Tooltip` autour d'un composant qui suit déjà son propre survol.
- Les indices `Selected` et `Active` sont corrigés automatiquement lorsqu'ils sont invalides.
- `ProgressBar.Value` est clampée, mais les champs d'entrée et les listes d'options restent à la responsabilité de l'appelant.
- `Select` estime l'espace vertical à partir des contraintes Gio du parent : dans une colonne complexe, le parent doit donc lui fournir des contraintes cohérentes.
- Il n'y a pas encore d'abstraction de données, de navigation, de validation de formulaire ou de gestionnaire de thème multiple : ces responsabilités restent dans l'application qui utilise AGUIK.

## 14. Ajouter un composant

Pour créer un nouveau composant :

1. décider s'il est stateless ou stateful ;
2. pour un composant interactif, conserver les widgets Gio dans sa structure ;
3. exposer une méthode `Layout(gtx, th)` si le composant porte de l'état ;
4. mesurer le contenu avant de calculer le cadre si la taille dépend du texte ;
5. réutiliser `NeoFrame`, `Text`, `Icon`, `glayout` et `interaction.ClickArea` au lieu de dupliquer leur logique ;
6. appeler `interaction.Invalidate` après les transitions d'état ;
7. retourner les dimensions complètes, y compris l'espace du cadre et de l'ombre ;
8. ajouter l'usage au showcase pour vérifier les contraintes, le survol, l'appui et les overlays.

Cette organisation garde le rendu visuel centralisé dans le thème et `NeoFrame`, tout en laissant à chaque composant la responsabilité de son état propre.
