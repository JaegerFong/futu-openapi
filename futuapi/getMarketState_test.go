package futuapi

import (
	"context"
	"futu-openapi/internal/pb/qotcommon"
	"testing"
)

func TestFutuAPI_GetMarketState(t *testing.T) {
	var m int32 = int32(qotcommon.QotMarket_QotMarket_CNSH_Security)
	lj := "601012"
	an := "603363"

	req := make([]*qotcommon.Security, 2)
	req[0] = &qotcommon.Security{
		Market: &m,
		Code:   &lj,
	}

	req[1] = &qotcommon.Security{
		Market: &m,
		Code:   &an,
	}

	apireq := NewFutuAPIT(1, "jinxtest")

	ctx := context.Background()
	apireq.Connect(ctx, "127.0.0.1:11111")
	rsp, err := apireq.GetMarketState(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	rsp2, err := apireq.GetSecuritySnapshot(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(rsp, rsp2)
}
