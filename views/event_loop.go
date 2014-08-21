package views

import (
	"time"
)

// ============================================================================
// Task

type Closure func()

type Task struct {
	closure Closure
	time    time.Time
}

func NewTask(closure Closure, time time.Time) *Task {
	task := new(Task)
	task.closure = closure
	task.time = time
	return task
}

// ============================================================================
// Task Queue
type task_node_t struct {
	data *Task
	prev *task_node_t
	next *task_node_t
}

type task_queue_t struct {
	head  task_node_t
	count int
}

func new_task_queue() *task_queue_t {
	t := new(task_queue_t)
	t.head.next = &t.head
	t.head.prev = &t.head
	return t
}

func (t *task_queue_t) Empty() bool {
	return t.head.prev == &t.head && t.head.next == &t.head
}

// Insert a new task node at the rear of the queue.
func (t *task_queue_t) Push(task *Task) {
	task_node := new(task_node_t)
	task_node.data = task

	task_node.prev = t.head.prev
	task_node.next = &t.head

	t.head.prev.next = task_node
	t.head.prev = task_node

	t.count++
}

func (t *task_queue_t) Pop() {
	if t.Empty() {
		return
	}

	t.head.next = t.head.next.next
	t.head.next.prev = &t.head

	t.count--
}

func (t *task_queue_t) Front() *Task {
	if t.Empty() {
		return nil
	}
	return t.head.next.data
}

func (t *task_queue_t) Count() int {
	return t.count
}

func (t *task_queue_t) debug_check() bool {
	var node *task_node_t
	var count = 0

	for node = t.head.prev; node != t.head.next; node = node.next {
		count++
	}

	if count == t.count {
		return true
	}
	return false
}

// ============================================================================
// Priority Task Queue
type priority_task_queue_t struct {
	data  []*Task
	count int
}

func new_priority_task_queue() *priority_task_queue_t {
	priority_task_queue := new(priority_task_queue_t)
	priority_task_queue.data = make([]*Task, 1, 100)
	return priority_task_queue
}

func (t *priority_task_queue_t) Empty() bool {
	return t.data == nil || len(t.data) <= 1
}

func (t *priority_task_queue_t) Top() *Task {
	if t.Empty() {
		return nil
	}
	return t.data[1]
}

func (t *priority_task_queue_t) Push(task *Task) {
	t.data = append(t.data, task)

	i0 := len(t.data) - 1
	i1 := int(i0 / 2)

	// The node in the tree is begin at 1.
	for i0 != 1 && t.data[i0].time.Before(t.data[i1].time) {
		t.data[i0], t.data[i1] = t.data[i1], t.data[i0]
		i0 = i1
		i1 = i0 / 2
	}

	t.count++
}

func (t *priority_task_queue_t) Pop() {
	// swap the root and the last element.
	n := len(t.data) - 1
	t.data[1] = t.data[n]

	// remove the last element.
	t.data = t.data[0:n] // [0, n)
	n = n - 1

	// p: parent, l: left child, r: right child.
	p := 1
	for {
		l := p * 2

		if l >= n {
			break
		}

		i := l
		r := l + 1
		if r < n && t.data[r].time.Before(t.data[l].time) {
			i = r
		}

		if t.data[i].time.Before(t.data[p].time) {
			t.data[p], t.data[l] = t.data[l], t.data[p]
			p = i
		} else {
			break
		}
	}

	t.count--
}

func (t *priority_task_queue_t) Count() int {
	return t.count
}

func (t *priority_task_queue_t) debug_check() bool {
	if len(t.data) == t.count {
		return true
	}
	return false
}

// ============================================================================
// Event loop

type EventLoop struct {
	should_quit        bool
	pending_task_queue *task_queue_t
	delayed_task_queue *priority_task_queue_t
	pump               event_pump_t
	recent_time        time.Time
	// pending_task_queue_lock sync.Lock
}

func (e *EventLoop) init() {
	e.pending_task_queue = new_task_queue()
	e.delayed_task_queue = new_priority_task_queue()
}

func (e *EventLoop) add_to_pending_task_queue(closure Closure, delay_misc int64) {
	time := time.Now().Add(time.Duration(delay_misc * 1000 * 1000))
	task := NewTask(closure, time)
	// Lock pending task queue.
	e.pending_task_queue.Push(task)
	// Unlock pending task queue.
}

func (e *EventLoop) PostTask(closure Closure) {
	e.add_to_pending_task_queue(closure, 0)
}

func (e *EventLoop) PostDelayedTask(closure Closure, delay_misc int64) {
	e.add_to_pending_task_queue(closure, delay_misc)
}

func (e *EventLoop) ShouldQuit() {
	e.should_quit = true
}

//
// EventPumpDelegate
//
func (e *EventLoop) DoWork() bool {
	now := time.Now()

	for !e.pending_task_queue.Empty() {
		// lock pending task queue.
		task := e.pending_task_queue.Front()
		e.pending_task_queue.Pop()
		// unlock pending task queue.

		if !task.time.After(now) {
			// TODO(nest): in case of the nest work. message_loop.cc 617
			task.closure()
			return true
		} else {
			e.delayed_task_queue.Push(task)
			// If we changed the topmost task, then it is time to reschedule.
			if e.delayed_task_queue.Top() == task {
				e.pump.ScheduleDelayedWork(e.delayed_task_queue.Top().time)
			}
		}
	}

	// nothing happened.
	return false
}

func (e *EventLoop) DoDelayedWork(next_delayed_work_time *time.Time) bool {
	if e.delayed_task_queue.Empty() {
		e.recent_time = time.Now()
		*next_delayed_work_time = e.recent_time
		return false
	}

	// When we "fall behind," there will be a lot of tasks in the delayed work
	// queue that are ready to run.  To increase efficiency when we fall behind,
	// we will only call Time::Now() intermittently, and then process all tasks
	// that are ready to run before calling it again.  As a result, the more we
	// fall behind (and have a lot of ready-to-run delayed tasks), the more
	// efficient we'll be at handling the tasks.
	next_run_time := e.delayed_task_queue.Top().time
	if next_run_time.After(e.recent_time) {
		// get a better view of Now();
		e.recent_time = time.Now()
		if next_run_time.After(e.recent_time) {
			*next_delayed_work_time = e.recent_time
			return false
		}
	}

	task := e.delayed_task_queue.Top()
	e.delayed_task_queue.Pop()

	if !e.delayed_task_queue.Empty() {
		*next_delayed_work_time = e.delayed_task_queue.Top().time
	}

	// DeferOrRunPendingTask()
	task.closure()

	return false
}

// ============================================================================

type UIEventLoop struct {
	EventLoop
}

func NewUIEventLoop() *UIEventLoop {
	u := new(UIEventLoop)
	u.EventLoop.init()
	u.pump = new_ui_event_pump(u)
	return u
}

func (u *UIEventLoop) Run() {
	u.pump.Run()
}

// ============================================================================
// Event loop mixup with ui events and others.

var g_current_ui_event_loop *UIEventLoop

func MainUIEventLoop() *UIEventLoop {
	if g_current_ui_event_loop == nil {
		g_current_ui_event_loop = NewUIEventLoop()
	}
	return g_current_ui_event_loop
}
