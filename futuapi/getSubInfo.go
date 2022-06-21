package futuapi

import (
	"context"
	"futu-openapi/internal/pb/qotgetsubinfo"
)

const (
	ProtoIDGetSubInfo = 3003
)

func (api *FutuAPI) GetSubInfo(ctx context.Context, isAllReq bool) (*qotgetsubinfo.Response, error) {
	req := &qotgetsubinfo.Request{
		C2S: &qotgetsubinfo.C2S{
			IsReqAllConn: &isAllReq,
		},
	}

	rsp := make(qotgetsubinfo.ResponseChan)
	err := api.req(ProtoIDGetSubInfo, req, rsp)
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
