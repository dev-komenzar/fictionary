package main

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/tuckKome/fictionary/data"
)

const (
	connect  string = "host=" + host + " port=" + port + " user=" + user + " dbname=" + dbname + " password=" + password + " sslmode=disable"
	host     string = "127.0.0.1"
	port     string = "5432"
	user     string = "tahoiya"
	dbname   string = "dbtahoiya"
	password string = "password"
)

//DB初期化
func dbInit() {
	db, err := gorm.Open("postgres", connect)
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&data.Kaitou{}, &data.Game{})
	defer db.Close()
}

//DBから一つ取り出す：回答ページで使用
func dbGetOne(id int) data.Game {
	db, err := gorm.Open("postgres", connect)
	if err != nil {
		panic("データベース開ず(dbInsert)")
	}
	defer db.Close()

	var game data.Game
	db.First(&game, id)
	return game
}

//DBから[]Kaitouを取り出す
func dbGetKaitou(id int) []data.Kaitou {
	db, err := gorm.Open("postgres", connect)
	if err != nil {
		panic("データベース開ず(dbInsert)")
	}
	defer db.Close()

	var kaitous []data.Kaitou
	db.Where("game_id = ?", id).Find(&kaitous)
	return kaitous
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

func dbGetGames() []data.Game {
	var games []data.Game

	db, err := gorm.Open("postgres", connect)
	if err != nil {
		panic("データベース開ず(dbInsert)")
	}
	defer db.Close()

	db.Find(&games)
	return games
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

func main() {

	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	router.StaticFile("/css/bootstrap.min.css", "./css/bootstrap.min.css")
	router.StaticFile("/js/bootstrap.min.js", "./js/bootstrap.min.js")

	dbInit()

	//はじめのページ：お題を入力：過去のお題
	router.GET("/", func(c *gin.Context) {
		var h = dbGetGames()
		c.HTML(200, "index.html", gin.H{"History": h})
	})

	router.POST("/newGame", func(c *gin.Context) {
		text := c.PostForm("odai")
		game := gameInit(text)
		//dbInsert(game, c) //DB：ゲームに登録

		db, err := gorm.Open("postgres", connect)
		if err != nil {
			panic("データベース開ず(dbInsert)")
		}
		defer db.Close()

		db.Create(&game)

		id := strconv.Itoa(int(game.ID))
		uri := "/games/" + id + "/new"
		c.Redirect(302, uri)
	})

	router.GET("/games/:id/new", func(c *gin.Context) {
		//idをint型に変換
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}

		game := dbGetOne(id)
		uri := "/games/" + n + "/new"
		c.HTML(200, "phase21.html", gin.H{"odai": game.Odai, "uri": uri})
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
		//INSERT
		db, err := gorm.Open("postgres", connect)
		if err != nil {
			panic("データベース開ず(dbInsert)")
		}
		defer db.Close()

		db.Create(&kaitou)

		uri := "/games/" + n + "/accepted"
		c.Redirect(302, uri)
	})

	//「回答受け付けました」ページを表示
	router.GET("/games/:id/accepted", func(c *gin.Context) {
		//idを取得
		n := c.Param("id")

		uri := "/games/" + n
		c.HTML(200, "phase22.html", gin.H{"uri": uri})
	})

	//「みんなの回答」ページを表示
	router.GET("/games/:id", func(c *gin.Context) {
		var numKaitou int

		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		game := dbGetOne(id)
		answers := dbGetKaitou(id)
		shuffle(answers)
		numKaitou = len(answers)
		c.HTML(200, "phase3.html", gin.H{
			"odai":         game.Odai,
			"countOfUsers": numKaitou,
			"kaitous":      answers,
		})
	})

	//起動
	router.Run()
}
