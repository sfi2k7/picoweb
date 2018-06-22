package picoweb

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

const (
	contentType = "application/json"
)

type PicoClient struct {
	Headers     map[string]string
	ContentType string
	r           *http.Request
}

func (pc *PicoClient) addHeaders(r *http.Request) {
	if pc.Headers != nil {
		for k, v := range pc.Headers {
			r.Header.Add(k, v)
		}
	}
}

func (pc *PicoClient) Get(url string) ([]byte, error) {
	c := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	pc.addHeaders(req)
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (pc *PicoClient) Post(url string, data []byte) ([]byte, error) {
	c := &http.Client{}
	rdr := bytes.NewReader(data)
	req, err := http.NewRequest("POST", url, rdr)
	if err != nil {
		return nil, err
	}
	pc.addHeaders(req)
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (pc *PicoClient) Put(url string, data []byte) ([]byte, error) {
	c := &http.Client{}
	rdr := bytes.NewReader(data)
	req, err := http.NewRequest("PUT", url, rdr)
	if err != nil {
		return nil, err
	}
	pc.addHeaders(req)
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func NewClient() *PicoClient {
	return &PicoClient{}
}
