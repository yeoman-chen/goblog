package admin

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"goblog/models"
	"goblog/util"
	"strconv"
	"strings"
	"time"
)

type ArticleController struct {
	baseController
}

func (this *ArticleController) List() {
	var (
		page       int
		pagesize   int = 10
		status     int
		offset     int
		list       []*models.Post
		post       models.Post
		searchtype string
		keyword    string
	)
	searchtype = this.GetString("searchtype")
	keyword = this.GetString("keyword")
	status, _ = this.GetInt("status")
	if page, _ := this.GetInt("page"); page < 1 {
		page = 1
	}
	offset = (page - 1) * pagesize

	query := post.Query().Filter("status", status)
	if keyword != "" {
		switch searchtype {
		case "title":
			query = query.Filter("title__icontains", keyword)
		case "author":
			query = query.Filter("author__icontains", keyword)
		case "tag":
			query = query.Filter("tag__icontains", keyword)
		}
	}

	count, _ := query.Count()
	if count > 0 {
		query.OrderBy("-is_top", "-id").Limit(pagesize, offset).All(&list)
	}

	for _, val := range list {
		fmt.Println(val.TagsLink)
	}
	this.Data["searchtype"] = searchtype
	this.Data["keyword"] = keyword
	this.Data["count_1"], _ = post.Query().Filter("status", 1).Count()
	this.Data["count_2"], _ = post.Query().Filter("status", 2).Count()
	this.Data["status"] = status
	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pagebar"] = util.NewPager(page, int(count), pagesize, fmt.Sprintf("/admin/article/list?status=%d&searchtype=%s&keyword=%s", status, searchtype, keyword), true).ToString()
	this.Display()
}

func (this *ArticleController) Add() {
	this.Data["posttime"] = this.getTime().Format("2006-01-02 15:04:05")
	this.Display()
}

func (this *ArticleController) Edit() {
	id, _ := this.GetInt("id")
	post := models.Post{Id: id}
	if post.Read() != nil {
		this.Abort("404")
	}
	post.Tags = strings.Trim(post.Tags, ",")
	this.Data["post"] = post
	this.Data["posttime"] = this.getTime().Format("2006-01-02 15:04:05")
	this.Display()
}

func (this *ArticleController) Save() {
	var (
		id      int    = 0
		title   string = strings.TrimSpace(this.GetString("title"))
		content string = this.GetString("content")
		color   string = strings.TrimSpace(this.GetString("color"))
		tags    string = strings.TrimSpace(this.GetString("tags"))
		urlname string = strings.TrimSpace(this.GetString("urlname"))
		timestr string = strings.TrimSpace(this.GetString("posttime"))
		status  int    = 0
		istop   int8   = 0
		urltype int8   = 0
		post    models.Post
	)
	if title == "" {
		this.showmsg("标题不能为空！")
	}
	id, _ = this.GetInt("id")
	status, _ = this.GetInt("status")

	if this.GetString("istop") == "1" {
		istop = 1
	}

	if this.GetString("urltype") == "1" {
		istop = 1
	}

	if status != 1 || status != 2 {
		status = 0
	}
	addtags := make([]string, 0)
	if tags != "" {
		tagarr := strings.Split(tags, ",")
		for _, v := range tagarr {
			if tag := strings.TrimSpace(v); tag != "" {
				exists := false
				for _, vv := range addtags {
					if vv == tag {
						exists = true
						break
					}
				}
				if exists == false {
					addtags = append(addtags, tag)
				}
			}
		}
	}
	if id < 1 {
		post.UserId = this.userid
		post.Author = this.username
		post.PostTime = this.getTime()
		post.UpdateTime = this.getTime()
		post.Insert()
	} else {
		post.Id = id
		if post.Read() != nil {
			goto RD
		}
		if post.Tags != "" {
			var tagpostobj models.TagPost
			var tagobj models.Tag
			oldtags := strings.Split(post.Tags, ",")
			//标签统计-1
			tagobj.Query().Filter("name__in", oldtags).Update(orm.Params{"count": orm.ColValue(orm.ColMinus, 1)})
			//删掉tag_post表的记录
			tagpostobj.Query().Filter("postid", post.Id).Delete()
		}

	}

	if len(addtags) > 0 {
		for _, v := range addtags {
			tag := models.Tag{Name: v}
			if tag.Read("name") == orm.ErrNoRows {
				tag.Count = 1
				tag.Insert()
			} else {
				tag.Count += 1
				tag.Update("Count")
			}
			tp := models.TagPost{TagId: tag.Id, PostId: post.Id, PostStatus: int8(status), PostTime: this.getTime()}
			tp.Insert()
		}

	}

	if posttime, err := time.Parse("2006-01-02 15:04:05", timestr); err == nil {
		post.PostTime = posttime
	} else {
		post.PostTime, _ = time.Parse("2006-01-02 15:04:05", post.PostTime.Format("2006-01-02 15:04:05"))
	}

	post.Title = title
	post.Color = color
	post.Content = content
	post.IsTop = istop
	post.Status = int8(status)
	post.Tags = strings.Join(addtags, ",")
	post.UrlName = urlname
	post.UrlType = urltype
	post.UpdateTime = this.getTime()
	post.Update("tags", "status", "title", "color", "is_top", "content", "url_name", "url_type", "update_time", "post_time")

RD:
	this.Redirect("/admin/article/list", 302)
}

func (this *ArticleController) Delete() {
	id, _ := this.GetInt("id")
	post := models.Post{Id: id}
	if post.Read() != nil {
		this.Abort("404")
	}
	post.Delete()
	this.Redirect("/admin/article/list", 302)
}

//批处理
func (this *ArticleController) Batch() {
	ids := this.GetStrings("ids[]")
	op := this.GetString("op")
	idarr := make([]int, 0)
	for _, val := range ids {
		fmt.Println("dadada:", val)
		if id, _ := strconv.Atoi(val); id > 0 {
			idarr = append(idarr, id)
		}
	}

	var post models.Post
	switch op {
	case "topub": //移到已发布
		post.Query().Filter("id__in", idarr).Update(orm.Params{"status": 0})
	case "todrafts": //移到草稿
		post.Query().Filter("id__in", idarr).Update(orm.Params{"status": 1})
	case "totrash": //移到回收站
		post.Query().Filter("id__in", idarr).Update(orm.Params{"status": 2})
	case "delete":
		for _, id := range idarr {
			obj := models.Post{Id: id}
			if obj.Read() == nil {
				obj.Delete()
			}
		}
	}
	this.Redirect(this.Ctx.Request.Referer(), 302)
}
