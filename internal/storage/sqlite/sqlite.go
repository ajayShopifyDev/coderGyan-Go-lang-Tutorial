package sqlite

import (
	"codersGyan/crud/internal/config"
	"codersGyan/crud/internal/types"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
	//_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {

	stmt, err := s.Db.Prepare("INSERT INTO students (name,email,age) VALUES (?,?,?)")
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)

	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (s *Sqlite) GetStudentbyID(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM students WHERE id =? LIMIT 1")

	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()

	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.ID, &student.Name, &student.Email, &student.Age)

	if err != nil {
		if err != sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("Student id is wrong %s", id)
		}
		return types.Student{}, fmt.Errorf("Query Error :%w", err)
	}

	return student, nil
}

func (s *Sqlite) GetAllStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id,name,email,age FROM students ")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []types.Student

	for rows.Next() {
		var student types.Student
		err = rows.Scan(&student.ID, &student.Name, &student.Email, &student.Age)

		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	return students, nil
}

func (s *Sqlite) UpdateStudentbyID(id int64, name string, email string, age int) (types.Student, error) {
	// 1. Check if student exists first
	_, err := s.GetStudentbyID(id) // Use your fixed GetStudentByID
	if err != nil {
		return types.Student{}, fmt.Errorf("student id %d not found", id)
	}

	// 2. Update only if exists
	query := `
        UPDATE students 
        SET name = ?, email = ?, age = ? 
        WHERE id = ?
    `
	result, err := s.Db.Exec(query, name, email, age, id)
	if err != nil {
		return types.Student{}, fmt.Errorf("update failed: %w", err)
	}

	// 3. Verify exactly 1 row was updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return types.Student{}, fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return types.Student{}, fmt.Errorf("student id %d not found", id)
	}

	// 4. Return updated student
	return s.GetStudentbyID(id)
}

func (s *Sqlite) DeleteStudentByID(id int64) (int64, error) {
	// 1. Verify student exists
	_, err := s.GetStudentbyID(id)
	if err != nil {
		return 0, fmt.Errorf("student id %d not found", id)
	}

	// 2. Delete the student
	stmt, err := s.Db.Prepare("DELETE FROM students WHERE id = ? ")
	if err != nil {
		return 0, fmt.Errorf("prepare delete failed: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return 0, fmt.Errorf("delete failed: %w", err)
	}

	// 3. Verify exactly 1 row deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return 0, fmt.Errorf("student id %d not found", id)
	}

	return rowsAffected, nil
}

// func(s *Sqlite) UpdateStudentbyID(id int64,name string, email string, age int)(types.Student,error){
// 	stmt,err:=s.Db.Prepare("SELECT * FROM students WHERE id=? LIMIT 1")

// 	if err != nil {
// 		return types.Student{},err
// 	}

// 	defer stmt.Close()

// 	var student types.Student
// 	err = stmt.QueryRow(id).Scan(&student.ID)

// 	if err != nil {
// 		if err != sql.ErrNoRows {
// 			return types.Student{}, fmt.Errorf("Student id is wrong %s", id)
// 		}
// 		return types.Student{}, fmt.Errorf("Query Error :%w", err)
// 	}

// 	newStmt, err := s.Db.Prepare("INSERT INTO students (name,email,age) VALUES (?,?,?)")
// 	if err != nil {
// 		return types.Student{}, err
// 	}

// 	defer newStmt.Close()

// 	updatedUser, err := stmt.Exec(name, email, age)

// 	if err != nil {
// 		return types.Student{}, err
// 	}

// 	return  updatedUser,nil
// }
