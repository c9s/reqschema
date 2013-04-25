reqschema
=============
As to handle the form values from http requests, we need to use strconv to
parse string format into the type that we need.

reqschema let you define the request param schema for parsing form field values, so
you don't have to do the conversion by yourself.

for example, you may define a simple struct, which describes the schema of a request type:

```go
type UserAuthRequest struct {
    UserId   int    `param:"id"`
    UserName string `param:"username"`
    Password string `param:"password"`
    ResourceJson string `param:"resource" decode:"json"`
}
```

then pass the struct to create a request schema handle object:

```go
rschema := Create(req, &UserAuthRequest{})
resourceData, err := rschema.Get("Resource")

var id int
id, err := rschema.Get("UserId")
```


