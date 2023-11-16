package model

import (
	"encoding/json"
	"rule_engine/global"
	"rule_engine/log"
	"rule_engine/redis"
)

type Base struct {
	ID        int    `db:"id" json:"id"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

type QueryID struct {
	ID int `form:"id" binding:"required"`
}

type QueryPage struct {
	Page int `form:"page" binding:"gte=1"`
	Size int `form:"size" binding:"gte=1"`
}

type QueryResult struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
	Pages int `json:"pages"`
}

type Message struct {
	Event string
}

func (m Message) Update() {
	data, _ := json.Marshal(&m)
	if err := redis.Publish(global.MessageChannel, data); err != nil {
		log.Errorf("hot update error:%s", err)
	}
}
