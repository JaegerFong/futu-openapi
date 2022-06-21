// doc: https://openapi.futunn.com/futu-api-doc/quote/get-capital-distribution.html
package futuapi

import (
	"context"
	"futu-openapi/internal/pb/qotcommon"
	"futu-openapi/internal/pb/qotgetcapitaldistribution"
)

const (
	ProtoIdGetCapitalDistribution = 3212
)

func (api *FutuAPI) GetCapitalDistribution(ctx context.Context, security *qotcommon.Security) (*qotgetcapitaldistribution.Response, error) {
	req := &qotgetcapitaldistribution.Request{
		C2S: &qotgetcapitaldistribution.C2S{
			Security: security,
		},
	}

	rsp := make(qotgetcapitaldistribution.ResponseChan)
	err := api.req(ProtoIdGetCapitalDistribution, req, rsp)
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
