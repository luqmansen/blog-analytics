package analytics

import "cloud.google.com/go/civil"

type Analytic struct {
	CreatedAt civil.DateTime `json:"created_at" bson:"created_at" msgpack:"code"`
	URL       string    `json:"url" bson:"url" msgpack:"url" validate:"empty=false & format=url"`
	IP        string    `json:"ip" bson:"ip" msgpack:"ip"`
	Info      IpInfo    `json:"info" bson:"info" msgpack:"info"`
}

type IpInfo struct {
	IP       string `json:"ip" bson:"ip" msgpack:"ip" validate:"format=ip"`
	City     string `json:"city" bson:"city" msgpack:"city"`
	Region   string `json:"region" bson:"region" msgpack:"region"`
	Country  string `json:"country" bson:"country" msgpack:"country"`
	Location string `json:"location" bson:"location" msgpack:"location"`
	Org      string `json:"org" bson:"org" msgpack:"org"`
	Timezone string `json:"timezone" bson:"timezone" msgpack:"timezone"`
}
