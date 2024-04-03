package domain

type Translator interface {
	Translate(source string, sourceLang, targetLang Language) (string, error)
}

type Language string

const (
	LanguageRU = "ru"
	LanguageEN = "en"
)
