package impl

import (
	"context"
	"github.com/oldbai555/lbtool/demo/pb"
)

type blogServerImpl struct {
	*pb.UnimplementedBlogServer
}

func (b blogServerImpl) GetBlog(ctx context.Context, req *pb.GetBlogReq) (*pb.GetBlogRsp, error) {
	var rsp pb.GetBlogRsp

	return &rsp, nil
}

var BlogServerImpl = &blogServerImpl{}

var _ pb.BlogServer = (*blogServerImpl)(nil)
