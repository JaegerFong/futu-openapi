package fundamentals

import (
	"context"
	"futu-openapi/internal/pb/qotcommon"
	"futu-openapi/module"
	"futu-openapi/module/gcs"
)

//
func GetPlateList(market int32, plateType int32) ([]*qotcommon.PlateInfo, error) {

	cli := module.GetFutuClient()
	ctx := context.Background()
	rsp, err := cli.GetPlateSet(ctx, market, plateType)
	if err != nil {
		return nil, gcs.NewErr(gcs.ErrInternal, err.Error())
	}

	if rsp.GetRetType() != gcs.RET_SUC {
		return nil, gcs.NewErr(int(rsp.GetErrCode()), rsp.GetRetMsg())
	}

	return rsp.S2C.GetPlateInfoList(), nil
}

func GetPlateSecurityList(market int32, code string, sort int32, ascend bool) ([]*qotcommon.SecurityStaticInfo, error) {
	plate := &qotcommon.Security{
		Market: &market,
		Code:   &code,
	}
	cli := module.GetFutuClient()
	ctx := context.Background()
	rsp, err := cli.GetPlateSecurity(ctx, plate, sort, ascend)
	if err != nil {
		return nil, gcs.NewErr(gcs.ErrInternal, err.Error())
	}

	if rsp.GetRetType() != gcs.RET_SUC {
		return nil, gcs.NewErr(int(rsp.GetErrCode()), rsp.GetRetMsg())
	}

	return rsp.S2C.StaticInfoList, nil
}
