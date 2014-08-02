package views

import (
	"log"
)

type table_t map[interface{}]interface{}
type msg_func_t func(interface{}, interface{}, string, ...interface{})

var g_observed_table table_t

func AddObserver(observed interface{}, condition string, observer interface{}, msg_func msg_func_t) {
	// Get condition table.
	val := g_observed_table[observed]
	var condition_table table_t
	if table, ok := val.(table_t); ok {
		condition_table = table
	}
	if condition_table == nil {
		condition_table = make(table_t)
		g_observed_table[observed] = condition_table
	}

	// Get observer table.
	val = condition_table[condition]
	var observer_table table_t
	if table, ok := val.(table_t); ok {
		observer_table = table
	}
	if observer_table == nil {
		observer_table = make(table_t)
		condition_table[condition] = observer_table
	}

	// Get observer.
	old_msg_func := observer_table[observer]
	if old_msg_func == &msg_func {
		log.Printf("WARNING: Overwrite the observer with the same msg.")
	} else {
		observer_table[observer] = msg
	}
}
