# Guide d'utilisation des composants AGUIK

Ce guide explique comment utiliser chaque composant du package `components`, quels champs configurer et quelles valeurs lire après une interaction.

AGUIK repose sur Gio. Tous les composants doivent être rendus pendant une frame avec un `layout.Context` et un `*theme.Theme`.

```go
func view(gtx layout.Context, th *theme.Theme, state *State) layout.Dimensions {
    return state.save.Layout(gtx, th)
}
```

## Règle essentielle : conserver l'état

Les composants interactifs doivent être créés une seule fois puis conservés entre les frames. Ne les instanciez pas dans la fonction de rendu.

```go
type State struct {
    saveButton components.Button
    name        components.Input
    terms       components.Checkbox
}

state := State{
    saveButton: components.Button{
        Text:    "Enregistrer",
        Variant: components.ButtonPrimary,
        Size:    components.SizeMD,
    },
}
```

Un composant comme `Button`, `Input`, `Dialog` ou `Toaster` contient de l'état Gio interne. Le recréer à chaque frame ferait perdre les clics, le focus, le texte saisi ou les timers.

## Variantes et tailles communes

### Variantes

Le type `components.ButtonVariant` est partagé par plusieurs composants :

```go
components.ButtonDefault
components.ButtonPrimary
components.ButtonSecondary
components.ButtonSuccess
components.ButtonWarning
components.ButtonError
```

La variante contrôle généralement la couleur principale. `ButtonError` utilise `theme.Colors.Danger`.

### Tailles

Le type `components.ButtonSize` est également partagé :

```go
components.SizeSM
components.SizeMD
components.SizeLG
```

La taille modifie principalement le padding, la taille des icônes et la taille des contrôles.

## Texte

### `Text`

`Text` est un widget sans état.

```go
components.Text(th, "Bonjour")
components.Text(th, "Erreur", components.WithColor(th.Colors.Danger))
components.Text(th, "Titre", components.WithSize(24))
```

Signature :

```go
components.Text(th *theme.Theme, value string, options ...components.TextOption) layout.Widget
```

Options :

| Option | Effet |
|---|---|
| `WithSize(sp)` | Définit la taille en sp. |
| `WithColor(color)` | Définit la couleur du texte. |
| `WithFont(font.Font)` | Remplace la police. |
| `WithAlignment(text.Alignment)` | Définit l'alignement du texte. |
| `WithMaxLines(n)` | Limite le nombre de lignes. |
| `WithWeight(font.Weight)` | Définit le poids, par exemple `font.Bold`. |
| `WithEllipsis()` | Limite à une ligne si nécessaire et ajoute une ellipse. |

Le `Theme` doit avoir un `Shaper` configuré pour afficher du texte.

### Variantes sémantiques

Ces fonctions préconfigurent `Text` :

```go
components.Title(th, "Titre principal")
components.Subtitle(th, "Sous-titre")
components.Body(th, "Texte courant")
components.Caption(th, "Information secondaire")
```

Elles acceptent les mêmes `TextOption`, qui sont appliquées après leur style par défaut.

## Icônes

### `Icon`

Dessine un `*widget.Icon`.

```go
components.Icon(th, homeIcon)
components.Icon(th, homeIcon,
    components.WithIconSize(32),
    components.WithIconColor(th.Colors.Primary),
)
```

Options :

- `WithIconSize(dp)` : taille carrée de l'icône en dp ;
- `WithIconColor(color)` : couleur de l'icône.

La taille par défaut est 24 dp.

### `IconText`

Compose une icône et un texte horizontalement.

```go
components.IconText(
    th,
    homeIcon,
    "Accueil",
    []components.IconOption{components.WithIconSize(18)},
    []components.TextOption{components.WithWeight(font.SemiBold)},
)
```

## Boutons

### `Button`

```go
button := components.Button{
    Text:     "Enregistrer",
    Icon:     checkIcon,
    Variant:  components.ButtonPrimary,
    Size:     components.SizeMD,
    Disabled: false,
}
```

Champs :

| Champ | Description |
|---|---|
| `Clickable` | État Gio interne. Ne pas recréer entre les frames. |
| `Text` | Libellé du bouton. |
| `Icon` | Icône optionnelle. |
| `Variant` | Couleur sémantique. |
| `Size` | `SizeSM`, `SizeMD` ou `SizeLG`. |
| `Disabled` | Désactive le bouton et sa hitbox. |
| `Clicked` | Vaut `true` pendant la frame où le clic est consommé. |

Utilisation :

```go
dims := state.saveButton.Layout(gtx, th)
if state.saveButton.Clicked {
    saveProfile()
}
return dims
```

`Clicked` n'est pas un état permanent. La logique métier doit être traitée immédiatement dans le rendu ou copiée dans un état applicatif.

### `IconButton`

Bouton dont le contenu est une icône.

```go
button := components.IconButton{
    Icon:          settingsIcon,
    Variant:       components.ButtonDefault,
    Size:          components.SizeMD,
    TooltipText:   "Paramètres",
    TooltipPlacement: components.TooltipRight,
}
```

Champs supplémentaires :

- `TooltipText` : texte de la bulle ;
- `TooltipPlacement` : `TooltipTop`, `TooltipBottom` ou `TooltipRight`.

Le tooltip est géré directement par `IconButton`. N'entourez pas ce composant avec `Tooltip`.

## Éléments visuels

### `Badge`

Étiquette non interactive pour un statut ou un tag.

```go
badge := components.Badge{
    Text:    "ACTIF",
    Variant: components.ButtonSuccess,
    Size:    components.SizeSM,
}
return badge.Layout(gtx, th)
```

Champs : `Text`, `Variant` et `Size`. Un badge ne possède pas de clic et ne réagit pas au survol.

### `Card`

`Card` enveloppe n'importe quel `layout.Widget` dans un cadre visuel.

```go
card := components.Card(
    th,
    glayout.Column(th.Spacing.SM, layout.Start,
        glayout.Rigid(components.Subtitle(th, "Profil")),
        glayout.Rigid(components.Text(th, "Contenu")),
    ),
    components.WithCardPadding(th.Spacing.LG),
)
return card(gtx)
```

Options :

- `WithCardPadding(dp)` : padding intérieur ;
- `WithCardBackground(color)` : couleur de surface.

Les valeurs par défaut viennent du thème.

### `Divider`

```go
horizontal := components.Divider{}
vertical := components.Divider{Vertical: true, Length: 120}
```

Champs :

- `Vertical` : orientation verticale si `true` ;
- `Thickness` : épaisseur en dp, par défaut la bordure du thème ;
- `Length` : longueur en dp, ou tout l'espace disponible si `0`.

## Saisie utilisateur

### `Input`

Champ de saisie sur une ligne.

```go
input := components.Input{
    Placeholder: "Adresse e-mail",
    Width:       280,
    Size:        components.SizeMD,
}
```

Champs :

| Champ | Description |
|---|---|
| `Editor` | `widget.Editor` Gio contenant la saisie. |
| `Placeholder` | Texte affiché quand le champ est vide. |
| `Width` | Largeur en dp, 200 dp par défaut. |
| `Size` | Taille du padding. |
| `Disabled` | Rend l'éditeur en lecture seule et désactive la hitbox. |
| `Password` | Masque les caractères saisis. |

Lire la valeur :

```go
value := state.email.Editor.Text()
```

### `TextArea`

Champ multiligne basé sur le même `widget.Editor`.

```go
bio := components.TextArea{
    Placeholder: "Présentez-vous",
    Width:       320,
    Lines:       4,
    Size:        components.SizeMD,
}
```

Champs supplémentaires :

- `Lines` : nombre de lignes visibles ; 4 par défaut si `0` ;
- `Editor`, `Placeholder`, `Width`, `Size` et `Disabled` ont le même rôle que pour `Input`.

Lire le texte avec `state.bio.Editor.Text()`.

## Contrôles booléens

### `Checkbox`

```go
terms := components.Checkbox{
    Label:    "J'accepte les conditions",
    Variant:  components.ButtonSuccess,
    Size:     components.SizeMD,
}
```

Champs :

- `Bool` : état Gio ;
- `Label` : texte associé ;
- `Variant` : couleur quand coché ;
- `Size` : taille de la case ;
- `Disabled` : désactive le contrôle.

Lire l'état :

```go
accepted := state.terms.Bool.Value
```

### `Switch`

```go
notifications := components.Switch{
    Variant: components.ButtonSuccess,
    Size:    components.SizeMD,
}
```

Champs : `Bool`, `Variant`, `Size` et `Disabled`.

Lire l'état avec `state.notifications.Bool.Value`.

## Sélections

### `RadioGroup`

Sélection exclusive parmi plusieurs options.

```go
plans := components.RadioGroup{
    Options: []components.RadioOption{
        {Label: "Gratuit"},
        {Label: "Pro"},
        {Label: "Entreprise"},
    },
    Selected: -1,
    Variant:  components.ButtonPrimary,
    Size:     components.SizeMD,
}
```

Champs :

- `Options` : liste de `RadioOption` ;
- `Selected` : index actif, `-1` si aucune sélection ;
- `Variant` et `Size` ;
- `Disabled`.

Chaque `RadioOption` contient un champ `Label`.

```go
if state.plans.Selected >= 0 {
    selectedPlan := state.plans.Options[state.plans.Selected].Label
}
```

### `Tabs`

Onglets horizontaux avec un seul onglet actif.

```go
tabs := components.Tabs{
    Items: []components.TabItem{
        {Label: "Général"},
        {Label: "Sécurité"},
        {Label: "Notifications"},
    },
    Variant: components.ButtonPrimary,
    Active:  0,
}
```

Champs : `Items`, `Active`, `Variant`, `Size` et `Disabled`.

Chaque `TabItem` contient un `Label`.

```go
state.tabs.Layout(gtx, th)
switch state.tabs.Active {
case 0:
    return generalPanel(gtx, th)
case 1:
    return securityPanel(gtx, th)
default:
    return notificationsPanel(gtx, th)
}
```

### `Select`

Dropdown à sélection unique.

```go
role := components.Select{
    Placeholder: "Choisir un rôle...",
    Options: []components.SelectOption{
        {Label: "Administrateur"},
        {Label: "Éditeur"},
        {Label: "Lecteur"},
    },
    Selected: -1,
    Width:    240,
    Size:     components.SizeMD,
}
```

Champs :

- `Options` : options disponibles ;
- `Selected` : index choisi, `-1` si aucune sélection ;
- `Placeholder` : texte avant sélection ;
- `Width` : largeur en dp, 200 dp par défaut ;
- `Size` et `Disabled`.

Chaque `SelectOption` contient un `Label`.

```go
state.role.Layout(gtx, th)
if state.role.Selected >= 0 {
    roleName := state.role.Options[state.role.Selected].Label
}
```

Le menu est affiché au-dessus ou en dessous du champ selon l'espace disponible dans les contraintes Gio du parent.

## Progression et chargement

### `ProgressBar`

Barre de progression déterminée.

```go
progress := components.ProgressBar{
    Value:   0.65,
    Variant: components.ButtonPrimary,
    Width:   240,
    Height:  12,
}
```

Champs :

- `Value` : progression attendue entre `0` et `1` ;
- `Variant` : couleur du remplissage ;
- `Width` : 200 dp par défaut ;
- `Height` : 12 dp par défaut.

Les valeurs inférieures à `0` sont ramenées à `0`, les valeurs supérieures à `1` à `1`.

### `Loader`

Spinner animé indéterminé.

```go
loader := components.Loader{
    Size:    24,
    Variant: components.ButtonPrimary,
}
return loader.Layout(gtx, th)
```

Champs : `Size` en dp, 24 dp par défaut, et `Variant`.

Le `Loader` doit rester persistant : il conserve son instant de démarrage et invalide les frames pour faire tourner l'animation.

## Navigation et aide

### `Sidebar`

Barre verticale d'icônes avec un item actif.

```go
sidebar := components.Sidebar{
    Items: []components.SidebarItem{
        {Icon: homeIcon, Tooltip: "Accueil"},
        {Icon: settingsIcon, Tooltip: "Paramètres"},
    },
    Active:  0,
    Variant: components.ButtonPrimary,
    Size:    components.SizeMD,
}
```

Champs : `Items`, `Active`, `Variant` et `Size`.

Chaque `SidebarItem` contient `Icon` et `Tooltip`. Le composant crée et conserve un `IconButton` par item.

```go
state.sidebar.Layout(gtx, th)
current := state.sidebar.Active
```

### `Tooltip`

Wrapper pour un widget non interactif :

```go
tooltip := components.Tooltip{
    Text:      "Accueil",
    Placement: components.TooltipBottom,
}

return tooltip.Layout(gtx, th, components.Icon(th, homeIcon))
```

Placements disponibles : `TooltipTop`, `TooltipBottom` et `TooltipRight`.

Utilisez ce wrapper pour `Icon`, `Text` ou un autre widget sans gestion de survol. Pour `IconButton`, utilisez plutôt `TooltipText` et `TooltipPlacement` directement sur le bouton.

## Notifications et fenêtres modales

### `Alert`

Alerte autonome avec titre, message et fermeture optionnelle.

```go
alert := components.Alert{
    Title:       "Enregistré",
    Message:     "Votre profil a été mis à jour.",
    Variant:     components.ButtonSuccess,
    Dismissible: true,
    MaxWidth:    320,
}
return alert.Layout(gtx, th)
```

Champs : `Icon`, `Title`, `Message`, `Variant`, `Dismissible` et `MaxWidth`.

Après le layout, `alert.Closed()` indique si l'utilisateur a fermé l'alerte. Un conteneur doit alors la retirer de son état.

### `Dialog`

Fenêtre modale centrée avec backdrop.

```go
dialog := components.Dialog{
    Title:                  "Supprimer le profil ?",
    Message:                "Cette action est irréversible.",
    Width:                  360,
    DismissOnBackdropClick: true,
    AutoClose:              true,
    Actions: []components.DialogAction{
        {Text: "Annuler", Variant: components.ButtonDefault},
        {Text: "Supprimer", Variant: components.ButtonError},
    },
}
```

Champs :

- `Open` : affiche ou masque le dialog ;
- `Title` et `Message` ;
- `Actions` : boutons d'action ;
- `Width` : 320 dp par défaut ;
- `DismissOnBackdropClick` : ferme sur clic hors de la carte ;
- `AutoClose` : ferme après une action ;
- `Clicked` : index de l'action cliquée pendant la frame courante.

Chaque `DialogAction` contient `Text` et `Variant`.

Utilisation recommandée :

```go
state.dialog.Layout(gtx, th)
if state.dialog.Clicked >= 0 {
    action := state.dialog.Actions[state.dialog.Clicked]
    switch action.Text {
    case "Supprimer":
        deleteProfile()
    case "Annuler":
        // Rien à faire.
    }
}
```

Un dialog doit être placé dans un `glayout.Stack` après le contenu principal pour apparaître au-dessus de celui-ci.

### `Toaster`

Pile d'alertes en bas à droite.

```go
toaster := components.Toaster{
    Gap:      8,
    Margin:   16,
    MaxWidth: 320,
}
```

Champs :

- `Gap` : espace entre les toasts, 8 dp par défaut ;
- `Margin` : marge aux bords, 16 dp par défaut ;
- `MaxWidth` : largeur maximale transmise aux alertes.

Ajouter un toast :

```go
state.toaster.Show(components.ToastConfig{
    Title:       "Enregistré",
    Message:     "Les modifications sont sauvegardées.",
    Variant:     components.ButtonSuccess,
    Dismissible: true,
    Duration:    4 * time.Second,
})
```

`ToastConfig` contient `Icon`, `Title`, `Message`, `Variant`, `Dismissible` et `Duration`. Une durée nulle désactive la fermeture automatique.

Le `Toaster` doit être layouté à chaque frame :

```go
state.toaster.Layout(gtx, th)
```


## Données tabulaires

### `Table`

`Table` est générique (`Table[T any]`) : chaque ligne est une valeur de ton propre type `T`. Contrairement aux autres composants, tu ne lui donnes pas directement du texte à afficher, mais des fonctions qui savent extraire ce texte depuis une ligne `T`.

```go
type Product struct {
    Name  string
    Price string
}

table := components.Table[Product]{
    Columns: []components.TableColumn[Product]{
        {
            Label:  "Nom",
            Weight: 2,
            Cell:   func(p Product) string { return p.Name },
            Less:   func(a, b Product) bool { return a.Name < b.Name },
        },
        {
            Label:  "Prix",
            Weight: 1,
            Cell:   func(p Product) string { return p.Price },
            Less:   func(a, b Product) bool { return a.Price < b.Price },
        },
    },
    Rows: []Product{
        {Name: "Clavier", Price: "49€"},
        {Name: "Souris", Price: "19€"},
    },
}
```

Champs :

| Champ | Description |
|---|---|
| `Columns` | Liste de `TableColumn[T]`, une par colonne. |
| `Rows` | Les données affichées, un `T` par ligne. |
| `SortColumn` | Index de la colonne actuellement triée, `-1` si aucune. |
| `SortDirection` | `SortNone`, `SortAsc` ou `SortDesc`. |
| `RowActions` | Fonction optionnelle construisant un widget à droite de chaque ligne (bouton supprimer, éditer...). `nil` = pas d'action. |

Chaque `TableColumn[T]` contient :

| Champ | Description |
|---|---|
| `Label` | Texte affiché dans l'en-tête. |
| `Weight` | Largeur relative de la colonne ; `0` = partage égal entre colonnes. |
| `Cell` | Extrait le texte à afficher pour une ligne donnée. |
| `Less` | Comparateur pour le tri. `nil` = colonne non triable (pas de flèche, clic ignoré). |

Le tri est géré entièrement par `Table` : cliquer un en-tête triable bascule `Aucun → Ascendant → Descendant → Aucun` et retrie `Rows` sur place, sans intervention de l'appelant.

### Ajouter une action par ligne

```go
table.RowActions = func(item Product, index int) layout.Widget {
    return func(gtx layout.Context) layout.Dimensions {
        var del components.IconButton
        del.Icon = trashIcon
        del.Variant = components.ButtonError
        del.Size = components.SizeSM

        dims := del.Layout(gtx, th)
        if del.Clicked {
            table.Rows = append(table.Rows[:index], table.Rows[index+1:]...)
        }
        return dims
    }
}
```

`index` correspond à la position de la ligne **au moment du rendu de cette frame** — si `Rows` est retrié entre deux frames, `index` reste cohérent avec l'ordre affiché à cet instant, pas avec un identifiant stable dans le temps. Pour des suppressions fiables sur des listes qui changent souvent, préfère un champ d'identifiant unique dans `T` plutôt que de te fier uniquement à `index`.

### Persistance de l'état

Comme `Button` ou `Select`, `Table` contient un état Gio interne (`Clickable` par en-tête pour le tri) — elle doit être un champ de ton `State`, jamais recréée dans le rendu :

```go
type State struct {
    products components.Table[Product]
}
```

```go
return state.products.Layout(gtx, th)
```


## Layouts AGUIK

Les composants sont souvent composés avec `glayout`.

```go
content := glayout.Column(th.Spacing.MD, layout.Start,
    glayout.Rigid(components.Title(th, "Profil")),
    glayout.Rigid(func(gtx layout.Context) layout.Dimensions {
        return state.name.Layout(gtx, th)
    }),
    glayout.Rigid(func(gtx layout.Context) layout.Dimensions {
        return state.save.Layout(gtx, th)
    }),
)
```

Helpers disponibles :

- `Rigid(widget)` : taille naturelle ;
- `Flexed(weight, widget)` : partage de l'espace restant ;
- `Row(spacing, alignment, children...)` : composition horizontale ;
- `Column(spacing, alignment, children...)` : composition verticale ;
- `Inset(space, widget)` : padding uniforme ;
- `Grid(spacing, columns, children...)` : grille à colonnes égales ;
- `Stack(children...)` : superposition plein écran ;
- `Spacer(space)` : espace fixe.

Exemple d'overlay complet :

```go
return glayout.Stack(
    mainContent,
    func(gtx layout.Context) layout.Dimensions {
        return state.dialog.Layout(gtx, th)
    },
    func(gtx layout.Context) layout.Dimensions {
        return state.toaster.Layout(gtx, th)
    },
)(gtx)
```

Le dernier enfant est dessiné au-dessus des précédents.

## Pattern recommandé pour une action

```go
func saveButton(th *theme.Theme, state *State) layout.Widget {
    return func(gtx layout.Context) layout.Dimensions {
        dims := state.save.Layout(gtx, th)
        if state.save.Clicked {
            state.lastAction = "Profil enregistré"
            state.toaster.Show(components.ToastConfig{
                Title:    "Enregistré",
                Variant:  components.ButtonSuccess,
                Duration: 3 * time.Second,
            })
        }
        return dims
    }
}
```

Le composant est layouté une fois, puis son événement est traité juste après. Les structures `Button`, `Dialog` et `Toaster` restent dans `State`.
