package models

import (
	"time"
)

type DBConfig struct {
	URI string
}

type FindRecords struct {
	DB         string
	Collection string
	DateStart  time.Time
	DateEnd    time.Time
	Fields     []string
}

type WeatherRecords struct {
	DateTime time.Time `csv:"Fecha" bson:"datetime" json:"datetime"`
	Temp     float64   `csv:"Temp" json:"temp"`
	Chill    float64   `csv:"Chill" bson:"chill" json:"chill"`
	Dew      float64   `csv:"Dew" bson:"dew" json:"dew"`
	Heat     float64   `csv:"Heat" bson:"heat" json:"heat"`
	Hum      float64   `csv:"Hum" bson:"hum" json:"hum"`
	WspdAvg  float64   `csv:"Wspdavg" bson:"wspd_avg" json:"wspd_avg"`
	WdirAvg  float64   `csv:"Wdiravg" bson:"wdir_avg" json:"wdir_avg"`
	Bar      float64   `csv:"Bar" bson:"bar" json:"bar"`
	Rain     float64   `csv:"Rain" bson:"rain" json:"rain"`
}

type EFMElectricField struct {
	DateTime      time.Time `bson:"datetime" json:"datetime"`
	Lightning     bool      `bson:"lightning" json:"lightning"`
	ElectricField float64   `bson:"electric_field" json:"electric_field"`
	Distance      uint8     `bson:"distance" json:"distance"`
	RotorFail     bool      `bson:"rotor_fail" json:"rotor_fail"`
}

type LineChart struct {
	Data []struct {
		Name string `bson:"name" json:"name"`
		Data []struct {
			X time.Time `bson:"x"`
			Y any       `bson:"y"`
		} `bson:"data" json:"data"`
	} `bson:"data" json:"data"`
}

type RadialRanges struct {
	Range1 int64 `bson:"range1" json:"range1"`
	Range2 int64 `bson:"range2" json:"range2"`
	Range3 int64 `bson:"range3" json:"range3"`
	Range4 int64 `bson:"range4" json:"range4"`
}

type RadialChart struct {
	Labels []string `json:"labels"`
	Series []int64  `json:"ranges"`
}
