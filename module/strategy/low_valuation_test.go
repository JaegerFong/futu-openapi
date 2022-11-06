// 1-3月 -50%
// 近10个交易日涨跌幅20%
// 近5个交易日资金为净流入
// 市盈率30以内
package strategy

import (
	"futu-openapi/module"
	"log"
	"testing"
)

func init() {
	module.DoInit()
}

func TestLowValuation(t *testing.T) {
	rsp, err := LowValuation()
	if err != nil {
		t.Fatal(err)
	}
	log.Fatal(rsp)
}
