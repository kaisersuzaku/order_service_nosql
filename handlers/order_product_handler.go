package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/kaisersuzaku/order_service_nosql/models"
	"github.com/kaisersuzaku/order_service_nosql/services"
)

type OrderProductHandler struct {
	ps services.IOrderProductService
}

func BuildOrderProductHandler(ps services.IOrderProductService) OrderProductHandler {
	return OrderProductHandler{
		ps: ps,
	}
}

type IOrderProductHandler interface {
	GetProductByID(w http.ResponseWriter, r *http.Request)
}

func (oph OrderProductHandler) OrderProduct(w http.ResponseWriter, r *http.Request) {
	var req models.OrderProductReq
	bodyByte, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respErr := models.GetInvalidPayloadResp()
		log.Println(respErr.ErrorCode, respErr.ErrorMessage, err)
		oph.respError(w, respErr)
		return
	}
	err = json.Unmarshal(bodyByte, &req)
	if err != nil {
		respErr := models.GetInvalidPayloadResp()
		log.Println(respErr.ErrorCode, respErr.ErrorMessage, err)
		oph.respError(w, respErr)
		return
	}
	_, err = govalidator.ValidateStruct(req)
	if err != nil {
		respErr := models.GetInvalidPayloadResp()
		log.Println(respErr.ErrorCode, respErr.ErrorMessage, err)
		oph.respError(w, respErr)
		return
	}
	resp, respErr := oph.ps.OrderProduct(r.Context(), req)
	if respErr.ErrorCode != "" {
		log.Println(respErr.ErrorCode, respErr.ErrorMessage)
		oph.respError(w, respErr)
		return
	}
	oph.respSuccess(w, resp)
}

func (oph OrderProductHandler) respSuccess(w http.ResponseWriter, data interface{}) {
	body, _ := json.Marshal(data)
	m := make(map[string]string)
	w.Header().Add("Content-type", "application/json")
	m["Content-type"] = "application/json"
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (oph OrderProductHandler) respError(w http.ResponseWriter, e models.RespError) {
	body, _ := json.Marshal(e)
	m := make(map[string]string)
	w.Header().Add("Content-type", "application/json")
	m["Content-type"] = "application/json"
	w.WriteHeader(http.StatusBadRequest)
	w.Write(body)
}
