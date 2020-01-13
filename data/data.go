package data

import (
	"github.com/jinzhu/gorm"
)

// 回答のDB　回答を集める
type Kaitou struct {
	gorm.Model
	User   string
	Answer string
	Note   string
	GameID uint
}

// ゲームのDB index>履歴　にも使う
type Game struct {
	gorm.Model
	Odai   string //お題
	Number uint //回答数
}
