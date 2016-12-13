package models

import (
	"blog/app/support"
	"encoding/json"
	"fmt"
	"time"

	"github.com/russross/blackfriday"
)

const (
	BLOG_STATUS_NORMAL  = 0 // 正常状态
	BLOG_STATUS_PENDING = 1 // 审核状态
	BLOG_TYPE_MD        = 0
	BLOG_TYPE_HTML      = 1
)

// Blogger model.
type Blogger struct {
	Id            int64     `xorm:"not null pk autoincr INT(11)"`
	Title         string    `xorm:"not null default '' VARCHAR(50)"`
	Content       string    `xorm:"not null TEXT"`
	CategoryId    int       `xorm:"INT(11)"`
	Passwd        string    `xorm:"VARCHAR(64)"`
	CreateTime    time.Time `xorm:"created"`
	CreateBy      int       `xorm:"not null INT(11)"`
	ReadCount     int64     `xorm:"default 0 BIGINT(20)"`
	LeaveCount    int64     `xorm:"default 0 BIGINT(20)"`
	UpdateTime    time.Time `xorm:"TIMESTAMP"`
	BackgroundPic string    `xorm:"VARCHAR(255)"`
	Type          int       `xorm:"INT(1)"`
	HtmlBak       string    `xorm:"TEXT"`
	Summary       string    `xorm:"VARCHAR(255)"`
	Status        int       `xrom:"INT(11)"`
}

// Get blogger list.
func (b *Blogger) FindList() ([]Blogger, error) {
	// get list data from cache.
	list := make([]Blogger, 0)
	res, _ := support.Cache.Get(support.SPY_BLOGGER_LIST).Result()

	if res != "" {
		err := json.Unmarshal([]byte(res), &list)
		if err == nil {
			return list, err
		}
	}
	// if list data is null in cache,get list data in db.
	err := support.Xorm.Find(&list)

	if err == nil {
		res, e1 := json.Marshal(&list)
		if e1 != nil {
			support.Cache.Set(support.SPY_BLOGGER_LIST, string(res), 0)
		}
	}

	return list, err
}

//New to Add new blogger.
func (b *Blogger) New() (int64, error) {
	blog := new(Blogger)
	blog.Title = b.Title
	blog.Content = b.Content
	blog.CreateBy = b.CreateBy
	blog.UpdateTime = time.Now()
	blog.Passwd = b.Passwd
	blog.CategoryId = b.CategoryId

	has, err := support.Xorm.InsertOne(blog)

	// refurbish cache.
	if err == nil {
		list := make([]Blogger, 0)
		err := support.Xorm.Find(&list)
		if err == nil {
			res, e1 := json.Marshal(&list)
			if e1 != nil {
				support.Cache.Set(support.SPY_BLOGGER_LIST, string(res), 0)
			}
		}
	}
	return has, err
}

// find blogger by id.
func (b *Blogger) FindById() (*Blogger, error) {

	blog := new(Blogger)
	// Get single blogger from cache.
	res, e1 := support.Cache.Get(support.SPY_BLOGGER_SINGLE + fmt.Sprintf("%d", b.Id)).Result()

	if e1 == nil {
		e2 := json.Unmarshal([]byte(res), &blog)
		if e2 == nil {
			return blog, nil
		}
	}
	// if cache not blogger data, find in db.
	_, err := support.Xorm.Id(b.Id).Get(blog)

	if err != nil {
		return blog, err
	}

	return blog, err
}

// Update blogger.
func (b *Blogger) Update() (bool, error) {
	has, err := support.Xorm.Id(b.Id).Update(&b)
	if err == nil {
		// refurbish cache.
		res, e1 := json.Marshal(&b)
		if e1 == nil {
			support.Cache.Del(support.SPY_BLOGGER_SINGLE + fmt.Sprintf("%d", b.Id))
			support.Cache.Set(support.SPY_BLOGGER_SINGLE+fmt.Sprintf("%d", b.Id), string(res), 0)
		}
	}
	return has > 0, err
}

// Delete blogger.
func (b *Blogger) Del() (bool, error) {

	has, err := support.Xorm.Id(b.Id).Delete(&b)

	if err == nil {
		// Delete cache.
		support.Cache.Del(support.SPY_BLOGGER_SINGLE + fmt.Sprintf("%d", b.Id))
	}

	return has > 0, err
}

func (b *Blogger) RenderContent() string {
	if b.Type == BLOG_TYPE_MD {
		return string(blackfriday.MarkdownCommon([]byte(b.Content)))
	}
	return b.Content
}

// GetSummary to cut out a part of blog content
func (b *Blogger) GetSummary() string {
	if b.Summary != "" {
		return b.Summary
	}
	if len(b.Content) < 300 {
		return b.Content
	}
	return b.Content[0:300]
}

// MainURL return the url of the blog
// TODO:Laily it is can be set as id, category, ident and so on
func (b *Blogger) MainURL() string {
	return fmt.Sprintf("/p/%d", b.Id)
}

func (b *Blogger) FindByCategory(categoryID int64) (*[]Blogger, error) {
	blogs := make([]Blogger, 0)
	err := support.Xorm.Where("category_id = ?", categoryID).Find(&blogs)
	return &blogs, err
}

func (b *Blogger) IsMD() bool {
	return b.Type == BLOG_TYPE_MD
}

// GetLatestBlog .
// 获取最新的博客
func (b *Blogger) GetLatestBlog(n int) []Blogger {
	blogs := make([]Blogger, 0)
	support.Xorm.Limit(n, 0).Find(&blogs)
	return blogs
}
