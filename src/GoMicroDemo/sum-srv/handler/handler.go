package handler

import (
	"context"
	"leslack/src/GoMicroDemo/proto/sum"
	"leslack/src/GoMicroDemo/sum-srv/service"
)

type handler struct {
}

func (h handler) GetSum(ctx context.Context, request *sum.SumRequest, response *sum.SumResponse) error {
	//panic("implement me")
	var inputs []int64
	for i := int64(0); i <= request.Input; i++ {
		inputs = append(inputs, i)
	}
	response.Output = service.GetSum(inputs...)
	return nil
}

func Handler() sum.SumHandler {
	return handler{}
}
