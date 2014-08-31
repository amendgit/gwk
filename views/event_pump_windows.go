package views

//
// @see message_pump_win.cc
//

import (
	. "gwk/sysc" // it's fine.
	"log"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"
)

// Message sent to get an additional time slice for pumping (processing) another
// task (a series of such messages creates a continuous task pump).
const kMsgHaveWork = WM_USER + 1

type ui_event_pump_t struct {
	message_wnd       Handle
	delayed_work_time time.Time
	delegate          event_pump_delegate_t
	should_quit       bool
	have_work_        int32
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
	u.should_quit = true
}

func (u *ui_event_pump_t) ScheduleWork() {
	if atomic.SwapInt32(&u.have_work_, 1) == 1 {
		return // someone else continued the pumping.
	}

	// make sure the message pump does some work for us.
	ret := PostMessage(u.message_wnd, kMsgHaveWork, uintptr(unsafe.Pointer(u)), 0)
	if ret {
		return // there was room in the Window Message queue.
	}

	// We have failed to insert a have-work message, so there is a chance that we
	// will starve tasks/timers while sitting in a nested message loop.  Nested
	// loops only look at Windows Message queues, and don't look at *our* task
	// queues, etc., so we might not get a time slice in such. :-(
	// We could abort here, but the fear is that this failure mode is plausibly
	// common (queue is full, of about 2000 messages), so we'll do a near-graceful
	// recovery.  Nested loops are pretty transient (we think), so this will
	// probably be recoverable.
	atomic.SwapInt32(&u.have_work_, 0) // clarify that we didn't really insert.
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
	if delay_msec < USER_TIMER_MINIMUM {
		delay_msec = USER_TIMER_MINIMUM
	}

	// Create a WM_TIMER event that will wake us up to check for any pending
	// timers (in case we are running within a nested, external sub-pump).
	ret := SetTimer(u.message_wnd, uintptr(unsafe.Pointer(u)), uint(delay_msec), 0)

	if ret != 0 {
		return
	}

	// If we can't set timers, we are in big trouble... but cross our fingers for
	// now.
	// log.Printf("SetTimer error")
}

func (u *ui_event_pump_t) handle_work_msg() {
	// Let whatever would have run had we not been putting messages in the queue
	// run now.  This is an attempt to make our dummy message not starve other
	// messages that may be in the Windows message queue.
	u.process_pump_replacement_msg()

	// Now give the delegate a chance to do some work.  He'll let us know if he
	// needs to do more work.
	if u.delegate.DoWork() {
		u.ScheduleWork()
	}
}

func (u *ui_event_pump_t) handle_timer_msg() {

}

func ui_event_loop_wnd_proc(hwnd Handle, msg uint32, warg uintptr, larg uintptr) uintptr {
	switch msg {
	case kMsgHaveWork:
		(*ui_event_pump_t)(unsafe.Pointer(warg)).handle_work_msg()
	case WM_TIMER:
		(*ui_event_pump_t)(unsafe.Pointer(warg)).handle_timer_msg()
	}
	return DefWindowProc(hwnd, msg, warg, larg)
}

func (u *ui_event_pump_t) init_message_wnd() {
	const kEventPumpClass = "www.ustc.edu.cn/sse/gwk/eventloop"

	var wc WNDCLASSEX
	wc.Size = uint32(unsafe.Sizeof(wc))
	wc.FnWndProc = syscall.NewCallback(ui_event_loop_wnd_proc)
	wc.HInstance = NULL
	wc.ClassName = syscall.StringToUTF16Ptr(kEventPumpClass)
	RegisterClassEx(&wc)

	u.message_wnd, _ =
		CreateWindowEx(0, syscall.StringToUTF16Ptr(kEventPumpClass), nil,
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
		delay = int(^uint(0) >> 1) // INT_MAX
	}

	result := MsgWaitForMultipleObjectsEx(0, nil, int32(delay),
		QS_ALLINPUT, MWMO_INPUTAVAILABLE)

	if result == WAIT_OBJECT_0 {
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
		var msg MSG
		queue_status := GetQueueStatus(QS_MOUSE)

		if (HIWORD(queue_status)&QS_MOUSE) != 0 &&
			!PeekMessage(&msg, NULL, WM_MOUSEFIRST, WM_MOUSELAST, PM_NOREMOVE) {
			WaitMessage()
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
	queue_status := GetQueueStatus(QS_SENDMESSAGE)
	if (HIWORD(queue_status) & QS_SENDMESSAGE) != 0 {
		sent_messages_in_queue = true
	}

	var msg MSG
	if PeekMessage(&msg, NULL, 0, 0, PM_REMOVE) {
		return u.process_msg(&msg)
	}

	return sent_messages_in_queue
}

func (u *ui_event_pump_t) process_msg(msg *MSG) bool {
	// While running our main message pump, we discard kMsgHaveWork messages.
	if msg.HWnd == u.message_wnd && msg.Message == kMsgHaveWork {
		return u.process_pump_replacement_msg()
	}

	TranslateMessage(msg)
	DispatchMessage(msg)

	return true
}

func (u *ui_event_pump_t) process_pump_replacement_msg() bool {
	var have_msg bool
	var msg MSG
	have_msg = PeekMessage(&msg, NULL, 0, 0, PM_REMOVE)

	// Since we discarded a kMsgHaveWork message, we must update the flag.
	old_have_work := atomic.SwapInt32(&u.have_work_, 0)
	if old_have_work == 0 {
		log.Printf("ERROR: we don't have work to do.")
	}

	// We don't need a special time slice if we didn't have_msg to process.
	if !have_msg {
		return false
	}

	// Guarantee we'll get another time slice in the case where we go into native
	// windows code.   This ScheduleWork() may hurt performance a tiny bit when
	// tasks appear very infrequently, but when the event queue is busy, the
	// kMsgHaveWork events get (percentage wise) rarer and rarer.
	u.ScheduleWork()
	return u.process_msg(&msg)
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
