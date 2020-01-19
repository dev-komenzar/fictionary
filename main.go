package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/tuckKome/fictionary/db"
	"github.com/tuckKome/fictionary/handler"
)

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

	router.POST("/newGame", handler.CreateGame(bot))

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
