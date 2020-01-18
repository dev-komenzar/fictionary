package handler

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgresql を使うためのライブラリ
	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/tuckKome/fictionary/data"
	"github.com/tuckKome/fictionary/db"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

//Gameの初期設定
func gameInit(text string) data.Game {
	var newGame data.Game
	newGame.Odai = text
	var now = time.Now()
	newGame.CreatedAt = now
	newGame.UpdatedAt = now
	return newGame
}

func lineInit(id string, typeOfSource string) data.Line {
	var newLine data.Line
	newLine.TalkID = id
	newLine.Type = typeOfSource
	var now = time.Now()
	newLine.CreatedAt = now
	newLine.UpdatedAt = now
	return newLine
}

//LineConnect : LINE bot 接続
func LineConnect() *linebot.Client {
	channelID := getEnv("CHANNEL_ID", "")
	channelSecret := getEnv("CHANNEL_SECRET", "")

	bot, err := linebot.New(channelSecret, channelID)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return bot
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

/*
//CreateGame は新しいゲームを作る
func CreateGame(c *gin.Context) {
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
}
*/

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
	db.InsertKaitou(kaitou)

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

func getNotNill(a string, b string, c string) string {
	if a != "" {
		return a
	} else if b != "" {
		return b
	} else {
		return c
	}
}

func MakeNewLine(events []*linebot.Event) {
	for _, event := range events {
		if event.Type == linebot.EventTypeJoin {
			userID := event.Source.UserID
			groupID := event.Source.GroupID
			roomID := event.Source.RoomID

			var line data.Line
			d := getNotNill(userID, groupID, roomID)
			if d == userID {
				line = lineInit(d, "user")
			} else if d == groupID {
				line = lineInit(d, "group")
			} else {
				line = lineInit(d, "room")
			}
			db.InsertLine(line) //DBにLINEからの情報が登録された

		}
	}
}
