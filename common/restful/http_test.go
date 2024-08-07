package restful

import (
	"net/url"
	"testing"
)

func Test_url(t *testing.T) {
	var uri = "http://abc:def@ww.a.com/b/c?a=1&b=3#111"
	u, err := url.Parse(uri)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", *u)
	t.Logf("%v", u.Hostname())
	t.Logf(u.Query().Encode())
	t.Logf("%s://%s%s", u.Scheme, u.Host, u.Path)
}

func TestClient_Get(t *testing.T) {
	var uri = "http://abc:def@ww.a.com/b/c?a=1&b=3#111"
	NewClient().Get(uri, nil)
}
