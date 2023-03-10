package userModel

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	UserTypeAdmin  = "admin"
	UserTypeNormal = "user"
)

type UserModel struct {
	User_id             primitive.ObjectID `bson:"_id,omitempty"`
	Name                string             `bson:"name"`
	Email               string             `bson:"email"`
	User_role           string             `bson:"user_type"`
	Account_deactivated bool
	Data                interface{}
}

type AdminUser struct {
	Password string `bson:"password"`
}
type NormalUser struct {
	ImageUrl string `bson:"imageurl"`
	PhoneNo  string `bson:"phoneno"`
}

func (user *UserModel) UnmarshalJSON(data []byte) (err error) {
	var dev map[string]interface{}
	var _err = json.Unmarshal(data, &dev)
	if _err != nil {
		panic(err)
	}
	discriminator, ok := dev["user_type"].(string)
	if !ok {
		return errors.New("invalid discriminator type")
	}
	user.Name = dev["name"].(string)
	user.Email = dev["email"].(string)
	// user.User_id = dev["_id"].(primitive.ObjectID)
	fmt.Println(user.User_id)
	switch discriminator {
	case UserTypeAdmin:
		user.User_role = "admin"
		user.Data = &AdminUser{
			Password: dev["password"].(string),
		}
	case UserTypeNormal:
		fmt.Println("user")
		user.User_role = "user"
		user.Data = NormalUser{}
	}
	return nil
}

func (user *UserModel) UnmarshalBSON(data []byte) (err error) {
	var doc bson.M
	var subdoc bson.M
	if err := bson.Unmarshal(data, &doc); err != nil {
		return err
	}
	user.Name = doc["name"].(string)
	user.Email = doc["email"].(string)
	user.User_id = doc["_id"].(primitive.ObjectID)
	subdoc = doc["data"].(primitive.M)
	switch doc["user_type"] {
	case UserTypeAdmin:
		fmt.Println("admin")
		user.User_role = "admin"
		user.Data = AdminUser{
			Password: subdoc["password"].(string),
		}
	case UserTypeNormal:
		fmt.Println("user")
		user.User_role = "user"
		user.Data = NormalUser{
			ImageUrl: subdoc["imageurl"].(string),
			PhoneNo:  subdoc["phoneno"].(string),
		}
	}
	return nil
}
