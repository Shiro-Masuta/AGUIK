package components

import (
	"image"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/Shiro-Masuta/AGUIK/theme"
)

// ToastConfig décrit un toast à afficher — l'appelant ne connaît jamais
// le widget.Clickable interne d'Alert, Toaster gère le cycle de vie.
type ToastConfig struct {
	Icon        *widget.Icon
	Title       string
	Message     string
	Variant     ButtonVariant
	Dismissible bool
	Duration    time.Duration // 0 = pas d'auto-dismiss (fermeture manuelle uniquement)
}

type toastEntry struct {
	alert    Alert
	deadline time.Time // zero si pas d'auto-dismiss
}

// Toaster empile des Alert en bas à droite, chacun avec son propre
// auto-dismiss optionnel. Champ persistant de l'état appelant, comme tout
// composant stateful du kit (jamais recréé à chaque frame).
type Toaster struct {
	Gap      float32 // dp entre deux toasts ; 0 -> defaultToasterGapDp
	Margin   float32 // dp depuis le bord ; 0 -> defaultToasterMarginDp
	MaxWidth float32 // transmis à chaque Alert

	entries []toastEntry
}

const (
	defaultToasterGapDp    = 8
	defaultToasterMarginDp = 16
)

// Show ajoute un toast à la pile. Safe à appeler depuis n'importe quel
// callback de clic pendant Layout — le rendu suit dès la frame courante
// si appelé avant que Toaster.Layout ne tourne, sinon à la suivante.
func (t *Toaster) Show(cfg ToastConfig) {
	e := toastEntry{
		alert: Alert{
			Icon:        cfg.Icon,
			Title:       cfg.Title,
			Message:     cfg.Message,
			Variant:     cfg.Variant,
			Dismissible: cfg.Dismissible,
			MaxWidth:    t.MaxWidth,
		},
	}
	if cfg.Duration > 0 {
		e.deadline = time.Now().Add(cfg.Duration)
	}
	t.entries = append(t.entries, e)
}

func (t *Toaster) Layout(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	now := time.Now()

	// Purge des toasts expirés AVANT tout dessin, comme tout état lu en
	// début de frame dans ce kit.
	live := t.entries[:0]
	for _, e := range t.entries {
		if !e.deadline.IsZero() && !now.Before(e.deadline) {
			continue
		}
		live = append(live, e)
	}
	t.entries = live

	gapDp := t.Gap
	if gapDp == 0 {
		gapDp = defaultToasterGapDp
	}
	marginDp := t.Margin
	if marginDp == 0 {
		marginDp = defaultToasterMarginDp
	}
	gap := gtx.Dp(unit.Dp(gapDp))
	margin := gtx.Dp(unit.Dp(marginDp))

	full := gtx.Constraints.Max
	y := full.Y - margin

	var nextWake time.Time
	for i := len(t.entries) - 1; i >= 0; i-- {
		e := &t.entries[i]

		mgtx := gtx
		mgtx.Constraints = layout.Constraints{Min: image.Point{}, Max: full}
		m := op.Record(mgtx.Ops)
		dims := e.alert.Layout(mgtx, th)
		call := m.Stop()

		if e.alert.Closed() {
			// Croix cliquée cette frame : retiré, pas dessiné.
			t.entries = append(t.entries[:i], t.entries[i+1:]...)
			continue
		}

		x := full.X - margin - dims.Size.X
		y -= dims.Size.Y
		off := op.Offset(image.Pt(x, y)).Push(gtx.Ops)
		call.Add(gtx.Ops)
		off.Pop()
		y -= gap

		if !e.deadline.IsZero() && (nextWake.IsZero() || e.deadline.Before(nextWake)) {
			nextWake = e.deadline
		}
	}

	// Réveille Gio pile au moment où le prochain toast doit expirer —
	// sans ça, l'auto-dismiss ne se déclencherait qu'au prochain
	// mouvement de souris fortuit, exactement comme le bug du Dialog.
	if !nextWake.IsZero() {
		gtx.Execute(op.InvalidateCmd{At: nextWake})
	}

	return layout.Dimensions{Size: full}
}
