package main

import (
	"auth/internal/audit"
	"auth/internal/cache"
	"auth/internal/config"
	"auth/internal/controller"
	localCache "auth/internal/infra/cache"
	"auth/internal/infra/casbin"
	"auth/internal/infra/db"
	"auth/internal/infra/redis"
	"auth/internal/repo"
	"auth/internal/service"
	"auth/internal/util"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MicahParks/keyfunc/v3"
)

func main() {
	// 读取配置并初始化结构化日志
	cfg := config.Load()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	// 启动阶段统一使用短超时上下文，避免外部依赖异常导致卡住
	bootCtx, bootCancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer bootCancel()

	// 初始化数据库（含自动迁移），失败即终止启动
	database, err := db.NewSQLite(cfg.SQLitePath, logger)
	if err != nil {
		logger.Error("init sqlite failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// 初始化 L1 本地缓存（所有缓存组件共享）
	l1, err := localCache.NewLocalCache(24 * time.Hour)
	if err != nil {
		logger.Error("init bigcache failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		if cerr := l1.Close(); cerr != nil {
			logger.Warn("close bigcache failed", slog.String("error", cerr.Error()))
		}
	}()

	// 初始化 Redis L2 缓存；不可用时自动降级为仅 L1
	redisClient := redis.NewRedisClient(bootCtx, cfg.RedisAddrs, cfg.RedisPassword, cfg.RedisDB, logger)
	if redisClient != nil {
		defer func() {
			if cerr := redisClient.Close(); cerr != nil {
				logger.Warn("close redis failed", slog.String("error", cerr.Error()))
			}

		}()
	}

	// 组装缓存与仓储，服务层通过接口访问，便于替换实现和测试
	emailMapCache := cache.NewEmailUserIDCache(l1, redisClient, logger)
	revokedMapCache := cache.NewRevokedAtUserIDCache(l1, redisClient, cfg.BufferSize, database, logger)
	emailBloom := cache.NewEmailBloomFilter(redisClient, logger, cfg.EmailBloomKey, cfg.EmailBloomErr, cfg.EmailBloomCap)
	tokenStoreCache := cache.NewTokenStore(l1, redisClient, logger)
	userByIDCache := cache.NewUserByIDCache(l1, redisClient, logger, cfg.UserByIDTTL, cfg.UserByIDNullTTL)

	userRepoInner := repo.NewUserRepoInner()
	userRepo := repo.NewUserRepoImpl(userRepoInner, userByIDCache, logger)

	jwks, err := keyfunc.NewDefault(cfg.KeyFuncURLs)
	if err != nil {
		panic(err)
	}
	// jwtManager := util.NewJWTManager(cfg.JWTIssuer, cfg.JWTAccessSecret, cfg.JWTRefreshSecret, cfg.AccessTTL, cfg.RefreshTTL)

	// 初始化审计日志通道，Kafka 不可用时会降级为本地日志输出
	auditLogger := audit.NewKafkaAuditLogger(cfg.KafkaBrokers, cfg.KafkaTopic, logger)
	defer func() {
		if cerr := auditLogger.Close(); cerr != nil {
			logger.Warn("close audit logger failed", slog.String("error", cerr.Error()))
		}
	}()

	// 初始化 casbin
	enforcer := casbin.InitCasbin(database, cfg.CasbinPath, redisClient)

	// 初始化两级 session 存储
	session1, _ := util.NewBigCacheSessionStore()
	session2 := util.NewRedisSessionStore(redisClient, "redis:session:")

	sessionStore := util.NewCompositeSessionStore(session1, *session2)
	// 组装服务层与 HTTP 层
	businessLimiter := service.NewKeyedRateLimiter(2, 6)
	authSvc := service.NewAuthServiceImpl(
		database,
		userRepo,
		emailMapCache,
		revokedMapCache,
		emailBloom,
		tokenStoreCache,
		sessionStore,
		auditLogger,
		businessLimiter,
		cfg.ClientSecret,
		cfg.EmailMapTTL,
		cfg.RevokeMapTTL,
		cfg.ResetTokenTTL,
		cfg.JWTIssuer,
		logger,
	)
	userSvc := service.NewUserServiceImpl(database, userRepo, emailMapCache, enforcer)

	// Handler 同时接收 client 凭据与 redirect 白名单，用于授权码流程安全校验
	handler := controller.NewHandler(authSvc, userSvc, cfg.HydraAdminURL, cfg.OAuthClients, cfg.OAuthRedirectURIs, cfg.SessionCookieMaxAge, cfg.HydraLoginRememberFor, cfg.HydraConsentRememberFor)
	router := controller.NewRouter(handler, jwks, sessionStore, redisClient, cfg.HttpRateLimitRPS, cfg.HttpRateLimitBurstlogger, logger)

	server := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	// 启动 HTTP 服务，主协程继续监听退出信号
	go func() {
		logger.Info("http server started", slog.String("addr", cfg.HTTPAddr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// 收到 SIGINT/SIGTERM 后执行优雅停机
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	logger.Info("shutdown signal received", slog.String("signal", sig.String()))

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger.Info("server stopped gracefully")
}
