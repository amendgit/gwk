package views

import (
	"log"
)

type table_t map[interface{}]interface{}

type observer_t interface{}
type subject_t interface{}
type msg_func_t func(subject_t, string, observer_t, ...interface{})

// Three table definition.
//
// from observer to msg_func.
type observer_table_t map[observer_t]msg_func_t

// from condition to observer_table
type condition_table_t map[string]observer_table_t

// from subject to condition table.
type subject_table_t map[subject_t]condition_table_t

var g_subject_table subject_table_t

func AddObserver(subject subject_t, condition string, observer observer_t, msg_func msg_func_t) {
	// Get condition table.
	condition_table := g_subject_table[subject]
	if condition_table == nil {
		condition_table = make(condition_table_t)
		g_subject_table[subject] = condition_table
	}

	// Get observer table.
	observer_table := condition_table[condition]
	if observer_table == nil {
		observer_table = make(observer_table_t)
		condition_table[condition] = observer_table
	}

	// Get observer.
	old_msg_func := observer_table[observer]
	if &old_msg_func == &msg_func {
		log.Printf("WARNING: Overwrite the observer with the same msg.")
	} else {
		observer_table[observer] = msg_func
	}
}

func RemoveObserver(subject subject_t, condition string, observer observer_t) {
	// Get condition table.
	condition_table := g_subject_table[subject]
	if condition_table == nil {
		condition_table = make(condition_table_t)
		g_subject_table[subject] = condition_table
	}

	// Get observer table.
	observer_table := condition_table[condition]
	if observer_table == nil {
		observer_table = make(observer_table_t)
		condition_table[condition] = observer_table
	}

	delete(observer_table, observer)
}

func NotifyObserver(subject subject_t, condition string, observer observer_t, args ...interface{}) {
	// Get condition table.
	condition_table := g_subject_table[subject]
	if condition_table == nil {
		condition_table = make(condition_table_t)
		g_subject_table[subject] = condition_table
	}

	// Get observer table.
	observer_table := condition_table[condition]
	if observer_table == nil {
		observer_table = make(observer_table_t)
		condition_table[condition] = observer_table
	}

	msg_func := observer_table[observer]
	if msg_func == nil {
		log.Print("WARNING: Notify on an empty msg function.")
		return
	} else {
		msg_func(subject, condition, observer, args...)
	}
}
