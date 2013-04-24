package reqschema
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
