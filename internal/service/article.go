package service

import (
	"context"
	"time"

	"github.com/tsukaychan/webook/internal/events"

	"github.com/tsukaychan/webook/internal/domain"
	"github.com/tsukaychan/webook/internal/repository"
	"github.com/tsukaychan/webook/pkg/logger"
)

var _ ArticleService = (*articleService)(nil)

//go:generate mockgen -source=./article.go -package=svcmocks -destination=mocks/article.mock.go ArticleService
type ArticleService interface {
	// author

	Save(ctx context.Context, atcl domain.Article) (int64, error)
	Publish(ctx context.Context, atcl domain.Article) (int64, error)
	Withdraw(ctx context.Context, id, authorId int64) error
	List(ctx context.Context, authorId int64,
		offset, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)

	// reader

	GetPublishedById(ctx context.Context, id, uid int64) (domain.Article, error)
	ListPub(ctx context.Context, start time.Time, offset, limit int) ([]domain.Article, error)
}

type articleService struct {
	articleRepo repository.ArticleRepository
	producer    events.Producer
	logger      logger.Logger
}

func NewArticleService(articleRepo repository.ArticleRepository, producer events.Producer, logger logger.Logger) ArticleService {
	return &articleService{
		articleRepo: articleRepo,
		producer:    producer,
		logger:      logger,
	}
}

func (svc *articleService) Save(ctx context.Context, atcl domain.Article) (int64, error) {
	atcl.Status = domain.ArticleStatusUnpublished
	if atcl.Id > 0 {
		err := svc.articleRepo.Update(ctx, atcl)
		return atcl.Id, err
	}

	return svc.articleRepo.Create(ctx, atcl)
}

func (svc *articleService) Publish(ctx context.Context, atcl domain.Article) (int64, error) {
	atcl.Status = domain.ArticleStatusPublished
	return svc.articleRepo.Sync(ctx, atcl)
}

func (svc *articleService) Withdraw(ctx context.Context, id, authorId int64) error {
	return svc.articleRepo.SyncStatus(ctx, id, authorId, domain.ArticleStatusPrivate)
}

func (svc *articleService) List(ctx context.Context, authorId int64, offset, limit int) ([]domain.Article, error) {
	return svc.articleRepo.List(ctx, authorId, offset, limit)
}

func (svc *articleService) GetById(ctx context.Context, id int64) (domain.Article, error) {
	return svc.articleRepo.GetById(ctx, id)
}

func (svc *articleService) GetPublishedById(ctx context.Context, id, uid int64) (domain.Article, error) {
	atcl, err := svc.articleRepo.GetPublishedById(ctx, id)
	if err == nil {
		go func() {
			er := svc.producer.ProduceReadEvent(events.ReadEvent{
				Aid: id,
				Uid: uid,
			})
			if er != nil {
				svc.logger.Error("send reader read event failed",
					logger.Int64("uid", uid),
					logger.Int64("aid", id),
					logger.Error(err))
			}
		}()
	}
	return atcl, err
}

func (svc *articleService) ListPub(ctx context.Context, start time.Time, offset, limit int) ([]domain.Article, error) {
	return svc.articleRepo.ListPub(ctx, start, offset, limit)
}
