package restful

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"github.com/pkg/errors"
)

const (
	contentTypeHeader = "Content-Type"
	acceptHeader      = "Accept"
	jsonType          = "application/json"
)

type HTTPRequest struct {
	URL    string
	Method string
	Query  url.Values
	Header map[string]string
	Data   io.Reader
}

func (c *Client) Get(reqURL string, resp interface{}) error {
	httpRequest, err := c.buildHTTPRequest(http.MethodGet, reqURL, nil, nil)
	if err != nil {
		return err
	}
	return c.DoRequest(httpRequest, resp)
}

func (c *Client) PostJson(reqURL string, data io.Reader, resp interface{}) error {
	httpRequest, err := c.buildHTTPRequest(http.MethodPost, reqURL, map[string]string{contentTypeHeader: jsonType}, data)
	if err != nil {
		return err
	}
	return c.DoRequest(httpRequest, resp)
}

func (c *Client) PutJson(reqURL string, data io.Reader, resp interface{}) error {
	httpRequest, err := c.buildHTTPRequest(http.MethodPut, reqURL, map[string]string{contentTypeHeader: jsonType}, data)
	if err != nil {
		return err
	}
	return c.DoRequest(httpRequest, resp)
}

func (c *Client) Delete(reqURL string, data io.Reader, resp interface{}) error {
	httpRequest, err := c.buildHTTPRequest(http.MethodDelete, reqURL, map[string]string{contentTypeHeader: jsonType}, data)
	if err != nil {
		return err
	}
	return c.DoRequest(httpRequest, resp)
}

func (c *Client) buildHTTPRequest(method, reqURL string, header map[string]string, data io.Reader) (HTTPRequest, error) {
	uri, err := url.ParseRequestURI(reqURL)
	if err != nil {
		return HTTPRequest{}, fmt.Errorf("invalid request url")
	}

	return HTTPRequest{
		URL:    fmt.Sprintf("%s://%s%s", uri.Scheme, uri.Host, uri.Path),
		Method: method,
		Header: header,
		Query:  uri.Query(),
		Data:   data,
	}, nil
}

func (c *Client) DoRequest(req HTTPRequest, data interface{}) error {
	defer func() {
		if err := recover(); err != nil {
			err = errors.WithStack(fmt.Errorf("%v", err))
		}
	}()

	tof := reflect.TypeOf(data)
	if tof.Kind() != reflect.Ptr {
		return errors.New("invalid data type")
	}
	switch req.Method {
	case http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete:
		// skip
	default:
		panic("unrecognized http request Method")
	}

	request, err := http.NewRequest(req.Method, req.URL, req.Data)
	if err != nil {
		return fmt.Errorf("http.NewRequest happend error: %s", err)
	}

	if req.Header != nil {
		for k, v := range req.Header {
			request.Header.Set(k, v)
		}
	}

	if len(req.Query) != 0 {
		q := request.URL.Query()
		for k, v := range req.Query {
			if !q.Has(k) {
				q.Set(k, v[0])
			}
		}
		request.URL.RawQuery = q.Encode()
	}

	if c.username != "" && c.password != "" {
		request.SetBasicAuth(c.username, c.password)
	}

	response, err := c.client.Do(request)
	if err != nil {
		return fmt.Errorf("client.Do happend error: %s", err)
	}

	if response.StatusCode >= 300 {
		return fmt.Errorf("unexpect status, code: %d, message: %s", response.StatusCode, response.Status)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll happend error: %s", err)
	}

	switch tof.Elem().Kind() {
	case reflect.String:
		reflect.ValueOf(data).Elem().SetString(string(body))
	default:
		if len(body) > 0 && data != nil {
			if err := json.Unmarshal(body, &data); err != nil {
				return fmt.Errorf("json.Unmarshal happend error: %s", err)
			}
		}
	}
	return nil
}
