// 获取快照
// doc: https://openapi.futunn.com/futu-api-doc/quote/get-market-snapshot.html
package futuapi

import (
	"context"
	"futu-openapi/internal/pb/qotcommon"
	"futu-openapi/internal/pb/qotgetsecuritysnapshot"
)

const (
	ProtoIDGetSecuritySnapshot = 3203
)

// 获取快照
func (api *FutuAPI) GetSecuritySnapshot(ctx context.Context, sl []*qotcommon.Security) (*qotgetsecuritysnapshot.Response, error) {
	req := &qotgetsecuritysnapshot.Request{
		C2S: &qotgetsecuritysnapshot.C2S{
			SecurityList: sl,
		},
	}

	rsp := make(qotgetsecuritysnapshot.ResponseChan)
	err := api.req(ProtoIDGetSecuritySnapshot, req, rsp)
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
