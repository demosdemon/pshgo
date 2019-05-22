package middleware

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/lars"

	"github.com/demosdemon/pshgo/cmd/serve/cpanic"
)

const (
	LogTimeFormat     = "01/Jan/2006:15:04:05 -0700"
	RequestContextKey = "github.com/pshgo/cmd/serve/middleware/RequestContextKey"
)

type Request struct {
	Start      time.Time         `json:"start"`
	ID         string            `json:"id"`
	Username   string            `json:"username,omitempty"`
	ClientIP   string            `json:"client_id"`
	RemoteAddr string            `json:"remote_addr"`
	Method     string            `json:"method"`
	URL        string            `json:"url"`
	Proto      string            `json:"proto"`
	Referrer   string            `json:"referrer,omitempty"`
	UserAgent  string            `json:"use_agent"`
	Host       string            `json:"host"`
	Headers    map[string]string `json:"headers"`
	Delay      time.Duration     `json:"delay,omitempty"`
	Status     int               `json:"status,omitempty"`
	Size       int64             `json:"size,omitempty"`
	Panic      *cpanic.Panic     `json:"panic,omitempty"`
}

func NewRequest(c lars.Context) {
	now := time.Now()
	req := c.Request()

	r := &Request{
		Start:      now,
		ID:         randomID(),
		ClientIP:   c.ClientIP(),
		RemoteAddr: req.RemoteAddr,
		Method:     req.Method,
		URL:        req.URL.String(),
		Proto:      req.Proto,
		Referrer:   req.Referer(),
		UserAgent:  req.UserAgent(),
		Host:       req.Host,
		Headers:    cloneHeaders(req.Header),
	}

	c.Set(RequestContextKey, r)
	c.Next()

	// TODO: save req for long-term logging
}

func randomID() string {
	var slug [16]byte
	_, _ = rand.Read(slug[:])
	s := base32.StdEncoding.EncodeToString(slug[:])
	s = strings.ToLower(s)
	s = strings.TrimRight(s, "=")
	return s
}

func cloneHeaders(h http.Header) map[string]string {
	rv := make(map[string]string, len(h))
	for k, v := range h {
		rv[k] = strings.Join(v, "; ")
	}
	return rv
}

// Update returns a new Request object with fields updated from the lars response.
func (r *Request) UpdateLARS(res *lars.Response) {
	r.Delay = time.Since(r.Start)
	r.Status = res.Status()
	r.Size = res.Size()
}

func (r Request) String() string {
	pieces := []string{
		r.ClientIP,
		r.ID,
		r.Username,
		r.Start.Format(LogTimeFormat),
		fmt.Sprintf("%s %s %s", r.Method, r.URL, r.Proto),
		r.Host,
		r.Delay.String(),
		fmt.Sprintf("%03d", r.Status),
		fmt.Sprintf("%d", r.Size),
		r.Referrer,
		r.UserAgent,
	}

	for idx, v := range pieces {
		switch v {
		case "", "0", "000":
			pieces[idx] = "-"
		}

		if strings.Index(v, " ") >= 0 {
			pieces[idx] = fmt.Sprintf("%q", v)
		}
	}

	return strings.Join(pieces, " ")
}
