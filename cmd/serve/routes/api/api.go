package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/demosdemon/pshgo/cmd/serve/server"
	"github.com/go-playground/lars"
	"github.com/go-playground/lars/middleware"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	StatusOK = http.StatusOK

	XForwardedProto    = "X-Forwarded-Proto"
	XForwardedProtocol = "X-Forwarded-Protocol"
	XForwardedSSL      = "X-Forwarded-Ssl"
)

func init() {
	server.RegisterConfigurator("/api", func(g lars.IRouteGroup) {
		g.Get("/ip", GetIP)
		g.Get("/uuid", GetUUID)
		g.Get("/headers", GetHeaders)
		g.Get("/user-agent", GetUserAgent)

		g.Get("/get", Anything)
		g.Post("/post", Anything)
		g.Put("/put", Anything)
		g.Patch("/patch", Anything)
		g.Delete("/delete", Anything)

		Any(g, "/anything", Anything)
		Any(g, "/anything/*", Anything)

		g.Get("/gzip", middleware.Gzip, Anything)

		g.Get("/redirect/:n", Redirect)
		g.Get("/stream/:n", Stream)
	})
}

func Any(g lars.IRouteGroup, p string, h ...lars.Handler) {
	g.Get(p, h...)
	g.Post(p, h...)
	g.Put(p, h...)
	g.Delete(p, h...)
	g.Patch(p, h...)
	g.Trace(p, h...)
}

func GetIP(c *server.Context) error {
	d := struct {
		Origin string `json:"origin"`
	}{
		Origin: c.ClientIP(),
	}

	return c.JSON(StatusOK, d)
}

func GetUUID(c *server.Context) error {
	v, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	d := struct {
		UUID uuid.UUID `json:"uuid"`
	}{
		UUID: v,
	}

	return c.JSON(StatusOK, d)
}

func GetHeaders(c *server.Context) error {
	return c.JSON(StatusOK, c.Request().Header)
}

func GetUserAgent(c *server.Context) error {
	d := struct {
		UserAgent string `json:"user-agent"`
	}{
		UserAgent: c.Request().UserAgent(),
	}

	return c.JSON(StatusOK, d)
}

func Anything(c *server.Context) error {
	return c.JSON(StatusOK, NewResponse(c))
}

func Redirect(c *server.Context) error {
	nString := c.Param("n")
	n, err := strconv.Atoi(nString)
	if err != nil {
		return err
	}

	if n == 1 {
		http.Redirect(c.Response(), c.Request(), "/api/get", http.StatusFound)
	} else {
		n--
		http.Redirect(c.Response(), c.Request(), fmt.Sprintf("/api/redirect/%d", n), http.StatusFound)
	}

	return nil
}

func Stream(c *server.Context) error {
	nString := c.Param("n")
	n, err := strconv.Atoi(nString)
	if err != nil {
		return err
	}

	if n > 100 {
		n = 100
	}

	type response struct {
		ID int `json:"id"`
		*BaseResponse
	}

	c.Response().Header().Set("Content-Type", "application/stream+json")
	r := response{
		BaseResponse: NewResponse(c),
		ID:           0,
	}

	c.Stream(func(w io.Writer) bool {
		rnd := time.Duration(rand.Int63n(int64(time.Second * 2)))
		c.Log().WithField("delay", rnd).Debug("sleeping")
		time.Sleep(time.Duration(rnd))

		enc := json.NewEncoder(w)
		err = enc.Encode(r)
		if err != nil {
			return false
		}
		r.ID++
		return r.ID < n
	})

	return err
}

type BaseResponse struct {
	URL     string              `json:"url"`
	Args    map[string][]string `json:"args"`
	Form    map[string][]string `json:"form"`
	Data    []byte              `json:"data"`
	Origin  string              `json:"origin"`
	Headers map[string][]string `json:"headers"`
	Files   map[string][][]byte `json:"files"`
	JSON    interface{}         `json:"json"`
	Method  string              `json:"method"`
}

func NewResponse(c lars.Context) *BaseResponse {
	req := c.Request()
	data := getData(req)

	var _json interface{}
	_ = json.Unmarshal(data, &_json)

	return &BaseResponse{
		URL:     getURL(req),
		Args:    getArgs(req),
		Form:    getForm(req),
		Data:    data,
		Origin:  c.ClientIP(),
		Headers: req.Header,
		Files:   getFiles(req),
		JSON:    _json,
		Method:  req.Method,
	}
}

func getURL(req *http.Request) string {
	scheme := req.Header.Get(XForwardedProto)
	if scheme == "" {
		scheme = req.Header.Get(XForwardedProtocol)
	}
	if scheme == "" && req.Header.Get(XForwardedSSL) == "on" {
		scheme = "https"
	}
	if scheme == "" {
		scheme = "http"
	}

	u := *req.URL
	u.Scheme = scheme
	u.Host = req.Host
	return u.String()
}

func getArgs(req *http.Request) map[string][]string {
	return req.URL.Query()
}

func getForm(req *http.Request) map[string][]string {
	_ = req.ParseMultipartForm(1 << 16)

	if req.MultipartForm != nil {
		return req.MultipartForm.Value
	}

	return req.Form
}

func getData(req *http.Request) []byte {
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(req.Body)

	// replace body with a copy that can be read again
	_ = req.Body.Close()
	req.Body = ioutil.NopCloser(&buf)

	return buf.Bytes()
}

func getFiles(req *http.Request) map[string][][]byte {
	_ = req.ParseMultipartForm(1 << 32)

	if req.MultipartForm == nil {
		return nil
	}

	rv := make(map[string][][]byte, len(req.MultipartForm.File))

	for k, v := range req.MultipartForm.File {
		rv[k] = make([][]byte, len(v))
		for idx, f := range v {
			fp, _ := f.Open()
			var buf bytes.Buffer
			_, _ = buf.ReadFrom(fp)
			_ = fp.Close()
			rv[k][idx] = buf.Bytes()
		}
	}

	return rv
}
