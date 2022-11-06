package module

import (
	"context"
	"futu-openapi/conf"
	"futu-openapi/futuapi"
)

var (
	futuApiClient *futuapi.FutuAPI
)

func initFutuConnect() {
	futuApiClient = futuapi.NewFutuAPIT(1, conf.CLIENT_ID)
	ctx := context.Background()
	futuApiClient.Connect(ctx, conf.CLIENT_ADDR)
}

func GetFutuClient() *futuapi.FutuAPI {
	return futuApiClient
}
