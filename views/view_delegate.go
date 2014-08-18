package views

type ViewDelegate interface {
	OnMouseEnter(event *MouseEvent)
	OnMouseLeave(event *MouseEvent)
	OnDraw(event *DrawEvent)
}
