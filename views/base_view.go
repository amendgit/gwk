package views

import (
	. "gwk/vango"
	. "image"
	"log"
)

type BaseView struct {
	id string

	canvas *Canvas

	x, y int // Relative to Parent.
	w, h int

	children []View
	parent   View

	uimap    UIMap
	layouter Layouter

	delegate ViewDelegate
}

func NewBaseView() *BaseView {
	var v = new(BaseView)
	return v
}

func (v *BaseView) ID() string {
	return v.id
}

func (v *BaseView) SetID(id string) {
	v.id = id
}

func (v *BaseView) AddChild(child View) {
	v.children = append(v.children, child)
	child.SetParent(v)
}

func (v *BaseView) Children() []View {
	return v.children
}

func (v *BaseView) Parent() View {
	return v.parent
}

func (v *BaseView) SetParent(parent View) {
	v.parent = parent
}

func (v *BaseView) Canvas() *Canvas {
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

func (v *BaseView) SetCanvas(canvas *Canvas) {
	v.canvas = canvas
}

func (v *BaseView) ToAbsPoint(pt Point) Point {
	if v.Parent() == nil {
		return pt
	}
	pt.X = pt.X + v.X()
	pt.Y = pt.Y + v.Y()
	return v.Parent().ToAbsPoint(pt)
}

func (v *BaseView) ToDevicePoint(pt Point) Point {
	pt = pt.Add(Pt(v.X(), v.Y()))
	return v.Parent().ToDevicePoint(pt)
}

func (v *BaseView) ToAbsRect(rect Rectangle) Rectangle {
	if v.Parent() == nil {
		return rect
	}
	rect.Min.X = rect.Min.X + v.X()
	rect.Min.Y = rect.Min.Y + v.Y()
	rect.Max.X = rect.Max.X + v.X()
	rect.Max.Y = rect.Max.Y + v.Y()
	return v.Parent().ToAbsRect(rect)
}

func (v *BaseView) ToDeviceRect(rect Rectangle) Rectangle {
	rect = rect.Add(Pt(v.X(), v.Y()))
	return v.Parent().ToDeviceRect(rect)
}

func update_rect(v View, rect Rectangle) {
	if v.Parent() == nil {
		v.ScheduleDrawRect(rect)
		return
	}
	rect = rect.Add(Pt(v.X(), v.Y()))
	update_rect(v.Parent(), rect)
}

func (v *BaseView) ScheduleDraw() {
	if v.Parent() != nil {
		update_rect(v.Parent(), v.Bounds())
	}
}

func (v *BaseView) ScheduleDrawRect(rect Rectangle) {
	if v.Parent() != nil {
		update_rect(v.Parent(), rect)
	}
}

func (v *BaseView) X() int {
	return v.x
}

func (v *BaseView) Y() int {
	return v.y
}

func (v *BaseView) W() int {
	return v.w
}

func (v *BaseView) H() int {
	return v.h
}

func (v *BaseView) XYWH() (x, y, w, h int) {
	x, y, w, h = v.x, v.y, v.w, v.h
	return
}

func (v *BaseView) SetXYWH(x, y, w, h int) {
	v.x, v.y, v.w, v.h = x, y, w, h
}

func (v *BaseView) Width() int {
	return v.W()
}

func (v *BaseView) SetWidth(width int) {
	v.w = width
}

func (v *BaseView) Height() int {
	return v.H()
}

func (v *BaseView) SetHeight(height int) {
	v.h = height
}

func (v *BaseView) Left() int {
	return v.X()
}

func (v *BaseView) SetLeft(left int) {
	v.x = left
}

func (v *BaseView) Top() int {
	return v.Y()
}

func (v *BaseView) SetTop(top int) {
	v.y = top
}

func (v *BaseView) Bounds() Rectangle {
	return Rect(v.x, v.y, v.x+v.w, v.y+v.h)
}

func (v *BaseView) LocalBounds() Rectangle {
	return Rect(0, 0, v.w, v.h)
}

func (v *BaseView) SetBounds(bounds Rectangle) {
	v.x, v.y = bounds.Min.X, bounds.Min.Y
	v.w, v.h = bounds.Dx(), bounds.Dy()
}

func (v *BaseView) Layouter() Layouter {
	return v.layouter
}

func (v *BaseView) SetLayouter(layouter Layouter) {
	v.layouter = layouter
}

func (v *BaseView) OnDraw(event *DrawEvent) {
	log.Printf("BaseView.OnDraw()")
}

func (v *BaseView) OnMouseEnter(event *MouseEvent) {
	if v.delegate == nil {
		return
	}
	v.delegate.OnMouseEnter(event)
}

func (v *BaseView) OnMouseLeave(event *MouseEvent) {
	log.Printf("BaseView.OnMouseLeave()")
}

func (v *BaseView) SetUIMap(ui UIMap) {
	v.uimap = ui
}

func (v *BaseView) UIMap() UIMap {
	return v.uimap
}

func (v *BaseView) MockUp(ui UIMap) {
	return
}

func (v *BaseView) SetDelegate(delegate ViewDelegate) {
	v.delegate = delegate
}

func (v *BaseView) Delegate() ViewDelegate {
	return v.delegate
}
