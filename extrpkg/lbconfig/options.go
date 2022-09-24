package lbconfig

import (
	"errors"
	"github.com/oldbai555/lb/extrpkg/lbconfig/domain"
)

type Option func(opt *options)

type options struct {
	dataSource domain.DataSource
	useLocal   bool
}

func WithDataSource(d domain.DataSource) Option {
	return func(opt *options) {
		opt.dataSource = d
	}
}

func WithUseLocal() Option {
	return func(opt *options) {
		opt.useLocal = true
	}
}

func newOptions(opts ...Option) (*options, error) {
	o := &options{
		dataSource: nil,
		useLocal:   false,
	}
	for _, opt := range opts {
		opt(o)
	}
	if o.dataSource == nil {
		return nil, errors.New("dataSource is nil")
	}
	return o, nil
}
