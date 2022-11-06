package futuapi

import (
	"context"
	"futu-openapi/internal/pb/qotcommon"
	"futu-openapi/internal/pb/qotsub"
)

const (
	ProtoIDSubscribe = 3001
)

func (api *FutuAPI) Subscribe() {

}

func (api *FutuAPI) UnSubscribe() {

}

func (api *FutuAPI) sub(ctx context.Context, sl []*qotcommon.Security, stl []int32, isSub, isReg bool, regPushRehabTypeList []int32) {
	req := &qotsub.Request{
		C2S: &qotsub.C2S{
			SecurityList:         sl,
			SubTypeList:          stl,
			IsSubOrUnSub:         &isSub,
			IsRegOrUnRegPush:     &isReg,
			RegPushRehabTypeList: regPushRehabTypeList,
		},
	}

	rsp := make(qotsub.ResponseChan)
	err := api.req(ProtoIDSubscribe, req, rsp)
	if err != nil {

	}
}
