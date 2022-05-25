package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var STANDARD_DEVIATION_MULTIPLIER = 2
var SMA_PERIOD = 100

type botStatistics struct {
	SMA         []float64
	FMA         []float64
	EMA         []float64
	UPPER_BB    []float64
	LOWER_BB    []float64
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
	stacks         map[string]float64
	charts         map[string]Charts
	stats          botStatistics
}

type Candle struct {
	pair   string
	date   int
	high   float64
	low    float64
	open   float64
	close  float64
	volume float64
}

type Charts struct {
	pair   []string
	date   []int
	high   []float64
	low    []float64
	open   []float64
	close  []float64
	volume []float64
}

func main() {

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
		charts:         make(map[string]Charts),
		stacks:         make(map[string]float64),
	}
	runBot(&botState)
}

func runBot(botState *botState) {
	stdin := bufio.NewReader(os.Stdin)

	for {
		text, _ := stdin.ReadString('\n')
		var string_cleared string = text
		string_cleared = strings.Replace(string_cleared, "\n", "", -1)
		parseInput(string_cleared, botState)
	}
}

func parseInput(text string, botState *botState) {
	stringSlice := strings.Split(text, " ")
	if strings.Compare(stringSlice[0], "settings") == 0 {
		update_settings(stringSlice[1], stringSlice[2], botState)
		fmt.Printf("\n\n%v\n\n", botState)

	}
	if strings.Compare(stringSlice[0], "update") == 0 {
		if stringSlice[1] == "game" {
			update_game(stringSlice[2], stringSlice[3], botState)
		}
	}
	if stringSlice[0] == "action" {
		handle_action(botState)
	}
}

func handle_action(botState *botState) {
	dollars := botState.stacks["USDT"]
	btc := botState.stacks["BTC"]
	fmt.Printf("LEN OF CLOSE => %d", len(botState.charts["USDT_BTC"].close))
	get_moving_average(botState, botState.charts["USDT_BTC"].close, botState.charts["USDT_BTC"].high, botState.charts["USDT_BTC"].low)
	compute_bollinger_bands(botState, botState.charts["USDT_BTC"].close)
	handle_signals(botState, botState.charts["USDT_BTC"].close)
	current_closing_price := botState.charts["USDT_BTC"].close[len(botState.charts["USDT_BTC"].close)-1]
	affordable := dollars / current_closing_price

	fmt.Printf("My stacks are USDT: %f and BTC: %f. The current closing price is %f . So I can afford %f", dollars, btc, current_closing_price, affordable)
}

func get_moving_average(botState *botState, close []float64, highs []float64, lows []float64) {
	get_slow_moving_average(botState, close, highs, lows)
}

func handle_signals(botState *botState, close []float64) {
	// if closes[-2] < self.UPPER_BB[len(self.UPPER_BB)-2] and closes[-1] > self.UPPER_BB[len(self.UPPER_BB)-1]:
	// 	self.actionOrder = "SELL"
	if len(botState.stats.UPPER_BB) > 1 {
		if close[len(close)-2] < botState.stats.UPPER_BB[len(botState.stats.UPPER_BB)-2] && close[len(close)-1] > botState.stats.UPPER_BB[len(botState.stats.UPPER_BB)-1] {
			botState.stats.actionOrder = "SELL"
		}
	}
	if len(botState.stats.LOWER_BB) > 1 {
		if close[len(close)-2] > botState.stats.LOWER_BB[len(botState.stats.LOWER_BB)-2] && close[len(close)-1] < botState.stats.LOWER_BB[len(botState.stats.LOWER_BB)-1] {
			botState.stats.actionOrder = "BUY"
		}
	}
}

func compute_bollinger_bands(botState *botState, close []float64) {
	get_lower_band(botState, close)
	get_upper_band(botState, close)
}

func get_lower_band(botState *botState, close []float64) {
	STOCK_SMA := botState.stats.SMA[len(botState.stats.SMA)-1]
	var deviationArray []float64
	i := len(close) - 1
	j := 0

	fmt.Println(close)
	for j != SMA_PERIOD+1 {
		deviationArray = append(deviationArray, close[i])
		i = i - 1
		j = j + 1
	}

	SMA_STANDARD_DEVIATION := StdDev(deviationArray)
	temp_LOWER_BB := STOCK_SMA - float64(STANDARD_DEVIATION_MULTIPLIER)*SMA_STANDARD_DEVIATION
	fmt.Println("temp_LOWER_BB")
	fmt.Println(temp_LOWER_BB)

	botState.stats.LOWER_BB = append(botState.stats.LOWER_BB, temp_LOWER_BB)
}

func get_upper_band(botState *botState, close []float64) {
	STOCK_SMA := botState.stats.SMA[len(botState.stats.SMA)-1]
	var deviationArray []float64
	i := len(close) - 1
	j := 0

	fmt.Println(close)
	for j != SMA_PERIOD+1 {
		deviationArray = append(deviationArray, close[i])
		i = i - 1
		j = j + 1
	}

	SMA_STANDARD_DEVIATION := StdDev(deviationArray)
	temp_UPPER_BB := STOCK_SMA + float64(STANDARD_DEVIATION_MULTIPLIER)*SMA_STANDARD_DEVIATION
	fmt.Println("temp_UPPER_BB")
	fmt.Println(temp_UPPER_BB)

	botState.stats.UPPER_BB = append(botState.stats.UPPER_BB, temp_UPPER_BB)
}

func get_slow_moving_average(botState *botState, close []float64, highs []float64, lows []float64) {
	i := len(close) - 1
	temp_SMA := 0.0
	TP := 0.0
	j := 0
	// Average closing value over the last SMA period
	for j != SMA_PERIOD {
		TP = (highs[i] + lows[i] + close[i]) / 3
		temp_SMA = temp_SMA + TP
		i = i - 1
		j = j + 1
	}

	fmt.Printf("\n\n\ntemp_SMA => %f\n\n\n", (temp_SMA / float64(SMA_PERIOD)))

	botState.stats.SMA = append(botState.stats.SMA, temp_SMA/float64(SMA_PERIOD))
}

func update_settings(key string, value string, botState *botState) {
	if strings.Compare(key, "timebank") == 0 {
		maxTimeBank, err := strconv.Atoi(value)
		handle_errors(err)
		timeBank, err := strconv.Atoi(value)
		handle_errors(err)
		botState.maxTimeBank = maxTimeBank
		botState.timeBank = timeBank
	}
	if strings.Compare(key, "time_per_move") == 0 {
		timePerMove, err := strconv.Atoi(value)
		handle_errors(err)
		botState.timePerMove = timePerMove
	}
	if strings.Compare(key, "candle_interval") == 0 {
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

}

func update_game(key string, value string, botState *botState) {
	if strings.Compare(key, "next_candles") == 0 {
		new_candles := strings.Split(value, ";")
		tmp_date := strings.Split(value, ",")
		date, err := strconv.Atoi(tmp_date[1])
		handle_errors(err)
		botState.date = date
		for _, candle_str := range new_candles {
			candle_infos := strings.Split(candle_str, ",")
			update_charts(candle_infos[0], candle_str, botState)
		}
	}
	if strings.Compare(key, "stacks") == 0 {
		tmp_date := strings.Split(value, ",")
		for _, candle_str := range tmp_date {
			candle_infos := strings.Split(candle_str, ":")
			update_stacks(candle_infos[0], candle_infos[1], botState)
		}
	}

}

func update_stacks(key string, value string, botState *botState) {
	valFloat, err := strconv.ParseFloat(value, 64)
	handle_errors(err)
	botState.stacks[key] = valFloat
}

func update_charts(pair string, new_candle_str string, botState *botState) {
	if len(botState.charts) == 0 {
		botState.charts[pair] = Charts{}
	}

	new_candle_obj := initCandle(botState.candleFormat, new_candle_str)

	addCandle(new_candle_obj, botState, pair)
}

func addCandle(candle Candle, botState *botState, pair string) {

	if entry, ok := botState.charts[pair]; ok {
		// Then we modify the copy
		entry.date = append(entry.date, candle.date)
		entry.open = append(entry.open, candle.open)
		entry.high = append(entry.high, candle.high)
		entry.low = append(entry.low, candle.low)
		entry.close = append(entry.close, candle.close)
		entry.volume = append(entry.volume, candle.volume)

		// Then we reassign map entry
		botState.charts[pair] = entry
	}
}

func initCandle(format []string, intel string) Candle {
	newCandle := Candle{}
	tmp := strings.Split(intel, ",")

	for i, key := range format {
		value := tmp[i]
		if key == "pair" {
			newCandle.pair = value
		}
		if strings.Compare(key, "date") == 0 {
			date, err := strconv.Atoi(value)
			handle_errors(err)
			newCandle.date = date
		}
		if key == "high" {
			high, err := strconv.ParseFloat(value, 64)
			handle_errors(err)
			newCandle.high = high
		}
		if key == "low" {
			low, err := strconv.ParseFloat(value, 64)
			handle_errors(err)
			newCandle.low = low
		}
		if key == "open" {
			open, err := strconv.ParseFloat(value, 64)
			handle_errors(err)
			newCandle.open = open
		}
		if key == "close" {
			close, err := strconv.ParseFloat(value, 64)
			handle_errors(err)
			newCandle.close = close
		}
		if key == "volume" {
			volume, err := strconv.ParseFloat(value, 64)
			handle_errors(err)
			newCandle.volume = volume
		}
	}

	return newCandle
}

func handle_errors(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

func Variance(xs []float64) float64 {
	if len(xs) == 0 {
		return math.NaN()
	} else if len(xs) <= 1 {
		return 0
	}

	// Based on Wikipedia's presentation of Welford 1962
	// (http://en.wikipedia.org/wiki/Algorithms_for_calculating_variance#Online_algorithm).
	// This is more numerically stable than the standard two-pass
	// formula and not prone to massive cancellation.
	mean, M2 := 0.0, 0.0
	for n, x := range xs {
		delta := x - mean
		mean += delta / float64(n+1)
		M2 += delta * (x - mean)
	}
	return M2 / float64(len(xs)-1)
}

func StdDev(xs []float64) float64 {
	return math.Sqrt(Variance(xs))
}
