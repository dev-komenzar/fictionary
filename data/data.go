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
	Base   int //シャッフルのための変数
	Votes  []Vote
}

// Game : ゲームのDB index>履歴　にも使う
type Game struct {
	gorm.Model
	Odai      string //お題
	Kaitous   []Kaitou
	Phase     string // accepting | playing | archive
	CreatedBy string //作った人
	Secret    string //合言葉
}

// Line : Line bot を友達追加したユーザー・招待したグループを保存
type Line struct {
	gorm.Model
	TalkID string //LINEのユーザー・グループ・ルームID
	Type   string //user OR group OR room
}

// Vote : 投票。Kaitou has many Votes
type Vote struct {
	gorm.Model
	KaitouID  int
	CreatedBy string
}

type Donation struct {
	gorm.Model
	Who      string
	HowMuch  int
	HowToPay string
}
