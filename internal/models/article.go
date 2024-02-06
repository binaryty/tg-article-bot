package models

import (
	"strconv"
	"time"
)

type Article struct {
	ID          int
	Title       string
	Link        string
	ThumbUrl    string
	CreatedAt   time.Time
	PublishedAt string
}

func (a Article) Id() string {
	return strconv.Itoa(a.ID)
}
