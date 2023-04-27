package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/c-robinson/iplib"
)

type ErrorResponse struct {
	StatusCode  int    `json:"statusCode"`
	Description string `json:"description"`
}

func NewErrorResponse(statusCode int, description string) ErrorResponse {
	return ErrorResponse{
		StatusCode:  statusCode,
		Description: description,
	}
}

func Contains(nets []iplib.Net4, toCheck iplib.Net4) bool {
	for _, n := range nets {
		if n.String() == toCheck.String() {
			return true
		}
	}
	return false
}

func WriteBody(w http.ResponseWriter, body any) {
	js, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write([]byte(js))
}

func Write400Error(w http.ResponseWriter) {
	writeError(w, 400, "Bad Request: the request body couldn't be parsed")
}

func Write500Error(w http.ResponseWriter) {
	writeError(w, 500, "Internal Server Error: unexpected error")
}

func Write406Error(w http.ResponseWriter, description string) {
	writeError(w, 406, description)
}

func writeError(w http.ResponseWriter, statusCode int, description string) {
	js, err := json.Marshal(NewErrorResponse(statusCode, description))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(statusCode)
		w.Write([]byte(js))
	}
}

func concat(args ...string) string {
	result := ""
	for _, s := range args {
		result = result + " " + s
	}
	return result
}

func ExecuteCommand(name string, args ...string) error {
	log.Printf("executing [%s %s]", name, concat(args...))
	_, err := exec.Command(name, args...).Output()
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	fmt.Printf("command %s successfully executed", name)
	return nil
}
