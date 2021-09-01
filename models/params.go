package models

type Params struct {
	Timestamp int64  `json:"timestamp", bson:"timestamp"`
	Key       string `json:"key", bson:"key"`
	Value     string `json:"value", bson:"value"`
	Limit     int    `json:"limit", bson:"limit"`
}

//IDS ...
type IDS struct {
	IDs []string `json:"ids,omitempty" bson:"ids,omitempty"`
}

type Count struct {
	Count int `json:"count", bson:"count"`
}
