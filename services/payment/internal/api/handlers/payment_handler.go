package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iwanlaudin/go-microservice/pkg/common/api"
	"github.com/iwanlaudin/go-microservice/pkg/common/helpers"
	"github.com/iwanlaudin/go-microservice/services/payment/internal/service"
)

type HandlerPayment struct {
	PaymentService *service.PaymentService
}

func NewPaymentHandler(paymentService *service.PaymentService) *HandlerPayment {
	return &HandlerPayment{
		PaymentService: paymentService,
	}
}

func (handler *HandlerPayment) FindAll(writer http.ResponseWriter, req *http.Request) {
	results, err := handler.PaymentService.FindAll(req.Context())
	if err != nil {
		api.NewAppResponse(err.Error(), http.StatusBadRequest).Err(writer)
		return
	}

	api.NewAppResponse("Succefully to find all payment", http.StatusOK).Ok(writer, results)
}

func (handler *HandlerPayment) FindById(writer http.ResponseWriter, req *http.Request) {
	paymentIdReq := chi.URLParam(req, "paymentId")
	paymentId, err := helpers.ConvertStringToUUID(paymentIdReq)
	helpers.PanicIfError(err)

	result, err := handler.PaymentService.FindById(req.Context(), paymentId)
	if err != nil {
		api.NewAppResponse(err.Error(), http.StatusBadRequest).Err(writer)
		return
	}

	api.NewAppResponse("Successfully to find payment", http.StatusOK).Ok(writer, result)
}

func (handler *HandlerPayment) Update(writer http.ResponseWriter, req *http.Request) {
	paymentIdReq := chi.URLParam(req, "paymentId")
	paymentId, err := helpers.ConvertStringToUUID(paymentIdReq)
	helpers.PanicIfError(err)

	err = handler.PaymentService.Update(req.Context(), paymentId)
	if err != nil {
		api.NewAppResponse(err.Error(), http.StatusBadRequest).Err(writer)
		return
	}

	api.NewAppResponse("Successfully update payment", http.StatusOK).Ok(writer, nil)
}
