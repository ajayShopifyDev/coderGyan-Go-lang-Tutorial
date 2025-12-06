package students

import (
	"codersGyan/crud/internal/types"
	"codersGyan/crud/internal/utils/response"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//w.Write([]byte("welcome to students API"))
		slog.Info("Student is creating")
		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("Empty Body..")))
			return
		}
		
		if err != nil{
			response.WriteJSON(w,http.StatusBadRequest,response.GeneralError(err))
			return 
		} 
   
		//validate RequestData

		 if errs:= validator.New().Struct(student); errs!=  nil{
			validateErrors := errs.(validator.ValidationErrors)
			response.WriteJSON(w,http.StatusBadRequest,response.ValidateError(validateErrors))
		 }


		response.WriteJSON(w, http.StatusAccepted, map[string]string{"status": "Ok","id":string(student.ID)})
	}
}
