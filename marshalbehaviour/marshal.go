package marshalbehaviour

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

/*
** JSON object example for request **
{
    "field1" : "apapun",
	"field2" : "apapun",
	"field3" : "apapun or omit",
    "field4" : "apapun or omit"
}
*/

type ObjToMarshal struct {
	// required, field will be usual string
	Field1 string `json:"field1" validate:"required"`
	// required, field will be pointer string
	Field2 *string `json:"field2" validate:"required"`
	// field will be usual string, if omitted, field will be an EMPTY string (default value for the specified type)
	Field3 string `json:"field3"`
	// field will be pointer string, if omitted, field will be a NIL pointer
	Field4 *string `json:"field4"`
}

func MarshalTest(c echo.Context) error {
	req := &ObjToMarshal{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, "error bind: "+err.Error())
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, "error validate: "+err.Error())
	}

	str := fmt.Sprintf("field type: field1 = '%+v', field2 = '%+v', field3 = '%+v', field4 = '%+v'", req.Field1, req.Field2, req.Field3, req.Field4)

	return c.JSON(http.StatusOK, str)
}
