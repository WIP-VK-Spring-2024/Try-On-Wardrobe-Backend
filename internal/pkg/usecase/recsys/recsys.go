package recsys

import (
	"context"
	"fmt"

	"try-on/internal/middleware"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

type Recsys struct {
	publisher  domain.Publisher[domain.RecsysRequest]
	subscriber domain.Subscriber[domain.RecsysResponse]
	redis      *redis.Client
	feed       domain.FeedRepository
}

func New(
	feed domain.FeedRepository,
	redis *redis.Client,
	publisher domain.Publisher[domain.RecsysRequest],
	subscriber domain.Subscriber[domain.RecsysResponse],
) domain.Recsys {
	return &Recsys{
		redis:      redis,
		feed:       feed,
		publisher:  publisher,
		subscriber: subscriber,
	}
}

func (rec Recsys) Close() {
	rec.publisher.Close()
}

func (rec Recsys) GetRecommendations(ctx context.Context, limit int, request domain.RecsysRequest) ([]domain.Post, error) {
	result, err := rec.redis.SRandMemberN(recsysSetKey(request.UserID), int64(request.SamplesAmount)).Result()
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, rec.publisher.Publish(ctx, request)
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
			args := make([]interface{}, 0, len(response.OutfitIds))
			for _, outfitId := range response.OutfitIds {
				args = append(args, outfitId)
			}

			result := rec.redis.SAdd(recsysSetKey(response.UserID), args...)
			if err := result.Err(); err != nil {
				logger.Errorw("recsys redis error", "error", err)
				return domain.ResultDiscard
			}

			return domain.ResultOk
		})
	}()

	return nil
}

func recsysSetKey(userId utils.UUID) string {
	return fmt.Sprintf("user:%s:recsys", userId.String())
}
