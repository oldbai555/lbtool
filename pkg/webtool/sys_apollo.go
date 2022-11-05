package webtool

import (
	"encoding/json"
	"github.com/oldbai555/lbtool/extpkg/lbconf"
	"github.com/oldbai555/lbtool/extpkg/lbconf/apollo"
	"github.com/oldbai555/lbtool/extpkg/lbconf/bconf"
	"github.com/oldbai555/lbtool/log"
)

type ApolloConf struct {
	AppId     string `json:"app_id"`
	NameSpace string `json:"name_space"`
	Address   string `json:"address"`
	Cluster   string `json:"cluster"`
	Secret    string `json:"secret"`
}

func initApollo(conf *ApolloConf) (bconf.Config, error) {
	c, err := apollo.NewApolloConfig(
		apollo.WithAppid(conf.AppId),
		apollo.WithNamespace(conf.NameSpace),
		apollo.WithAddr(conf.Address),
		apollo.WithCluster(conf.Cluster),
		apollo.WithSecret(conf.Secret),
	)
	if err != nil {
		log.Errorf("err is : %v", err)
		return nil, err
	}

	config, err := lbconf.NewConfig(lbconf.WithDataSource(c))
	if err != nil {
		log.Errorf("err is : %v", err)
		return nil, err
	}

	err = config.Load()
	if err != nil {
		log.Errorf("err is : %v", err)
		return nil, err
	}

	if err = config.Watch(func(path string, v bconf.Val) {
		log.Infof("path %s val %+v\n", path, v)
	}); err != nil {
		log.Errorf("err is : %v", err)
		return nil, err
	}
	log.Infof("init apollo successfully\n")

	return config, nil
}

func getJson4Apollo(conf bconf.Config, key string, out interface{}) error {
	re, err := conf.Get(key)
	if err != nil {
		log.Errorf("err is : %v", err)
		return err
	}
	marshal, err := json.Marshal(re)
	if err != nil {
		log.Errorf("err is : %v", err)
		return err
	}
	err = json.Unmarshal(marshal, out)
	if err != nil {
		log.Errorf("err is : %v", err)
		return err
	}
	return nil
}
