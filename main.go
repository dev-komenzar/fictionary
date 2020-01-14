package main

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/tuckKome/fictionary/data"
)

//DB初期化
func dbInit() {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=tahoiya dbname=tahoiya password=password")
	if err != nil {
		panic("データベース開ず(dbInit)")
	}
	db.AutoMigrate(&data.Kaitou{}, &data.Game{})
	defer db.Close()
}

//DBに追加
func dbInsert(new interface{}, c *gin.Context) {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=tahoiya dbname=tahoiya password=password")
	if err != nil {
		panic("データベース開ず(dbInsert)")
	}
	defer db.Close()

	if db.NewRecord(&new) == false {
		panic("すでにデータが存在します。")
	} else {
		db.Create(&new)
		if db.NewRecord(&new) == false {
			log.Printf("History Recorded\n")
		} else {
			log.Println("History not created") //エラー内容を表示したい http://gorm.io/ja_JP/docs/error_handling.html
			c.HTML(200, "error.html", gin.H{"message": "問題が作成されませんでした。もう一度試してください"})
		}
	}
}

//DBから一つ取り出す：回答ページで使用
func dbGetOne(id int) data.Game {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=tahoiya dbname=tahoiya password=password")
	if err != nil {
		panic("データベース開ず(dbInsert)")
	}
	defer db.Close()

	var game data.Game
	db.First(&game, id)
	return game
}

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

func makeHistory() []string {
	var games []data.Game
	var history []string

	db, err := gorm.Open("postgres", "host=localhost port=5432 user=tahoiya dbname=tahoiya password=password")
	if err != nil {
		panic("データベース開ず(dbInsert)")
	}
	defer db.Close()

	db.Find(&games)
	for i := range games {
		var odai = games[i].Odai
		history = append(history, odai)
	}
	return history
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

func main() {

	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	dbInit()

	//はじめのページ：お題を入力：過去のお題
	router.GET("/", func(c *gin.Context) {
		var h = makeHistory()
		c.HTML(200, "index.html", gin.H{"History": h})
	})

	router.POST("/newGame", func(c *gin.Context) {
		text := c.PostForm("odai")
		newGame := gameInit(text)
		dbInsert(newGame, c) //DB：ゲームに登録

		id := strconv.Itoa(int(newGame.ID))
		uri := "/games/" + id + "/kaitou"
		c.Redirect(302, uri)
	})

	router.GET("/games/:id/kaitou", func(c *gin.Context) {
		//idをint型に変換
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}

		odai := dbGetOne(id)
		uri := "/games/" + n + "/new"
		c.HTML(200, "phase21.html", gin.H{"odai": odai, "uri": uri})
	})

	router.POST("/games/:id/new", func(c *gin.Context) {
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
		dbInsert(kaitou, c)

		uri := "/games/" + n + "/accepted"
		c.Redirect(302, uri)
	})

	//「回答受け付けました」ページ
	router.GET("/games/:id/accepted", func(c *gin.Context) {
		//idをint型に変換
		n := c.Param("id")

		uri := "/games/" + n
		c.HTML(200, "phase22.html", gin.H{"uri": uri})
	})
}
