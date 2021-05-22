package services

import (
	"context"
	"log"

	"github.com/kaisersuzaku/order_service_nosql/models"
	"github.com/kaisersuzaku/order_service_nosql/repo"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderProductService struct {
	pr repo.IProductRepo
}

func BuildOrderProductService(pr repo.IProductRepo) OrderProductService {
	return OrderProductService{
		pr: pr,
	}
}

type IOrderProductService interface {
	OrderProduct(ctx context.Context, req models.OrderProductReq) (resp models.OrderProductResp, err models.RespError)
}

func (ops OrderProductService) OrderProduct(ctx context.Context, req models.OrderProductReq) (resp models.OrderProductResp, err models.RespError) {
	var product models.Product
	tx, e := ops.pr.StartSession(context.Background())
	if e != nil {
		log.Println(e)
		err = models.GetRequestFailed()
		return
	}
	defer tx.EndSession(context.Background())

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		product.ID = req.ID
		res, eCallback := ops.pr.Update(sessCtx, product, req.Qty)
		if eCallback != nil {
			ops.pr.Rollback(sessCtx, tx)
			err = models.GetRequestFailed()
			log.Println(eCallback)
			return nil, eCallback
		}

		if res.ModifiedCount == 0 {
			ops.pr.Rollback(sessCtx, tx)
			err = models.GetStockLessThanRequest()
			log.Println(err.ErrorMessage)
			return nil, eCallback
		}

		if res.ModifiedCount == 1 {
			// logics hits to vendor payment
			// will rollback substract if payment failed

			// logics saving to coll payment
			// will rollback substract if payment failed
			// and send refund
			log.Println("Processed Payment logic")
		}

		eCallback = ops.pr.Commit(sessCtx, tx)
		if eCallback != nil {
			err = models.GetRequestFailed()
			log.Println(eCallback)
			return nil, eCallback
		}
		err = models.RespError{}

		return nil, nil
	}
	_, e = tx.WithTransaction(ctx, callback)
	if e != nil {
		log.Println("e", e)
		err = models.GetRequestFailed()
		log.Println(e)
		return
	}
	resp.Status = "OK"
	return
}
