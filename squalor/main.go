package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/square/squalor"
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
	_db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/employees?charset=utf8&parseTime=true")
	panicOnError(err)
	defer _db.Close()

	db, err := squalor.NewDB(_db)
	panicOnError(err)

	employee := &Employee{}
	employeeModel, err := db.BindModel("employees", employee)
	panicOnError(err)

	err = db.Get(employee, 10001)
	panicOnError(err)
	fmt.Printf("%v\n", employee)

	q := employeeModel.Select(employeeModel.All()).Where(employeeModel.C("emp_no").Gte(10000)).Limit(10)
	var results []Employee
	err = db.Select(&results, q)
	panicOnError(err)
	fmt.Printf("results: %v\n", results)

	q = employeeModel.Select(employeeModel.C("emp_no").Max())
	var maxID []int
	err = db.Select(&maxID, q)
	panicOnError(err)
	fmt.Printf("maxID: %v\n", maxID[0])

	employee.EmpNo = maxID[0] + 1
	err = db.Insert(employee)
	panicOnError(err)
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
