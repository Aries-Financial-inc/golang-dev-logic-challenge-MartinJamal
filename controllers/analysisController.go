package controllers

import (
	"encoding/json"
	"math"
	"net/http"
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

// AnalysisResponse represents the data structure of the analysis result
type AnalysisResponse struct {
	XYValues        []XYValue `json:"xy_values"`
	MaxProfit       float64   `json:"max_profit"`
	MaxLoss         float64   `json:"max_loss"`
	BreakEvenPoints []float64 `json:"break_even_points"`
}

// XYValue represents a pair of X and Y values
type XYValue struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func AnalysisHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var contracts []OptionsContract
	err := json.NewDecoder(r.Body).Decode(&contracts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(contracts) > 4 {
		http.Error(w, "Too many options contracts; maximum is 4", http.StatusBadRequest)
		return
	}

	xyValues := calculateXYValues(contracts)
	maxProfit := calculateMaxProfit(contracts)
	maxLoss := calculateMaxLoss(contracts)
	breakEvenPoints := calculateBreakEvenPoints(contracts)

	response := AnalysisResponse{
		XYValues:        xyValues,
		MaxProfit:       maxProfit,
		MaxLoss:         maxLoss,
		BreakEvenPoints: breakEvenPoints,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func calculateXYValues(contracts []OptionsContract) []XYValue {
	// For simplicity, let's assume the underlying price range is between 0 and 2 * max strike price
	var minStrike, maxStrike float64
	for i, contract := range contracts {
		if i == 0 || contract.StrikePrice < minStrike {
			minStrike = contract.StrikePrice
		}
		if i == 0 || contract.StrikePrice > maxStrike {
			maxStrike = contract.StrikePrice
		}
	}

	var xyValues []XYValue
	for price := minStrike * 0.5; price <= maxStrike*1.5; price += (maxStrike * 2) / 100 {
		profit := 0.0
		for _, contract := range contracts {
			profit += calculateProfit(contract, price)
		}
		xyValues = append(xyValues, XYValue{X: price, Y: profit})
	}

	return xyValues
}

func calculateProfit(contract OptionsContract, price float64) float64 {
	profit := 0.0
	switch contract.Type {
	case "call":
		if contract.LongShort == "long" {
			profit = math.Max(0, price-contract.StrikePrice) - contract.Ask
		} else {
			profit = contract.Bid - math.Max(0, price-contract.StrikePrice)
		}
	case "put":
		if contract.LongShort == "long" {
			profit = math.Max(0, contract.StrikePrice-price) - contract.Ask
		} else {
			profit = contract.Bid - math.Max(0, contract.StrikePrice-price)
		}
	}
	return profit
}

func calculateMaxProfit(contracts []OptionsContract) float64 {
	var maxProfit float64
	for price := 0.0; price <= 2*maxStrikePrice(contracts); price += 1.0 {
		profit := 0.0
		for _, contract := range contracts {
			profit += calculateProfit(contract, price)
		}
		if profit > maxProfit {
			maxProfit = profit
		}
	}
	return maxProfit
}

func calculateMaxLoss(contracts []OptionsContract) float64 {
	var maxLoss float64
	for price := 0.0; price <= 2*maxStrikePrice(contracts); price += 1.0 {
		profit := 0.0
		for _, contract := range contracts {
			profit += calculateProfit(contract, price)
		}
		if profit < maxLoss {
			maxLoss = profit
		}
	}
	return maxLoss
}

func calculateBreakEvenPoints(contracts []OptionsContract) []float64 {
	var breakEvenPoints []float64
	var previousProfit float64
	for price := 0.0; price <= 2*maxStrikePrice(contracts); price += 1.0 {
		profit := 0.0
		for _, contract := range contracts {
			profit += calculateProfit(contract, price)
		}
		if price != 0.0 && ((previousProfit < 0 && profit >= 0) || (previousProfit > 0 && profit <= 0)) {
			breakEvenPoints = append(breakEvenPoints, price)
		}
		previousProfit = profit
	}
	return breakEvenPoints
}

func maxStrikePrice(contracts []OptionsContract) float64 {
	maxStrike := 0.0
	for _, contract := range contracts {
		if contract.StrikePrice > maxStrike {
			maxStrike = contract.StrikePrice
		}
	}
	return maxStrike
}
