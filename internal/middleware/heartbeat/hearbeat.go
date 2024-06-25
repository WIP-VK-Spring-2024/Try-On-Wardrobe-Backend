package heartbeat

import (
	"net/http"

	"try-on/internal/generated/proto/centrifugo"
	"try-on/internal/pkg/common"

	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type Dependencies struct {
	DB         *pgxpool.Pool
	Centrifugo *grpc.ClientConn
	Redis      *redis.Client
}

//easyjson:json
type heartbeatResponse struct {
	DB         string
	Centrifugo string
	Redis      string
}

func Heartbeat(deps Dependencies) func(*fiber.Ctx) error {
	centrifugoClient := centrifugo.NewCentrifugoApiClient(deps.Centrifugo)

	return func(ctx *fiber.Ctx) error {
		err := deps.DB.Ping(ctx.UserContext())
		if err != nil {
			return ctx.Status(http.StatusServiceUnavailable).
				JSON(&heartbeatResponse{
					DB: err.Error(),
				})
		}

		resp, err := centrifugoClient.Info(ctx.UserContext(), &centrifugo.InfoRequest{})
		if err != nil {
			return ctx.Status(http.StatusServiceUnavailable).
				JSON(&heartbeatResponse{
					Centrifugo: err.Error(),
				})
		}
		if resp.Error != nil {
			return ctx.Status(http.StatusServiceUnavailable).
				JSON(&heartbeatResponse{
					Centrifugo: resp.Error.Message,
				})
		}

		result, err := deps.Redis.Ping().Result()
		if err != nil {
			return ctx.Status(http.StatusServiceUnavailable).
				JSON(&heartbeatResponse{
					Redis: result,
				})
		}

		return ctx.SendString(common.EmptyJson)
	}
}
