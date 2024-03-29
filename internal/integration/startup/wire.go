//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/tsukaychan/webook/internal/repository"
	articleCache "github.com/tsukaychan/webook/internal/repository/cache/article"
	captchaCache "github.com/tsukaychan/webook/internal/repository/cache/captcha"
	cache "github.com/tsukaychan/webook/internal/repository/cache/interactive"
	"github.com/tsukaychan/webook/internal/repository/dao"
	articleDao "github.com/tsukaychan/webook/internal/repository/dao/article"
	"github.com/tsukaychan/webook/internal/service"
	"github.com/tsukaychan/webook/internal/web"
	ijwt "github.com/tsukaychan/webook/internal/web/jwt"
	"github.com/tsukaychan/webook/ioc"
)

var thirdProvider = wire.NewSet(InitRedis, InitTestDB, InitLog)

var userSvcProvider = wire.NewSet(
	service.NewUserService,
	repository.NewCachedUserRepository,
	dao.NewGORMUserDAO,
	userCache.NewUserRedisCache,
)

var articleSvcProvider = wire.NewSet(
	service.NewArticleService,
	repository.NewCachedArticleRepository,
	articleDao.NewGORMArticleDAO,
	articleCache.NewRedisArticleCache,
)

var interactiveSvcProvider = wire.NewSet(
	service.NewInteractiveService,
	repository.NewCachedInteractiveRepository,
	dao.NewGORMInteractiveDAO,
	cache.NewRedisInteractiveCache,
)

func InitWebServer() *gin.Engine {
	wire.Build(
		thirdProvider,

		userSvcProvider,
		articleSvcProvider,
		interactiveSvcProvider,

		web.NewUserHandler,
		web.NewArticleHandler,
		web.NewOAuth2Handler,

		service.NewCaptchaService,
		repository.NewCachedCaptchaRepository,
		captchaCache.NewCaptchaRedisCache,

		ioc.InitSMSService,
		InitPhantomWechatService,
		InitWechatHandlerConfig,
		ijwt.NewRedisJWTHandler,

		ioc.InitMiddlewares,
		ioc.InitLimiter,

		ioc.InitWebServer,
	)
	return gin.Default()
}

func InitArticleHandler(atclDao articleDao.ArticleDAO) *web.ArticleHandler {
	wire.Build(
		thirdProvider,
		interactiveSvcProvider,
		userSvcProvider,
		service.NewArticleService,
		repository.NewCachedArticleRepository,
		articleCache.NewRedisArticleCache,
		web.NewArticleHandler,
	)
	return &web.ArticleHandler{}
}

func InitInteractiveService() service.InteractiveService {
	wire.Build(thirdProvider, interactiveSvcProvider)
	return service.NewInteractiveService(nil, nil)
}

func InitUserSvc() service.UserService {
	wire.Build(
		thirdProvider,
		userSvcProvider,
	)
	return service.NewUserService(nil, nil)
}

func InitJwtHdl() ijwt.Handler {
	// wire.Build(thirdProvider, ijwt.NewRedisJWTHandler)
	return ijwt.NewRedisJWTHandler(nil)
}
