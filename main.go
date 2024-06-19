package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID       int
	Username string
	Password string
}

func main() {
	var wg sync.WaitGroup
	var inputUSername, inputPassword string
	var storedpassword string
	// 自信のパソコン（localhost:3306）にgo_sample_appというユーザー名でログイン
	// アクセスするデータベースはsample_app_data
	db, err := sql.Open("mysql","go_sample_app:Wtatsumi0317@tcp(localhost:3306)/sample_app_data")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// データベース内(sample_app_data)に、テーブルがなければ作成
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS users
		(
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(50) NOT NULL,
		password VARCHAR(50) NOT NULL
		)`)
	if err != nil {
		log.Fatal(err)
	}
	wg.Add(1)
	txtCh := make(chan string)
	dateCh := make(chan string)
	go Greet(txtCh, dateCh, &wg)
	fmt.Print("Username: ")
	// ターミナルにUsernameと表示した後に、fmt.Scanで標準入力を求め、Scanしたその入力値を、
	// 引数に取った変数のアドレス先に格納する
	fmt.Scan(&inputUSername)
	// ⇒fmt.Printとセット
	fmt.Print("Password: ")
	fmt.Scan(&inputPassword)
	// 一応SQLインジェクションを対策してみる
	err = db.QueryRow("SELECT password FROM users WHERE username = ?",inputUSername).Scan(&storedpassword)
	if err != nil {
		fmt.Println("failed. Please retry\n")
		return
	}
	if inputPassword == storedpassword {
		fmt.Println("Welcome!!\n")
		fmt.Printf("%s、%sさん。本日は%sです。",<-txtCh,inputUSername,<-dateCh)
	} else {
		fmt.Println("failed. Invalid Password")
	}
	wg.Wait()
}

func Greet(txtch chan <- string, datech chan <- string, wg *sync.WaitGroup){
	defer wg.Done()//channelをクローズしないと、gorutineが終わらので、一生実行されない。
	defer close(txtch)
	defer close(datech)
	currentTime := time.Now()
	if ( 5 < currentTime.Hour() && currentTime.Hour() < 11 ) {
		msg_to_channel(txtch, "おはようございます")
	}else if( 11 <= currentTime.Hour() && currentTime.Hour() < 18 ){
		msg_to_channel(txtch, "こんにちは")
	}else {
		msg_to_channel(txtch, "こんばんは")
	}
	msg_to_channel(datech, currentTime.Format("2006-01-02"))
}

func msg_to_channel(ch chan <- string, msg string){
	ch <- msg
}