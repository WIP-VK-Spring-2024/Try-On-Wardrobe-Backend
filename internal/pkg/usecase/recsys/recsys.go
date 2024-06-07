package recsys

import (
	"context"
	"fmt"

	"try-on/internal/middleware"
	"try-on/internal/pkg/app_errors"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/repository/ml"
	"try-on/internal/pkg/utils"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

type Recsys struct {
	publisher    domain.Publisher[domain.RecsysRequest]
	subscriber   domain.Subscriber[domain.RecsysResponse]
	redis        *redis.Client
	feed         domain.FeedRepository
	availability domain.AvailabilityChecker
}

func New(
	feed domain.FeedRepository,
	redis *redis.Client,
	publisher domain.Publisher[domain.RecsysRequest],
	subscriber domain.Subscriber[domain.RecsysResponse],
) domain.Recsys {
	return &Recsys{
		redis:        redis,
		feed:         feed,
		publisher:    publisher,
		subscriber:   subscriber,
		availability: ml.NewAvailabilityChecker(),
	}
}

func (rec Recsys) Close() {
	rec.publisher.Close()
}

func (rec Recsys) makeRecsysRequest(ctx context.Context, request domain.RecsysRequest) error {
	key := recsysFlagKey(request.UserID)
	err := rec.redis.Get(key).Err()
	if err != redis.Nil {
		return nil
	}

	err = rec.redis.Set(key, true, 0).Err()
	if err != nil {
		return err
	}
	return rec.publisher.Publish(ctx, request)
}

func (rec Recsys) GetRecommendations(ctx context.Context, limit int, request domain.RecsysRequest) ([]domain.Post, error) {
	cfg := middleware.Config(ctx).ModelsHealth

	isAvailable, err := rec.availability.IsAvailable(cfg.Recsys, ctx)
	if err != nil {
		middleware.GetLogger(ctx).Warnw("recsys", "error", err)
		return nil, app_errors.ErrModelUnavailable
	}
	if !isAvailable {
		return nil, app_errors.ErrModelUnavailable
	}

	redisKey := recsysSetKey(request.UserID)

	result, err := rec.redis.SPopN(redisKey, int64(limit)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	fmt.Println("Got results from redis:", result)

	if len(result) == 0 {
		return nil, rec.makeRecsysRequest(ctx, request)
	}

	resultUUIDs := make([]utils.UUID, 0, len(result))
	for _, elem := range result {
		uuid, err := utils.ParseUUID(elem)
		if err != nil {
			return nil, err
		}
		resultUUIDs = append(resultUUIDs, uuid)
	}

	return rec.feed.GetPostsByOutfitIds(request.UserID, resultUUIDs)
}

func (rec Recsys) ListenResults(logger *zap.SugaredLogger) error {
	ctx := middleware.WithLogger(context.Background(), logger)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorw("recsys results", "error", err)
			}
		}()

		rec.subscriber.Listen(ctx, func(response *domain.RecsysResponse) domain.Result {
			if !utils.HttpOk(response.StatusCode) {
				logger.Warnw("rabbit", "queue", "recsys", "error", response.Message)
				err := rec.redis.Del(recsysFlagKey(response.UserID)).Err()
				if err != nil {
					logger.Errorw("recsys redis error", "error", err)
				}
				return domain.ResultDiscard
			}

			args := make([]interface{}, 0, len(response.OutfitIds))
			for _, outfitId := range response.OutfitIds {
				args = append(args, outfitId.String())
			}

			_, err := rec.redis.TxPipelined(func(pipeline redis.Pipeliner) error {
				pipeline.SAdd(recsysSetKey(response.UserID), args...)
				pipeline.Del(recsysFlagKey(response.UserID))
				return nil
			})
			if err != nil {
				logger.Errorw("recsys redis error", "error", err)
				return domain.ResultDiscard
			}

			return domain.ResultOk
		})
	}()

	return nil
}

func recsysSetKey(userId utils.UUID) string {
	return fmt.Sprintf("user:%s:recsys:items", userId.String())
}

func recsysFlagKey(userId utils.UUID) string {
	return fmt.Sprintf("user:%s:recsys:flag", userId.String())
}
