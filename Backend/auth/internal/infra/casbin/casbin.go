package casbin

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func InitCasbin(db *gorm.DB, modelPath string, rdb redis.UniversalClient) *casbin.SyncedEnforcer {
	// GORM 适配器
	a, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		log.Fatalf("failed to init casbin adapter: %v", err)
	}
	fmt.Println("Current working directory:", modelPath)
	// 创建 Enforcer
	e, err := casbin.NewSyncedEnforcer(modelPath, a)
	if err != nil {
		log.Fatalf("failed to create enforcer: %v", err)
	}

	// 初始化 Watcher
	w, err := rediswatcher.NewWatcher("", rediswatcher.WatcherOptions{
		PubClient: rdb,
		SubClient: rdb,
		Channel:   "casbin_sync",
	})

	if err != nil {
		log.Printf("Casbin: 无法启用 Redis Watcher (%v)，降级为本地模式", err)
	} else {
		e.SetWatcher(w)
		w.SetUpdateCallback(func(string) { e.LoadPolicy() })
		log.Println("Casbin: Redis Watcher 成功复用现有连接")
	}

	// 4. 加载初始策略
	if err := e.LoadPolicy(); err != nil {
		log.Fatalf("failed to load policy: %v", err)
	}

	return e
}
