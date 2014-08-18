package views

type UIMap map[string]interface{}

func (u UIMap) Bool(key string) (bool, bool) {
	val := u[key]
	if rv, ok := val.(bool); ok {
		return rv, true
	}
	return false, false
}

func (u UIMap) String(key string) (string, bool) {
	val := u[key]
	if rv, ok := val.(string); ok {
		return rv, true
	}
	return "", false
}

func (u UIMap) Int(key string) (int, bool) {
	val := u[key]
	if rv, ok := val.(int); ok {
		return rv, true
	}
	return 0, false
}

func (u UIMap) UIMap(key string) (UIMap, bool) {
	val := u[key]
	if rv, ok := val.(UIMap); ok {
		return rv, true
	}
	return nil, false
}

func (u UIMap) UIMaps(key string) ([]UIMap, bool) {
	val := u[key]
	if rv, ok := val.([]UIMap); ok {
		return rv, true
	}
	return nil, false
}

func (u UIMap) View(key string) (View, bool) {
	val := u[key]
	if rv, ok := val.(View); ok {
		return rv, true
	}
	return nil, false
}
