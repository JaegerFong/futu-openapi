package fundamentals

import (
	"context"
	"errors"
	"futu-openapi/internal/pb/qotcommon"
	"futu-openapi/internal/pb/qotgetcapitaldistribution"
	"futu-openapi/internal/pb/qotgetcapitalflow"
	"futu-openapi/internal/pb/qotgetsecuritysnapshot"
	"futu-openapi/module"
	"futu-openapi/module/gcs"
)

type DistrAndFlow struct {
	Flow         *qotgetcapitalflow.S2C
	Distribution *qotgetcapitaldistribution.S2C
}

func GetCapitalDistrAndFlow(market int32, code string, period qotcommon.PeriodType, begin, end string) (*DistrAndFlow, error) {
	security := &qotcommon.Security{
		Market: &market,
		Code:   &code,
	}
	cli := module.GetFutuClient()
	ctx := context.Background()
	flowRsp, err := cli.GetCapitalFlow(ctx, security, int32(period), begin, end)
	if err != nil {
		return nil, err
	}

	if flowRsp.GetRetType() != gcs.RET_SUC {
		return nil, errors.New(flowRsp.GetRetMsg())
	}

	result := &DistrAndFlow{
		Flow: flowRsp.S2C,
	}

	distrRsp, err := cli.GetCapitalDistribution(ctx, security)
	if err != nil {
		return nil, err
	}

	if distrRsp.GetRetType() != gcs.RET_SUC {
		return nil, errors.New(distrRsp.GetRetMsg())

	}
	result.Distribution = distrRsp.S2C
	return result, nil
}

func GetSnapShot(market []*qotcommon.Security) ([]*qotgetsecuritysnapshot.Snapshot, error) {
	cli := module.GetFutuClient()
	ctx := context.Background()
	rsp, err := cli.GetSecuritySnapshot(ctx, market)
	if err != nil {
		return nil, err
	}

	if rsp.GetRetType() != gcs.RET_SUC {
		return nil, errors.New(rsp.GetRetMsg())
	}

	return rsp.S2C.GetSnapshotList(), nil
}
