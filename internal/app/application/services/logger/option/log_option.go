package option

type LogOption func(map[string]any)

func Any(key string, value any) LogOption {
	return func(m map[string]any) {
		m[key] = value
	}
}

func Error(err error) LogOption {
	return func(m map[string]any) {
		m["error"] = err.Error()
	}
}
