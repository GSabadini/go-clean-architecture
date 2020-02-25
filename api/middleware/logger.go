package middleware

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

//Logger armazena a estrutura de logger para entrada e saídas da API
type Logger struct {
	log *logrus.Logger
}

//NewLogger constrói um logger com suas dependências
func NewLogger(log *logrus.Logger) Logger {
	return Logger{log: log}
}

//Logging cria logs de entrada e saída da API
func (l Logger) Logging(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	_, err := getRequestPayload(r)
	if err != nil {
		l.log.Warningln("error when getting payload")
	}

	//@TODO CORRIGIR BODY
	l.log.WithFields(logrus.Fields{
		"key": "api_request",
		//"payload":     body,
		"url":         r.URL.Path,
		"http_method": r.Method,
	}).Info("request received by the API")

	next.ServeHTTP(w, r)

	l.log.WithFields(logrus.Fields{
		"key":         "api_response",
		"url":         r.URL.Path,
		"http_method": r.Method,
	}).Info("response returned from the API")
}

func getRequestPayload(r *http.Request) (string, error) {
	if r.Body == nil {
		return "", errors.New("body not defined")
	}

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	// re-adiciona o payload ao buffer da request para ser lido por outros handlers
	r.Body = ioutil.NopCloser(bytes.NewBuffer(payload))

	return strings.TrimSpace(string(payload)), nil
}
