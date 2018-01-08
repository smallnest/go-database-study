package main

import (
	"log"
	"time"

	"github.com/kr/pretty"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Employee struct {
	EmpNo     int       `db:"emp_no,omitempty"`
	BirthDate time.Time `db:"birth_date,omitempty"`
	FirstName string    `db:"first_name,omitempty"`
	LastName  string    `db:"last_name,omitempty"`
	Gender    string    `db:"gender,omitempty"`
	HireDate  time.Time `db:"hire_date,omitempty"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db, err := sqlx.Connect("mysql", "root@tcp(127.0.0.1:3306)/employees?charset=utf8&parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	employees := []Employee{}
	err = db.Select(&employees, "SELECT * FROM employees limit 10")
	if err != nil {
		log.Fatalln(err)
	}
	pretty.Println(employees)
}
