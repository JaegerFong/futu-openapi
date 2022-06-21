package futuapi

import (
	"context"
	"futu-openapi/internal/pb/qotgetplateset"
)

const (
	ProtoIDGetPlateSet = 3204
)

func (api *FutuAPI) GetPlateSet(ctx context.Context, market, plateSetType int32) (*qotgetplateset.Response, error) {
	req := &qotgetplateset.Request{
		C2S: &qotgetplateset.C2S{
			Market:       &market,
			PlateSetType: &plateSetType,
		},
	}

	rsp := make(qotgetplateset.ResponseChan)
	err := api.req(ProtoIDGetPlateSet, req, rsp)
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
