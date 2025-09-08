package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	js = append(js, '\n') // this is interesting

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {

	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	/*
		important note:- firstly r.Body is io.ReadCloser (its an interface)
		the function http.MaxBytesReader returns a io.ReadCloser (interface) I don't know the concrete type
		now we can't read more than the specified bytes from r.Body
		if we do it will automatically send the response 413 http.StatusRequestEntityTooLarge (413)
		that's why we need to provide w, i don't know exactly which interface type w satisfies, but you can write to it

	*/

	decoder := json.NewDecoder(r.Body)
	/*
		we have the decoder, here, its important to note that .decode() doesn't decode the entire request but only decodes the first json object
		so instead of directly applying .Decode(//dst) we must have this decoder so that we can call .Decode(//dst) more than once
	*/
	err := decoder.Decode(dst)
	if err != nil {
		var (
			syntaxError         *json.SyntaxError
			unmarshalTypeError  *json.UnmarshalTypeError
			invalidUnmarshalErr *json.InvalidUnmarshalError
		)

		switch {
		// Badly-formed JSON with location info.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		// Some syntax errors surface as io.ErrUnexpectedEOF.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		// Wrong JSON type for a field or at an offset.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		// Empty body.
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// Programmer error: passed a non-pointer or nil pointer into Decode.
		case errors.As(err, &invalidUnmarshalErr):
			panic(err)

		// Anything else: bubble up as-is.
		default:
			return err
		}
	}

	return nil
}
