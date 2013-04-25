package reqschema
import "reflect"
import "net/http"
import "strconv"
import "fmt"
import "encoding/json"

const defaultMaxMemory = 32 << 20 // 32MB

/*
Synopsis

import "reqschema"

type UserAuthRequest struct {
	UserId	 int `param:"id"`
	UserName string `param:"username"`
	Password string `param:"password"`
	Resource string `param:"resource" decode:"json"`
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



/*
Get the form value and convert the form value string to the type of given value,
then use reflect to set the value back.

The given value must be addressable, so that we can set value via reflect.

Returns (found, error)
*/
func (self * RequestSchema) GetParamTo(paramName string, value interface{}) (bool,error) {
	// get the value from pointer (Elem())
	valueValue := reflect.ValueOf(value).Elem()
	valueType := valueValue.Type()

	if ! valueValue.CanSet() {
		panic("the value can not be set.")
	}

	if requestValues , ok := self.Request.Form[ paramName ]; ok && len(requestValues) > 0 {
		newValue, err := parseStringByType(requestValues[0], valueType)
		if err != nil {
			return true, err
		}
		switch t := newValue.(type) {
		default:
			return true, fmt.Errorf( "unsupported type %s" , t)
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
		return true, nil
	}

	// value not found, should not be error
	return false, nil
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

func (self * RequestSchema) GetParamName(name string) string {
	field := self.GetField(name)
	if field != nil {
		paramName := field.Tag.Get("param")
		if paramName != "" {
			return paramName
		}
	}
	return ""
}

func (self * RequestSchema) Has(name string) bool {
	paramName := self.GetParamName(name)
	if requestValues, ok := self.Request.Form[ paramName ]; ok && len(requestValues) > 0 {
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

	paramName := field.Tag.Get("param")
	if paramName == "" {
		return nil, nil
	}

	// found value in form
	if requestValues , ok := self.Request.Form[ paramName ]; ok && len(requestValues) > 0 {
		value, err := parseStringByType(requestValues[0], field.Type)

		if err != nil {
			return nil, err
		}

		if field.Type.Name() == "string" {
			decodeType := field.Tag.Get("decode")
			if decodeType == "json" {
				var data interface{}
				err := json.Unmarshal( []byte(value.(string)) , &data)
				if err != nil {
					return nil, nil
				}
				return data, nil
			}
		}
		return value, nil
	}
	return nil, nil
}

