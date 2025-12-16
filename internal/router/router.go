package router

import (
	"git-register-project/internal/Database/redis"
	handlogin "git-register-project/internal/handler/hand_login/login"
	handler_comfirm_register "git-register-project/internal/handler/hand_register/confirm_register"
	handler_register "git-register-project/internal/handler/hand_register/register"
	repository "git-register-project/internal/repository/interface"
	serversmtp "git-register-project/internal/server_smtp"
	"git-register-project/internal/servise/login"
	comfirm_register "git-register-project/internal/servise/register/confirm_register"
	"git-register-project/internal/servise/register/register"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine,
	redisClient *redis.Redis,
	repo repository.UserRegister,
	loginRepo repository.UserLogin,
	generate repository.TokenGenerator) {
	// Создаём SMTP-сервис
	mailSender := serversmtp.NewSMTPSender()

	// Создаём сервис регистрации (с бизнес-логикой)
	registerService := register.NewRegisterService(repo, mailSender, redisClient)

	// Создаём хендлер с сервисом
	registerHandler := handler_register.NewRegister(registerService)

	// Создаём сервис подтверждения регистрации (с бизнес-логикой)
	confirmService := comfirm_register.NewConfirmRegisterService(repo, mailSender, redisClient)

	// Создаём хендлер с сервисом
	confirmRegisterHandler := handler_comfirm_register.NewConfirmRegister(confirmService)

	// Создаём сервис авторизации (с бизнес-логикой)
	authService := login.NewLoginService(loginRepo, mailSender, redisClient, generate)

	// Создаём хендлер с сервисом
	loginHandler := handlogin.NewLogin(authService)

	public := r.Group("/api")
	{
		public.POST("/register", registerHandler.Register)
		public.POST("/confirm_register", confirmRegisterHandler.Confirm_register)
		public.POST("/login", loginHandler.Login)
	}
}
