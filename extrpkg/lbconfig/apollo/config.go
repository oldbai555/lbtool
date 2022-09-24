package apollo

import (
	"fmt"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/oldbai555/lb/extrpkg/lbconfig/domain"
	"strings"
)

var _ domain.DataSource = (*apolloConfig)(nil)

type apolloConfig struct {
	client  agollo.Client
	options *options
}

func NewApolloConfig(opts ...Option) (domain.DataSource, error) {
	options := NewOptions(opts...)
	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return &config.AppConfig{
			AppID:          options.appid,
			Cluster:        options.cluster,
			NamespaceName:  options.namespace,
			IP:             options.addr,
			IsBackupConfig: options.isBackupConfig,
			Secret:         options.secret,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	conf := &apolloConfig{
		client:  client,
		options: options,
	}
	return conf, nil
}

func (c *apolloConfig) Load() ([]*domain.Data, error) {
	data := make([]*domain.Data, 0)
	for _, v := range strings.Split(c.options.namespace, ",") {
		data = append(data, c.loadNameSpace(v)...)
	}
	return data, nil
}

func (c *apolloConfig) loadNameSpace(namespace string) []*domain.Data {
	var list []*domain.Data
	c.client.GetConfigCache(namespace).Range(func(key, value interface{}) bool {
		list = append(list, &domain.Data{Key: fmt.Sprintf("%s", key), Val: value})
		return true
	})
	return list
}

func (c *apolloConfig) Watch() (domain.DataWatcher, error) {
	return newWatcher(c), nil
}
