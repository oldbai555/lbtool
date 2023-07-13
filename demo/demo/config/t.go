package config

type Mgr struct {
	m map[string]string
}

var M = &Mgr{
	m: make(map[string]string),
}

func (m *Mgr) Get(key string) string {
	return m.m[key]
}

func (m *Mgr) Set(key, value string) {
	m.m[key] = value
}
