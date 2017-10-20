package admin

import (
	"goblog/models"
	"goblog/util"
	"strconv"
	"strings"
)

type AccountController struct {
	baseController
}

func (this *AccountController) Login() {
	if this.userid > 0 {
		this.Redirect("/admin", 302)
	}
	if this.GetString("dosubmit") == "yes" {
		account := strings.TrimSpace(this.GetString("account"))
		password := strings.TrimSpace(this.GetString("password"))
		remember := this.GetString("remember")
		if account != "" && password != "" {
			var user models.User
			user.UserName = account
			if user.Read("user_name") != nil || user.Password != util.Md5([]byte(password)) {
				this.Data["errmsg"] = "帐号或密码错误"
			} else if user.Active == 0 {
				this.Data["errmsg"] = "该帐号未激活"
			}

			user.LoginCount += 1
			user.LastIp = this.getClientIp()
			user.LastLogin = this.getTime()

			user.Update()
			authkey := util.Md5([]byte(this.getClientIp() + "|" + user.Password))
			if remember == "yes" {
				this.Ctx.SetCookie("auth", strconv.Itoa(user.Id)+"|"+authkey, 7*86400)
			} else {
				this.Ctx.SetCookie("auth", strconv.Itoa(user.Id)+"|"+authkey)
			}

			this.Redirect("/admin", 302)
		}
	}
	this.TplName = this.moduleName + "/account/login.html"
}

func (this *AccountController) Logout() {
	this.Ctx.SetCookie("auth", "")
	this.Redirect("/admin/login", 302)
}

//todo
func (this *AccountController) Profile() {
	user := models.User{Id: this.userid}

	if err := user.Read(); err != nil {
		this.showmsg(err.Error())
	}
	//todo
	if this.isPost() {

	}
	this.Data["user"] = user
	this.Display()
}
