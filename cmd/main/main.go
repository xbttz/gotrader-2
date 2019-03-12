package main

import (
	"gotrader/internal/central"
	"gotrader/internal/display"
	"gotrader/internal/logic"
)

func main() {
	for {
		daemonize()
	}
}

func daemonize() {
	central.InitFlag()
	display.HelloMsg(central.Asset())

	central.TelegramSend(
		display.HelloMsg(central.Asset()),
	)

	typeOrder := logic.CandleRunner()
	central.CreateOrder(typeOrder)

	if typeOrder == "Buy" {
		central.ClosePositionProfitBuy()
	} else {
		central.ClosePositionProfitSell()
	}
	central.GetProfit()
}
