package pshgo

import (
	"fmt"
	"net/url"
	"strings"
)

type Relationships map[string][]Relationship

type Relationship struct {
	Cluster  string     `json:"cluster"`
	Fragment string     `json:"fragment"`
	Host     string     `json:"host"`
	Hostname string     `json:"hostname"`
	IP       string     `json:"ip"`
	Password string     `json:"password"`
	Path     string     `json:"path"`
	Port     int        `json:"port"`
	Public   bool       `json:"public"`
	Query    JSONObject `json:"query"`
	Rel      string     `json:"rel"`
	Scheme   string     `json:"scheme"`
	Service  string     `json:"service"`
	SSL      JSONObject `json:"ssl"`
	Type     string     `json:"type"`
	Username string     `json:"username"`
}

func (r Relationship) URL(user, query bool) string {
	var b strings.Builder
	b.WriteString(r.Scheme)
	b.WriteString("://")

	if user && r.Username != "" {
		b.WriteString(r.Username)
		if r.Password != "" {
			b.WriteString(":")
			b.WriteString(r.Password)
		}
		b.WriteString("@")
	}

	b.WriteString(r.Host)
	if r.Port > 0 {
		_, _ = fmt.Fprintf(&b, ":%d", r.Port)
	}

	if r.Path != "" {
		b.WriteString("/")
		b.WriteString(r.Path)
	}

	if query {
		values := make(url.Values, len(r.Query))
		for k, v := range r.Query {
			switch v := v.(type) {
			case []string:
				values[k] = v
			case JSONArray:
				values[k] = make([]string, len(v))
				for idx, val := range v {
					values[k][idx] = fmt.Sprint(val)
				}
			case string:
				values[k] = []string{v}
			case *string:
				values[k] = []string{*v}
			default:
				values[k] = []string{fmt.Sprint(v)}
			}
		}

		if len(values) > 0 {
			if r.Path == "" {
				b.WriteString("/")
			}

			b.WriteString("?")
			b.WriteString(values.Encode())
		}
	}

	return b.String()
}
