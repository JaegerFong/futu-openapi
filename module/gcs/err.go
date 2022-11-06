package gcs

const (
	ErrInternal = iota + 50000
)

type Err struct {
	code int
	msg  string
}

func NewErr(c int, m string) *Err {
	return &Err{
		code: c,
		msg:  m,
	}
}

func (e *Err) Error() string {
	return e.msg
}
