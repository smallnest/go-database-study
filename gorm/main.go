package main

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
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

	db, err := gorm.Open("mysql", "root@tcp(127.0.0.1:3306)/employees?charset=utf8&parseTime=true")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Employee{})

	// Read
	var employee Employee
	db.First(&employee, "emp_no = ?", "10001")

	// Create
	employee.EmpNo = 60000
	db.Create(&employee)

	// Update
	db.Model(&employee).Update("hire_date", time.Now())

	// Delete - delete product
	db.Delete(&employee)
}
