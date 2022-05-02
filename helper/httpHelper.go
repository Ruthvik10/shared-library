package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type envelope map[string]interface{}

func readJSON(w http.ResponseWriter, r *http.Request, dest interface{}) error {

	maxBytes := 1048576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(dest)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("the body contains badly formed json (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("the body contains badly formed json")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("type mismatch for field %q in the request", unmarshalTypeError.Field)
			}
			return fmt.Errorf("type mismatch in the request at %d", unmarshalTypeError.Offset)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		case errors.Is(err, io.EOF):
			return errors.New("the body cannot be empty")

		default:
			return err
		}
	}

	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must contain only a single JSON value")
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, data envelope, statusCode int, headers http.Header) error {
	res, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	for k, v := range headers {
		w.Header()[k] = v
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(res)
	return nil
}
