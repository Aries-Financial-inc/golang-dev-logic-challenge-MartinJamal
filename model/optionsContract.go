package model

import (
	"time"
)

// OptionsContract represents the data structure of an options contract
type OptionsContract struct {
	Type           string    `json:"type"`            // call or put
	StrikePrice    float64   `json:"strike_price"`    // strike price
	Bid            float64   `json:"bid"`             // bid price
	Ask            float64   `json:"ask"`             // ask price
	ExpirationDate time.Time `json:"expiration_date"` // expiration date
	LongShort      string    `json:"long_short"`      // long or short
}
