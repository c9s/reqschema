package reqschema
import "testing"
import "net/http"
import "net/url"



func TestRequestSchema(t *testing.T) {
	req := &http.Request{Method: "GET"}
	req.URL, _ = url.Parse("http://www.google.com/search?q=foo&id=13&id64=64")
	if q := req.FormValue("q"); q != "foo" {
		t.Errorf(`req.FormValue("q") = %q, want "foo"`, q)
	}

	type GoogleSearchRequest struct {
		Query string `field:"q"`
		Id64  int64  `field:"id64"`
		Id    int    `field:"id"`
	}

	rschema := Create(req, &GoogleSearchRequest{})
	query, err := rschema.Get("Query")

	if err != nil {
		t.Error(err)
	}

	if query != "foo" {
		t.Log(query)
		t.Error("Not foo")
	}

	id , err := rschema.Get("Id")
	if err != nil {
		t.Error(err)
	}
	if id.(int64) != 13 {
		t.Errorf("%d is not 13",id)
	}


	id64 , err := rschema.Get("Id64")
	if err != nil {
		t.Error(err)
	}
	if id64.(int64) != 64 {
		t.Errorf("%d should be 64", id64)
	}
	_ = rschema
	_ = query
	_ = err
}



