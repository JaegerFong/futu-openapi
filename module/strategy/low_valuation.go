// 1-3月 -50%
// 近10个交易日涨跌幅20%
// 近5个交易日资金为净流入
// 市盈率30以内
package strategy

import (
	"context"
	"futu-openapi/futuapi"
	"futu-openapi/internal/pb/qotstockfilter"
	"futu-openapi/module"
	"futu-openapi/module/fundamentals"
	"futu-openapi/module/gcs"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func LowValuation() ([]*qotstockfilter.StockData, error) {
	isNoFilter := false
	cli := module.GetFutuClient()
	ctx := context.Background()
	basicFilter := make([]*qotstockfilter.BaseFilter, 0, 1)
	var ttmMin float64 = 15
	var ttmMax float64 = 30
	ttmName := int32(qotstockfilter.StockField_StockField_PeTTM)
	ttmFilter := &qotstockfilter.BaseFilter{
		FieldName:  &ttmName,
		FilterMin:  &ttmMin,
		FilterMax:  &ttmMax,
		IsNoFilter: &isNoFilter,
	}
	basicFilter = append(basicFilter, ttmFilter)

	acfilter := make([]*qotstockfilter.AccumulateFilter, 0, 2)
	var day60Min float64 = -60
	var day60Max float64 = -30
	var day60 int32 = 60
	day60ChangeRateFilter := &qotstockfilter.AccumulateFilter{
		FieldName:  (*int32)(qotstockfilter.AccumulateField_AccumulateField_ChangeRate.Enum()),
		FilterMin:  &day60Min,
		FilterMax:  &day60Max,
		IsNoFilter: &isNoFilter,
		Days:       &day60,
	}
	acfilter = append(acfilter, day60ChangeRateFilter)

	var day10Min float64 = 2
	var day10Max float64 = 10
	var day10 int32 = 10
	day10ChangeRateFilter := &qotstockfilter.AccumulateFilter{
		FieldName:  (*int32)(qotstockfilter.AccumulateField_AccumulateField_ChangeRate.Enum()),
		FilterMin:  &day10Min,
		FilterMax:  &day10Max,
		IsNoFilter: &isNoFilter,
		Days:       &day10,
	}
	acfilter = append(acfilter, day10ChangeRateFilter)

	filter := &futuapi.StockFilter{
		BaseFilterList:   basicFilter,
		AccumulateFilter: acfilter,
	}

	var errgroup []error
	filterData := make([]*qotstockfilter.StockData, 0, 200)
	var page int32 = 1
	var num int32 = 200
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second * 5)
		var begin int32 = (page - 1) * num
		rsp, err := cli.StockFilter(ctx, begin, num, gcs.MARKET_SH, filter)
		if err != nil {
			// error
			errgroup = append(errgroup, gcs.NewErr(gcs.ErrInternal, err.Error()))
		}

		if rsp.GetRetType() != gcs.RET_SUC {
			// error
			errgroup = append(errgroup, gcs.NewErr(int(rsp.GetErrCode()), rsp.GetRetMsg()))
		}

		s2c := rsp.GetS2C()
		count := s2c.GetAllCount()
		if count == 0 {
			page++
			continue
		}

		stockData := s2c.GetDataList()
		filterData = append(filterData, stockData...)
		page++
	}

	printer := message.NewPrinter(language.English)
	nowTime := time.Now()
	beginDate := nowTime.AddDate(0, 0, -12).Format("2006-01-02")
	endDate := nowTime.AddDate(0, 0, -1).Format("2006-01-02")
	for _, v := range filterData {
		stockName := v.GetName()
		baseData := v.GetSecurity()
		market := baseData.GetMarket()
		stockCode := baseData.GetCode()
		distrAndFlow, err := fundamentals.GetCapitalDistrAndFlow(market, stockCode, gcs.PERIOD_DAY, beginDate, endDate)
		if err != nil {
			log.Errorf("GetCapitalDistrAndFlow Error %s", err)
			continue
		}

		// distribution
		distr := distrAndFlow.Distribution
		superIn := distr.GetCapitalInSuper()
		bigIn := distr.GetCapitalInBig()
		midIn := distr.GetCapitalInMid()
		smallIn := distr.GetCapitalInSmall()
		superOut := distr.GetCapitalOutSuper()
		bigOut := distr.GetCapitalOutBig()
		midOut := distr.GetCapitalOutMid()
		smallOut := distr.GetCapitalOutSmall()
		realSuperIn := superIn - superOut
		fRealSuperIn := printer.Sprintf("%.2f", realSuperIn)
		realBigIn := bigIn - bigOut
		fRealBigIn := printer.Sprintf("%.2f", realBigIn)
		realMidIn := midIn - midOut
		fRealMidIn := printer.Sprintf("%.2f", realMidIn)
		realSmallIn := smallIn - smallOut
		fRealSmallIn := printer.Sprintf("%.2f", realSmallIn)
		logrus.Infof("『%s %s RealIn』Super:%s Big:%s Mid:%s Small:%s", stockName, stockCode, fRealSuperIn, fRealBigIn, fRealMidIn, fRealSmallIn)
		// flow
		flowList := distrAndFlow.Flow.FlowItemList
		for _, fv := range flowList {
			flowTime := fv.GetTime()
			inflow := fv.GetInFlow()
			fInflow := printer.Sprintf("%.2f", inflow)
			logrus.Infof("『%s %s』%s ---- %s", stockName, stockCode, fInflow, flowTime)
		}
	}

	return filterData, nil
}

// BaseFilter
// fieldName: StockField_PeTTM
// filterMin: 15
// filterMax: 30
// isNoFilter: false

// AccumulateFilter
// fieldName: AccumulateField_ChangeRate
// filterMin: -60
// filterMax: -50
// isNoFilter: false
// days: 60

// AccumulateFilter
// fieldName: AccumulateField_ChangeRate
// filterMin: 0
// filterMax: +20
// isNoFilter: false
// days: 10

// GetCapitalDistrAndFlow
// market code sort ascend
