// doc: https://openapi.futunn.com/futu-api-doc/quote/get-capital-flow.html
package futuapi

import (
	"context"
	"futu-openapi/internal/pb/qotcommon"
	"futu-openapi/internal/pb/qotgetcapitalflow"
)

const (
	ProtoIdGetCapitalFlow = 3211
)

func (api *FutuAPI) GetCapitalFlow(ctx context.Context, security *qotcommon.Security, periodType int32, begin, end string) (*qotgetcapitalflow.Response, error) {
	req := &qotgetcapitalflow.Request{
		C2S: &qotgetcapitalflow.C2S{
			Security:   security,
			PeriodType: &periodType,
			BeginTime:  &begin,
			EndTime:    &end,
		},
	}

	rsp := make(qotgetcapitalflow.ResponseChan)
	err := api.req(ProtoIdGetCapitalFlow, req, rsp)
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
