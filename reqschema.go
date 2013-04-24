package reqschema
import "reflect"
import "net/http"
import "fmt"
import "strconv"


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

func (self * RequestSchema) Get(fieldName string) (interface{}, error) {
	// Get The Field By Reflect
	fieldType, found := self.TypeOfSchema.FieldByName(fieldName)

	if ! found {
		return nil, nil
	}

	// valueType := self.ValueOfSchema.FieldByName(fieldName)

	if self.Request.Form == nil {
		// parse form
		self.Request.ParseMultipartForm(defaultMaxMemory)
	}

	fmt.Println(fieldType);
	fmt.Println(fieldType.Name);
	fmt.Println(fieldType.Type);

	// found value in form
	if requestValues , ok := self.Request.Form[ fieldName ]; ok && len(requestValues) > 0 {
		requestValue  := requestValues[0]

		var returnValue interface{}
		var err error

		switch fieldType.Type.Name() {
		case "int":
			returnValue , err = strconv.ParseInt(requestValue, 0, 0)
		case "int8":
			returnValue , err = strconv.ParseInt(requestValue, 0, 8)
		case "int32":
			returnValue , err = strconv.ParseInt(requestValue, 0, 32)
		case "int64":
			returnValue , err = strconv.ParseInt(requestValue, 0, 64)
		case "float32":
			returnValue , err = strconv.ParseFloat(requestValue, 32)
		case "float64":
			returnValue , err = strconv.ParseFloat(requestValue, 64)
		case "string":
			returnValue = requestValue
			err = nil
		}
		if err != nil {
			return nil, err
		}
		return returnValue, nil
	}
	return nil, nil
}

