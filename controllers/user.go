package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"time"
)

type User struct {
	Ctx iris.Context
	Session *sessions.Session
}

func UserMVC(app *mvc.Application) {
	userSession := sessions.New(sessions.Config{
		Cookie:                      "xupt_session_id",
		Expires:                     24 * 5 * time.Hour,
	})
	app.Register(userSession.Start).Handle(new(User))
}

func (u *User) Get() ResponseFmtData {
	return ResponseFmtData{
		Message: "成功获取用户信息",
		Code:    0,
		Data:    nil,
	}
}

// PostLogin 登录
func (u *User) PostLogin() ResponseFmtData{
	u.Session.Set("xupt_session_id", "login_success" )
	return ResponseFmtData{
		Message: "登录成功",
		Code: 1,
		Data: nil,
	}
}