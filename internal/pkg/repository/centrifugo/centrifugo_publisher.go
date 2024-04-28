package centrifugo

import (
	"context"

	"try-on/internal/generated/proto/centrifugo"
	"try-on/internal/middleware"
	"try-on/internal/pkg/domain"

	"github.com/mailru/easyjson"
	"google.golang.org/grpc"
)

type CentrifugoPublisher struct {
	centrifugo centrifugo.CentrifugoApiClient
}

func New(conn grpc.ClientConnInterface) domain.ChannelPublisher[easyjson.Marshaler] {
	return &CentrifugoPublisher{
		centrifugo: centrifugo.NewCentrifugoApiClient(conn),
	}
}

func (h CentrifugoPublisher) Publish(ctx context.Context, channel string, message easyjson.Marshaler) error {
	logger := middleware.GetLogger(ctx)

	payload, _ := easyjson.Marshal(message)

	logger.Infow("centrifugo", "channel", channel, "payload", string(payload))

	centrifugoResp, err := h.centrifugo.Publish(
		ctx,
		&centrifugo.PublishRequest{
			Channel: channel,
			Data:    payload,
		},
	)

	switch {
	case err != nil:
		logger.Errorw("centrifugo", "err", err)
		return err
	case centrifugoResp.Error != nil:
		logger.Errorw("centrifugo", "err", centrifugoResp.Error.Message)
	}
	return nil
}
