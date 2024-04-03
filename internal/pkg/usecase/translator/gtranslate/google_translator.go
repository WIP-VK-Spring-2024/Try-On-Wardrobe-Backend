package gtranslate

import (
	"try-on/internal/pkg/domain"

	google "github.com/gilang-as/google-translate"
)

type GoogleTranslator struct{}

func (g *GoogleTranslator) Translate(source string, sourceLang, targetLang domain.Language) (string, error) {
	translated, err := google.Translator(google.Translate{
		Text: source,
		From: string(sourceLang),
		To:   string(targetLang),
	})
	if err != nil {
		return "", err
	}
	return translated.Text, nil
}
