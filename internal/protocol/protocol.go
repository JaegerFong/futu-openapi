package protocol

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"futu-openapi/internal/tcp"
	"io"
	"net"
	"sync"

	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

const (
	PFT_PROTO = iota // protobuf
	PFT_JSON         // json
)

// byte: 2+4+2+2+4+4+20+8= 46
type protoHeader struct {
	HeaderFlag   [2]byte  // 包头标志，默认为 FT
	ProtoID      uint32   // 协议 ID
	ProtoFmtType uint8    // 协议格式类型，0 为 Protobuf 格式，1 为 Json 格式 ,推荐protobuf
	ProtoVer     uint8    // 协议版本，用于迭代兼容，目前填 0
	SerialNo     uint32   // 包序列号，用于对应请求包和回包，要求递增
	BodyLen      uint32   // 包体长度
	BodySHA1     [20]byte // 包体原始数据(解密后)的 SHA1 哈希值
	Reserved     [8]byte  // 保留 8 字节扩展
}

type FutuEncoder struct {
	proto uint32
	seq   uint32
	msg   proto.Message
}

func NewFutuEncoder(p, s uint32, m proto.Message) *FutuEncoder {
	return &FutuEncoder{
		proto: p,
		seq:   s,
		msg:   m,
	}
}

func (e *FutuEncoder) WriteTo(c net.Conn) error {
	body, err := proto.Marshal(e.msg)
	if err != nil {
		return err
	}

	h := protoHeader{
		HeaderFlag:   [2]byte{'F', 'T'},
		ProtoID:      e.proto,
		ProtoFmtType: PFT_PROTO,
		ProtoVer:     0,
		SerialNo:     e.seq,
		BodyLen:      uint32(len(body)),
	}

	bodySha1 := sha1.Sum(body)
	h.BodySHA1 = bodySha1

	var buf bytes.Buffer
	err = binary.Write(&buf, binary.LittleEndian, &h)
	if err != nil {
		return err
	}

	_, err = buf.Write(body)
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(c)
	if err != nil {
		return err
	}

	return nil
}

type FutuDecoder struct {
	reg *Registry
}

func NewFutuDecoder(reg *Registry) *FutuDecoder {
	return &FutuDecoder{
		reg: reg,
	}
}

// 从连接中读取数据
func (d *FutuDecoder) ReadFrom(c net.Conn) (tcp.Handler, error) {

	var h protoHeader
	err := binary.Read(c, binary.LittleEndian, &h)
	if err != nil {
		return nil, err
	}

	if h.HeaderFlag != [2]byte{'F', 'T'} {
		return nil, errors.New("header flag error")
	}

	bodyLen := h.BodyLen
	body := make([]byte, bodyLen)
	_, err = io.ReadFull(c, body)
	if err != nil {
		return nil, err
	}

	bodySha1 := sha1.Sum(body)
	for i, c := range bodySha1 {
		if h.BodySHA1[i] != c {
			return nil, errors.New("SHA1 not match")
		}
	}

	// 获取连接
	ch := d.reg.Get(h.ProtoID, h.SerialNo)
	return &handler{
		ch:   ch,
		body: body,
	}, nil

}

type RespChan interface {
	Send(b []byte) error
	Close()
}

const (
	CHAN_TYPE_GET      uint32 = 1
	CHAN_TYPE_CALLBACK uint32 = 2
)

type Registry struct {
	m  map[string]RespChan
	mu sync.Mutex
}

func NewRegistry() *Registry {
	return &Registry{
		m: make(map[string]RespChan),
	}
}

func (reg *Registry) Close() {

}

func (reg *Registry) Regist(protoId, serial uint32, ch RespChan) error {
	k := reg.rkey(protoId, serial)

	reg.mu.Lock()
	defer reg.mu.Unlock()

	reg.m[k] = ch
	return nil
}

func (reg *Registry) UnRegist() {}

func (reg *Registry) Get(protoId, serial uint32) RespChan {
	k := reg.rkey(protoId, serial)

	reg.mu.Lock()
	defer reg.mu.Unlock()

	respChan := reg.m[k]
	return respChan
}

func (reg *Registry) rkey(protoId, serial uint32) string {
	return fmt.Sprintf("%d-%d", protoId, serial)
}

type handler struct {
	ch     RespChan
	proto  uint32
	serial uint32
	body   []byte
}

func (h *handler) Handle() error {
	log.Infof("handle: proto=%d serial=%d", h.proto, h.serial)
	// handle
	err := h.ch.Send(h.body)
	if err != nil {
		log.Errorf("send fail: proto=%d serial=%d", h.proto, h.serial)
		return err
	}

	return nil
}
