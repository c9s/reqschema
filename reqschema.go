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

	if r.Form == nil {
		// parse form
		r.ParseMultipartForm(defaultMaxMemory)
	}
	return &RequestSchema{ Request: r, Schema: schema, ValueOfSchema: valueOf, TypeOfSchema: typeOf }
}

func (self * RequestSchema) GetTo(name string, value interface{}) error {
	// get the value from pointer (Elem())
	valueValue := reflect.ValueOf(value).Elem()
	valueType := valueValue.Type()

	if ! valueValue.CanSet() {
		panic("the value can not be set.")
	}

	// valueType.Type.Name()
	if requestValues , ok := self.Request.Form[ name ]; ok && len(requestValues) > 0 {
		newValue, err := parseStringByType(requestValues[0], valueType)
		if err != nil {
			return err
		}
		switch t := newValue.(type) {
		default:
			panic( fmt.Sprintf("unsupported type %s" , t))
		case int32:
			valueValue.SetInt( int64(newValue.(int32)) )
		case int64:
			valueValue.SetInt(newValue.(int64))
		case int:
			valueValue.SetInt( int64(newValue.(int)) )
		case string:
			valueValue.SetString(newValue.(string))
		case float64:
			valueValue.SetFloat(newValue.(float64))
		case float32:
			valueValue.SetFloat( float64(newValue.(float32)) )
		}
	}
	return nil
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




func (self * RequestSchema) GetField(name string) *reflect.StructField {
	field, found := self.TypeOfSchema.FieldByName(name)
	if found {
		return &field
	}
	return nil
}

func (self * RequestSchema) GetFieldName(name string) string {
	field := self.GetField(name)
	if field != nil {
		fieldName := field.Tag.Get("field")
		if fieldName != "" {
			return fieldName
		}
	}
	return ""
}

func (self * RequestSchema) Has(name string) bool {
	fieldName := self.GetFieldName(name)
	if requestValues, ok := self.Request.Form[ fieldName ]; ok && len(requestValues) > 0 {
		return true
	}
	return false
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

	// found value in form
	if requestValues , ok := self.Request.Form[ fieldName ]; ok && len(requestValues) > 0 {
		return parseStringByType(requestValues[0], field.Type)
	}
	return nil, nil
}

