package internal

import (
	"database/sql"
	"encoding/json"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestConnectDatabase(t *testing.T) {
	db, err := sql.Open("sqlite3", "memory")
	if err != nil {
		t.Fatalf("Failed to create memory database: %v", err)
	}
	defer db.Close()

	app := &App{}

	err = app.connectDatabase("memory")
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
}

func TestCreateTable(t *testing.T) {
	db, err := sql.Open("sqlite3", "memory")
	if err != nil {
		t.Fatalf("Failed to create memory database: %v", err)
	}
	defer db.Close()

	app := &App{db: db}

	err = app.createTable()
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
}

func TestGetTasksDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	app := &App{db: db}

	rows := sqlmock.NewRows([]string{"id", "title", "completed"}).
		AddRow(1, "Task 1", true).
		AddRow(2, "Task 2", false)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM tasks")).WillReturnRows(rows)

	tasks, err := app.getTasks()
	if err != nil {
		t.Fatalf("Failed to get tasks: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, but got %d", len(tasks))
	}

	if tasks[0].Id != 1 || tasks[0].Title != "Task 1" || tasks[0].Completed != true {
		t.Errorf("Unexpected content for task 1")
	}

	if tasks[1].Id != 2 || tasks[1].Title != "Task 2" || tasks[1].Completed != false {
		t.Errorf("Unexpected content for task 2")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestAddTasksDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	app := &App{db: db}

	result := sqlmock.NewResult(1, 1)

	mock.ExpectPrepare("INSERT INTO tasks(.+)").ExpectExec().WillReturnResult(result)

	title := "New Task"
	completed := true
	data, err := app.addTask(title, completed)
	if err != nil {
		t.Fatalf("Failed to add task: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(data, &response); err != nil {
		t.Fatalf("Failed to unmarshal response JSON: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestDeleteTaskDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	app := &App{db: db}

	result := sqlmock.NewResult(1, 1)

	mock.ExpectPrepare("DELETE FROM tasks WHERE id = \\?;").ExpectExec().WillReturnResult(result)

	id := 123
	data, err := app.deleteTask(id)
	if err != nil {
		t.Fatalf("Failed to delete task: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(data, &response); err != nil {
		t.Fatalf("Failed to unmarshal response JSON: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}

	data, err = app.deleteTask(id)
	if err == nil {
		t.Fatalf("error delete non existed task: %v", err)
	}
}

func TestUpdateTaskDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	app := &App{db: db}

	task := Task{
		Id:        123,
		Title:     "Updated Task",
		Completed: true,
	}

	expectedAffectedRows := int64(1)
	result := sqlmock.NewResult(expectedAffectedRows, expectedAffectedRows)

	mock.ExpectPrepare("UPDATE tasks SET title = \\?, completed = \\? WHERE id = \\?;").ExpectExec().WillReturnResult(result)

	data, err := app.updateTask(task)
	if err != nil {
		t.Fatalf("Failed to update task: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(data, &response); err != nil {
		t.Fatalf("Failed to unmarshal response JSON: %v", err)
	}

	task = Task{
		Id:        1233333333,
		Title:     "Updated Task",
		Completed: true,
	}
	data, err = app.updateTask(task)
	if err == nil {
		t.Fatalf("error update non exiest task: %v", err)
	}
}