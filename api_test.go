package namecheap

import (
	"testing"
)

func makeClient(t *testing.T) *Client {
	client, err := NewClient("user", "apiuser", "secret", "128.0.0.1", true)

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	return client
}

func TestClient_NewRequest(t *testing.T) {
	c := makeClient(t)

	body := map[string]interface{}{
		"foo": "bar",
		"baz": "bar",
	}
	req, err := c.NewRequest(body)
	if err != nil {
		t.Fatalf("bad: %v", err)
	}

	if req.URL.String() != "https://api.sandbox.namecheap.com/xml.response" {
		t.Fatalf("bad base url: %v", req.URL.String())
	}

	if req.Method != "POST" {
		t.Fatalf("bad method: %v", req.Method)
	}
}
