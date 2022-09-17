package zipkin

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	openzipkin "github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
)

var (
	zkReporter reporter.Reporter
	zkTracer   opentracing.Tracer
)

// Req 链路追踪
type Req struct {
	// ServiceName 当前服务名称，用于注册到zipkin
	ServiceName string `json:"serviceName"`
	// ServerAddress 当前服务地址
	ServerAddress string `json:"serverAddress"`
	// ZipkinAddr zipkin的服务地址
	ZipkinAddr string `json:"zipkinAddr"`
}

// InitZipkinTracer 初始化zipkin客户端，并将服务注册到zipkin
func InitZipkinTracer(req Req) (opentracing.Tracer, error) {

	zkReporter = zipkinHTTP.NewReporter(req.ZipkinAddr)

	endpoint, err := openzipkin.NewEndpoint(req.ServiceName, req.ServerAddress)
	if err != nil {
		hlog.Fatalf("unable to create local endpoint: %+v\n", err)
		return nil, err
	}

	nativeTracer, err := openzipkin.NewTracer(zkReporter, openzipkin.WithTraceID128Bit(true), openzipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		hlog.Fatalf("unable to create tracer: %+v\n", err)
		return nil, err
	}

	zkTracer = zipkinot.Wrap(nativeTracer)
	opentracing.SetGlobalTracer(zkTracer)

	// 将tracer注入到 Hertz 的中间件中
	//h.Use(func(c context.Context, ctx *app.RequestContext) {
	//	span := zkTracer.StartSpan(ctx.FullPath())
	//	defer span.Finish()
	//	ctx.Next(c)
	//})

	return zkTracer, nil
}

func Close() {
	err := zkReporter.Close()
	if err != nil {
		hlog.Errorf("err is : %v", err)
	}
}
