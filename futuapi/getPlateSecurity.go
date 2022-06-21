package futuapi

import (
	"context"
	"futu-openapi/internal/pb/qotcommon"
	"futu-openapi/internal/pb/qotgetplatesecurity"
)

const (
	ProtoIDGetPlateSecurity = 3205
)

func (api *FutuAPI) GetPlateSecurity(ctx context.Context, plate *qotcommon.Security, sortField int32, ascend bool) (*qotgetplatesecurity.Response, error) {
	req := &qotgetplatesecurity.Request{
		C2S: &qotgetplatesecurity.C2S{
			Plate:     plate,
			SortField: &sortField,
			Ascend:    &ascend,
		},
	}

	rsp := make(qotgetplatesecurity.ResponseChan)
	err := api.req(ProtoIDGetPlateSecurity, req, rsp)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-rsp:
		if !ok {
			return nil, ErrChannelClosed
		}

		return resp, nil
	}
}
