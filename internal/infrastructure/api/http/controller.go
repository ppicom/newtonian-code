package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	usecases "github.com/ppicom/newtonian/internal/application/use_cases"
)

type Controller struct {
	transferMoneyUseCase *usecases.TransferMoneyUseCase
}

func NewController(transferMoneyUseCase *usecases.TransferMoneyUseCase) *Controller {
	return &Controller{
		transferMoneyUseCase: transferMoneyUseCase,
	}
}

func (c *Controller) SetupRoutes(router *Router) {
	api := router.Engine().Group("/api/v1")
	{
		api.POST("/transfer", c.TransferMoney)
	}
}

func (c *Controller) TransferMoney(ctx *gin.Context) {
	from := ctx.PostForm("from")
	to := ctx.PostForm("to")
	amount, err := strconv.Atoi(ctx.PostForm("amount"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	err = c.transferMoneyUseCase.Execute(from, to, amount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Transfer successful"})
}
