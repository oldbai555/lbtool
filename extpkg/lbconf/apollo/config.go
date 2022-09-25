package apollo

import (
	"fmt"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/oldbai555/lbtool/extpkg/lbconf/bconf"
	"strings"
)

var _ bconf.DataSource = (*apolloConfig)(nil)

type apolloConfig struct {
	client  agollo.Client
	options *options
}

func NewApolloConfig(opts ...Option) (bconf.DataSource, error) {
	newOpts := NewOptions(opts...)
	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return &config.AppConfig{
			AppID:          newOpts.appid,
			Cluster:        newOpts.cluster,
			NamespaceName:  newOpts.namespace,
			IP:             newOpts.addr,
			IsBackupConfig: newOpts.isBackupConfig,
			Secret:         newOpts.secret,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	conf := &apolloConfig{
		client:  client,
		options: newOpts,
	}
	return conf, nil
}

func (c *apolloConfig) Load() ([]*bconf.Data, error) {
	data := make([]*bconf.Data, 0)
	for _, v := range strings.Split(c.options.namespace, ",") {
		data = append(data, c.loadNameSpace(v)...)
	}
	return data, nil
}

func (c *apolloConfig) loadNameSpace(namespace string) []*bconf.Data {
	var list []*bconf.Data
	c.client.GetConfigCache(namespace).Range(func(key, value interface{}) bool {
		list = append(list, &bconf.Data{Key: fmt.Sprintf("%s", key), Val: value})
		return true
	})
	return list
}

func (c *apolloConfig) Watch() (bconf.DataWatcher, error) {
	return newWatcher(c), nil
}
