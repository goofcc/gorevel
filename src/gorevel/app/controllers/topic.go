package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/robfig/revel"
	"github.com/robfig/revel/cache"

	"gorevel/app/models"
	"gorevel/app/routes"
)

type Topic struct {
	Application
}

func (c Topic) New() revel.Result {
	title := "发表新帖"

	return c.Render(title)
}

func (c Topic) NewPost(topic models.Topic, category int64) revel.Result {
	topic.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Topic.New())
	}

	topic.User = models.User{Id: c.user().Id}
	topic.Category = models.Category{Id: category}

	aff, _ := engine.Insert(&topic)
	if aff > 0 {
		c.Flash.Success("发表新帖成功")
		cache.Flush()
	} else {
		c.Flash.Error("发表新帖失败")
	}

	return c.Redirect(routes.Topic.Index(1))
}

// 帖子详细
func (c Topic) Show(id int64) revel.Result {
	topic := new(models.Topic)
	str := strconv.Itoa(int(id))

	if err := cache.Get("topic"+str, &topic); err != nil {
		has, _ := engine.Id(id).Get(topic)
		if !has {
			return c.NotFound("帖子不存在")
		}
		go cache.Set("topic"+str, topic, cache.FOREVER)
	}

	topic.Hits += 1
	engine.Id(id).Cols("hits").Update(topic)

	replies := getReplies(id)

	title := topic.Title
	return c.Render(title, topic, replies)
}

// 回复帖子
func (c Topic) Reply(id int64, content string) revel.Result {
	c.Validation.Required(content).Message("请填写回复内容")
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Topic.Show(id))
	}

	aff, _ := engine.Insert(&models.Reply{
		Topic:   models.Topic{Id: id},
		User:    models.User{Id: c.user().Id},
		Content: content,
	})

	if aff > 0 {
		engine.Exec("UPDATE topic SET replies = replies + 1 WHERE id = ?", id)
		cache.Flush()
	} else {
		c.Flash.Error("发表回复失败")
	}

	return c.Redirect(routes.Topic.Show(id))
}

func (c Topic) Edit(id int64) revel.Result {
	var topic models.Topic
	has, _ := engine.Id(id).Get(&topic)
	if !has {
		return c.NotFound("帖子不存在")
	}

	if user := c.user(); user.Id != topic.User.Id {
		return c.Forbidden("抱歉，您没有权限")
	}

	c.bindVars(Vars{
		"title": "编辑帖子",
		"topic": topic,
	})

	return c.RenderTemplate("topic/New.html")
}

func (c Topic) EditPost(id int64, topic models.Topic, category int64) revel.Result {
	user := c.user()
	if has, _ := engine.Where("id = ? AND user_id = ?", id, user.Id).Get(&models.Topic{}); !has {
		return c.Forbidden("抱歉，您没有权限")
	}

	topic.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Topic.Edit(id))
	}

	topic.Category = models.Category{Id: category}

	aff, _ := engine.Id(id).Cols("title", "category_id", "content").Update(&topic)
	if aff > 0 {
		c.Flash.Success("编辑帖子成功")
		cache.Flush()
	} else {
		c.Flash.Error("编辑帖子失败")
	}

	return c.Redirect(routes.Topic.Show(id))
}

// 帖子列表
func (c Topic) Index(page int) revel.Result {
	title := "最近发表"
	topics, pagination := GetTopics(page, "", "created", routes.Topic.Index(page))

	return c.Render(title, topics, pagination)
}

func (c Topic) Hot(page int) revel.Result {
	title := "最多点击"
	topics, pagination := GetTopics(page, "", "hits", routes.Topic.Hot(page))

	c.bindVars(Vars{
		"title":      title,
		"topics":     topics,
		"pagination": pagination,
	})

	return c.RenderTemplate("topic/Index.html")
}

func (c Topic) Good(page int) revel.Result {
	title := "好帖推荐"
	topics, pagination := GetTopics(page, "good = true", "created", routes.Topic.Good(page))

	c.bindVars(Vars{
		"title":      title,
		"topics":     topics,
		"pagination": pagination,
	})

	return c.RenderTemplate("topic/Index.html")
}

func (c Topic) SetGood(id int64) revel.Result {
	aff, _ := engine.Id(id).Cols("good").Update(&models.Topic{Good: true})
	if aff > 0 {
		return c.RenderJson(map[string]bool{"status": true})
	}

	return c.RenderJson(map[string]bool{"status": false})
}

// 帖子分类查询，帖子列表按时间排序
func (c Topic) Category(id int64, page int) revel.Result {
	title := "最近发表"
	topics, pagination := GetTopics(page, fmt.Sprintf("category_id = %d", id), "created", routes.Topic.Category(id, page))

	c.bindVars(Vars{
		"title":      title,
		"topics":     topics,
		"pagination": pagination,
	})

	return c.RenderTemplate("topic/Index.html")
}

func GetTopics(page int, where string, order string, url string) ([]models.Topic, *Pagination) {
	if page < 1 {
		page = 1
		url = url[:strings.Index(url, "=")+1] + "1"
	}

	var topics []models.Topic
	var pagination *Pagination

	if page > models.CachePageSize {
		// 当前页超出缓存页，从数据库取数据
		topics, pagination = queryDb(page, where, order, url)

	} else {
		// 当前页在缓存中，从缓存中取数据
		if err := cache.Get("topics"+url, &topics); err != nil {
			// 缓存中没有找到，查数据库
			topics, pagination = queryDb(page, where, order, url)
			if len(topics) == 0 {
				return topics, pagination
			}
			// 写入缓存
			go cache.Set("topics"+url, topics, cache.FOREVER)
			go cache.Set("pagination"+url, pagination, cache.FOREVER)

		} else {
			// 从缓存中取分页
			if err := cache.Get("pagination"+url, &pagination); err != nil {
				revel.ERROR.Println(err)
			}
		}
	}

	return topics, pagination
}

func queryDb(page int, where string, order string, url string) ([]models.Topic, *Pagination) {
	var topics []models.Topic
	var pagination *Pagination
	var rows int64
	if where == "" {
		rows, _ = engine.Count(&models.Topic{})
		err := engine.Omit("Content").Desc(order).Limit(ROWS_PER_PAGE, (page-1)*ROWS_PER_PAGE).Find(&topics)
		if err != nil {
			revel.ERROR.Println(err)
		}
	} else {
		rows, _ = engine.Where(where).Count(&models.Topic{})
		err := engine.Where(where).Omit("Content").Desc(order).Limit(ROWS_PER_PAGE, (page-1)*ROWS_PER_PAGE).Find(&topics)
		if err != nil {
			revel.ERROR.Println(err)
		}
	}

	pagination = NewPagination(page, int(rows), url[:strings.Index(url, "=")+1])

	return topics, pagination
}

func getReplies(id int64) []models.Reply {
	var replies []models.Reply
	str := strconv.Itoa(int(id))

	if err := cache.Get("replies"+str, &replies); err != nil {
		engine.Where("topic_id = ?", id).Find(&replies)
		go cache.Set("replies"+str, replies, cache.FOREVER)
	}

	return replies
}
