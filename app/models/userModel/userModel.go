package userModel

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	UserTypeAdmin   = "admin"
	UserTypeStudent = "student"
	UserTypeFaculty = "faculty"
)

type UserModel struct {
	Name      string `bson:"name"`
	Email     string `bson:"email"`
	User_role string `bson:"user_type"`
	Data      interface{}
}

type AdminModel struct {
	Password string `bson:"password"`
}
type FacultyModel struct {
	Faculty_id        string `bson:"faculty_id"`
	FacultyDepartment string `bson:"Department,omitempty"`
}
type StudentModel struct {
	Roll_number string `bson:"roll_number"`
	Year        string `bson:"year"`
	Department  string `bson:"department,omitempty"`
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
	switch discriminator {
	case UserTypeAdmin:
		user.User_role = "admin"
		user.Data = &AdminModel{
			Password: dev["password"].(string),
		}
	case UserTypeStudent:
		fmt.Println("student")
		user.User_role = "student"
		user.Data = StudentModel{
			Roll_number: dev["roll_number"].(string),
			Year:        dev["year"].(string),
			Department:  dev["department"].(string),
		}
	case UserTypeFaculty:
		user.User_role = "student"
		user.Data = FacultyModel{
			Faculty_id:        dev["faculty_id"].(string),
			FacultyDepartment: dev["department"].(string),
		}
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
	subdoc = doc["data"].(primitive.M)
	switch doc["user_type"] {
	case UserTypeAdmin:
		fmt.Println("admin")
		user.User_role = "admin"
		user.Data = AdminModel{
			Password: subdoc["password"].(string),
		}
	case UserTypeStudent:
		fmt.Println("student")
		user.User_role = "student"
		user.Data = StudentModel{
			Roll_number: subdoc["roll_number"].(string),
			Year:        subdoc["year"].(string),
			Department:  subdoc["department"].(string),
		}
	case UserTypeFaculty:
		fmt.Println("faculty")
		user.User_role = "student"
		user.Data = FacultyModel{
			Faculty_id:        subdoc["faculty_id"].(string),
			FacultyDepartment: subdoc["department"].(string),
		}

	}
	return nil
}
