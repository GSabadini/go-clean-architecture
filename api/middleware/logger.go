package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"

	"github.com/pkg/errors"
	"github.com/urfave/negroni"
)

//Logger armazena a estrutura de logger para entrada e saídas da API
type Logger struct {
	logger logger.Logger
}

//NewLoggerFactory constrói um Logger com suas dependências
func NewLogger(log logger.Logger) Logger {
	return Logger{logger: log}
}

//Execute cria logs de entrada e saída da API
func (l Logger) Execute(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	const (
		logKey      = "logger_middleware"
		requestKey  = "api_request"
		responseKey = "api_response"
	)

	body, err := getRequestPayload(r)
	if err != nil {
		l.logger.WithFields(logger.Fields{
			"key":         logKey,
			"http_status": http.StatusBadRequest,
			"error":       err.Error(),
		}).Errorf("error when getting payload")

		return
	}

	l.logger.WithFields(logger.Fields{
		"key":         requestKey,
		"payload":     body,
		"url":         r.URL.Path,
		"http_method": r.Method,
	}).Infof("started handling request")

	next.ServeHTTP(w, r)

	end := time.Since(start).Seconds()
	res := w.(negroni.ResponseWriter)
	l.logger.WithFields(logger.Fields{
		"key":           responseKey,
		"url":           r.URL.Path,
		"http_method":   r.Method,
		"http_status":   res.Status(),
		"response_time": end,
	}).Infof("completed handling request")
}

func getRequestPayload(r *http.Request) (string, error) {
	if r.Body == nil {
		return "", errors.New("body not defined")
	}

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", errors.Wrap(err, "error read body")
	}

	// re-adiciona o payload ao buffer da request
	r.Body = ioutil.NopCloser(bytes.NewBuffer(payload))

	return strings.TrimSpace(string(payload)), nil
}
