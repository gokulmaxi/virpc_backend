package batchModel

import "go.mongodb.org/mongo-driver/bson/primitive"

type BatchModel struct {
	BatchName        string
	BatchDescription string
	Startdate        string
	Enddate          string
	Totaldays        string
	ImageId          primitive.ObjectID
	CpuLimit         int
	AddFeatures      []string
}
