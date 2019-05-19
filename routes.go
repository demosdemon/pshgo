package pshgo

import (
	"encoding/json"
	"net/url"

	"github.com/sirupsen/logrus"
)

type (
	Routes map[url.URL]Route
)

func (r Routes) MarshalJSON() ([]byte, error) {
	logrus.Trace("Routes.MarshalJSON")
	intermediate := make(map[string]Route, len(r))

	for kURL, v := range r {
		intermediate[kURL.String()] = v
	}

	return json.Marshal(intermediate)
}

func (r *Routes) UnmarshalJSON(text []byte) error {
	logrus.Trace("Routes.UnmarshalJSON")
	var intermediate map[string]Route
	err := json.Unmarshal(text, &intermediate)
	if err != nil {
		return err
	}

	*r = make(Routes, len(intermediate))
	for k, v := range intermediate {
		kURL, err := url.Parse(k)
		if err != nil {
			return err
		}

		(*r)[*kURL] = v
	}

	return nil
}
