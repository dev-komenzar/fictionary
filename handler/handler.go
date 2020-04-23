package handler

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgresql を使うためのライブラリ

	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	"github.com/tuckKome/fictionary/data"
	"github.com/tuckKome/fictionary/db"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

//Gameの初期設定 POST "/newGame"で使う
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
	channelSecret := getEnv("CHANNEL_SECRET", "")
	channelToken := getEnv("CHANNEL_ACCESS_TOKEN", "")

	bot, err := linebot.New(channelSecret, channelToken)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return bot
}

//TwitterConnect : twitter bot 接続
func TwitterConnect() *twitter.Client {
	consumerKey := getEnv("TWITTER_CONSUMER_KEY", "")
	consumerSecret := getEnv("TWITTE_CONSUMER_SECRET", "")
	accessToken := getEnv("TWITTER_ACCESS_TOKEN", "")
	accessSecret := getEnv("TWITTER_ACCESS_SECRET", "")

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	// http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)
	return client
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
	//[0,1,2,...,k-1]を用意
	k := len(a)
	arr := make([]int, k)
	for i := 0; i < k; i++ {
		arr[i] = i
	}

	//Fisher–Yates シャッフル
	rand.Seed(time.Now().UnixNano())
	for i := k - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		tmp := arr[i]
		arr[i] = arr[j]
		arr[j] = tmp
	}

	//シャッフルされたarr を[]Kaitou に入れる
	for i := range a {
		a[i].Base = arr[i]
	}
}

//Index ははじめのページを取得
func Index(c *gin.Context) {
	var h = db.GetGames()
	c.HTML(200, "index.html", gin.H{"History": h})
}

//GetKaitou は回答フォームを取得する
func GetKaitou(c *gin.Context) {
	//idをint型に変換
	n := c.Param("id")
	id, err := strconv.Atoi(n)
	if err != nil {
		panic(err)
	}

	game := db.GetGame(id)
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

	a := db.GetKaitous(id)
	shuffle(a)
	db.UpdateKaitous(a)

	uri := "/games/" + n + "/accepted"
	c.Redirect(302, uri)
}

//GetAccepted はAcceptedページを取得
func GetAccepted(c *gin.Context) {
	//idを取得
	n := c.Param("id")
	id, err := strconv.Atoi(n)
	if err != nil {
		panic(err)
	}

	game := db.GetGame(id)
	uri := "/games/" + n
	uri2 := "/games/" + n + "/new"
	c.HTML(200, "phase22.html", gin.H{"odai": game.Odai, "uri": uri, "uri2": uri2})
}

//GetList は回答一覧を取得
func GetList(c *gin.Context) {
	//idを取得
	n := c.Param("id")
	id, err := strconv.Atoi(n)
	if err != nil {
		panic(err)
	}

	game := db.GetGame(id)

	a := db.GetKaitous(id) //回答一覧を取得

	//Kaitou.Base　で並び替える
	sort.SliceStable(a, func(i, j int) bool { return a[i].Base < a[j].Base })

	k := len(a) //coutOfUsers のため
	c.HTML(200, "phase3.html", gin.H{
		"odai":         game.Odai,
		"countOfUsers": k,
		"kaitous":      a,
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

//CreateGame は「*linebot.Client型を引数にした」新しいゲームを作る関数
func CreateGame(bot *linebot.Client, twitterClient *twitter.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		text := c.PostForm("odai")
		lineUse := c.PostForm("checkLine")
		twitterUse := c.PostForm("checkTwitter")
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

		//LINE bot・twitterに投げる
		if getEnv("GIN_MODE", "debug") == "release" {
			url := getEnv("HOST_ADDRESS", "localhost:8080") + uri
			message := fmt.Sprintf("お題は「%s」\nこのURLから回答してね\n%s", text, url)
			if lineUse == "on" {
				var lines []data.Line
				db.Find(&lines)

				for i := range lines {
					to := lines[i].TalkID
					if _, err := bot.PushMessage(to, linebot.NewTextMessage(message)).Do(); err != nil {
						log.Fatal(err)
					}
				}
			}
			if twitterUse == "on" {
				// Send a Tweet
				_, _, err := twitterClient.Statuses.Update(message, nil)
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		c.Redirect(302, uri)
	}
}

//CreateLine は「*linebot.Clientを引数にする」ユーザー・グループIDをDBに登録するhandler
func CreateLine(bot *linebot.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(events) //jsonを確認したい
		for _, event := range events {
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
			switch event.Type {
			case linebot.EventTypeJoin:
				db.InsertLine(line) //DBにLINEからの情報が登録された
				fmt.Printf("New Line created. ID %s \n", string(line.ID))
			case linebot.EventTypeLeave:
				db.DeleteLine(line)
				fmt.Printf("ID %s is deleted \n", string(line.ID))
			case linebot.EventTypeUnfollow:
				db.DeleteLine(line)
				fmt.Printf("ID %s is deleted \n", string(line.ID))
			}

		}
	}
}
