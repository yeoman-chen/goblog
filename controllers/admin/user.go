package admin

import (
	"github.com/astaxie/beego/validation"
	"goblog/models"
	"goblog/util"
	"strings"
)

type UserController struct {
	baseController
}

func (this *UserController) List() {
	var page int
	var pagesize int = 10
	var list []*models.User
	var user models.User

	if page, _ = this.GetInt("page"); page < 1 {
		page = 1
	}
	offset := (page - 1) * pagesize
	count, _ := user.Query().Count()
	if count > 0 {
		user.Query().OrderBy("-id").Limit(pagesize, offset).All(&list)
	}

	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pagebar"] = util.NewPager(page, int(count), pagesize, "/admin/user/list", true).ToString()
	this.Data["pagebar"] = util.NewPager(page, int(count), pagesize, "/admin/user/list", true).ToString()
	this.Display()
}

func (this *UserController) Add() {
	input := make(map[string]string)
	errmsg := make(map[string]string)

	if this.isPost() {
		username := strings.TrimSpace(this.GetString("username"))
		password := strings.TrimSpace(this.GetString("password"))
		password2 := strings.TrimSpace(this.GetString("password2"))
		email := strings.TrimSpace(this.GetString("email"))
		active, _ := this.GetInt("active")

		input["username"] = username
		input["password"] = password
		input["password2"] = password2
		input["email"] = email

		valid := validation.Validation{}
		if v := valid.Required(username, "username"); !v.Ok {
			errmsg["username"] = "请输入用户名"
		} else if v := valid.MaxSize(username, 15, "username"); !v.Ok {
			errmsg["username"] = "用户名长度不能大于15个字符"
		} else if v := valid.MinSize(username, 3, "username"); !v.Ok {
			errmsg["username"] = "用户名长度不能小于3个字符"
		}

		if v := valid.Required(password, "password"); !v.Ok {
			errmsg["password"] = "请输入密码"
		}
		if v := valid.Required(password2, "password2"); !v.Ok {
			errmsg["password2"] = "请输入确认密码"
		} else if password != password2 {
			errmsg["password2"] = "两次输入的密码不一致"
		}

		if v := valid.Required(email, "email"); !v.Ok {
			errmsg["email"] = "请输入邮箱"
		} else if v := valid.Email(email, "email"); !v.Ok {
			errmsg["email"] = "Email无效"
		}
		if active > 1 {
			active = 1
		} else {
			active = 0
		}

		if len(errmsg) == 0 {
			var user models.User
			user.UserName = username
			user.Password = util.Md5([]byte(password))
			user.Email = email
			user.Active = int8(active)
			if err := user.Insert(); err != nil {
				this.showmsg(err.Error())
			}
			this.Redirect("/admin/user/list", 302)
		}
	}
	this.Data["input"] = input
	this.Data["errmsg"] = errmsg
	this.Display()
}

func (this *UserController) Edit() {
	id, _ := this.GetInt("id")
	user := models.User{Id: id}
	if err := user.Read(); err != nil {
		this.showmsg("用户不存在")
	}
	errmsg := make(map[string]string)

	if this.isPost() {
		password := strings.TrimSpace(this.GetString("password"))
		password2 := strings.TrimSpace(this.GetString("password2"))
		email := strings.TrimSpace(this.GetString("email"))
		active, _ := this.GetInt("active")
		valid := validation.Validation{}

		if password != "" {
			if v := valid.Required(password2, "password2"); !v.Ok {
				errmsg["password2"] = "请再次输入密码"
			} else if password != password2 {
				errmsg["password2"] = "两次输入的密码不一致"
			} else {
				user.Password = util.Md5([]byte(password))
			}
		}

		if v := valid.Required(email, "email"); !v.Ok {
			errmsg["email"] = "请输入email"
		} else if v := valid.Email(email, "email"); !v.Ok {
			errmsg["email"] = "无效的email"
		} else {
			user.Email = email
		}
		if active > 1 {
			user.Active = 1
		} else {
			user.Active = 0
		}

		if len(errmsg) == 0 {
			user.Update()
			this.Redirect("/admin/user/list", 302)
		}
	}
	this.Data["user"] = user
	this.Data["errmsg"] = errmsg
	this.Display()
}

func (this *UserController) Delete() {
	id, _ := this.GetInt("id")
	if id == 1 {
		this.showmsg("不能删除id为1的用户")
	}
	user := models.User{Id: id}
	if user.Read() == nil {
		user.Delete()
	}
	this.Redirect("/admin/user/list", 302)
}
