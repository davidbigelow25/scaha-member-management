package model

import (
	"errors"
	"github.com/jinzhu/gorm"
	"html"
	"strings"
	"time"
)

type Person struct {
	gorm.Model
	FirstName  string    `gorm:"size:45;not null;`
	LastName  string     `gorm:"size:45;not null;`
	Email string    	 `gorm:"size:45;not null;`
	Phone  string    	 `gorm:"size:14;not null;"`
}

//
// Hey, lets be smart and filter out all the garbaded that can come it
// make it html safeish
func (p *Person) Prepare() {
	p.ID = 0
	p.FirstName = html.EscapeString(strings.TrimSpace(p.FirstName))
	p.LastName = html.EscapeString(strings.TrimSpace(p.LastName))
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

//
// Here are all the validation rules we can apply here for the given object
//
func (p *Person) Validate() error {

	if p.FirstName == "" {
		return errors.New("First Name Required")
	}
	if p.LastName == "" {
		return errors.New("Last Name")
	}
	return nil
}

/*func (p *Person) Save(db *gorm.DB) (*Person, error) {

	var err error
	err = db.Debug().Model(&Person{}).Create(&p).Error
	if err != nil {
		return &Person{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&Person{}).Where("id = ?", p).Save(&p).Error
		if err != nil {
			return &Person{}, err
		}
	}
	return p, nil
}

func (p *Post) FindAllPosts(db *gorm.DB) (*[]Post, error) {
	var err error
	posts := []Post{}
	err = db.Debug().Model(&Post{}).Limit(100).Find(&posts).Error
	if err != nil {
		return &[]Post{}, err
	}
	if len(posts) > 0 {
		for i, _ := range posts {
			err := db.Debug().Model(&User{}).Where("id = ?", posts[i].AuthorID).Take(&posts[i].Author).Error
			if err != nil {
				return &[]Post{}, err
			}
		}
	}
	return &posts, nil
}

func (p *Post) FindPostByID(db *gorm.DB, pid uint64) (*Post, error) {
	var err error
	err = db.Debug().Model(&Post{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}

func (p *Post) UpdateAPost(db *gorm.DB) (*Post, error) {

	var err error

	err = db.Debug().Model(&Post{}).Where("id = ?", p.ID).Updates(Post{Title: p.Title, Content: p.Content, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}

func (p *Post) DeleteAPost(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Post{}).Where("id = ? and author_id = ?", pid, uid).Take(&Post{}).Delete(&Post{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Post not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}*/