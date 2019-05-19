package pshgo

import (
	"fmt"
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
		b.WriteString("?")
		first := true
		for k, v := range r.Query {
			if !first {
				b.WriteString("&")
			}
			first = false
			_, _ = fmt.Fprintf(&b, "%s=%v", k, v)
		}
	}

	return b.String()
}
