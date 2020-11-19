package omaha

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
)

func (r *Request) Send(url string) (*Response, error) {
	b := &bytes.Buffer{}
	if err := xml.NewEncoder(b).Encode(r); err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "text/xml", b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseResponse(resp.Body)
}

func ParseResponse(r io.Reader) (*Response, error) {
	v := &Response{}
	if err := xml.NewDecoder(r).Decode(v); err != nil {
		return nil, err
	}
	return v, nil
}

func (r *Response) GetApp(id string) *AppResponse {
	for i := range r.Apps {
		a := &r.Apps[i]
		if a.ID == id {
			return a
		}
	}
	return nil
}

func (u *UpdateCheckResponse) PayloadURLs(p *Package) ([]string, error) {
	ref, err := url.Parse(p.Name)
	if err != nil {
		return nil, err
	}

	urls := make([]string, len(u.URLs))
	for i, v := range u.URLs {
		base, err := url.Parse(v.CodeBase)
		if err != nil {
			return nil, err
		}
		urls[i] = base.ResolveReference(ref).String()
	}
	return urls, nil
}

func (m *Manifest) GetAction(event string) *Action {
	for i := range m.Actions {
		a := &m.Actions[i]
		if a.Event == event {
			return a
		}
	}
	return nil
}
