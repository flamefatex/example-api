package external

import (
	"context"
	"fmt"
	"time"

	protos_ep "github.com/flamefatex/protos/goout/example"
	protos_v2_ext "github.com/flamefatex/protos/goout/example-api/v2/external"
	"github.com/golang/protobuf/ptypes"
)

type exampleHandler struct {
}

func NewExampleHandler() *exampleHandler {
	return &exampleHandler{}
}

func (h *exampleHandler) All(ctx context.Context, req *protos_v2_ext.ExampleAllRequest) (resp *protos_v2_ext.ExampleAllResponse, err error) {
	thisTime, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		err = fmt.Errorf("ptypes.TimestampProto err:%w", err)
		return
	}

	resp = &protos_v2_ext.ExampleAllResponse{
		Examples: []*protos_v2_ext.Example{
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
