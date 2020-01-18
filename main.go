package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/tuckKome/fictionary/data"
	"github.com/tuckKome/fictionary/db"
	"github.com/tuckKome/fictionary/handler"
)

//Gameの初期設定 POST "/newGame"で使う
func gameInit(text string) data.Game {
	var newGame data.Game
	newGame.Odai = text
	var now = time.Now()
	newGame.CreatedAt = now
	newGame.UpdatedAt = now
	return newGame
}

//POST "/newGame"で使う
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/css", "./css")
	router.Static("/js", "./js")
	router.Static("/resources", "./resources")

	db.Init()

	bot := handler.LineConnect()

	//はじめのページ：お題を入力：過去のお題
	router.GET("/", handler.Index)

	router.POST("/newGame", func(c *gin.Context) {
		text := c.PostForm("odai")
		lineUse := c.PostForm("checkLine")
		fmt.Println(lineUse)
		game := gameInit(text)

		connect := db.ArgInit()
		db, err := gorm.Open("postgres", connect)
		if err != nil {
			panic("データベース開ず(CreateGame)")
		}
		defer db.Close()

		db.Create(&game)

		id := strconv.Itoa(int(game.ID))
		uri := "/games/" + id + "/new"

		if getEnv("GIN_MODE", "debug") == "release" {
			if lineUse == "on" {
				var lines []data.Line
				db.Find(&lines)

				lineMessage := fmt.Sprintf("このURLから回答してね\n%s", uri)
				for i := range lines {
					to := lines[i].TalkID
					if _, err := bot.PushMessage(to, linebot.NewTextMessage(lineMessage)).Do(); err != nil {
						log.Fatal(err)
					}
				}
			}
		}
		c.Redirect(302, uri)
	})

	router.GET("/games/:id/new", handler.GetKaitou)

	router.POST("/games/:id/new", handler.CreateKaitou)

	//「回答受け付けました」ページを表示
	router.GET("/games/:id/accepted", handler.GetAccepted)

	//「みんなの回答」ページを表示
	router.GET("/games/:id", handler.GetList)

	//LINE bot からのwebhookを受ける
	router.POST("/line", func(c *gin.Context) {
		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(events) //jsonを確認したい
		handler.MakeNewLine(events)
	})

	//起動
	router.Run()
}
