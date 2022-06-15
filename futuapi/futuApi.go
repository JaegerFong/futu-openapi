// 基础API
package futuapi

import (
	"context"
	"errors"
	"futu-openapi/internal/pb/common"
	"futu-openapi/internal/pb/initconnect"
	"futu-openapi/internal/protocol"
	"futu-openapi/internal/tcp"
	"sync"

	"google.golang.org/protobuf/proto"
)

var (
	ErrInterrupted   = errors.New("process is interrupted")
	ErrChannelClosed = errors.New("channel is closed")
)

const (
	ProtoIDInitConnect = 1001 //InitConnect
)

type FutuAPI struct {
	clientVer  int32
	clientID   string
	recvNotify bool
	encAlgo    common.PacketEncAlgo
	protoFmt   common.ProtoFmt
	lang       string

	serial uint32
	mu     sync.Mutex

	// 数据接受注册表
	reg *protocol.Registry

	// 建立的连接
	conn *tcp.Conn
}

func NewFutuAPIT(ver int32, clientId string) *FutuAPI {
	return &FutuAPI{
		clientVer: ver,
		clientID:  clientId,
		encAlgo:   common.PacketEncAlgo_PacketEncAlgo_None,
		protoFmt:  common.ProtoFmt_ProtoFmt_Protobuf,
		lang:      "golang",
		mu:        sync.Mutex{},
		reg:       protocol.NewRegistry(),
	}
}

func (api *FutuAPI) Connect(ctx context.Context, addr string) {
	api.initConnect(ctx, addr)
}

func (api *FutuAPI) initConnect(ctx context.Context, addr string) (*initconnect.Response, error) {
	de := protocol.NewFutuDecoder(api.reg)
	c, err := tcp.Dial(addr, de)
	if err != nil {
		return nil, err
	}

	api.conn = c
	req := &initconnect.Request{
		C2S: &initconnect.C2S{
			ClientVer:           &api.clientVer,
			ClientID:            &api.clientID,
			RecvNotify:          &api.recvNotify,
			PacketEncAlgo:       (*int32)(&api.encAlgo),
			PushProtoFmt:        (*int32)(&api.protoFmt),
			ProgrammingLanguage: &api.lang,
		},
	}

	rsp := make(initconnect.ResponseChan)
	err = api.req(ProtoIDInitConnect, req, rsp)
	if err != nil {
		return nil, err
	}

	// 监听channel
	select {
	case <-ctx.Done():
		return nil, ErrInterrupted
	case resp, ok := <-rsp:
		if !ok {
			return nil, ErrChannelClosed
		}

		// TODO: 保存连接信息

		return resp, nil
	}
}

func (api *FutuAPI) req(proto uint32, req proto.Message, rsp protocol.RespChan) error {
	// 注册
	ser := api.serialNo()
	api.reg.Regist(proto, ser, rsp)
	// 发送请求
	en := protocol.NewFutuEncoder(proto, ser, req)
	err := api.conn.Send(en)
	if err != nil {
		return err
	}

	return nil
}

func (api *FutuAPI) serialNo() uint32 {
	api.mu.Lock()
	defer api.mu.Unlock()
	api.serial++
	return api.serial
}
