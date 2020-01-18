package db

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgresqlを使うためのライブラリ
	"github.com/tuckKome/fictionary/data"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

//ArgInt はOpenの引数を作る
func ArgInit() string {
	host := getEnv("FICTIONARY_DATABASE_HOST", "127.0.0.1")
	port := getEnv("FICTIONARY_PORT", "5432")
	user := getEnv("FICTIONARY_USER", "tahoiya")
	dbname := getEnv("FICTIONARY_DB_NAME", "dbtahoiya")
	password := getEnv("FICTIONARY_DB_PASS", "password")
	sslmode := getEnv("FICTIONARY_SSLMODE", "disable")

	dbinfo := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		user,
		password,
		host,
		port,
		dbname,
		sslmode,
	)
	return dbinfo
}

//DB初期化
func Init() {
	connect := ArgInit()
	db, err := gorm.Open("postgres", connect)
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&data.Kaitou{}, &data.Game{}, &data.Line{})
	defer db.Close()
}

//DBから一つ取り出す：回答ページで使用
func GetOne(id int) data.Game {
	connect := ArgInit()
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
func GetKaitou(id int) []data.Kaitou {
	connect := ArgInit()
	db, err := gorm.Open("postgres", connect)
	if err != nil {
		panic("データベース開ず(dbInsert)")
	}
	defer db.Close()

	var kaitous []data.Kaitou
	db.Where("game_id = ?", id).Find(&kaitous)
	return kaitous
}

func GetGames() []data.Game {
	var games []data.Game

	connect := ArgInit()
	db, err := gorm.Open("postgres", connect)
	if err != nil {
		panic("データベース開ず(dbInsert)")
	}
	defer db.Close()

	db.Find(&games)
	return games
}

//InsertKaitou : DBに新しいkaitouを追加
func InsertKaitou(kaitou data.Kaitou) {
	connect := ArgInit()
	db, err := gorm.Open("postgres", connect)
	if err != nil {
		panic("データベース開ず(dbInsert)")
	}
	defer db.Close()

	db.Create(&kaitou)
}

//InsertLine : DBに新しいlineを追加
func InsertLine(line data.Line) {
	connect := ArgInit()
	db, err := gorm.Open("postgres", connect)
	if err != nil {
		panic("データベース開ず(InsertLine)")
	}
	defer db.Close()

	db.Create(&line)
}
