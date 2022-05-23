package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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
	candleFormat   []string
	candlesTotal   int
	candlesGiven   int
	initialStack   int
	transactionFee float64
	date           int
	buyingPrice    int
	stacks         stacksDict
	charts         chartsDict
	stats          botStatistics
}

type stacksDict struct {
	btc  string
	usdt string
}

type chartsDict struct {
	usdt_btc string
}

func main() {

	stacksDict := stacksDict{
		btc:  "BTC",
		usdt: "USDT",
	}

	chartsDict := chartsDict{
		usdt_btc: "USDT_BTC",
	}

	botState := botState{
		timeBank:       0,
		maxTimeBank:    0,
		timePerMove:    0,
		candleInterval: 0,
		candleFormat:   []string{},
		candlesTotal:   0,
		candlesGiven:   0,
		initialStack:   0,
		transactionFee: 0,
		date:           0,
		buyingPrice:    0,
		stacks:         stacksDict,
		charts:         chartsDict,
	}
	runBot(botState)
}

func runBot(botState botState) {
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		var string_cleared string = text
		string_cleared = strings.Replace(string_cleared, "\n", "", -1)
		parseInput(string_cleared, botState)
	}
}

func parseInput(text string, botState botState) {
	stringSlice := strings.Split(text, " ")
	if strings.Compare(stringSlice[0], "settings") == 0 {
		update_settings(stringSlice[1], stringSlice[2], botState)
	}
	if strings.Compare(stringSlice[0], "update") == 0 {
		if stringSlice[1] == "game" {
			update_game(stringSlice[2], stringSlice[3], botState)
		}
	}
	if strings.Compare(stringSlice[0], "action") == 0 {
		print("action target")
	}
}

func update_settings(key string, value string, botState botState) {
	fmt.Println(key)
	fmt.Println(value)
	if strings.Compare(key, "timebank") == 0 {
		fmt.Println("im here")

		maxTimeBank, err := strconv.Atoi(value)
		handle_errors(err)
		timeBank, err := strconv.Atoi(value)
		handle_errors(err)
		botState.maxTimeBank = maxTimeBank
		botState.timeBank = timeBank
	}
	if strings.Compare(key, "time_per_move") == 0 {
		fmt.Println("im here2")

		timePerMove, err := strconv.Atoi(value)
		handle_errors(err)
		botState.timePerMove = timePerMove
	}
	if strings.Compare(key, "candle_interval") == 0 {
		fmt.Println("im here3")

		candleInterval, err := strconv.Atoi(value)
		handle_errors(err)
		botState.candleInterval = candleInterval
	}
	if strings.Compare(key, "candle_format") == 0 {
		candleFormat := strings.Split(value, ",")
		botState.candleFormat = candleFormat
	}
	if strings.Compare(key, "candles_total") == 0 {
		candles_total, err := strconv.Atoi(value)
		handle_errors(err)
		botState.candlesTotal = candles_total
	}
	if strings.Compare(key, "candles_given") == 0 {
		candlesGiven, err := strconv.Atoi(value)
		handle_errors(err)
		botState.candlesGiven = candlesGiven
	}
	if strings.Compare(key, "initial_stack") == 0 {
		initialStack, err := strconv.Atoi(value)
		handle_errors(err)
		botState.initialStack = initialStack
	}
	if strings.Compare(key, "transaction_fee_percent") == 0 {
		transactionFee, err := strconv.ParseFloat(value, 64)
		handle_errors(err)
		botState.transactionFee = transactionFee
	}
	fmt.Printf("%+v\n", botState)
}

func update_game(key string, value string, botState botState) {
	print("coucou1")
	print(key)

	if strings.Compare(key, "next_candles") == 0 {
		print("coucou2")
		new_candles := strings.Split(value, ";")
		tmp_date := strings.Split(value, ",")
		date, err := strconv.Atoi(tmp_date[1])
		handle_errors(err)
		botState.date = date
		for index, candle_str := range new_candles {
			fmt.Println(index)
			fmt.Println(candle_str)
			candle_infos := strings.Split(candle_str, ",")
			update_charts(candle_infos, candle_str)
		}

	}

}

func update_charts(candle_infos []string, candle_str string) {

}

func handle_errors(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
