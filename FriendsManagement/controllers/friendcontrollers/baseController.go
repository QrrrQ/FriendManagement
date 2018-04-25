package friendcontrollers

import (
	"FriendsManagement/utils/logger"
	"regexp"

	"github.com/astaxie/beego"
)

type errorJSON struct {
	Success   bool      `json:"success"`
	ErrorBody errorBody `json:"errorbody"`
}

type errorBody struct {
	StatusCode int    `json:"statuscode"`
	Message    string `json:"message"`
}

type successJSON struct {
	Success    bool        `json:"success"`
	ResultBody interface{} `json:"-"`
}

//BaseController ...
type BaseController struct {
	beego.Controller
}

func (c *BaseController) echoError(statusCode int, msg string) {
	c.Ctx.Output.Header("Server", "FriendsManagement V0.1")
	// msg := fmt.Sprintf(f, v)
	data := errorJSON{
		Success:   false,
		ErrorBody: errorBody{StatusCode: statusCode, Message: msg},
	}
	c.Data["json"] = data
	c.ServeJSON()
	logger.LogWithDepth(logger.LevelDebug, 4, `Response: { "statusCode": %d, "message": %s }`, statusCode, msg)
}

func (c *BaseController) echoSuccess(resultBody interface{}) {
	c.Ctx.Output.Header("Server", "FriendsManagement V0.1")
	data := successJSON{
		Success:    true,
		ResultBody: resultBody,
	}
	c.Data["json"] = data
	c.ServeJSON()
	logger.LogWithDepth(logger.LevelDebug, 4, `Response: { "success": true}`)
}

func (c *BaseController) echoData(resultBody interface{}) {
	c.Ctx.Output.Header("Server", "FriendsManagement V0.1")
	c.Data["json"] = resultBody
	c.ServeJSON()
	logger.LogWithDepth(logger.LevelDebug, 4, `Response: { "success": true}`)
}

func (c *BaseController) IsEmail(str string) bool {
	var b bool
	b, _ = regexp.MatchString(`^([a-z0-9_\.-]+)@([\da-z\.-]+)\.([a-z\.]{2,6})$`, str)
	return b
}
