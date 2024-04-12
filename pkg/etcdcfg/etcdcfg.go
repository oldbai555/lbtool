package etcdcfg

import (
	"encoding/json"
	"fmt"
	"github.com/oldbai555/lbtool/log"
	"os"
	"runtime"
	"sync"
	"time"
)

// Node etcd 的节点
type Node struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Config struct {
	Node             []Node `json:"node"`
	ConnectTimeoutMs int32  `json:"connect_timeout_ms"`
	OpTimeoutMs      int32  `json:"op_timeout_ms"`
	epList           []string
	Auth             Auth `json:"auth"`
}

func (c *Config) GetEndpointList() []string {
	return c.epList
}
func (c *Config) GetOpTimeout() time.Duration {
	return time.Duration(c.OpTimeoutMs) * time.Millisecond
}
func (c *Config) GetConnectTimeout() time.Duration {
	return time.Duration(c.ConnectTimeoutMs) * time.Millisecond
}

var ConfigPath string
var cfg *Config
var lastLoadTime int64
var mu sync.RWMutex

const loadIntervalSec = 30
const defaultConnectTimeoutMs = 800
const defaultOpTimeoutMs = 5000

func init() {
	if runtime.GOOS == "windows" {
		ConfigPath = `c:/work/etcd.json`
	} else {
		ConfigPath = `/etc/work/etcd.json`
	}
}

func SetConfigPath(path string) {
	ConfigPath = path
}

func GetConfig() *Config {
	now := time.Now().Unix()
	var nl *Config
	mu.RLock()
	if cfg != nil && lastLoadTime+loadIntervalSec >= now {
		nl = cfg
	}
	mu.RUnlock()
	if nl != nil {
		return nl
	}
	mu.Lock()
	defer mu.Unlock()
	// double check, may be update by other goroutine
	if cfg != nil && lastLoadTime+loadIntervalSec >= now {
		return cfg
	}
	defer func() {
		lastLoadTime = now
	}()

	// read file
	dat, err := os.ReadFile(ConfigPath)
	if err != nil {
		log.Errorf("load etcd config file error, path %s, err %s", ConfigPath, err)
		goto OUT
	}
	nl = &Config{}
	err = json.Unmarshal(dat, nl)
	if err != nil {
		log.Errorf("unmarshal etcd config file error, path %s, err %s", ConfigPath, err)
		goto OUT
	}
	if nl.ConnectTimeoutMs <= 0 {
		nl.ConnectTimeoutMs = defaultConnectTimeoutMs
	}
	if nl.OpTimeoutMs <= 0 {
		nl.OpTimeoutMs = defaultOpTimeoutMs
	}
	if len(nl.Node) == 0 {
		log.Errorf("node empty, path %s", ConfigPath)
		goto OUT
	}
	for _, n := range nl.Node {
		if n.Ip == "" || n.Port == 0 {
			log.Errorf("invalid node, node %+v", n)
			goto OUT
		}
		nl.epList = append(nl.epList, fmt.Sprintf("http://%s:%d", n.Ip, n.Port))
	}
	// swap
	cfg = nl
OUT:
	if cfg != nil {
		// return old one
		return cfg
	}
	return &Config{}
}
