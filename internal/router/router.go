package router

import (
	handler_register "git-register-project/internal/handler/hand_register"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	public := r.Group("/api")
	{
		public.POST("/register", handler_register.Register)
		public.POST("confirm_register", handler_register.Confirm_register)
	}
}
