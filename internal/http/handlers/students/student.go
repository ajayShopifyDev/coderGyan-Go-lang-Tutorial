package students

import (
	"codersGyan/crud/internal/storage"
	"codersGyan/crud/internal/types"
	"codersGyan/crud/internal/utils/response"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//w.Write([]byte("welcome to students API"))
		slog.Info("Student is creating")
		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("Empty Body..")))
			return
		}

		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//validate RequestData

		if errs := validator.New().Struct(student); errs != nil {
			validateErrors := errs.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidateError(validateErrors))
			return
		}

		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)
		slog.Info("user created sucessfully", slog.String("userId", fmt.Sprint(lastId)))
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJSON(w, http.StatusAccepted, map[string]int64{"id": lastId})
	}
}

func GETbyID(storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("getting a student", slog.String("id:", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		student, err := storage.GetStudentbyID(intId)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJSON(w, http.StatusOK,student)
	}
}
