package main

import (
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

	twitterClient := handler.TwitterConnect()

	//はじめのページ：お題を入力：過去のお題
	router.GET("/", handler.Index)

	router.POST("/newGame", handler.CreateGame(bot, twitterClient))

	router.GET("/games/:id/new", handler.GetKaitou)

	router.POST("/games/:id/new", handler.CreateKaitou)

	//「回答受け付けました」ページを表示
	router.GET("/games/:id/accepted", handler.GetAccepted)

	//「みんなの回答」ページを表示
	router.GET("/games/:id", handler.GetList)

	//「投票」機能
	router.POST("/games/:id", handler.CreateVote)

	//LINE bot からのwebhookを受ける
	router.POST("/line", handler.CreateLine(bot))

	//起動
	router.Run()
}
