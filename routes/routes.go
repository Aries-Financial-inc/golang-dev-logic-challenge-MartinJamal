package routes

import (
	"JamalMartin/golang-dev-logic-challenge-MartinJamal/model"
	"fmt"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AnalysisResult structure for the response body
type AnalysisResult struct {
	GraphData       []GraphPoint `json:"graph_data"`
	MaxProfit       float64      `json:"max_profit"`
	MaxLoss         float64      `json:"max_loss"`
	BreakEvenPoints []float64    `json:"break_even_points"`
}

// GraphPoint structure for X & Y values of the risk & reward graph
type GraphPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/analyze", func(c *gin.Context) {
		var contracts []model.OptionsContract

		if err := c.ShouldBindJSON(&contracts); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(contracts) > 4 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Too many options contracts; maximum is 4"})
			return
		}

		graphData := calculateGraphData(contracts)
		fmt.Println(graphData)
		maxProfit := calculateMaxProfit(contracts)
		fmt.Println("Max Profit - ", maxProfit)
		maxLoss := calculateMaxLoss(contracts)
		fmt.Println("Max Loss - ", maxLoss)
		breakEvenPoints := calculateBreakEvenPoints(contracts)
		fmt.Println("breakEvenPoints - ", breakEvenPoints)

		response := AnalysisResult{
			GraphData:       graphData,
			MaxProfit:       maxProfit,
			MaxLoss:         maxLoss,
			BreakEvenPoints: breakEvenPoints,
		}

		c.JSON(http.StatusOK, response)
	})

	return router
}

func calculateGraphData(contracts []model.OptionsContract) []GraphPoint {
	var minStrike, maxStrike float64
	for i, contract := range contracts {
		if i == 0 || contract.StrikePrice < minStrike {
			minStrike = contract.StrikePrice
		}
		if i == 0 || contract.StrikePrice > maxStrike {
			maxStrike = contract.StrikePrice
		}
	}

	var graphData []GraphPoint
	for price := minStrike * 0.5; price <= maxStrike*1.5; price += (maxStrike * 2) / 100 {
		profit := 0.0
		for _, contract := range contracts {
			profit += calculateProfit(contract, price)
		}
		graphData = append(graphData, GraphPoint{X: price, Y: profit})
	}

	return graphData
}

// calculateProfit calculates the profit/loss for an options contract at a given underlying price.
func calculateProfit(contract model.OptionsContract, underlyingPrice float64) float64 {
	switch contract.Type {
	case "Call":
		return calculateCallProfit(contract, underlyingPrice)
	case "Put":
		return calculatePutProfit(contract, underlyingPrice)
	default:
		return 0.0 // Handle unsupported contract types gracefully
	}
}

// calculateCallProfit calculates profit/loss for a call option.
func calculateCallProfit(contract model.OptionsContract, underlyingPrice float64) float64 {
	if underlyingPrice <= contract.StrikePrice {
		// If underlying price is less than or equal to strike price, option expires worthless
		return -contract.Ask
	}
	// Profit calculation for a long call option
	return (underlyingPrice - contract.StrikePrice - contract.Ask)
}

// calculatePutProfit calculates profit/loss for a put option.
func calculatePutProfit(contract model.OptionsContract, underlyingPrice float64) float64 {
	if underlyingPrice >= contract.StrikePrice {
		// If underlying price is greater than or equal to strike price, option expires worthless
		return -contract.Ask
	}
	// Profit calculation for a long put option
	return (contract.StrikePrice - underlyingPrice - contract.Ask)
}

func calculateMaxProfit(contracts []model.OptionsContract) float64 {
	if len(contracts) == 0 {
		return 0.0
	}

	maxProfit := math.Inf(-1) // Initialize to negative infinity

	// Iterate over a range of potential prices based on strike prices
	minStrikePrice := minStrikePrice(contracts)
	maxStrikePrice := maxStrikePrice(contracts)

	for price := minStrikePrice; price <= maxStrikePrice; price += 0.5 { // Adjust step as needed
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

func calculateMaxLoss(contracts []model.OptionsContract) float64 {
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

func calculateBreakEvenPoints(contracts []model.OptionsContract) []float64 {
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

func maxStrikePrice(contracts []model.OptionsContract) float64 {
	if len(contracts) == 0 {
		return math.Inf(-1) // Return negative infinity if there are no contracts
	}

	maxPrice := contracts[0].StrikePrice
	for _, contract := range contracts {
		if contract.StrikePrice > maxPrice {
			maxPrice = contract.StrikePrice
		}
	}
	return maxPrice
}

func minStrikePrice(contracts []model.OptionsContract) float64 {
	if len(contracts) == 0 {
		return math.Inf(1) // Return positive infinity if there are no contracts
	}

	minPrice := contracts[0].StrikePrice
	for _, contract := range contracts {
		if contract.StrikePrice < minPrice {
			minPrice = contract.StrikePrice
		}
	}
	return minPrice
}
