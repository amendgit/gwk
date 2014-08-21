package views

//
// @see message_pump_win.cc
//

import (
	"gwk/sysc"
	"log"
	"syscall"
	"time"
	"unsafe"
)

const kMsgHaveWork = sysc.WM_USER + 1

type ui_event_pump_t struct {
	message_wnd       sysc.Handle
	delayed_work_time time.Time
	delegate          event_pump_delegate_t
	should_quit       bool
}

func new_ui_event_pump(delegate event_pump_delegate_t) *ui_event_pump_t {
	u := new(ui_event_pump_t)
	if u.message_wnd == 0 {
		u.init_message_wnd()
	}

	u.delegate = delegate
	u.should_quit = false

	return u
}

//
// event_pump_t methods
//
func (u *ui_event_pump_t) Run() {
	for {
		more_work_is_plausible := u.process_next_ui_event()
		if u.should_quit {
			break
		}

		more_work := u.delegate.DoWork()
		more_work_is_plausible = more_work_is_plausible || more_work
		if u.should_quit {
			break
		}

		var next_delayed_work_time time.Time
		more_delayed_work := u.delegate.DoDelayedWork(&next_delayed_work_time)
		more_work_is_plausible = more_work_is_plausible || more_delayed_work
		// If we did not process any delayed work, then we can assume that our
		// existing WM_TIMER if any will fire when delayed work should run.  We
		// don't want to disturb that timer if it is already in flight.  However,
		// if we did do all remaining delayed work, then lets kill the WM_TIMER.
		// if (more_work_is_plausible && delayed_work_time_.is_null())
		//   KillTimer(message_hwnd_, reinterpret_cast<UINT_PTR>(this));
		if u.should_quit {
			break
		}

		if more_work_is_plausible {
			continue
		}

		// more_work_is_plausible = state_->delegate->DoIdleWork();
		//     if (state_->should_quit)
		// 	       break;

		// Wait (sleep) until we have work to do again.
		u.wait_for_work()
	}
}

func (u *ui_event_pump_t) Quit() {

}

func (u *ui_event_pump_t) ScheduleWork() {

}

func (u *ui_event_pump_t) ScheduleDelayedWork(delayed_work_time time.Time) {
	//
	// We would *like* to provide high resolution timers.  Windows timers using
	// SetTimer() have a 10ms granularity.  We have to use WM_TIMER as a wakeup
	// mechanism because the application can enter modal windows loops where it
	// is not running our MessageLoop; the only way to have our timers fire in
	// these cases is to post messages there.
	//
	// To provide sub-10ms timers, we process timers directly from our run loop.
	// For the common case, timers will be processed there as the run loop does
	// its normal work.  However, we *also* set the system timer so that WM_TIMER
	// events fire.  This mops up the case of timers not being able to work in
	// modal message loops.  It is possible for the SetTimer to pop and have no
	// pending timers, because they could have already been processed by the
	// run loop itself.
	//
	// We use a single SetTimer corresponding to the timer that will expire
	// soonest.  As new timers are created and destroyed, we update SetTimer.
	// Getting a spurrious SetTimer event firing is benign, as we'll just be
	// processing an empty timer queue.
	//
	u.delayed_work_time = delayed_work_time
	delay_msec := u.get_current_delay()
	if delay_msec < sysc.USER_TIMER_MINIMUM {
		delay_msec = sysc.USER_TIMER_MINIMUM
	}

	// Create a WM_TIMER event that will wake us up to check for any pending
	// timers (in case we are running within a nested, external sub-pump).
	ret := sysc.SetTimer(u.message_wnd, uintptr(unsafe.Pointer(u)), uint(delay_msec), 0)

	if ret != 0 {
		return
	}

	// If we can't set timers, we are in big trouble... but cross our fingers for
	// now.
	// log.Printf("SetTimer error")
}

func (u *ui_event_pump_t) handle_work_msg() {

}

func (u *ui_event_pump_t) handle_timer_msg() {

}

func ui_event_loop_wnd_proc(hwnd sysc.Handle, msg uint32, warg uintptr, larg uintptr) uintptr {
	switch msg {
	case kMsgHaveWork:
		(*ui_event_pump_t)(unsafe.Pointer(warg)).handle_work_msg()
	case sysc.WM_TIMER:
		(*ui_event_pump_t)(unsafe.Pointer(warg)).handle_timer_msg()
	}
	return sysc.DefWindowProc(hwnd, msg, warg, larg)
}

func (u *ui_event_pump_t) init_message_wnd() {
	const kEventPumpClass = "www.ustc.edu.cn/sse/gwk/eventloop"

	var wc sysc.WNDCLASSEX
	wc.Size = uint32(unsafe.Sizeof(wc))
	wc.FnWndProc = syscall.NewCallback(ui_event_loop_wnd_proc)
	wc.HInstance = sysc.NULL
	wc.ClassName = syscall.StringToUTF16Ptr(kEventPumpClass)
	sysc.RegisterClassEx(&wc)

	u.message_wnd, _ =
		sysc.CreateWindowEx(0, syscall.StringToUTF16Ptr(kEventPumpClass), nil,
			0, 0, 0, 0, 0, 0, 0, 0, 0)
	if u.message_wnd == 0 {
		log.Printf("ERROR: create eventloop message_wnd failed.")
	}
}

func (u *ui_event_pump_t) get_current_delay() int {
	if u.delayed_work_time.Nanosecond() <= 0 {
		return -1
	}

	timeout := u.delayed_work_time.Sub(time.Now())

	// timeout is nanosecond to millisecond
	var delay int = int(int64(timeout) / (1000 * 1000))
	if delay < 0 {
		delay = 0
	}

	return delay
}

func (u *ui_event_pump_t) wait_for_work() {
	delay := u.get_current_delay()
	if delay < 0 {
		delay = int(^uint(0) >> 1)
	}

	result := sysc.MsgWaitForMultipleObjectsEx(0, nil, int32(delay),
		sysc.QS_ALLINPUT, sysc.MWMO_INPUTAVAILABLE)

	if result == sysc.WAIT_OBJECT_0 {
		// A WM_* message is available.
		// If a parent child relationship exists between windows across threads
		// then their thread inputs are implicitly attached.
		// This causes the MsgWaitForMultipleObjectsEx API to return indicating
		// that messages are ready for processing (Specifically, mouse messages
		// intended for the child window may appear if the child window has
		// capture).
		// The subsequent PeekMessages call may fail to return any messages thus
		// causing us to enter a tight loop at times.
		// The WaitMessage call below is a workaround to give the child window
		// some time to process its input messages.
		var msg sysc.MSG
		queue_status := sysc.GetQueueStatus(sysc.QS_MOUSE)
		if (sysc.HIWORD(queue_status)&sysc.QS_MOUSE) != 0 &&
			!sysc.PeekMessage(&msg, sysc.NULL, sysc.WM_MOUSEFIRST, sysc.WM_MOUSELAST, sysc.PM_NOREMOVE) {
			sysc.WaitMessage()
		}
		return
	}

}

func (u *ui_event_pump_t) process_next_ui_event() bool {
	// If there are sent messages in the queue then PeekMessage internally
	// dispatches the message and returns false. We return true in this
	// case to ensure that the message loop peeks again instead of calling
	// MsgWaitForMultipleObjectsEx again.
	var sent_messages_in_queue = false
	queue_status := sysc.GetQueueStatus(sysc.QS_SENDMESSAGE)
	if (sysc.HIWORD(queue_status) & sysc.QS_SENDMESSAGE) != 0 {
		sent_messages_in_queue = true
	}

	var msg sysc.MSG
	if sysc.PeekMessage(&msg, sysc.NULL, 0, 0, sysc.PM_REMOVE) {
		return u.process_message(&msg)
	}

	return sent_messages_in_queue
}

func (u *ui_event_pump_t) process_message(msg *sysc.MSG) bool {
	sysc.TranslateMessage(msg)
	sysc.DispatchMessage(msg)
	return true
}

// bool MessagePumpForUI::ProcessMessageHelper(const MSG& msg) {
//   TRACE_EVENT1("base", "MessagePumpForUI::ProcessMessageHelper",
//                "message", msg.message);
//   if (WM_QUIT == msg.message) {
//     // Repost the QUIT message so that it will be retrieved by the primary
//     // GetMessage() loop.
//     state_->should_quit = true;
//     PostQuitMessage(static_cast<int>(msg.wParam));
//     return false;
//   }

//   // While running our main message pump, we discard kMsgHaveWork messages.
//   if (msg.message == kMsgHaveWork && msg.hwnd == message_hwnd_)
//     return ProcessPumpReplacementMessage();

//   if (CallMsgFilter(const_cast<MSG*>(&msg), kMessageFilterCode))
//     return true;

//   WillProcessMessage(msg);

//   if (!message_filter_->ProcessMessage(msg)) {
//     if (state_->dispatcher) {
//       if (!state_->dispatcher->Dispatch(msg))
//         state_->should_quit = true;
//     } else {
//       TranslateMessage(&msg);
//       DispatchMessage(&msg);
//     }
//   }

//   DidProcessMessage(msg);
//   return true;
// }
