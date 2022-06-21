package futuapi

import (
	"context"
	"fmt"
	"futu-openapi/internal/pb/qotcommon"
	"futu-openapi/internal/pb/qotstockfilter"
	"net/http"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

const (
	CLIENT_ID   = "JINX"
	CLIENT_ADDR = "127.0.0.1:11111"
)

func TestFutuAPI_initConnect(t *testing.T) {
	c := NewFutuAPIT(1, CLIENT_ID)
	ctx := context.Background()
	r, err := c.initConnect(ctx, CLIENT_ADDR)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(r)
}

type indexHandler struct {
	content string
}

func (ih *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, ih.content)
}

func TestHttp(t *testing.T) {
	http.Handle("/", &indexHandler{content: "jinx"})
	http.ListenAndServe(":8081", nil)
}

func TestFutuAPI_GetCapitalDistribution(t *testing.T) {
	var m int32 = int32(qotcommon.QotMarket_QotMarket_CNSH_Security)
	ljcode := "601012"
	lj := &qotcommon.Security{
		Market: &m,
		Code:   &ljcode,
	}

	apireq := NewFutuAPIT(1, CLIENT_ID)

	ctx := context.Background()
	apireq.Connect(ctx, CLIENT_ADDR)
	rsp, err := apireq.GetCapitalDistribution(ctx, lj)
	if err != nil {
		t.Fatal(err)
	}

	logrus.Info(rsp)
}

func TestFutuAPI_GetCapitalFlow(t *testing.T) {
	var m int32 = int32(qotcommon.QotMarket_QotMarket_CNSH_Security)
	ljcode := "601012"
	lj := &qotcommon.Security{
		Market: &m,
		Code:   &ljcode,
	}
	begin := "2022-06-01"
	end := "2022-06-20"

	apireq := NewFutuAPIT(1, CLIENT_ID)

	ctx := context.Background()
	apireq.Connect(ctx, CLIENT_ADDR)
	rsp, err := apireq.GetCapitalFlow(ctx, lj, int32(qotcommon.PeriodType_PeriodType_DAY), begin, end)
	if err != nil {
		t.Fatal(err)
	}

	logrus.Info(rsp)
}

func TestFutuAPI_GetMarketState(t *testing.T) {
	var m int32 = int32(qotcommon.QotMarket_QotMarket_CNSH_Security)
	var sz int32 = int32(qotcommon.QotMarket_QotMarket_CNSZ_Security)
	lj := "601012"
	an := "603363"
	yj := "159865"

	req := make([]*qotcommon.Security, 3)
	req[0] = &qotcommon.Security{
		Market: &m,
		Code:   &lj,
	}

	req[1] = &qotcommon.Security{
		Market: &m,
		Code:   &an,
	}

	req[2] = &qotcommon.Security{
		Market: &sz,
		Code:   &yj,
	}

	apireq := NewFutuAPIT(1, CLIENT_ID)

	ctx := context.Background()
	apireq.Connect(ctx, CLIENT_ADDR)
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

func TestFutuAPI_StockFilter(t *testing.T) {
	var begin int32 = 1
	var num int32 = 10
	var tor float64 = 20
	isnofilter := false

	market := int32(qotcommon.QotMarket_QotMarket_CNSH_Security)
	var trunoverday int32 = 5
	turnoverrate := &qotstockfilter.AccumulateFilter{
		FieldName:  (*int32)(qotstockfilter.AccumulateField_AccumulateField_TurnoverRate.Enum()),
		FilterMin:  &tor,
		IsNoFilter: &isnofilter,
		Days:       &trunoverday,
	}

	var changedays int32 = 5
	changerate := &qotstockfilter.AccumulateFilter{
		FieldName:  (*int32)(qotstockfilter.AccumulateField_AccumulateField_ChangeRate.Enum()),
		FilterMin:  &tor,
		IsNoFilter: &isnofilter,
		Days:       &changedays,
	}

	afl := make([]*qotstockfilter.AccumulateFilter, 0, 2)
	afl = append(afl, turnoverrate)
	afl = append(afl, changerate)
	filter := &StockFilter{
		AccumulateFilter: afl,
	}

	ctx := context.Background()
	apireq := NewFutuAPIT(1, CLIENT_ID)
	apireq.Connect(ctx, CLIENT_ADDR)
	rsp, err := apireq.StockFilter(ctx, begin, num, market, filter)
	if err != nil {
		t.Fatal(err)
	}

	stocklist := rsp.S2C.GetDataList()
	for _, v := range stocklist {
		sec := v.Security
		fmt.Printf("market:%d code:%s name:%s \n", sec.GetMarket(), sec.GetCode(), v.GetName())
	}
	t.Log(rsp)
}

func TestFutuAPI_GetPlateSecurity(t *testing.T) {
	var m int32 = int32(qotcommon.QotMarket_QotMarket_CNSH_Security)
	bdtcode := "BK0194"
	bdt := &qotcommon.Security{
		Market: &m,
		Code:   &bdtcode,
	}

	apireq := NewFutuAPIT(1, CLIENT_ID)

	ctx := context.Background()
	apireq.Connect(ctx, CLIENT_ADDR)
	rsp, err := apireq.GetPlateSecurity(ctx, bdt, int32(qotcommon.SortField_SortField_Code), true)
	if err != nil {
		t.Fatal(err)
	}

	logrus.Info(rsp)
}

func TestFutuAPI_GetPlateSet(t *testing.T) {
	apireq := NewFutuAPIT(1, CLIENT_ID)

	ctx := context.Background()
	apireq.Connect(ctx, CLIENT_ADDR)
	rsp, err := apireq.GetPlateSet(ctx, int32(qotcommon.QotMarket_QotMarket_CNSH_Security), int32(qotcommon.PlateSetType_PlateSetType_All))
	if err != nil {
		t.Fatal(err)
	}

	plateList := rsp.S2C.GetPlateInfoList()
	for _, v := range plateList {
		if strings.Contains(*v.Name, "猪") {
			p := v.GetPlate()
			fmt.Printf("name:%s,market:%d,code:%s", *v.Name, p.GetMarket(), p.GetCode())
		}
	}

	// BK0194 BK0365

	logrus.Info(rsp)
}

func TestFutuAPI_FilterStockPlanA(t *testing.T) {
	var begin int32 = 1
	var num int32 = 10
	isnofilter := false

	market := int32(qotcommon.QotMarket_QotMarket_CNSH_Security)

	var consecutive5 int32 = 5
	p1 := &qotstockfilter.PatternFilter{
		FieldName:         (*int32)(qotstockfilter.PatternField_PatternField_MAAlignmentLong.Enum()),
		IsNoFilter:        &isnofilter,
		KlType:            (*int32)(qotcommon.KLType_KLType_Day.Enum()),
		ConsecutivePeriod: &consecutive5,
	}

	ptf := make([]*qotstockfilter.PatternFilter, 0, 2)
	ptf = append(ptf, p1)
	filter := &StockFilter{
		PatternFilter: ptf,
	}

	ctx := context.Background()
	apireq := NewFutuAPIT(1, CLIENT_ID)
	apireq.Connect(ctx, CLIENT_ADDR)
	rsp, err := apireq.StockFilter(ctx, begin, num, market, filter)
	if err != nil {
		t.Fatal(err)
	}

	stocklist := rsp.S2C.GetDataList()
	sl := make([]*qotcommon.Security, 0, len(stocklist))
	for _, v := range stocklist {
		sec := v.Security
		fmt.Printf("market:%d code:%s name:%s \n", sec.GetMarket(), sec.GetCode(), v.GetName())
		i := &qotcommon.Security{
			Market: v.Security.Market,
			Code:   v.Security.Code,
		}
		sl = append(sl, i)
	}

	sldetail, err := apireq.GetSecuritySnapshot(ctx, sl)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(rsp, sldetail)
}

func TestFutuAPI_GetSubInfo(t *testing.T) {
	ctx := context.Background()
	apireq := NewFutuAPIT(1, CLIENT_ID)
	apireq.Connect(ctx, CLIENT_ADDR)
	resp, err := apireq.GetSubInfo(ctx, true)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resp)
}
