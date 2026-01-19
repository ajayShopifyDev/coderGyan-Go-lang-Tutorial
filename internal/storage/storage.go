package storage

import "codersGyan/crud/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentbyID(id int64) (types.Student, error)
	GetAllStudents()([]types.Student,error)
	UpdateStudentbyID(id int64,name string, email string, age int)(types.Student,error)
	DeleteStudentByID(id int64)(int64,error)
}
