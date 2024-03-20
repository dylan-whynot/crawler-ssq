package ssq

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var DB *sql.DB

var sql1 = "INSERT INTO ssq(id,date,week,red_numbers,red_number1,red_number2,red_number3,red_number4,red_number5,red_number6,blue,sales,pool_amount) values(?,?,?,?,?,?,?,?,?,?,?,?,?)"
var sql2 = "INSERT INTO ssq_prizegrade(code,number,people_number,money) values(?,?,?,?)"

func init() {
	db, err := sql.Open("sqlite3", "D:\\sqlite\\db\\ssq.db")
	DB = db
	if err != nil {
		log.Fatalln("打开数据库异常", err)
		return
	}
}
func CloseDB() {
	DB.Close()
}
func InsertDatas(ssqs []Ssq) {
	if len(ssqs) == 0 {
		return
	}
	prepare1, err := DB.Prepare(sql1)
	if err != nil {
		log.Fatalln("数据库PrePare 异常", err)
		return
	}
	defer prepare1.Close()
	prepare2, err := DB.Prepare(sql2)
	if err != nil {
		log.Fatalln("数据库PrePare 2异常", err)
		return
	}
	defer prepare2.Close()
	tx, err := DB.Begin()
	if err != nil {
		log.Fatalln("数据库事务声明异常", err)
	}
	defer func() {
		if e := recover(); e != nil {
			tx.Rollback()
			log.Fatalln("执行写入异常", e)
		}
	}()
	for _, v := range ssqs {
		_, err := prepare1.Exec(v.Id, v.Date, v.Week, v.Red_numbers, v.Red_number1, v.Red_number2, v.Red_number3, v.Red_number4, v.Red_number5, v.Red_number6, v.Blue, v.Sales, v.Pool_amount)
		if err != nil {
			log.Println("写入ssq异常", v)
			panic(err)
		}
		for _, v2 := range v.Prizegrades {
			_, err := prepare2.Exec(v2.code, v2.number, v2.people_number, v2.money)
			if err != nil {
				log.Fatalln("写入ssq_prizegrade异常", v2)
				panic(err)
			}
		}
	}
	tx.Commit()
}
