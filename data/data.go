package data

import (
	"github.com/jinzhu/gorm"
)

// Kaitou : 回答のDB　回答を集める
type Kaitou struct {
	gorm.Model
	User   string
	Answer string
	Note   string
	GameID uint
	Base   int
}

// Game : ゲームのDB index>履歴　にも使う
type Game struct {
	gorm.Model
	Odai string //お題
}

// Line : LINE bot を友達追加したユーザー・招待したグループを保存
type Line struct {
	gorm.Model
	TalkID string //LINEのユーザー・グループ・ルームID
	Type   string //user OR group OR room
}
