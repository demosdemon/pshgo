package pshgo_test

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"

	. "github.com/demosdemon/pshgo"
)

type encoder interface {
	Encode(v interface{}) error
}

type decoder interface {
	Decode(v interface{}) error
}

func jsonEncoder(w io.Writer) encoder {
	return json.NewEncoder(w)
}

func yamlEncoder(w io.Writer) encoder {
	return yaml.NewEncoder(w)
}

func jsonDecoder(r io.Reader) decoder {
	return json.NewDecoder(r)
}

func yamlDecoder(r io.Reader) decoder {
	return yaml.NewDecoder(r)
}

func TestDuration_MarshalText(t *testing.T) {
	cases := []struct {
		name      string
		duration  Duration
		mkEncoder func(w io.Writer) encoder
		wantW     string
		wantErr   bool
	}{
		{
			name:      "json",
			duration:  Duration{Duration: time.Minute * 3},
			mkEncoder: jsonEncoder,
			wantW:     "\"3m0s\"\n",
		},
		{
			name:      "yaml",
			duration:  Duration{Duration: time.Second * 30},
			mkEncoder: yamlEncoder,
			wantW:     "30s\n",
		},
	}

	t.Parallel()
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			var buf bytes.Buffer
			enc := c.mkEncoder(&buf)
			err := enc.Encode(c.duration)
			assert.True(t, (err != nil) == c.wantErr)
			assert.Equal(t, c.wantW, string(buf.Bytes()))
		})
	}
}

func TestDuration_UnmarshalText(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		mkDecoder func(r io.Reader) decoder
		want      Duration
		wantErr   bool
	}{
		{
			name:      "json",
			input:     `"30m20s"`,
			mkDecoder: jsonDecoder,
			want:      Duration{Duration: time.Minute*30 + time.Second*20},
		},
		{
			name:      "yaml",
			input:     `3h0s`,
			mkDecoder: yamlDecoder,
			want:      Duration{Duration: time.Hour * 3},
		},
	}

	t.Parallel()
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			buf := bytes.NewReader([]byte(c.input))
			dec := c.mkDecoder(buf)
			var got Duration
			err := dec.Decode(&got)
			assert.True(t, (err != nil) == c.wantErr)
			assert.Equal(t, c.want, got)
		})
	}
}
