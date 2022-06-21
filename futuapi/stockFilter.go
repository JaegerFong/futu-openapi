package futuapi

// 条件选股
// doc: https://openapi.futunn.com/futu-api-doc/quote/get-stock-filter.html
import (
	"context"
	"futu-openapi/internal/pb/qotcommon"
	"futu-openapi/internal/pb/qotstockfilter"
)

const (
	ProtoIDStockFilter = 3215
)

func (api *FutuAPI) StockFilter(ctx context.Context, begin, num, market int32, filter *StockFilter) (*qotstockfilter.Response, error) {
	req := &qotstockfilter.Request{
		C2S: &qotstockfilter.C2S{
			Begin:  &begin,
			Num:    &num,
			Market: &market,
		},
	}

	if filter != nil {
		if filter.plate != nil {
			req.C2S.Plate = filter.plate
		}

		if len(filter.BaseFilterList) > 0 {
			req.C2S.BaseFilterList = filter.BaseFilterList
		}

		if len(filter.AccumulateFilter) > 0 {
			req.C2S.AccumulateFilterList = filter.AccumulateFilter
		}

		if len(filter.FinancialFilter) > 0 {
			req.C2S.FinancialFilterList = filter.FinancialFilter
		}

		if len(filter.PatternFilter) > 0 {
			req.C2S.PatternFilterList = filter.PatternFilter
		}

		if len(filter.CustomIndicatorFilter) > 0 {
			req.C2S.CustomIndicatorFilterList = filter.CustomIndicatorFilter
		}
	}

	rsp := make(qotstockfilter.ResponseChan)
	err := api.req(ProtoIDStockFilter, req, rsp)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-rsp:
		{
			if !ok {
				return nil, ErrChannelClosed
			}

			return resp, nil
		}
	}
}

type StockFilter struct {
	plate                 *qotcommon.Security
	BaseFilterList        []*qotstockfilter.BaseFilter
	AccumulateFilter      []*qotstockfilter.AccumulateFilter
	FinancialFilter       []*qotstockfilter.FinancialFilter
	PatternFilter         []*qotstockfilter.PatternFilter
	CustomIndicatorFilter []*qotstockfilter.CustomIndicatorFilter
}
