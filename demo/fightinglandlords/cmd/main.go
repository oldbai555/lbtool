package main

//
//import (
//	"github.com/oldbai555/lbtool/demo/fightinglandlords"
//	"github.com/oldbai555/lbtool/log"
//	"github.com/oldbai555/lbtool/pkg/baix/baix"
//	"github.com/oldbai555/lbtool/pkg/baix/iface"
//	"github.com/urfave/cli/v2"
//	"os"
//	"strings"
//)
//
//func main() {
//	app := cli.NewApp()
//	app.Name = "baiX"
//	app.Version = "v0.0.1"
//	app.Description = "lb tcp server"
//	app.Action = serve
//	err := app.Run(os.Args)
//	if err != nil {
//		log.Errorf("err:%v", err)
//		return
//	}
//	return
//}
//
//func serve(c *cli.Context) error {
//	s := baix.NewServer()
//
//	initRouter(s)
//
//	return s.Serve(c.Context)
//}
//
//func initRouter(s iface.IServer) {
//	for k, router := range fightinglandlords.RouterMap {
//		s.RegisterRouter(strings.ToLower(k), router)
//	}
//}
