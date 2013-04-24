package reqschema
import "reflect"
import "net/http"
import "strconv"
import "fmt"


const defaultMaxMemory = 32 << 20 // 32MB

/*
Synopsis

import "reqschema"

type UserAuthRequest struct {
	UserId	 `field:id		 type:integer`
	UserName `field:username type:string`
	Password `field:password type:string`
}

func userAuthRequestHandler( w http.ResponseWriter, r * http.Request)
{
	params := UserAuthRequest{}
	reqschema.Init(r, &params)
}
*/

type RequestSchema struct {
	Request *http.Request
	Schema  interface{}
	TypeOfSchema reflect.Type
	ValueOfSchema reflect.Value
}

func Create(r * http.Request, schema interface{}) (*RequestSchema) {
	valueOf := reflect.ValueOf(schema).Elem()
	typeOf := valueOf.Type()
	return &RequestSchema{ Request: r, Schema: schema, ValueOfSchema: valueOf, TypeOfSchema: typeOf }
}


func (self * RequestSchema) GetTo(fieldName string, value interface{}) {
	// get the value from pointer (Elem())
	valueType := reflect.TypeOf(value)

	fmt.Println( valueType )

	_ = valueType
}

func parseStringByType(value string, typeInfo reflect.Type) (interface{}, error) {
	switch typeInfo.Name() {
	case "int":
		return strconv.ParseInt(value, 0, 0)
	case "int8":
		return strconv.ParseInt(value, 0, 8)
	case "int32":
		return strconv.ParseInt(value, 0, 32)
	case "int64":
		return strconv.ParseInt(value, 0, 64)
	case "float32":
		return strconv.ParseFloat(value, 32)
	case "float64":
		return strconv.ParseFloat(value, 64)
	case "string":
		return value, nil
	}
	return nil, nil
}


func (self * RequestSchema) Get(name string) (interface{}, error) {
	// Get The Field By Reflect
	field, found := self.TypeOfSchema.FieldByName(name)

	if ! found {
		return nil, nil
	}

	fieldName := field.Tag.Get("field")
	if fieldName == "" {
		return nil, nil
	}

	// valueType := self.ValueOfSchema.FieldByName(name)

	if self.Request.Form == nil {
		// parse form
		self.Request.ParseMultipartForm(defaultMaxMemory)
	}

	// found value in form
	if requestValues , ok := self.Request.Form[ fieldName ]; ok && len(requestValues) > 0 {
		return parseStringByType(requestValues[0], field.Type)
	}
	return nil, nil
}

