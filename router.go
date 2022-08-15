package main

import (
	"ArvanVoucher/controllers"
	"github.com/gin-gonic/gin"
)

func GetRouter(port string) (router *gin.Engine) {
	router = gin.Default()
	router.ForwardedByClientIP = true
	router.RedirectFixedPath = true

	api := router.Group("api")
	{
		voucher := api.Group("voucher")
		{
			voucher.POST("", controllers.AddVoucher)
			voucher.POST("register", controllers.RegisterVoucher)
			voucher.GET("report/:code", controllers.GetVoucherReport)
		}
	}

	router.Run("0.0.0.0:" + port)
	return
}