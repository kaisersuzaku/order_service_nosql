package repo

import (
	"context"

	"github.com/kaisersuzaku/order_service_nosql/models"
	"github.com/kaisersuzaku/order_service_nosql/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	productCollName string = "product"
	dbName          string = "online_store"
)

type ProductRepo struct {
	db *utils.DataStore
}

func BuildProductRepo(db *utils.DataStore) ProductRepo {
	return ProductRepo{
		db: db,
	}
}

type IProductRepo interface {
	StartSession(ctx context.Context) (mongo.Session, error)
	Rollback(ctx mongo.SessionContext, tx mongo.Session) error
	Commit(ctx context.Context, tx mongo.Session) error
	Update(ctx context.Context, product models.Product, reqQty int) (res *mongo.UpdateResult, err error)
}

func (p ProductRepo) StartSession(ctx context.Context) (mongo.Session, error) {
	return p.db.C.StartSession()
}

func (p ProductRepo) Rollback(ctx mongo.SessionContext, tx mongo.Session) error {
	return tx.AbortTransaction(ctx)
}

func (p ProductRepo) Commit(ctx context.Context, tx mongo.Session) error {
	return tx.CommitTransaction(ctx)
}

func (p ProductRepo) Update(ctx context.Context, product models.Product, reqQty int) (res *mongo.UpdateResult, err error) {
	id, _ := primitive.ObjectIDFromHex(product.ID)
	filter := bson.M{
		"_id":   id,
		"stock": bson.M{"$gte": reqQty},
	}
	data := bson.D{
		{"$inc", bson.D{{"stock", -reqQty}}},
	}
	coll := p.db.C.Database(dbName).Collection(productCollName)
	return coll.UpdateOne(ctx, filter, data)
}
