package delivery

// import (
// 	"context"

// 	"try-on/internal/generated/proto/centrifugo"
// 	"try-on/internal/pkg/domain"

// 	"github.com/mailru/easyjson"
// 	"go.uber.org/zap"
// )

// type CentrifugoHandler struct {
// 	logger     *zap.SugaredLogger
// 	centrifugo centrifugo.CentrifugoApiClient
// }

// func (h CentrifugoHandler) Publish(message easyjson.Marshaler, channel string) error {
// 	payload, _ := easyjson.Marshal(message)

// 	h.logger.Infow("centrifugo", "channel", channel, "payload", string(payload))

// 	centrifugoResp, err := h.centrifugo.Publish(
// 		context.Background(),
// 		&centrifugo.PublishRequest{
// 			Channel: channel,
// 			Data:    payload,
// 		},
// 	)

// 	switch {
// 	case err != nil:
// 		h.logger.Errorw(err.Error())
// 	case centrifugoResp.Error != nil:
// 		h.logger.Errorw(centrifugoResp.Error.Message)
// 	default:
// 		return nil
// 	}

// 	return nil
// }
