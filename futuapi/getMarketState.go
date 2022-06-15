// 获取市场标的状态
// doc: https://openapi.futunn.com/futu-api-doc/quote/get-market-state.html
package futuapi

import (
	"context"
	"futu-openapi/internal/pb/qotcommon"
	"futu-openapi/internal/pb/qotgetmarketstate"
)

const (
	ProtoIDGetMarketState = 3223
)

// 标的的市场状态
func (api *FutuAPI) GetMarketState(ctx context.Context, sl []*qotcommon.Security) (*qotgetmarketstate.Response, error) {
	req := &qotgetmarketstate.Request{
		C2S: &qotgetmarketstate.C2S{
			SecurityList: sl,
		},
	}

	rsp := make(qotgetmarketstate.ResponseChan)
	err := api.req(ProtoIDGetMarketState, req, rsp)
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
