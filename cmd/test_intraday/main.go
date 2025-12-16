package main

import "github.com/SwanHtetAungPhyo/alphavintage"

func main() {
	// Skip intraday (premium) and test daily data directly
	alphavintage.TestSingleDayFromDaily()
}
