package api

import (
	"fmt"
	"github.com/ethereum-blockchain-service/service/api/controller"
	"github.com/ethereum-blockchain-service/tools/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Start() {
	db.Init()
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", 8080),
		Handler:        setupRouter(),
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}

}

func setupRouter() *gin.Engine {

	gin.SetMode(gin.DebugMode)
	r := gin.New()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, "token")

	r.Use(gin.Recovery(), cors.New(config))

	controller := controller.NewController()
	route := r.Group("/api/v1")
	route.GET("/blocks", controller.Block.GetBlocksByLimit)
	route.GET("/blocks/:id", controller.Block.GetBlocksById)
	route.GET("/transaction/:txHash", controller.Block.GetTransactionByTxHash)
	return r
}
