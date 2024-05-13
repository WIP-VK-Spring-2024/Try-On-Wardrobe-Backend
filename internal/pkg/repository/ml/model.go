package ml

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"try-on/internal/middleware"
	"try-on/internal/pkg/domain"

	"github.com/mailru/easyjson"
)

type ModelAvailabilityChecker struct{}

func NewAvailabilityChecker() domain.MlModel {
	return &ModelAvailabilityChecker{}
}

func (m ModelAvailabilityChecker) IsAvailable(model string, ctx context.Context) (bool, error) {
	cfg := middleware.Config(ctx).ModelsHealth

	req, err := http.NewRequest(http.MethodGet, cfg.Endpoint+model, nil)
	if err != nil {
		return false, err
	}

	req.Header.Add(cfg.TokenHeader, cfg.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	fmt.Println("Got response", string(body))

	modelResp := domain.ModelHealthResponse{}
	err = easyjson.Unmarshal(body, &modelResp)
	if err != nil {
		return false, err
	}

	return modelResp.IsListening, nil
}
