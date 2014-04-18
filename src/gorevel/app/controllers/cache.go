package controllers

import (
	"strconv"
	"strings"

	"github.com/robfig/revel/cache"

	"gorevel/app/models"
	"gorevel/app/routes"
)

func cacheUpdateHits(id int64) {
	go updateHits(id, routes.Topic.Index(0))
	go updateHits(id, routes.Topic.Hot(0))
	go updateHits(id, routes.Topic.Good(0))
}

func updateHits(id int64, url string) {
	str := "topics" + url[:strings.Index(url, "=")+1]
	for i := 1; i <= models.CachePageSize; i++ {
		url = str + strconv.Itoa(i)
		var topics []models.Topic
		if err := cache.Get(url, &topics); err == nil {
			for index, topic := range topics {
				if topic.Id == id {
					topics[index].Hits += 1
					cache.Replace(url, topics, cache.FOREVER)
					return
				}
			}
		}
	}
}

func cacheUpdateReplies(id int64) {
	go updateReplies(id, routes.Topic.Index(0))
	go updateReplies(id, routes.Topic.Hot(0))
	go updateReplies(id, routes.Topic.Good(0))
}

func updateReplies(id int64, url string) {
	str := "topics" + url[:strings.Index(url, "=")+1]
	for i := 1; i <= models.CachePageSize; i++ {
		url = str + strconv.Itoa(i)
		var topics []models.Topic
		if err := cache.Get(url, &topics); err == nil {
			for index, topic := range topics {
				if topic.Id == id {
					topics[index].Replies += 1
					cache.Replace(url, topics, cache.FOREVER)
					return
				}
			}
		}
	}
}
