package gcs

import (
	"futu-openapi/internal/pb/common"
	"futu-openapi/internal/pb/qotcommon"
)

// 市场
const (
	MARKET_SH = int32(qotcommon.QotMarket_QotMarket_CNSH_Security)
	MARKET_SZ = int32(qotcommon.QotMarket_QotMarket_CNSZ_Security)
	MARKET_HK = int32(qotcommon.QotMarket_QotMarket_HK_Security)
)

//
const (
	PLATE_ALL      = int32(qotcommon.PlateSetType_PlateSetType_All)
	PLATE_INDUSTRY = int32(qotcommon.PlateSetType_PlateSetType_Industry)
	PLATE_REGION   = int32(qotcommon.PlateSetType_PlateSetType_Region)
	PLATE_CONCEPT  = int32(qotcommon.PlateSetType_PlateSetType_Concept)
	PLATE_OTHER    = int32(qotcommon.PlateSetType_PlateSetType_Other)
)

// 周期
const (
	PERIOD_REALTIME = qotcommon.PeriodType_PeriodType_INTRADAY
	PERIOD_DAY      = qotcommon.PeriodType_PeriodType_DAY
	PERIOD_WEEK     = qotcommon.PeriodType_PeriodType_WEEK
	PERIOD_MONTH    = qotcommon.PeriodType_PeriodType_MONTH
)

// 返回结果
const (
	RET_SUC  = int32(common.RetType_RetType_Succeed)
	RET_FAIL = int32(common.RetType_RetType_Failed)
)
