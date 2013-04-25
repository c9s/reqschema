package reqschema
import "testing"
import "net/http"
import "net/url"
import "encoding/json"

func TestRequestSchemaWithJsonDecoder(t *testing.T) {
	resource := map[string]string {
		"name": "John",
	}
	jsonBytes, err := json.Marshal(resource)

	if err != nil {
		t.Error(err)
	}

	req := &http.Request{Method: "GET"}
	req.URL, _ = url.Parse("http://www.google.com/search?resource=" + string(jsonBytes) )
	if resource := req.FormValue("resource"); resource == "" {
		t.Error("resource param not found")
	}

	type ResourceRequest struct {
		Resource string `param:"resource" decode:"json"`
	}

	rschema := Create(req, &ResourceRequest{})
	resourceData, err := rschema.Get("Resource")

	t.Log(resourceData)

	if _, ok := resourceData.(map[string]interface{})[ "name" ] ; ! ok {
		t.Error("name not found")
	}

	if err != nil {
		t.Error(err)
	}
	_ = resourceData
}


func TestRequestSchema(t *testing.T) {
	req := &http.Request{Method: "GET"}
	req.URL, _ = url.Parse("http://www.google.com/search?q=foo&id=13&id64=64")
	if q := req.FormValue("q"); q != "foo" {
		t.Errorf(`req.FormValue("q") = %q, want "foo"`, q)
	}

	type GoogleSearchRequest struct {
		Query string `param:"q"`
		Id64  int64  `param:"id64"`
		Id    int    `param:"id"`
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

	var id2 int64 = 0
	found, err := rschema.GetParamTo("id",&id2)
	if err != nil {
		t.Error(err)
	}
	if ! found {
		t.Error("id not found")
	}
	t.Log(id2)
	if id2 != 13 {
		t.Error("%d is not equal to 13", id2)
	}


	_ = rschema
	_ = query
	_ = err
}



