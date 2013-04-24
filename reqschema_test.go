package reqschema
import "testing"
import "net/http"
import "net/url"



func TestRequestSchema(t *testing.T) {
	req := &http.Request{Method: "GET"}
	req.URL, _ = url.Parse("http://www.google.com/search?q=foo&id=13")
	if q := req.FormValue("q"); q != "foo" {
		t.Errorf(`req.FormValue("q") = %q, want "foo"`, q)
	}

	type GoogleSearchRequest struct {
		Query string `field:q`
		Id64  int64  `field:id64`
		Id    int    `field:id`
	}

	rschema := Create(req, &GoogleSearchRequest{})
	rschema.Get("Query")
	rschema.Get("Id")
	_ = rschema
}



