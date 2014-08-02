package views

import (
	. "gwk/vango"
	. "image"
	"log"
)

type View struct {
	id string

	canvas *Canvas

	x, y int // Relative to Parent.
	w, h int

	children []Viewer
	parent   Viewer

	uimap    UIMap
	layouter Layouter
}

func NewView() *View {
	var v = new(View)
	return v
}

func (v *View) ID() string {
	return v.id
}

func (v *View) SetID(id string) {
	v.id = id
}

func (v *View) AddChild(child Viewer) {
	v.children = append(v.children, child)
	child.SetParent(v)
}

func (v *View) Children() []Viewer {
	return v.children
}

func (v *View) Parent() Viewer {
	return v.parent
}

func (v *View) SetParent(parent Viewer) {
	v.parent = parent
}

func (v *View) Canvas() *Canvas {
	if v.canvas == nil {
		v.canvas = NewCanvas(v.W(), v.H())
	}
	canvas_bounds := v.canvas.Bounds()
	if canvas_bounds.Dx() < v.W() || canvas_bounds.Dy() < v.H() {
		v.canvas = NewCanvas(v.W(), v.H())
	} else {
		return v.canvas.SubCanvas(Rect(0, 0, v.W(), v.H()))
	}
	return v.canvas
}

func (v *View) SetCanvas(canvas *Canvas) {
	v.canvas = canvas
}

func (v *View) ToAbsPoint(pt Point) Point {
	if v.Parent() == nil {
		return pt
	}
	pt.X = pt.X + v.X()
	pt.Y = pt.Y + v.Y()
	return v.Parent().ToAbsPoint(pt)
}

func (v *View) ToDevicePoint(pt Point) Point {
	pt = pt.Add(Pt(v.X(), v.Y()))
	return v.Parent().ToDevicePoint(pt)
}

func (v *View) ToAbsRect(rect Rectangle) Rectangle {
	if v.Parent() == nil {
		return rect
	}
	rect.Min.X = rect.Min.X + v.X()
	rect.Min.Y = rect.Min.Y + v.Y()
	rect.Max.X = rect.Max.X + v.X()
	rect.Max.Y = rect.Max.Y + v.Y()
	return v.Parent().ToAbsRect(rect)
}

func (v *View) ToDeviceRect(rect Rectangle) Rectangle {
	rect = rect.Add(Pt(v.X(), v.Y()))
	return v.Parent().ToDeviceRect(rect)
}

func update_rect(v Viewer, rect Rectangle) {
	if v.Parent() == nil {
		v.ScheduleDrawRect(rect)
		return
	}
	rect = rect.Add(Pt(v.X(), v.Y()))
	update_rect(v.Parent(), rect)
}

func (v *View) ScheduleDraw() {
	if v.Parent() != nil {
		update_rect(v.Parent(), v.Bounds())
	}
}

func (v *View) ScheduleDrawRect(rect Rectangle) {
	if v.Parent() != nil {
		update_rect(v.Parent(), rect)
	}
}

func (v *View) X() int {
	return v.x
}

func (v *View) Y() int {
	return v.y
}

func (v *View) W() int {
	return v.w
}

func (v *View) H() int {
	return v.h
}

func (v *View) XYWH() (x, y, w, h int) {
	x, y, w, h = v.x, v.y, v.w, v.h
	return
}

func (v *View) SetXYWH(x, y, w, h int) {
	v.x, v.y, v.w, v.h = x, y, w, h
}

func (v *View) Width() int {
	return v.W()
}

func (v *View) SetWidth(width int) {
	v.w = width
}

func (v *View) Height() int {
	return v.H()
}

func (v *View) SetHeight(height int) {
	v.h = height
}

func (v *View) Left() int {
	return v.X()
}

func (v *View) SetLeft(left int) {
	v.x = left
}

func (v *View) Top() int {
	return v.Y()
}

func (v *View) SetTop(top int) {
	v.y = top
}

func (v *View) Bounds() Rectangle {
	return Rect(v.x, v.y, v.x+v.w, v.y+v.h)
}

func (v *View) LocalBounds() Rectangle {
	return Rect(0, 0, v.w, v.h)
}

func (v *View) SetBounds(bounds Rectangle) {
	v.x, v.y = bounds.Min.X, bounds.Min.Y
	v.w, v.h = bounds.Dx(), bounds.Dy()
}

func (v *View) Layouter() Layouter {
	return v.layouter
}

func (v *View) SetLayouter(layouter Layouter) {
	v.layouter = layouter
}

func (v *View) OnDraw(event *DrawEvent) {
	log.Printf("View.OnDraw()")
}

func (v *View) OnMouseEnter(event *MouseEvent) {
	log.Printf("View.OnMouseEnter()")
}

func (v *View) OnMouseLeave(event *MouseEvent) {
	log.Printf("View.OnMouseLeave()")
}

func (v *View) SetUIMap(ui UIMap) {
	v.uimap = ui
}

func (v *View) UIMap() UIMap {
	return v.uimap
}

func (v *View) MockUp(ui UIMap) {
	return
}
