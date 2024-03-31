package heartbeat

import (
	"net/http"

	"try-on/internal/pkg/common"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type Dependencies struct {
	DB         *pgxpool.Pool
	Centrifugo *grpc.ClientConn
}

//easyjson:json
type heartbeatResponse struct {
	DB         string
	Centrifugo string
}

func Hearbeat(deps Dependencies) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		err := deps.DB.Ping(ctx.UserContext())
		if err != nil {
			return ctx.Status(http.StatusServiceUnavailable).
				JSON(&heartbeatResponse{
					DB: err.Error(),
				})
		}

		if deps.Centrifugo.GetState() != connectivity.Ready {
			return ctx.Status(http.StatusServiceUnavailable).
				JSON(&heartbeatResponse{
					Centrifugo: deps.Centrifugo.GetState().String(),
				})
		}

		return ctx.SendString(common.EmptyJson)
	}
}
