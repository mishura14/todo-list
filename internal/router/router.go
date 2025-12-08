package router

import (
	"git-register-project/internal/Database/redis"
	handler_register "git-register-project/internal/handler/hand_register"
	"git-register-project/internal/repository"
	serversmtp "git-register-project/internal/server_smtp"
	"git-register-project/internal/servise/register"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, redisClient *redis.Redis, repo repository.UserRegister) {
	// Создаём SMTP-сервис
	mailSender := serversmtp.NewSMTPSender()

	// Создаём сервис регистрации (с бизнес-логикой)
	registerService := register.NewRegisterService(repo, mailSender, redisClient)

	// Создаём хендлер с сервисом
	registerHandler := handler_register.NewRegister(registerService)

	public := r.Group("/api")
	{
		public.POST("/register", registerHandler.Register)
		public.POST("/confirm_register", registerHandler.Confirm_register)
	}
}
