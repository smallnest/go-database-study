package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	empNo     int
	birthDate time.Time
	firstName string
	lastName  string
	gender    string
	hireDate  time.Time
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 0. drivers: https://github.com/golang/go/wiki/SQLDrivers
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/employees?charset=utf8&parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 1. ping
	if err := db.Ping(); err != nil {
		fmt.Println("can't connect mysql:", err)
	}

	// 2. query
	testQuery(db)

	// 3. prepared statement
	testStatement(db)

	// 4. QueryRow
	err = db.QueryRow("select first_name from employees where emp_no = ?", 10001).Scan(&firstName)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(firstName)

	// 5. exec
	testExec(db)

	// 6. transaction
	testTransaction(db)

	// 7. testNULL
	testNULL(db)

	// 8. testConnectingPool
	testNULL(db)
}

func testQuery(db *sql.DB) {

	rows, err := db.Query("select * from employees limit 10")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&empNo, &birthDate, &firstName, &lastName, &gender, &hireDate)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(empNo, birthDate.Format("2006-01-02"), firstName, lastName, gender, hireDate.Format("2006-01-02"))
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func testStatement(db *sql.DB) {
	stmt, err := db.Prepare("select * from employees where emp_no > ? limit 10")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(10000)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&empNo, &birthDate, &firstName, &lastName, &gender, &hireDate)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(empNo, birthDate.Format("2006-01-02"), firstName, lastName, gender, hireDate.Format("2006-01-02"))
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

}

func testExec(db *sql.DB) {
	var maxID int
	err := db.QueryRow("select max(emp_no) from employees").Scan(&maxID)
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("INSERT INTO employees(emp_no,birth_date,first_name,last_name,gender,hire_date) VALUES(?,?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(maxID+1, time.Now().Add(-20*365*24*time.Hour), "xin", "wong", "F", time.Now())
	if err != nil {
		log.Fatal(err)
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID = %d, affected = %d\n", lastID, rowCnt)
}

func testTransaction(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("select max(emp_no) from employees")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		_, err = stmt.Exec()
		if err != nil {
			log.Fatal(err)
		}
	}
	if err = stmt.Close(); err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("commit tx")
}

func testNULL(db *sql.DB) {
	var firstName sql.NullString
	err := db.QueryRow("select first_name from employees where emp_no = ?", 10001).Scan(&firstName)
	if err != nil {
		log.Fatal(err)
	}

	if firstName.Valid {
		log.Println(firstName.String)
	} else {
		log.Println("null value")
	}

}

func testConnectingPool(db *sql.DB) {
	db.SetConnMaxLifetime(8 * time.Hour)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)

	err := db.QueryRow("select first_name from employees where emp_no = ?", 10001).Scan(&firstName)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(firstName)
}
