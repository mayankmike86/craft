package models

type Params struct {
	Timestamp int64  `json:"timestamp", bson:"timestamp"`
	Key       string `json:"key", bson:"key"`
	Value     string `json:"value", bson:"value"`
}

//IDS ...
type IDS struct {
	IDs []string `json:"ids,omitempty" bson:"ids,omitempty"`
}
