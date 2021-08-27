package models

// CraftConfiguration ...
type CraftConfiguration struct {
	ID         string `json:"id", bson:"id"`
	Attributes []AttributesBase
	Timestamp  int64 `json:"timestamp", bson:"timestamp"`
}

// AttributesBase ...
type AttributesBase struct {
	Key   string `json:"key", bson:"key"`
	Value string `json:"value", bson:"value"`
}
