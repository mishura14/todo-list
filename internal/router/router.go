package router

import (
	"git-register-project/internal/Database/redis"
	handler_comfirm_register "git-register-project/internal/handler/hand_register/confirm_register"
	handler_register "git-register-project/internal/handler/hand_register/register"
	repository "git-register-project/internal/repository/interface"
	serversmtp "git-register-project/internal/server_smtp"
	comfirm_register "git-register-project/internal/servise/register/confirm_register"
	"git-register-project/internal/servise/register/register"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, redisClient *redis.Redis, repo repository.UserRegister) {
	// Создаём SMTP-сервис
	mailSender := serversmtp.NewSMTPSender()

	// Создаём сервис регистрации (с бизнес-логикой)
	registerService := register.NewRegisterService(repo, mailSender, redisClient)

	// Создаём хендлер с сервисом
	registerHandler := handler_register.NewRegister(registerService)
	confirmService := comfirm_register.NewConfirmRegisterService(repo, mailSender, redisClient)
	confirmRegisterHandler := handler_comfirm_register.NewConfirmRegister(confirmService)

	public := r.Group("/api")
	{
		public.POST("/register", registerHandler.Register)
		public.POST("/confirm_register", confirmRegisterHandler.Confirm_register)
	}
}
