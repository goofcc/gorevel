package controllers

import (
	"fmt"
	"github.com/coocood/qbs"
	"github.com/robfig/revel"
	"gorevel/app/models"
	"gorevel/app/routes"
	"time"
)

type Topic struct {
	Application
}

// 帖子列表
func (c *Topic) Index(page int) revel.Result {
	title := "最近发表"

	categories := getCategories(c.q)
	topics, pagination := models.GetTopics(c.q, page, "", "", "created", routes.Topic.Index(page))

	return c.Render(title, topics, pagination, categories)
}

func (c *Topic) Hot(page int) revel.Result {
	title := "最多点击"

	categories := getCategories(c.q)
	topics, pagination := models.GetTopics(c.q, page, "", "", "hits", routes.Topic.Hot(page))

	c.Render(title, topics, pagination, categories)
	return c.RenderTemplate("topic/Index.html")
}

func (c *Topic) Good(page int) revel.Result {
	title := "好帖推荐"

	categories := getCategories(c.q)
	topics, pagination := models.GetTopics(c.q, page, "good", true, "created", routes.Topic.Good(page))

	c.Render(title, topics, pagination, categories)
	return c.RenderTemplate("topic/Index.html")
}

// 帖子分类查询，帖子列表按时间排序
func (c *Topic) Category(id int64, page int) revel.Result {
	title := "最近发表"

	categories := getCategories(c.q)
	topics, pagination := models.GetTopics(c.q, page, "category_id", id, "created", routes.Topic.Category(id, page))

	c.Render(title, topics, pagination, categories)
	return c.RenderTemplate("topic/Index.html")
}

func (c *Topic) New() revel.Result {
	title := "发表新帖"
	categories := getCategories(c.q)

	return c.Render(title, categories)
}

func (c *Topic) NewPost(topic models.Topic) revel.Result {
	topic.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Topic.New())
	}

	topic.UserId = c.RenderArgs["user"].(*models.User).Id
	if topic.Save(c.q) {
		c.Flash.Success("发表新帖成功")
	} else {
		c.Flash.Error("发表新帖失败")
	}

	return c.Redirect(routes.Topic.Index(1))
}

// 帖子详细
func (c *Topic) Show(id int64) revel.Result {
	topic := findTopicById(c.q, id)
	if topic.Id == 0 {
		return c.NotFound("帖子不存在")
	}

	topic.Hits += 1

	type Topic struct {
		Hits    int
		Created time.Time
	}

	t := new(Topic)
	t.Hits = topic.Hits
	t.Created = topic.Created
	c.q.WhereEqual("id", id).Update(t)

	replies := getReplies(c.q, id)
	categories := getCategories(c.q)

	title := topic.Title
	return c.Render(title, topic, replies, categories)
}

// 回复帖子
func (c *Topic) Reply(id int64, content string) revel.Result {
	c.Validation.Required(content).Message("请填写回复内容")
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Topic.Show(id))
	}

	reply := new(models.Reply)
	reply.TopicId = id
	reply.UserId = c.RenderArgs["user"].(*models.User).Id
	reply.Content = content

	if !reply.Save(c.q) {
		c.Flash.Error("发表回复失败")
	}

	return c.Redirect(routes.Topic.Show(id))
}

func (c *Topic) Edit(id int64) revel.Result {
	title := "编辑帖子"

	topic := findTopicById(c.q, id)
	if topic.Id == 0 {
		return c.NotFound("帖子不存在")
	}

	categories := getCategories(c.q)

	c.Render(title, topic, categories)
	return c.RenderTemplate("topic/New.html")
}

func (c *Topic) EditPost(id int64, topic models.Topic) revel.Result {
	topic.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Topic.Edit(id))
	}

	t := findTopicById(c.q, id)

	if t.Id == 0 {
		return c.NotFound("帖子不存在")
	} else if t.UserId != c.RenderArgs["user"].(*models.User).Id {
		return c.Forbidden("流氓啊")
	}

	t.Title = topic.Title
	t.CategoryId = topic.CategoryId
	t.Content = topic.Content

	if t.Save(c.q) {
		c.Flash.Success("编辑帖子成功")
	} else {
		c.Flash.Error("编辑帖子失败")
	}

	return c.Redirect(routes.Topic.Show(id))
}

func getReplies(q *qbs.Qbs, id int64) []*models.Reply {
	var replies []*models.Reply
	q.WhereEqual("topic_id", id).FindAll(&replies)

	return replies
}

func findTopicById(q *qbs.Qbs, id int64) *models.Topic {
	topic := new(models.Topic)
	err := q.WhereEqual("topic.id", id).Find(topic)
	if err != nil {
		fmt.Println(err)
	}

	return topic
}
