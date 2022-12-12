package proto_emicklei

import (
	"fmt"
	"os"

	"github.com/emicklei/proto"
)

func GenDemo() {
	reader, _ := os.Open("test.proto")
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, _ := parser.Parse()

	proto.Walk(definition,
		proto.WithService(handleService),
		proto.WithMessage(handleMessage))
}

func handleService(s *proto.Service) {
	fmt.Println(s.Name)
}

func handleMessage(m *proto.Message) {
	lister := new(optionLister)
	for _, each := range m.Elements {
		each.Accept(lister)
	}
	fmt.Println(m.Name)
}

type optionLister struct {
	proto.NoopVisitor
}

func (l optionLister) VisitOption(o *proto.Option) {
	fmt.Println(o.Name)
}
