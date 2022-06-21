package qotgetcapitalflow

import "google.golang.org/protobuf/proto"

type ResponseChan chan *Response

func (ch ResponseChan) Send(b []byte) error {
	resp := new(Response)
	if err := proto.Unmarshal(b, resp); err != nil {
		return err
	}

	ch <- resp
	return nil
}

func (ch ResponseChan) Close() {
	close(ch)
}
