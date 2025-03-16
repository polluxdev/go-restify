package webapi

type WebApiResponseSuccess struct {
	Code   int    `bson:"code" json:"code"`
	Status string `bson:"status" json:"status"`
}

type WebApiResponseFailed struct {
	Code   int         `bson:"code" json:"code"`
	Status string      `bson:"status" json:"status"`
	Data   interface{} `bson:"data" json:"data"`
}
