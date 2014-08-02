package views

import (
	. "gwk/vango"
)

type Layer struct {
	canvas        Canvas
	delegate_view Viewer
}

func (l *Layer) ScheduleDraw() {

}
