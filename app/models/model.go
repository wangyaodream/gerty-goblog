package models

import "github.com/wangyaodream/gerty-goblog/pkg/types"

type BaseModel struct {
	ID uint64
}

func (bm BaseModel) GetStringID() string {
	return types.Uint64ToString(bm.ID)
}
