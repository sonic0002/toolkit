package mocker

import "reflect"

// RestoreMock restores the mocked functions or variables
func RestoreMock(args ...interface{}) func() {
	origArgs := map[interface{}]reflect.Value{}

	for _, name := range args {
		value := reflect.ValueOf(name)
		if value.Kind() != reflect.Ptr {
			panic("unsupported value")
		}

		pv := reflect.New(value.Type().Elem())
		pv.Elem().Set(reflect.Indirect(value))
		origArgs[name] = pv
	}

	return func() {
		for key, value := range origArgs {
			reflect.ValueOf(key).Elem().Set(value.Elem())
		}
	}
}
