package fundamentals

import (
	"futu-openapi/internal/pb/qotcommon"
	"futu-openapi/module/gcs"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestGetPlateSecurityList(t *testing.T) {
	// BK0365
	code := "BK0194"
	l, err := GetPlateSecurityList(gcs.MARKET_SH, code, int32(qotcommon.SortField_SortField_Code),
		true)
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range l {
		basic := v.GetBasic()
		name := basic.GetName()
		sec := basic.GetSecurity()
		code := sec.GetCode()
		logrus.Infof("[%s-%s]", name, code)
	}
}

func TestGetPlateList(t *testing.T) {
	pl, err := GetPlateList(gcs.MARKET_SH, gcs.PLATE_ALL)
	if err != nil {
		t.Fatal(err)
	}

	hc := 0
	sc := 0
	for _, v := range pl {
		name := v.GetName()
		plate := v.GetPlate()
		code := plate.GetCode()
		market := plate.GetMarket()
		mn := "沪"
		if market == gcs.MARKET_SZ {
			mn = "深"
			sc++
		} else if market == gcs.MARKET_SH {
			hc++
		}
		logrus.Infof("[%s:%s-%s]", mn, name, code)
	}

	logrus.Infof("count:%d,沪:%d,深:%d", len(pl), hc, sc)
}
