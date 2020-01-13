package main

import (
	"log"
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
func dbInsert(newGame data.Game, c *gin.Context) {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=tahoiya dbname=tahoiya password=password")
	if err != nil {
		panic("データベース開ず(dbInsert)")
	}
	defer db.Close()

	if db.NewRecord(&newGame) == false {
		panic("すでにデータが存在します。")
	} else {
		db.Create(&newGame)
		if db.NewRecord(&newGame) == false {
			log.Printf("History %d Recorded\n", newGame.ID)
		} else {
			log.Println("History not created") //エラー内容を表示したい http://gorm.io/ja_JP/docs/error_handling.html
			c.HTML(200, "error.html", gin.H{"message": "問題が作成されませんでした。もう一度試してください"})
		}
	}
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

func main() {

	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	dbInit()

	//はじめのページ：お題を入力：過去のお題
	var h = makeHistory()

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{"History": h})
	})

	router.POST("/new", func(c *gin.Context) {
		text := c.PostForm("odai")
		newGame := gameInit(text)
		dbInsert(newGame, c)
		c.HTML(200, "phase21.html", gin.H{"odai": text})
	})
}
