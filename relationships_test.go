package pshgo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/demosdemon/pshgo"
)

func strptr(s string) *string { return &s }

func TestRelationship_URL(t *testing.T) {
	cases := []struct {
		name        string
		rel         Relationship
		user, query bool
		want        string
	}{
		{
			name: "zero",
			want: "://",
		},
		{
			name: "no user",
			rel: Relationship{
				Scheme:   "https",
				Username: "main",
				Password: "main",
				Host:     "test.example",
			},
			user:  false,
			query: false,
			want:  "https://test.example",
		},
		{
			name: "user no password",
			rel: Relationship{
				Scheme:   "https",
				Username: "main",
				Host:     "test.example",
			},
			user:  true,
			query: false,
			want:  "https://main@test.example",
		},
		{
			name: "user password",
			rel: Relationship{
				Scheme:   "https",
				Username: "main",
				Password: "main",
				Host:     "test.example",
			},
			user:  true,
			query: false,
			want:  "https://main:main@test.example",
		},
		{
			name: "port",
			rel: Relationship{
				Scheme:   "https",
				Username: "main",
				Host:     "test.example",
				Port:     3030,
			},
			user:  false,
			query: false,
			want:  "https://test.example:3030",
		},
		{
			name: "path",
			rel: Relationship{
				Scheme:   "https",
				Username: "main",
				Host:     "test2.example",
				Port:     3040,
				Path:     "v2",
			},
			user:  false,
			query: false,
			want:  "https://test2.example:3040/v2",
		},
		{
			name: "empty query",
			rel: Relationship{
				Scheme: "https",
				Host:   "test.example",
				Path:   "v1",
			},
			user:  true,
			query: true,
			want:  "https://test.example/v1",
		},
		{
			name: "string array query",
			rel: Relationship{
				Scheme: "postgresql",
				Host:   "test.example",
				Path:   "main",
				Query: JSONObject{
					"sslmode": []string{"disable"},
				},
			},
			user:  false,
			query: true,
			want:  "postgresql://test.example/main?sslmode=disable",
		},
		{
			name: "json array query",
			rel: Relationship{
				Scheme: "http",
				Host:   "example.com",
				Path:   "test",
				Query: JSONObject{
					"example": JSONArray{
						true,
						false,
						"test",
					},
				},
			},
			user:  false,
			query: true,
			want:  "http://example.com/test?example=true&example=false&example=test",
		},
		{
			name: "string query",
			rel: Relationship{
				Scheme: "mongodb",
				Host:   "cloud.local",
				Path:   "main",
				Query: JSONObject{
					"is_master": "true",
				},
			},
			user:  false,
			query: true,
			want:  "mongodb://cloud.local/main?is_master=true",
		},
		{
			name: "string ptr query",
			rel: Relationship{
				Scheme: "http",
				Host:   "test.example",
				Query: JSONObject{
					"test": strptr("foobar"),
				},
			},
			user:  false,
			query: true,
			want:  "http://test.example/?test=foobar",
		},
		{
			name: "default query",
			rel: Relationship{
				Scheme: "http",
				Host:   "test.example",
				Query: JSONObject{
					"is_master": true,
				},
			},
			user:  false,
			query: true,
			want:  "http://test.example/?is_master=true",
		},
	}

	t.Parallel()
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			s := c.rel.URL(c.user, c.query)
			assert.Equal(t, c.want, s)
		})
	}
}
