package views

type Toolbar struct {
	BaseView
}

func NewToolbar() *Toolbar {
	toolbar := new(Toolbar)
	return toolbar
}

func (t *Toolbar) OnDraw(event *DrawEvent) {
	event.Canvas.DrawColor(19, 19, 19)
}
