package qotgetsubinfo

import "google.golang.org/protobuf/proto"

type ResponseChan chan *Response

func (ch ResponseChan) Send(b []byte) error {
	rsp := new(Response)
	if err := proto.Unmarshal(b, rsp); err != nil {
		return err
	}

	ch <- rsp
	return nil
}

func (ch ResponseChan) Close() {
	close(ch)
}
