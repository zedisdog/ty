package database

import (
	"github.com/zedisdog/ty/generate/snowflake"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
)

var _ callbacks.BeforeCreateInterface = (*SnowflakeID)(nil)

type CommonField struct {
	SnowflakeID
	CreatedAt int64 `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64 `json:"updated_at" gorm:"autoUpdateTime"`
}

type SnowflakeID struct {
	ID uint64 `json:"id,string" gorm:"primaryKey;autoIncrement:false"`
}

func (s *SnowflakeID) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == 0 {
		s.ID, err = snowflake.NextID()
	}
	return
}
