package views

type UIBinder struct {
	object0 interface{}
	key0    string
	object1 interface{}
	key1    string
}

func UIBind(object0 interface{}, key0 string, object1 interface{}, key1 string) {
	binder := &UIBinder{
		object0: object0,
		key0:    key0,
		object1: object1,
		key1:    key1,
	}

	AddObserver(object0, key0, binder,
		func(subject subject_t, condition string, observer observer_t, args ...interface{}) {
			var binder *UIBinder
			var ok bool
			if binder, ok = observer.(*UIBinder); !ok {
				return
			}
			key0_changed(binder)
		})

	AddObserver(object1, key1, binder,
		func(subject subject_t, condition string, observer observer_t, args ...interface{}) {
			var binder *UIBinder
			var ok bool
			if binder, ok = observer.(*UIBinder); !ok {
				return
			}
			key1_changed(binder)
		})
}

func key0_changed(binder *UIBinder) {
	// table1.setValueForKeyPath(binder.key1, value0)
}

func key1_changed(binder *UIBinder) {
	// table0.setValueForKeyPath(binder.key0, value0)
}
