package db

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/tuckKome/fictionary/data"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

//Openの引数を作る
func ArgInit() string {
	host := getEnv("FICTIONARY_DATABASE_HOST", "127.0.0.1")
	port := getEnv("FICTIONARY_PORT", "5432")
	user := getEnv("FICTIONARY_USER", "tahoiya")
	dbname := getEnv("FICTIONARY_DB_NAME", "dbtahoiya")
	password := getEnv("FICTIONARY_DB_PASS", "password")

	dbinfo := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=require",
		user,
		password,
		host,
		port,
		dbname,
	)
	return dbinfo
}

//DB初期化
func DBInit() {
	connect := ArgInit()
	db, err := gorm.Open("postgres", connect)
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&data.Kaitou{}, &data.Game{})
	defer db.Close()
}

//DBから一つ取り出す：回答ページで使用
func DBGetOne(id int) data.Game {
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
func DBGetKaitou(id int) []data.Kaitou {
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

func DBGetGames() []data.Game {
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

func DBInsert(kaitou data.Kaitou) {
	connect := ArgInit()
	db, err := gorm.Open("postgres", connect)
	if err != nil {
		panic("データベース開ず(dbInsert)")
	}
	defer db.Close()

	db.Create(&kaitou)
}
