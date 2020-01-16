package handler

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/tuckKome/fictionary/data"
	"github.com/tuckKome/fictionary/db"
)

//Gameの初期設定
func gameInit(text string) data.Game {
	var newGame data.Game
	newGame.Odai = text
	newGame.Number = 0
	var now = time.Now()
	newGame.CreatedAt = now
	newGame.UpdatedAt = now
	return newGame
}

func makeAns(name string, ans string, id uint) data.Kaitou {
	var kaitou data.Kaitou
	kaitou.User = name
	kaitou.Answer = ans
	kaitou.GameID = id
	var now = time.Now()
	kaitou.CreatedAt = now
	kaitou.UpdatedAt = now

	return kaitou
}

func shuffle(a []data.Kaitou) {
	rand.Seed(time.Now().UnixNano())
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

//Index ははじめのページを取得
func Index(c *gin.Context) {
	var h = db.GetGames()
	c.HTML(200, "index.html", gin.H{"History": h})
}

//CreateGame は新しいゲームを作る
func CreateGame(c *gin.Context) {
	text := c.PostForm("odai")
	game := gameInit(text)
	//dbInsert(game, c) //DB：ゲームに登録

	connect := db.ArgInit()
	db, err := gorm.Open("postgres", connect)
	if err != nil {
		panic("データベース開ず(dbInsert)")
	}
	defer db.Close()

	db.Create(&game)

	id := strconv.Itoa(int(game.ID))
	uri := "/games/" + id + "/new"
	c.Redirect(302, uri)
}

//GetKaitou は回答フォームを取得する
func GetKaitou(c *gin.Context) {
	//idをint型に変換
	n := c.Param("id")
	id, err := strconv.Atoi(n)
	if err != nil {
		panic(err)
	}

	game := db.GetOne(id)
	uri := "/games/" + n + "/new"
	c.HTML(200, "phase21.html", gin.H{"odai": game.Odai, "uri": uri})
}

//CreateKaitou は回答を作る
func CreateKaitou(c *gin.Context) {
	//idをuint型に変換
	n := c.Param("id")
	id, err := strconv.Atoi(n)
	if err != nil {
		panic(err)
	}
	iduint := uint(id)

	name := c.PostForm("name")
	ans := c.PostForm("answer")

	kaitou := makeAns(name, ans, iduint)
	//INSERT
	db.Insert(kaitou)

	uri := "/games/" + n + "/accepted"
	c.Redirect(302, uri)
}

//GetAccepted はAcceptedページを取得
func GetAccepted(c *gin.Context) {
	//idを取得
	n := c.Param("id")

	uri := "/games/" + n
	uri2 := "/games/" + n + "/new"
	c.HTML(200, "phase22.html", gin.H{"uri": uri, "uri2": uri2})
}

//GetList は回答一覧を取得
func GetList(c *gin.Context) {
	var numKaitou int

	n := c.Param("id")
	id, err := strconv.Atoi(n)
	if err != nil {
		panic(err)
	}
	game := db.GetOne(id)
	answers := db.GetKaitou(id)
	shuffle(answers)
	numKaitou = len(answers)
	c.HTML(200, "phase3.html", gin.H{
		"odai":         game.Odai,
		"countOfUsers": numKaitou,
		"kaitous":      answers,
	})
}
