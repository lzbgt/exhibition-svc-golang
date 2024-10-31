// Author: Bruce Lu
// Email: lzbgt_AT_icloud.com

package models

import (
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Base struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	CreateTime time.Time `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"autoUpdateTime"`
	IsDeleted  bool      `json:"is_deleted" gorm:"default:false"`
	IsActive   bool      `json:"is_active" gorm:"default:true"`
}

type UserLogin struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type ExibitionInput struct {
	Title       string          `json:"title,omitempty"`
	Description string          `json:"description,omitempty"`
	Location    string          `json:"location,omitempty"`
	Creator     json.RawMessage `json:"creator,omitempty" gorm:"type:json"`
	Sponsors    json.RawMessage `json:"sponsors,omitempty" gorm:"type:json"`
	Videos      json.RawMessage `json:"videos,omitempty" gorm:"type:json"`
	Images      json.RawMessage `json:"images,omitempty" gorm:"type:json"`
	StartTime   time.Time       `json:"start_time" gorm:"default:null"`
	EndTime     time.Time       `json:"end_time" gorm:"default:null"`
}

type Exibition struct {
	Base
	ExibitionInput
}

type ExUserInput struct {
	Eid      int    `json:"eid,omitempty" gorm:"index"`
	Name     string `json:"name"`
	Uname    string `json:"uname" gorm:"unique"`
	Password string `json:"password"`
	Title    string `json:"title"`
	Mobile   string `json:"mobile"`
}

type ExUser struct {
	Base
	ExUserInput
}

type ExCatalogInput struct {
	Pid         int             `json:"pid" gorm:"index,default:0"`
	Eid         int             `json:"eid" gorm:"index"`
	Name        string          `json:"name"`
	NameEn      string          `json:"name_en"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Images      json.RawMessage `json:"images" gorm:"type:json"`
	Videos      json.RawMessage `json:"videos" gorm:"type:json"`
	Banner      json.RawMessage `json:"banner" gorm:"type:json"`
}

type ExCatalog struct {
	Base
	ExCatalogInput
	RootId int `json:"root_id"`
}

// Claims struct
type Claims struct {
	Username string `json:"username"`
	UserId   int    `json:"user_id"`
	Eid      int    `json:"eid"`
	Perms    string `json:"perms"`
	jwt.RegisteredClaims
}

type BatchUserInput struct {
	NamePrefix string `json:"name_prefix"`
	IndexRange []int  `json:"index_range" binding:"required,gt=1,dive"`
}

type ExRateInput struct {
	Eid  int     `json:"eid"`
	Rate float64 `json:"rate" binding:"required,gt=0" gorm:"default:5"`
	Uid  int     `json:"uid" gorm:"uniqueIndex:idx_rate_user_per_item"`
	Iid  int     `json:"iid" gorm:"uniqueIndex:idx_rate_user_per_item"`
}

type ExRate struct {
	Base
	ExRateInput
	//Item ExItem `json:"item" gorm:"foreignKey:Iid;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type ExAmountInput struct {
	Eid    int `json:"eid"`
	Amount int `json:"amount" gorm:"default:0"`
	Uid    int `json:"uid" gorm:"uniqueIndex:idx_amount_user_per_item"`
	Iid    int `json:"iid" gorm:"uniqueIndex:idx_amount_user_per_item"`
}

type ExAmount struct {
	Base
	ExAmountInput
	//Item ExItem `json:"item" gorm:"foreignKey:Iid;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type ExItemInput struct {
	Eid         int             `json:"eid"`
	Name        string          `json:"name" binding:"required"`
	Description string          `json:"description"`
	Thumbnails  json.RawMessage `json:"thumbnails" gorm:"type:json"`
	Images      json.RawMessage `json:"images" gorm:"type:json"`
	Videos      json.RawMessage `json:"videos" gorm:"type:json"`
	Cid         int             `json:"cid"`
}

type ExItem struct {
	Base
	ExItemInput
	//Catalog ExCatalog `json:"catalog" gorm:"foreignKey:Cid;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type ExCommentInput struct {
	Eid     int    `json:"eid"`
	Uid     int    `json:"uid" gorm:"uniqueIndex:idx_cmt_user_per_item"`
	Iid     int    `json:"iid" gorm:"uniqueIndex:idx_cmt_user_per_item"`
	Content string `json:"content" binding:"required"`
}

type ExComment struct {
	Base
	ExCommentInput
	//Item ExItem `json:"item" gorm:"foreignKey:Iid;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type RatingDistribution struct {
	Category string
	Count    int64
	Percent  float64
}
