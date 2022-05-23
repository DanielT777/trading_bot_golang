package main

type stacksDic struct {
}

type chartsDic struct {
}

type botStatistics struct {
	SMA         []int
	FMA         []int
	EMA         []int
	UPPER_BB    []int
	LOWER_BB    []int
	actionOrder string
}

type botState struct {
	timeBank       int
	maxTimeBank    int
	timePerMove    int
	candleInterval int
	candleFormat   []int
	candlesTotal   int
	candlesGiven   int
	initialStack   int
	transactionFee float64
	date           int
	buyingPrice    int
	stacks         stacksDic
	charts         chartsDic
	stats          botStatistics
}

type bot struct {
	botState botState
	runBot   runBot
}

type runBot func()

func init() {

}

func main() {

}
