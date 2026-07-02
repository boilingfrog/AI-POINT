// Command server 启动金币钱包 HTTP 服务。main 只负责装配各层与启动，不含业务逻辑。
package main

import (
	"log"
	"net/http"
	"time"

	"demo/coinwallet/dao"
	"demo/coinwallet/service"
	httpapi "demo/coinwallet/transport/http"
)

func main() {
	store := dao.NewStore()
	svc := service.New(store)
	handler := httpapi.NewHandler(svc)

	cfg := loadConfig()
	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second, // 防慢速头攻击，也满足审查的超时要求
	}

	log.Printf("coinwallet listening on %s", cfg.Addr)
	log.Fatal(srv.ListenAndServe())
}
