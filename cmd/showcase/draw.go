package main

import (
	"gioui.org/layout"
	"gioui.org/op/paint"

	"github.com/Shiro-Masuta/AGUIK/components"
	"github.com/Shiro-Masuta/AGUIK/glayout"
	"github.com/Shiro-Masuta/AGUIK/theme"
)

type themed interface {
	Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions
}

// bind adapte un composant AGUIK en layout.Widget.
func bind(th *theme.Theme, c themed) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions { return c.Layout(gtx, th) }
}

func (s *state) layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	paint.Fill(gtx.Ops, th.Colors.Background)

	main := glayout.Inset(th.Spacing.LG,
		glayout.Row(th.Spacing.LG, layout.Start,
			glayout.Rigid(bind(th, &s.sidebar)),
			glayout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return s.page(gtx, th)
			}),
		),
	)

	// Ordre = z-order : le dernier enfant est dessiné au-dessus.
	return glayout.Stack(
		main,
		func(gtx layout.Context) layout.Dimensions {
			dims := s.dialog.Layout(gtx, th)
			if s.dialog.Clicked >= 0 && s.dialog.Clicked < len(s.dialog.Actions) {
				switch s.dialog.Actions[s.dialog.Clicked].Text {
				case "Supprimer":
					s.toast("Profil supprimé", "C'est fait.", components.ButtonError)
				case "Annuler":
					s.toast("Annulé", "Rien n'a changé.", components.ButtonDefault)
				}
			}
			return dims
		},
		bind(th, &s.toaster),
	)(gtx)
}

func (s *state) page(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	switch s.sidebar.Active {
	case 1:
		return s.formPage(gtx, th)
	case 2:
		return s.tablePage(gtx, th)
	default:
		return s.componentsPage(gtx, th)
	}
}

func (s *state) componentsPage(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return glayout.Column(th.Spacing.MD, layout.Start,
		glayout.Rigid(components.Title(th, "Composants")),
		glayout.Rigid(glayout.Row(th.Spacing.SM, layout.Middle,
			glayout.Rigid(badge(th, "ACTIF", components.ButtonSuccess)),
			glayout.Rigid(badge(th, "BETA", components.ButtonWarning)),
			glayout.Rigid(badge(th, "ERREUR", components.ButtonError)),
		)),
		glayout.Rigid(bind(th, &components.Divider{})),
		glayout.Rigid(glayout.Row(th.Spacing.SM, layout.Middle,
			glayout.Rigid(bind(th, &s.switchNotif)),
			glayout.Rigid(components.Body(th, "Notifications")),
		)),
		glayout.Rigid(glayout.Row(th.Spacing.MD, layout.Middle,
			glayout.Rigid(bind(th, &s.progress)),
			glayout.Rigid(func(gtx layout.Context) layout.Dimensions {
				dims := s.progressBtn.Layout(gtx, th)
				if s.progressBtn.Clicked {
					s.progress.Value += 0.1
					if s.progress.Value > 1 {
						s.progress.Value = 0
					}
				}
				return dims
			}),
			glayout.Rigid(bind(th, &s.loader)),
		)),
		glayout.Rigid(bind(th, &s.role)),
		glayout.Rigid(func(gtx layout.Context) layout.Dimensions {
			dims := s.openDialog.Layout(gtx, th)
			if s.openDialog.Clicked {
				s.dialog.Open = true
			}
			return dims
		}),
		glayout.Rigid(components.Caption(th, "Dernière action : "+s.lastAction)),
	)(gtx)
}

func (s *state) formPage(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return glayout.Column(th.Spacing.MD, layout.Start,
		glayout.Rigid(components.Title(th, "Formulaire")),
		glayout.Rigid(bind(th, &s.name)),
		glayout.Rigid(bind(th, &s.pass)),
		glayout.Rigid(bind(th, &s.bio)),
		glayout.Rigid(bind(th, &s.plan)),
		glayout.Rigid(bind(th, &s.terms)),
		glayout.Rigid(func(gtx layout.Context) layout.Dimensions {
			dims := s.save.Layout(gtx, th)
			if s.save.Clicked {
				switch {
				case !s.terms.Bool.Value:
					s.toast("Conditions", "Accepte les conditions d'abord.", components.ButtonWarning)
				default:
					s.toast("Enregistré", "Bonjour "+s.name.Editor.Text()+" !", components.ButtonSuccess)
				}
			}
			return dims
		}),
	)(gtx)
}

func (s *state) tablePage(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	return glayout.Column(th.Spacing.MD, layout.Start,
		glayout.Rigid(components.Title(th, "Table")),
		glayout.Rigid(components.Caption(th, "Clique un en-tête pour trier")),
		glayout.Rigid(bind(th, &s.products)),
	)(gtx)
}

func badge(th *theme.Theme, text string, v components.ButtonVariant) layout.Widget {
	b := components.Badge{Text: text, Variant: v, Size: components.SizeSM}
	return bind(th, &b)
}
