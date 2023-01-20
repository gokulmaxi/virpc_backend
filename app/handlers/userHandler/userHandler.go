package userHandler

import (
	"app/database"
	"app/models/userModel"
	"app/utilities"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func signup(c *fiber.Ctx) error {
	var user = userModel.UserModel{}
	req := c.Body()
	json.Unmarshal(req, &user)
	coll := database.Instance.Db.Collection("users")
	_, err := coll.InsertOne(context.TODO(), user)
	if err != nil {
		panic(err)
	}
	return c.SendString("signed in")
}
func signin(c *fiber.Ctx) error {
	// create new session
	var results userModel.UserModel
	var req map[string]interface{}
	res := make(map[string]interface{})
	sess, err := utilities.Store.Get(c)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(c.Body(), &req)
	if err != nil {
		panic(err)
	}
	filter := bson.D{{"email", req["user"].(string)}}
	coll := database.Instance.Db.Collection("users")
	err = coll.FindOne(context.TODO(), filter).Decode(&results)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			res["data"] = "User Not found"
			data, _ := json.Marshal(res)
			return c.Send(data)
		}
	}
	if results.User_role == userModel.UserTypeAdmin {
		if req["password"] != results.Data.(userModel.AdminUser).Password {
			res["data"] = "password is incorrect"
		} else {
			res["name"] = results.Name
			res["role"] = results.User_role
			res["is_deactivated"] = results.Account_deactivated
		}
	} else {
		fmt.Println(results.Data)
		if results.Data.(userModel.NormalUser).ImageUrl == "" {
			res["registered"] = false
		} else {
			res["registered"] = true
		}
		res["name"] = results.Name
		res["role"] = results.User_role
		res["is_deactivated"] = results.Account_deactivated

	}
	// res["account_deactivated"] = results.Account_deactivated
	fmt.Printf("%+v\n", results)
	sess.Set("name", results.Name)
	sess.Set("role", results.User_role)
	if err := sess.Save(); err != nil {
		panic(err)
	}
	data, err := json.Marshal(res)
	return c.Send(data)
}
func signout(c *fiber.Ctx) error {
	sess, err := utilities.Store.Get(c)
	if err != nil {
		panic(err)
	}
	sess.Delete("name")
	sess.Delete("role")
	if err := sess.Save(); err != nil {
		panic(err)
	}
	return c.SendString("logged out")
}
func Register(_route fiber.Router) {
	_route.Post("/signup", signup)
	_route.Post("/login", signin)
	_route.Get("/logout", signout)
}
