package admin

import (
	"goblog/models"
	"goblog/util"
	"strconv"
	"strings"
)

type TagController struct {
	baseController
}

func (this *TagController) Index() {
	act := this.GetString("act")
	switch act {
	case "batch":
		this.batch()
	default:
		this.list()
	}

}

func (this *TagController) list() {
	var page int
	var pagesize int = 10
	var list []*models.Tag
	var tag models.Tag

	if page, _ := this.GetInt("page"); page < 1 {
		page = 1
	}
	offset := (page - 1) * pagesize
	count, _ := tag.Query().Count()
	if count > 0 {
		tag.Query().OrderBy("-count").Limit(pagesize, offset).All(&list)
	}

	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pagebar"] = util.NewPager(page, int(count), pagesize, "/admin/tag", true).ToString()
	this.Display("tag/list")
}

//批量操作
func (this *TagController) batch() {
	ids := this.GetStrings("ids[]")
	op := this.GetString("op")

	idarr := make([]int, 0)
	for _, v := range ids {
		if id, _ := strconv.Atoi(v); id > 0 {
			idarr = append(idarr, id)
		}
	}
	var tag models.Tag
	switch op {
	case "upcount": //更新统计
		for _, id := range idarr {
			tag := models.Tag{Id, id}
			if tag.Read() == nil {
				tag.UpCount()
			}
		}
	case "merge": //合并到
		toname := strings.TrimSpace(this.GetString("toname"))
		if toname != "" && len(idarr) > 0 {
			tag := new(models.Tag)
			tag.Name = toname
			if tag.Read("name") != nil {
				tag.Count = 0
				tag.Insert()
			}
			for _, id := range idarr {
				obj := models.Tag{Id: id}
				if obj.Read() == nil {
					obj.MergeTo(tag)
					obj.Delete()
				}
			}
			tag.UpCount()
		}
	case "delete": //删除
		for id, _ := range idarr {
			tag := models.Tag{Id, id}
			if tag.Read() == nil {
				tag.Delete()
			}
		}
	}
	this.Redirect("/admin/tag", 302)
}