package action

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gsabadini/go-bank-transfer/api/response"
	"github.com/gsabadini/go-bank-transfer/domain"
	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
	"github.com/gsabadini/go-bank-transfer/usecase"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/pkg/errors"
)

//Transfer armazena as dependências de uma transferência
type Transfer struct {
	logger  logger.Logger
	usecase usecase.TransferUseCase
}

//NewTransfer constrói uma transferência com suas dependências
func NewTransfer(usecase usecase.TransferUseCase, log logger.Logger) Transfer {
	return Transfer{usecase: usecase, logger: log}
}

//Store é um handler para criação de transferência
func (t Transfer) Store(w http.ResponseWriter, r *http.Request) {
	const logKey = "create_transfer"

	var input usecase.TransferInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		t.logError(
			logKey,
			"error when decoding json",
			http.StatusBadRequest,
			err,
		)

		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}
	defer r.Body.Close()

	if err := validate(input); len(err) > 0 {
		t.logError(
			logKey,
			"input invalid",
			http.StatusBadRequest,
			errors.New("validate"),
		)

		response.NewMessagesError(err, http.StatusBadRequest).Send(w)
		return
	}

	result, err := t.usecase.Store(input)
	if err != nil {
		switch err {
		case domain.ErrInsufficientBalance:
			t.logError(
				logKey,
				"insufficient balance",
				http.StatusUnprocessableEntity,
				err,
			)

			response.NewError(err, http.StatusUnprocessableEntity).Send(w)
			return
		default:
			t.logError(
				logKey,
				"error when creating a new transfer",
				http.StatusInternalServerError,
				err,
			)

			response.NewError(err, http.StatusInternalServerError).Send(w)
			return
		}
	}
	t.logSuccess(logKey, "success create transfer", http.StatusCreated)

	response.NewSuccess(result, http.StatusCreated).Send(w)
}

//Index é um handler para retornar a lista de transferências
func (t Transfer) Index(w http.ResponseWriter, _ *http.Request) {
	const logKey = "index_transfer"

	result, err := t.usecase.FindAll()
	if err != nil {
		t.logError(
			logKey,
			"error when returning the transfer list",
			http.StatusInternalServerError,
			err,
		)

		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}
	t.logSuccess(logKey, "success when returning transfer list", http.StatusOK)

	response.NewSuccess(result, http.StatusOK).Send(w)
}

func validate(input usecase.TransferInput) []string {
	var messages []string

	var errAccountsEquals = errors.New("account origin equals destination account")
	if input.AccountOriginID == input.AccountDestinationID {
		messages = append(messages, errAccountsEquals.Error())
	}

	translator := en.New()
	uni := ut.New(translator, translator)
	trans, found := uni.GetTranslator("en")
	if !found {
		log.Fatal("translator not found")
	}
	v := validator.New()
	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		log.Fatal(err)
	}

	//_ = v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
	//	return ut.Add("required", "{0}xx is a required field!", true) // see universal-translator for details
	//}, func(ut ut.Translator, fe validator.FieldError) string {
	//	t, _ := ut.T("required", fe.Field())
	//	return t
	//})

	var err = v.Struct(input)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			messages = append(messages, strings.ToLower(e.Translate(trans)))
		}
	}

	return messages
}

func (t Transfer) logSuccess(key string, message string, httpStatus int) {
	t.logger.WithFields(logger.Fields{
		"key":         key,
		"http_status": httpStatus,
	}).Infof(message)
}

func (t Transfer) logError(key string, message string, httpStatus int, err error) {
	t.logger.WithFields(logger.Fields{
		"key":         key,
		"http_status": httpStatus,
		"error":       err.Error(),
	}).Errorf(message)
}
