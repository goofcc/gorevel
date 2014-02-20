package controllers

import (
	"fmt"
	"strings"

	"github.com/robfig/revel"

	"gorevel/app/models"
	"gorevel/app/routes"
)

type Topic struct {
	Application
}

// 帖子列表
func (c Topic) Index(page int) revel.Result {
	title := "最近发表"
	categories := getCategories()
	topics, pagination := getTopics(page, "", "created", routes.Topic.Index(page))

	return c.Render(title, topics, pagination, categories)
}

func (c Topic) Hot(page int) revel.Result {
	title := "最多点击"
	categories := getCategories()
	topics, pagination := getTopics(page, "", "hits", routes.Topic.Hot(page))

	c.Vars(Vars{
		"title":      title,
		"topics":     topics,
		"pagination": pagination,
		"categories": categories,
	})

	return c.RenderTemplate("topic/Index.html")
}

func (c Topic) Good(page int) revel.Result {
	title := "好帖推荐"
	categories := getCategories()
	topics, pagination := getTopics(page, "good = true", "created", routes.Topic.Good(page))

	c.Vars(Vars{
		"title":      title,
		"topics":     topics,
		"pagination": pagination,
		"categories": categories,
	})

	return c.RenderTemplate("topic/Index.html")
}

// 帖子分类查询，帖子列表按时间排序
func (c Topic) Category(id int64, page int) revel.Result {
	title := "最近发表"
	categories := getCategories()
	topics, pagination := getTopics(page, fmt.Sprintf("category_id = %d", id), "created", routes.Topic.Category(id, page))

	c.Vars(Vars{
		"title":      title,
		"topics":     topics,
		"pagination": pagination,
		"categories": categories,
	})

	return c.RenderTemplate("topic/Index.html")
}

func (c Topic) New() revel.Result {
	title := "发表新帖"
	categories := getCategories()

	return c.Render(title, categories)
}

func (c Topic) NewPost(topic models.Topic, category int64) revel.Result {
	topic.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Topic.New())
	}

	topic.User = models.User{Id: c.userId}
	topic.Category = models.Category{Id: category}

	aff, _ := engine.Insert(&topic)
	if aff > 0 {
		c.Flash.Success("发表新帖成功")
	} else {
		c.Flash.Error("发表新帖失败")
	}

	return c.Redirect(routes.Topic.Index(1))
}

// 帖子详细
func (c Topic) Show(id int64) revel.Result {
	topic := new(models.Topic)
	has, _ := engine.Id(id).Get(topic)

	if !has {
		return c.NotFound("帖子不存在")
	}

	topic.Hits += 1
	engine.Id(id).Cols("hits").Update(topic)

	replies := getReplies(id)
	categories := getCategories()

	title := topic.Title
	return c.Render(title, topic, replies, categories)
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
		User:    models.User{Id: c.userId},
		Content: content,
	})
	if aff == 0 {
		c.Flash.Error("发表回复失败")
	}

	return c.Redirect(routes.Topic.Show(id))
}

func (c Topic) Edit(id int64) revel.Result {
	title := "编辑帖子"
	categories := getCategories()

	var topic models.Topic
	has, _ := engine.Id(id).Get(&topic)
	if !has {
		return c.NotFound("帖子不存在")
	}

	c.Vars(Vars{
		"title":      title,
		"topic":      topic,
		"categories": categories,
	})

	return c.RenderTemplate("topic/New.html")
}

func (c Topic) EditPost(id int64, topic models.Topic, category int64) revel.Result {
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
	} else {
		c.Flash.Error("编辑帖子失败")
	}

	return c.Redirect(routes.Topic.Show(id))
}

func getTopics(page int, where string, order string, url string) ([]models.Topic, *Pagination) {
	page -= 1
	if page < 0 {
		page = 0
	}

	var topics []models.Topic
	var rows int64
	if where == "" {
		rows, _ = engine.Count(&models.Topic{})
		err := engine.Omit("Content").Desc(order).Limit(ItemsPerPage, page*ItemsPerPage).Find(&topics)
		if err != nil {
			revel.ERROR.Println(err)
		}
	} else {
		rows, _ = engine.Where(where).Count(&models.Topic{})
		err := engine.Where(where).Omit("Content").Desc(order).Limit(ItemsPerPage, page*ItemsPerPage).Find(&topics)
		if err != nil {
			revel.ERROR.Println(err)
		}
	}

	url = url[:strings.Index(url, "=")+1]
	pagination := NewPagination(page, int(rows), url)

	return topics, pagination
}

func getReplies(id int64) []models.Reply {
	var replies []models.Reply
	engine.Where("topic_id = ?", id).Find(&replies)

	return replies
}
