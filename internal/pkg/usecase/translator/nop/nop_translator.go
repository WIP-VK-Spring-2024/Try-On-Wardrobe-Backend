package nop

import "try-on/internal/pkg/domain"

type NopTranslator struct{}

func (tr *NopTranslator) Translate(source string, _, _ domain.Language) (string, error) {
	return source, nil
}
