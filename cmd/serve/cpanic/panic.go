package cpanic

import (
	"bytes"
	"fmt"
	"runtime"
	"time"
)

type Handler func(p *Panic)

type Panic struct {
	Time  time.Time   `json:"time"`
	Value interface{} `json:"value"`
	Trace string      `json:"trace"`
}

func Recover(fn Handler) {
	if fn == nil {
		return
	}
	if r := recover(); r != nil {
		var trace [1 << 16]byte
		n := runtime.Stack(trace[:], true)
		p := &Panic{
			Time:  time.Now(),
			Value: r,
			Trace: string(trace[:n]),
		}
		fn(p)
	}
}

func (p Panic) MarshalText() ([]byte, error) {
	var buf bytes.Buffer
	var err error

	_, err = fmt.Fprintf(&buf, "panic: %v\n\n", p.Value)
	if err != nil {
		return nil, err
	}

	_, err = fmt.Fprint(&buf, p.Trace)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (p Panic) String() string {
	b, _ := p.MarshalText()
	return string(b)
}
