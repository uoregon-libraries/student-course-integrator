package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/uoregon-libraries/student-course-integrator/src/sci-server"
)

func main() {
	sciserver.Run()
}
