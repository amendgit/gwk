package views

import (
	"gwk/sysc"
	"log"
	"syscall"
	"unsafe"
)

const kMsgHaveWork = sysc.WM_USER + 1

type UIEventLoop struct {
	EventLoop
	message_wnd sysc.Handle
}

func NewUIEventLoop() *UIEventLoop {
	u := new(UIEventLoop)
	u.init()
	if u.message_wnd == 0 {
		u.init_message_wnd()
	}
	return u
}

func (u *UIEventLoop) handle_work_msg() {

}

func (u *UIEventLoop) handle_timer_msg() {

}

func ui_event_loop_wnd_proc(hwnd sysc.Handle, msg uint32, warg uintptr, larg uintptr) uintptr {
	switch msg {
	case kMsgHaveWork:
		(*UIEventLoop)(unsafe.Pointer(warg)).handle_work_msg()
	case sysc.WM_TIMER:
		(*UIEventLoop)(unsafe.Pointer(warg)).handle_timer_msg()
	}
	return sysc.DefWindowProc(hwnd, msg, warg, larg)
}

func (u *UIEventLoop) init_message_wnd() {
	const kUIEventLoopClass = "www.ustc.edu.cn/sse/gwk/eventloop"

	var wc sysc.WNDCLASSEX
	wc.Size = uint32(unsafe.Sizeof(wc))
	wc.FnWndProc = syscall.NewCallback(ui_event_loop_wnd_proc)
	wc.HInstance = sysc.NULL
	wc.ClassName = syscall.StringToUTF16Ptr(kUIEventLoopClass)
	sysc.RegisterClassEx(&wc)

	u.message_wnd, _ =
		sysc.CreateWindowEx(0, syscall.StringToUTF16Ptr(kUIEventLoopClass), nil,
			0, 0, 0, 0, 0, 0, 0, 0, 0)
	if u.message_wnd == 0 {
		log.Printf("Create eventloop message HWND failed.")
	}
}

func (u *UIEventLoop) Run() {
	for {
		u.process_next_ui_event()
		if u.should_quit {
			break
		}

		u.do_work()
		if u.should_quit {
			break
		}

		u.do_delayed_work()
		if u.should_quit {
			break
		}
	}

	// message_pump_win.cc 268
	// u.WaitForMoreWork()
}

func (u *UIEventLoop) process_next_ui_event() {
	var msg sysc.MSG
	// for sysc.GetMessage(&msg, sysc.NULL, 0, 0) != 0 {
	for sysc.PeekMessage(&msg, sysc.NULL, 0, 0, sysc.PM_REMOVE) {
		sysc.TranslateMessage(&msg)
		sysc.DispatchMessage(&msg)
	}
}
