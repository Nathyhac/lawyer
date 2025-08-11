package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Envelop map[string]interface{}

func WriteJson(w http.ResponseWriter, status int, data Envelop) error {
	js, err := json.MarshalIndent(data, "", "")

	if err != nil {
		return err
	}

	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func IdReader(r *http.Request) (int64, error) {
	IdParam := chi.URLParam(r, "id")
	if IdParam == "" {
		return 0, errors.New("id cannot be empty")
	}

	parsedId, err := strconv.ParseInt(IdParam, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error happend :%v ", err)
	}

	return parsedId, nil

}
