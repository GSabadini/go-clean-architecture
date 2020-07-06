package validator

import (
	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type goPlayground struct {
	validator *validator.Validate
	translate ut.Translator
	log       logger.Logger
	err       error
	msg       []string
}

//NewGoPlayground constrói uma instância do validator GoPlayground
func NewGoPlayground(log logger.Logger) Validator {
	language := en.New()
	uni := ut.New(language, language)
	translate, found := uni.GetTranslator("en")
	if !found {
		log.Fatalln("translator not found")
	}

	v := validator.New()

	if err := en_translations.RegisterDefaultTranslations(v, translate); err != nil {
		log.Fatalln("translator not found")
	}

	return &goPlayground{validator: v, translate: translate, log: log}
}

func (g *goPlayground) Validate(i interface{}) error {
	if len(g.msg) > 0 {
		g.msg = nil
	}

	g.err = g.validator.Struct(i)
	if g.err != nil {
		return g.err
	}

	return nil
}

func (g *goPlayground) Messages() []string {
	if g.err != nil {
		for _, err := range g.err.(validator.ValidationErrors) {
			g.msg = append(g.msg, err.Translate(g.translate))
		}
	}

	return g.msg
}
