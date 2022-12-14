package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var students = []Student{
	{Id: 1, Name: "Sude", Class: "1-c", Teacher: "Mehmet"},
	{Id: 2, Name: "Atakan", Class: "1-b", Teacher: "Kemal"},
}

type Student struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Class   string `json:"class"`
	Teacher string `json:"teacher"`
}

func listStudents(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, students)
}
func listStudentsById(context *gin.Context) {
	var studentByUserId Student
	err := context.BindJSON(&studentByUserId)
	fmt.Println()
	for i := 0; i < len(students); i++ {
		if students[i].Id == studentByUserId.Id && err == nil {
			context.IndentedJSON(http.StatusOK, gin.H{"message": "Student has been found", "student information": students[i]})
			return
		}
	}
}
func createStudent(context *gin.Context) {
	var studentByUser Student
	err := context.BindJSON(&studentByUser)

	if err == nil && studentByUser.Id != 0 && studentByUser.Class != "" && studentByUser.Name != "" && studentByUser.Teacher != "" {
		students = append(students, studentByUser)
		addStudentOnDatabase(&studentByUser)
		context.IndentedJSON(http.StatusCreated, gin.H{"message": "Student has been created", "student_id": studentByUser.Id})
		return
	} else {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Student has not been created"})
		return
	}
}

func addStudentOnDatabase(studentByUser *Student) {

	db, err := sql.Open("mysql", "root:12345@tcp(localhost:3306)/studentApi")

	if err != nil {
		fmt.Println("Error validating sql.Open arguments")
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Error verifying connection with db.Ping")
	}

	insert, err := db.Prepare("INSERT INTO `studentApi`.`student` (`student_id`,`student_name`,`student_class`,`student_teacher`) VALUES (?,?,?,?); ")
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
	res, err := insert.Exec(studentByUser.Id, studentByUser.Name, studentByUser.Class, studentByUser.Teacher)

	rowsAffec, _ := res.RowsAffected()
	if err != nil || rowsAffec != 1 {
		fmt.Println("Error inserting row:", err)
		return
	}

	lastInserted, _ := res.LastInsertId()
	rowsAffected, _ := res.RowsAffected()
	fmt.Println("ID of last row inserted:", lastInserted)
	fmt.Println("number of rows affected:", rowsAffected)

	fmt.Println("Successful Connection to Database!")
}

func getStudentByID(int_id int) (*Student, error) {
	for i, s := range students {
		if s.Id == int_id {
			return &students[i], nil
		}

	}
	return nil, errors.New("Student cannot be found!")
}

func getStudent(context *gin.Context) {
	str_id := context.Param("id")
	int_id, err := strconv.Atoi(str_id)

	if err != nil {
		panic(err)
	}

	student, err := getStudentByID(int_id)
	if err == nil {
		context.IndentedJSON(http.StatusOK, student)
	} else {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Student cannot be found!"})
	}

}
func main() {

	router := gin.Default()
	router.GET("/students", listStudents)
	router.POST("/students", createStudent)
	router.GET("/studentById", listStudentsById)
	router.GET("/students/:id", getStudent)
	router.Run("localhost:9090")

}
