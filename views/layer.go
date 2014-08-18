package views

import (
	. "gwk/vango"
)

type Layer struct {
	canvas        Canvas
	delegate_view View
}

func (l *Layer) ScheduleDraw() {

}
