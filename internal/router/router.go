package router

import (
	"git-register-project/internal/Database/redis"
	handler_register "git-register-project/internal/handler/hand_register"
	"git-register-project/internal/repository"
	serversmtp "git-register-project/internal/server_smtp"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, redisClient *redis.Redis, repo repository.UserRepository) {

	mailSender := serversmtp.NewSMTPSender()

	registerHandler := handler_register.NewRegister(redisClient, repo, mailSender)

	public := r.Group("/api")
	{
		public.POST("/register", registerHandler.Register)
		public.POST("/confirm_register", registerHandler.Confirm_register)
	}
}
