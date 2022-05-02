package errorhandler

import (
	"log"
	"net/http"

	"github.com/Ruthvik10/shared-library/data"
	"github.com/Ruthvik10/shared-library/helper"
)

func logError(logger *log.Logger, r *http.Request, err error) {
	logger.Println(err)
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, logger *log.Logger, statusCode int, errMsg interface{}) {
	env := data.Envelope{"error": errMsg}
	err := helper.WriteJSON(w, env, statusCode, nil)
	if err != nil {
		logError(logger, r, err)
	}
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, logger *log.Logger, err error) {
	logError(logger, r, err)
	msg := "the server encountered an problem and could not process your request"
	ErrorResponse(w, r, logger, http.StatusInternalServerError, msg)
}

func BadRequestErrorResponse(w http.ResponseWriter, r *http.Request, logger *log.Logger, err error) {
	logError(logger, r, err)
	msg := err.Error()
	ErrorResponse(w, r, logger, http.StatusBadRequest, msg)
}

func NotFoundErrorResponse(w http.ResponseWriter, r *http.Request, logger *log.Logger, err error) {
	logError(logger, r, err)
	msg := "could not find the requested resource"
	ErrorResponse(w, r, logger, http.StatusNotFound, msg)
}
