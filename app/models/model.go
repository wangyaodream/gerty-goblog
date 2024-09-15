package models

import (
	"time"

	"github.com/wangyaodream/gerty-goblog/pkg/types"
)

type BaseModel struct {
	ID uint64 `gorm:"column:id;primarykey;autoIncreament;not null"`

	CreatedAt time.Time `gorm:"column:created_at;index"`
	UpdatedAt time.Time `gorm:"column:updated_at;index"`
}

func (bm BaseModel) GetStringID() string {
	return types.Uint64ToString(bm.ID)
}
