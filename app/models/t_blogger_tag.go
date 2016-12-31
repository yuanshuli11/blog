package models

import (
	"blog/app/support"
	"fmt"
	"strings"
	"log"
)

const(
	TABLE_TAG = "t_tag"
)

//BloggerTag model
type BloggerTag struct {
	Id     int    `xorm:"not null pk autoincr INT(11)"`
	Type   int    `xorm:"not null INT(11)"`
	Name   string `xorm:"not null VARCHAR(20)"`
	Ident  string
	Parent int    `xorm:"INT(11)`
}

func (t *BloggerTag) TableName() string {
	return "t_tag"
}

// Query all tag
// 查找所有 tag
func (b *BloggerTag) ListAll() ([]BloggerTag, error) {
	bt := make([]BloggerTag, 0)
	err := support.Xorm.Find(&bt)
	return bt, err
}

// 根据 id 获取标签
func (b *BloggerTag) GetByID(id int64) (*BloggerTag, error) {
	tag := new(BloggerTag)
	has, err := support.Xorm.Id(id).Get(tag)
	if has {
		return tag, nil
	}
	return nil, err
}

// 根据 ident 获取标签
func (b *BloggerTag) GetByIdent(ident string) (*BloggerTag, error) {
	tag := &BloggerTag{}
	has, err := support.Xorm.Where("ident = ?", ident).Get(tag)
	if has {
		return tag, nil
	}
	return nil, err
}

// Add new tag
func (b *BloggerTag) New() (bool, error) {

	bt := new(BloggerTag)
	bt.Type = b.Type
	bt.Name = b.Name
	bt.Type = b.Parent
	has, err := support.Xorm.InsertOne(bt)

	return has > 0, err
}

// FindBlogCount to get count of blog related to this tag
// 查询标签关联的文章数目
func (t *BloggerTag) FindBlogCount() {

}

// QueryTags to Search for tag
// 根据用户输入的单词匹配 tag
func (t *BloggerTag) QueryTags(str string) ([]map[string][]byte, error) {
	sql := "SELECT name,id FROM t_tag WHERE name LIKE \"%" + str + "%\" ORDER BY LENGTH(name)-LENGTH(\"" + str + "\") ASC LIMIT 10"
	//sql := "SELECT name FROM t_tag"
	ress, err := support.Xorm.Query(sql)
	fmt.Println("res: ", ress)
	if err != nil {
		fmt.Println("err: ", err)
		return ress, err
	}
	return ress, nil
}

// 更新标签
func (t *BloggerTag) Update() bool{
	if t.Id <= 0{
		return false
	}
	support.Xorm.Id(t.Id).Update(t)
	return true
}


// 删除标签
func (t *BloggerTag) Delete(ids []string){
	log.Println("tags",ids)
	sql := "DELETE FROM "+TABLE_TAG+" WHERE id in ("+strings.Join(ids,",")+")"
	support.Xorm.Exec(sql)
}
