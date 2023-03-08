package gorm

import (
	"github.com/zedisdog/ty/generate/snowflake"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
)

var _ callbacks.BeforeCreateInterface = (*CommonField)(nil)

type CommonField struct {
	ID        uint64 `json:"id,string" gorm:"primaryKey"`
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (a *CommonField) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == 0 {
		a.ID, err = snowflake.NextID()
	}
	return
}
