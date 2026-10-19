package main

import (
	"strconv"
	"time"

	"gioui.org/widget"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/Shiro-Masuta/AGUIK/components"
	"github.com/Shiro-Masuta/AGUIK/components/table"
)

type product struct {
	Name     string
	Category string
	Price    int // euros
}

// state regroupe tous les composants stateful : créé une fois, jamais dans draw.
type state struct {
	sidebar components.Sidebar
	toaster components.Toaster
	dialog  components.Dialog

	// Page composants
	switchNotif components.Switch
	progress    components.ProgressBar
	progressBtn components.Button
	loader      components.Loader
	role        components.Select
	openDialog  components.Button

	// Page formulaire
	name  components.Input
	pass  components.Input
	bio   components.TextArea
	terms components.Checkbox
	plan  components.RadioGroup
	save  components.Button

	// Page table
	products table.Table[product]

	lastAction string
}

func mustIcon(data []byte) *widget.Icon {
	ic, err := widget.NewIcon(data)
	if err != nil {
		panic(err)
	}
	return ic
}

func newState() *state {
	return &state{
		sidebar: components.Sidebar{
			Items: []components.SidebarItem{
				{Icon: mustIcon(icons.ActionHome), Tooltip: "Composants"},
				{Icon: mustIcon(icons.ActionSettings), Tooltip: "Formulaire"},
				{Icon: mustIcon(icons.ActionList), Tooltip: "Table"},
			},
			Active:  0,
			Variant: components.ButtonPrimary,
			Size:    components.SizeMD,
		},
		toaster: components.Toaster{Gap: 8, Margin: 16, MaxWidth: 320},
		dialog: components.Dialog{
			Title:                  "Supprimer le profil ?",
			Message:                "Cette action est irréversible.",
			Width:                  360,
			DismissOnBackdropClick: true,
			AutoClose:              true,
			Actions: []components.DialogAction{
				{Text: "Annuler", Variant: components.ButtonDefault},
				{Text: "Supprimer", Variant: components.ButtonError},
			},
		},

		switchNotif: components.Switch{Variant: components.ButtonSuccess, Size: components.SizeMD},
		progress:    components.ProgressBar{Value: 0.3, Variant: components.ButtonPrimary, Width: 240},
		progressBtn: components.Button{Text: "+10 %", Variant: components.ButtonSecondary, Size: components.SizeSM},
		loader:      components.Loader{Size: 28, Variant: components.ButtonPrimary},
		role: components.Select{
			Placeholder: "Choisir un rôle...",
			Options: []components.SelectOption{
				{Label: "Administrateur"}, {Label: "Éditeur"}, {Label: "Lecteur"},
			},
			Selected: -1,
			Width:    240,
			Size:     components.SizeMD,
		},
		openDialog: components.Button{Text: "Ouvrir le dialog", Variant: components.ButtonError, Size: components.SizeMD},

		name:  components.Input{Placeholder: "Nom", Width: 280, Size: components.SizeMD},
		pass:  components.Input{Placeholder: "Mot de passe", Width: 280, Size: components.SizeMD, Password: true},
		bio:   components.TextArea{Placeholder: "Présente-toi", Width: 360, Lines: 4, Size: components.SizeMD},
		terms: components.Checkbox{Label: "J'accepte les conditions", Variant: components.ButtonSuccess, Size: components.SizeMD},
		plan: components.RadioGroup{
			Options:  []components.RadioOption{{Label: "Gratuit"}, {Label: "Pro"}, {Label: "Entreprise"}},
			Selected: -1,
			Variant:  components.ButtonPrimary,
			Size:     components.SizeMD,
		},
		save: components.Button{Text: "Enregistrer", Variant: components.ButtonPrimary, Size: components.SizeMD},

		products: table.Table[product]{
			Columns: []table.TableColumn[product]{
				{Label: "Nom", Weight: 2, Cell: func(p product) string { return p.Name },
					Less: func(a, b product) bool { return a.Name < b.Name }},
				{Label: "Catégorie", Weight: 2, Cell: func(p product) string { return p.Category },
					Less: func(a, b product) bool { return a.Category < b.Category }},
				{Label: "Prix", Weight: 1, Cell: func(p product) string { return strconv.Itoa(p.Price) + "€" },
					Less: func(a, b product) bool { return a.Price < b.Price }},
			},
			Rows: []product{
				{"Clavier", "Périphérique", 49},
				{"Souris", "Périphérique", 19},
				{"Écran 27\"", "Affichage", 239},
				{"Webcam", "Vidéo", 59},
			},
			SortColumn: -1,
		},

		lastAction: "Aucune action",
	}
}

func (s *state) toast(title, msg string, v components.ButtonVariant) {
	s.lastAction = title
	s.toaster.Show(components.ToastConfig{
		Title:       title,
		Message:     msg,
		Variant:     v,
		Dismissible: true,
		Duration:    3 * time.Second,
	})
}
