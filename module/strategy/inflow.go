package strategy

import (
	"context"
	"futu-openapi/futuapi"
	"futu-openapi/internal/pb/qotstockfilter"
	"futu-openapi/module"
	"futu-openapi/module/fundamentals"
	"futu-openapi/module/gcs"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func NInFlow(day int32) ([]*qotstockfilter.StockData, error) {
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
	filter := &futuapi.StockFilter{
		BaseFilterList: basicFilter,
	}

	var errgroup []error
	filterData := make([]*qotstockfilter.StockData, 0, 200)
	var page int32 = 1
	var num int32 = 200
	for i := 0; i < 5; i++ {
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
		time.Sleep(time.Second * 2)
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
