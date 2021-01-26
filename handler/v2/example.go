package v2

import (
	"context"
	"fmt"
	"time"

	protos_ep "github.com/flamefatex/protos/goout/example"
	protos_v2 "github.com/flamefatex/protos/goout/example-api/v2"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/flamefatex/example-api/service/example"
)

type exampleHandler struct {
	exampleSvc example.ExampleSvc
}

func NewExampleHandler() *exampleHandler {
	return &exampleHandler{
		exampleSvc: example.ExampleSvcInstance(),
	}
}

func (h *exampleHandler) All(ctx context.Context, req *protos_v2.ExampleAllRequest) (resp *protos_v2.ExampleAllResponse, err error) {
	thisTime, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		err = fmt.Errorf("ptypes.TimestampProto err:%w", err)
		return
	}

	resp = &protos_v2.ExampleAllResponse{
		Examples: []*protos_v2.Example{
			{
				Id:          1,
				Name:        "example",
				Type:        protos_ep.ExampleType_EXAMPLE_TYPE_ACCESS,
				Description: "样例",
				CreateTime:  thisTime,
				UpdateTime:  thisTime,
			},
		},
	}
	return
}

func (h *exampleHandler) List(ctx context.Context, req *protos_v2.ExampleListRequest) (resp *protos_v2.ExampleListResponse, err error) {
	return
}
func (h *exampleHandler) Get(ctx context.Context, req *protos_v2.ExampleGetRequest) (resp *protos_v2.ExampleGetResponse, err error) {
	return
}
func (h *exampleHandler) Create(ctx context.Context, req *protos_v2.Example) (emp *empty.Empty, err error) {
	return
}
func (h *exampleHandler) Update(ctx context.Context, req *protos_v2.Example) (emp *empty.Empty, err error) {
	return
}
func (h *exampleHandler) Delete(ctx context.Context, req *protos_v2.ExampleDeleteRequest) (emp *empty.Empty, err error) {
	return nil, nil
}
