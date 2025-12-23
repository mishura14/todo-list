package router

import (
	"git-register-project/internal/Database/redis"
	"git-register-project/internal/handler"
	handlogin "git-register-project/internal/handler/hand_login/login"
	handrefresh "git-register-project/internal/handler/hand_refresh"
	handler_comfirm_register "git-register-project/internal/handler/hand_register/confirm_register"
	handler_register "git-register-project/internal/handler/hand_register/register"
	"git-register-project/internal/middleware/auth"
	repository "git-register-project/internal/repository/interface"
	serversmtp "git-register-project/internal/server_smtp"
	"git-register-project/internal/servise/login"
	refreshtoken "git-register-project/internal/servise/refreshToken"
	comfirm_register "git-register-project/internal/servise/register/confirm_register"
	"git-register-project/internal/servise/register/register"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	r *gin.Engine,
	redisClient *redis.Redis,
	userRepo repository.UserRegister,
	loginRepo repository.UserLogin,
	refreshRepo repository.Refreshtoken,
	tokenGen repository.TokenGenerator,
) {
	mailSender := serversmtp.NewSMTPSender()

	registerService := register.NewRegisterService(userRepo, mailSender, redisClient)
	registerHandler := handler_register.NewRegister(registerService)

	confirmService := comfirm_register.NewConfirmRegisterService(userRepo, mailSender, redisClient)
	confirmHandler := handler_comfirm_register.NewConfirmRegister(confirmService)

	authService := login.NewLoginService(loginRepo, mailSender, redisClient, tokenGen)
	loginHandler := handlogin.NewLogin(authService)

	refreshService := refreshtoken.NewRefreshService(tokenGen, redisClient, refreshRepo)
	refreshHandler := handrefresh.NewHandlerRefresh(refreshService)

	public := r.Group("/api")
	{
		public.POST("/register", registerHandler.Register)
		public.POST("/confirm_register", confirmHandler.Confirm_register)
		public.POST("/login", loginHandler.Login)
	}

	protected := r.Group("/api/auth")
	protected.Use(auth.AuthMiddleware())
	{
		protected.GET("/get", handler.GetHandler)
		protected.POST("/refresh", refreshHandler.Refresh)
	}
}
