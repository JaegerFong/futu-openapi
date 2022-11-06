package fundamentals

import (
	"futu-openapi/internal/pb/qotcommon"
	"futu-openapi/module"
	"futu-openapi/module/gcs"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func init() {
	module.DoInit()
}

func TestGetCapitalDistrAndFlow(t *testing.T) {
	today := time.Now().Format("2006-01-02")
	res, err := GetCapitalDistrAndFlow(gcs.MARKET_SH, "601012", gcs.PERIOD_REALTIME, today, today)
	if err != nil {
		t.Fatal(err)
	}

	logrus.Info("Flow")
	flow := res.Flow.FlowItemList
	for _, v := range flow {
		logrus.Infof("[UpdateTime:%s - InFlow:%.2f - Main:%.2f]Super:%.2f Big:%.2f Mid:%f Small:%.2f", v.GetTime(), v.GetInFlow(), v.GetMainInFlow(), v.GetSuperInFlow(), v.GetBigInFlow(), v.GetMidInFlow(), v.GetSmlInFlow())
	}

	distr := res.Distribution
	superIn := distr.GetCapitalInSuper()
	bigIn := distr.GetCapitalInBig()
	midIn := distr.GetCapitalInMid()
	smallIn := distr.GetCapitalInSmall()

	superOut := distr.GetCapitalOutSuper()
	bigOut := distr.GetCapitalOutBig()
	midOut := distr.GetCapitalOutMid()
	smallOut := distr.GetCapitalOutSmall()

	logrus.Info("Distribution In")
	logrus.Infof("Super:%.2f Big:%.2f Mid:%.2f Small:%.2f", superIn, bigIn, midIn, smallIn)

	logrus.Info("Distribution Out")
	logrus.Infof("Super:%.2f Big:%.2f Mid:%.2f Small:%.2f", superOut, bigOut, midOut, smallOut)

	realSuperIn := superIn - superOut
	realBigIn := bigIn - bigOut
	realMidIn := midIn - midOut
	realSmallIn := smallIn - smallOut

	logrus.Info("Real In")
	logrus.Infof("Super:%.2f Big:%.2f Mid:%.2f Small:%.2f", realSuperIn, realBigIn, realMidIn, realSmallIn)
}

func TestGetSnapShot(t *testing.T) {
	m := int32(gcs.MARKET_SH)
	code := "601012"
	markets := make([]*qotcommon.Security, 0, 1)
	markets = append(markets, &qotcommon.Security{
		Market: &m,
		Code:   &code,
	})

	snapShot, err := GetSnapShot(markets)
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range snapShot {
		basic := v.GetBasic()
		turnOver := basic.GetTurnover()
		turnOverRate := basic.GetTurnoverRate()
		amplitude := basic.GetAmplitude()
		avgPrice := basic.GetAvgPrice()
		volumnRatio := basic.GetVolumeRatio()
		logrus.Infof("TurnOver:%f TurnOverRate:%f Amplitude:%f AvgPrice:%f VolumnRatio:%f", turnOver, turnOverRate, amplitude, avgPrice, volumnRatio)
	}
}
