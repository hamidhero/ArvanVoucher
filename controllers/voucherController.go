package controllers

import (
	"ArvanVoucher/requests"
	"ArvanVoucher/resources"
	"ArvanVoucher/services"
	"ArvanVoucher/services/queue"
	"ArvanVoucher/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

func AddVoucher(c *gin.Context) {
	var req requests.AddVoucherRequest
	output := utils.NewOutput()

	if e := c.ShouldBindBodyWith(&req, binding.JSON); e != nil {
		utils.SetError(e, c, &output, http.StatusBadRequest, http.StatusBadRequest)
		return
	}

	s := services.VoucherService{ServiceDB: services.DB, Redis: services.REDIS, Ctx: services.CTX}
	voucherId, e := s.AddVoucher(req)
	if e != nil {
		utils.SetError(e, c, &output, http.StatusBadRequest, http.StatusBadRequest)
		return
	}

	output.Data = resources.AddVoucherResource{VoucherId: *voucherId}
	c.JSON(http.StatusOK, output)
	return
}

func RegisterVoucher(c *gin.Context) {
	var req requests.RegisterVoucherRequest
	output := utils.NewOutput()

	if e := c.ShouldBindBodyWith(&req, binding.JSON); e != nil {
		utils.SetError(e, c, &output, http.StatusBadRequest, http.StatusBadRequest)
		return
	}

	mobile, e := req.Mobile.Int64()
	if e != nil {
		utils.SetError(e, c, &output, http.StatusBadRequest, http.StatusBadRequest)
		return
	}

	s := queue.QueueService{
		Redis:     services.REDIS,
		ServiceDB: services.DB,
		Ctx:       services.CTX,
	}
	sendData := requests.AddToQueueRequest{
		Name: req.VoucherCode,
		Data: mobile,
	}
	//call sender service to check validation and sent to queue to update user wallet
	if e := s.AddToQueue(sendData); e != nil {
		utils.SetError(e, c, &output, http.StatusBadRequest, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, output)
	return
}

func GetVoucherReport(c *gin.Context) {
	output := utils.NewOutput()

	voucherCode := c.Param("code")
	if voucherCode == "" {
		utils.SetError(errors.New("NotValidVoucher"), c, &output, http.StatusBadRequest, http.StatusBadRequest)
		return
	}

	s := services.VoucherService{ServiceDB: services.DB}
	//call voucher report service to get the list
	res, e := s.GetVoucherReport(voucherCode)
	if e != nil {
		utils.SetError(e, c, &output, http.StatusBadRequest, http.StatusBadRequest)
		return
	}

	output.Data = res
	c.JSON(http.StatusOK, output)
	return

}
