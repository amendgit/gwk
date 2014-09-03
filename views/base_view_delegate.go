package views

type BaseViewDelegate struct {
	on_mouse_enter func(*MouseEvent)
	on_mouse_leave func(*MouseEvent)
	on_draw        func(*DrawEvent)
}

func NewBaseViewDelegate() *BaseViewDelegate {
	return new(BaseViewDelegate)
}

func (d *BaseViewDelegate) InitWithUIMap(delegate UIMap) *BaseViewDelegate {
	void_ptr := delegate["on_mouse_enter"]
	if void_ptr != nil {
		if on_mouse_enter, ok := void_ptr.(func(*MouseEvent)); ok {
			d.on_mouse_enter = on_mouse_enter
		}
	}

	void_ptr = delegate["on_mouse_leave"]
	if void_ptr != nil {
		if on_mouse_leave, ok := void_ptr.(func(*MouseEvent)); ok {
			d.on_mouse_leave = on_mouse_leave
		}
	}

	void_ptr = delegate["on_draw"]
	if void_ptr != nil {
		if on_draw, ok := void_ptr.(func(*DrawEvent)); ok {
			d.on_draw = on_draw
		}
	}

	return d
}

func (d *BaseViewDelegate) OnMouseEnter(event *MouseEvent) {
	if d.on_mouse_enter == nil {
		return
	}
	d.on_mouse_enter(event)
}

func (d *BaseViewDelegate) OnMouseLeave(event *MouseEvent) {
	if d.on_mouse_leave == nil {
		return
	}
	d.on_mouse_leave(event)
}

func (d *BaseViewDelegate) OnDraw(event *DrawEvent) {
	if d.on_draw == nil {
		return
	}
	d.on_draw(event)
}
