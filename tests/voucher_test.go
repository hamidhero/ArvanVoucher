package tests

import (
	"ArvanVoucher/requests"
	"ArvanVoucher/services"
	"ArvanVoucher/services/queue"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestAddVoucher(t *testing.T) {
	services.ConnectDB()
	s := services.VoucherService{ServiceDB: services.DB}

	tt := time.Now().Add(time.Hour * 24)
	tm := requests.DateTime{tt}
	req := requests.AddVoucherRequest{
		VoucherCode:    "t1",
		Amount:         1000,
		Limit:          0,
		ExpirationDate: &tm,
	}
	_, e := s.AddVoucher(req)
	assert.NotEqual(t, nil, e)
}

func TestGetVoucherReport(t *testing.T) {
	services.ConnectDB()
	services.ConnectRedis()
	s := services.VoucherService{ServiceDB: services.DB}

	res, e := s.GetVoucherReport("t1")
	assert.Equal(t, nil, e)
	assert.NotEqual(t, 0, len(res))
}

func TestAddToVoucherUsersQueue(t *testing.T) {
	services.ConnectDB()
	services.ConnectRedis()
	services.ConnectRabbitMq()
	s := queue.QueueService{
		Rabbit:    services.RABBIT,
		ServiceDB: services.DB,
	}

	req := requests.AddToQueueRequest{
		Name: "test",
		Data: nil,
	}
	e := s.AddToVoucherUsersQueue(req)
	assert.Equal(t, nil, e)
}

func handler(c *gin.Context) {
	var info requests.RegisterVoucherRequest
	if err := c.ShouldBindJSON(&info); err != nil {
		log.Panic(err)
	}
	fmt.Println(info)
	c.Writer.Write([]byte(`{"status": 200}`))
}

func TestRegisterVoucher(t *testing.T) {
	services.ConnectDB()
	services.ConnectRedis()
	services.ConnectRabbitMq()

	payload := strings.NewReader(`{
    "mobile": 9376019080,
    "voucher_code": "t11"}`)

	path := "/api/voucher/register"
	router := gin.Default()
	router.POST(path, handler)
	req, _ := http.NewRequest("POST", path, payload)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	t.Logf("status: %d", w.Code)
	t.Logf("response: %s", w.Body.String())

	assert.Equal(t, http.StatusOK, w.Code)
}
