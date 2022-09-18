package lbconfig

import (
	"context"
	"errors"
	"github.com/oldbai555/lb/extrpkg/lbconfig/hconf"
	"github.com/spf13/viper"
)

type config struct {
	opts    *options
	watcher hconf.DataWatcher
	viper   *viper.Viper
}
type HConfig interface {
	Load() error
	Get(key string) (HVal, error)
	Watch(event WatchEvent) error
	Close() error
}

type WatchEvent func(path string, v HVal)

func NewHConfig(opts ...Option) (HConfig, error) {
	options, err := newOptions(opts...)
	if err != nil {
		return nil, err
	}
	return &config{
		opts:  options,
		viper: viper.New(),
	}, nil
}

func (c *config) Load() error {
	kvs, err := c.opts.dataSource.Load()
	if err != nil {
		return nil
	}
	for _, v := range kvs {
		c.viper.Set(v.Key, v.Val)
	}
	return nil
}

func (c *config) Get(key string) (HVal, error) {
	return c.viper.Get(key), nil
}

func (c *config) Watch(event WatchEvent) error {
	var err error
	if c.watcher, err = c.opts.dataSource.Watch(); err != nil {
		return err
	}
	go c.watch(event)
	return nil
}

func (c *config) watch(event WatchEvent) {
	for {
		kvs, err := c.watcher.Change()
		if errors.Is(err, context.Canceled) {
			return
		}
		if err != nil {
			continue
		}
		for _, v := range kvs {
			c.viper.Set(v.Key, v.Val)
			event(v.Key, v.Val)
		}
	}
}

func (c *config) Close() error {
	c.viper = viper.New()
	if c.watcher != nil {
		return c.watcher.Close()
	}
	return nil
}
