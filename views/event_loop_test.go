package views

import (
	"testing"
)

func TestTaskQueue(t *testing.T) {
	q := new_task_queue()
	q.Push(&Task{})
	q.Push(&Task{})
	q.Push(&Task{})
	q.Pop()
	q.Pop()
	q.Pop()
	if !q.Empty() {
		t.Fatalf("task queue push/pop failed")
	}
}

func TestPriorityQueue(t *testing.T) {
	pq := new_priority_task_queue()
	pq.Push(&Task{})
	pq.Push(&Task{})
	pq.Push(&Task{})
	pq.Pop()
	pq.Pop()
	pq.Pop()
	if pq.Count() != 0 {
		t.Fatalf("priority task queue push/pop failed")
	}
}
