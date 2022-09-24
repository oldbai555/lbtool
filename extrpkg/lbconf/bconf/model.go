package bconf

// DataSource 数据源
type DataSource interface {
	Load() ([]*Data, error)
	Watch() (DataWatcher, error)
}

// DataWatcher 数据监听器
type DataWatcher interface {
	Change() ([]*Data, error)
	Close() error
}

// Data 数据
type Data struct {
	Key string
	Val interface{}
}

// LbConfig 配置接入方应当提供的接口
type LbConfig interface {
	Load() error
	Get(key string) (Val, error)
	Watch(event WatchEvent) error
	Close() error
}

// Val 值
type Val interface{}

// WatchEvent 监听事件
type WatchEvent func(path string, v Val)
