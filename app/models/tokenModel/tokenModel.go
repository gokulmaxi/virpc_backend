package tokenModel

import "go.mongodb.org/mongo-driver/bson/primitive"

type TokenModel struct {
	Title       string
	Description string
	CreateDate  string
	Status      string
	category    string
	UserId      primitive.ObjectID
}
