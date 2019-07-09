package central

import (
	"encoding/json"
	"log"
	"net/url"
	"time"

	"github.com/thiago-scherrer/gotrader/internal/api"
	cvt "github.com/thiago-scherrer/gotrader/internal/convert"
	"github.com/thiago-scherrer/gotrader/internal/display"
	rd "github.com/thiago-scherrer/gotrader/internal/reader"
)

// Order path to use on API Request
const oph string = "/api/v1/order"

// Position path to use on API Request
const poh string = "/api/v1/position?"

// Basic path to use on API Request
const ith string = "/api/v1/instrument?"

// Laverage path to use on API Request
const lth = "/api/v1/position/leverage"

// A random number to make a sleep before staring a new round
const tlp = 50

// parserAmount unmarshal a r API to return the wallet amount
func parserAmount(data []byte) int {
	ap := rd.APISimple()
	err := json.Unmarshal(data, &ap)
	if err != nil {
		log.Println("Error to get Amount: ", err)
	}
	return ap.Amount
}

// lastPrice unmarshal a r API to return the last price
func lastPrice(d []byte) float64 {
	ap := rd.APIArray()
	var r float64

	err := json.Unmarshal(d, &ap)
	if err != nil {
		log.Println("Error to get last price: ", err)
	}
	for _, v := range ap[:] {
		r = v.LastPrice
	}
	return r
}

func makeOrder(orderType string) string {
	ap := rd.APISimple()
	hfl := cvt.IntToString(rd.Hand())
	ast := rd.Asset()
	prc := cvt.FloatToString(Price())
	rtp := "POST"

	u := url.Values{}
	u.Set("symbol", ast)
	u.Add("side", orderType)
	u.Add("orderQty", hfl)
	u.Add("price", prc)
	u.Add("ordType", "Limit")
	data := cvt.StringToBytes(u.Encode())

	for {
		glt := api.ClientRobot(rtp, oph, data)
		err := json.Unmarshal(glt, &ap)
		if err != nil {
			log.Println("Error to make a order:", err)
			time.Sleep(time.Duration(5) * time.Second)
		} else {
			return ap.OrderID
		}
	}
}

// GetPosition get the actual open possitions
func GetPosition() float64 {
	ap := rd.APIArray()
	var r float64
	pth := poh + `filter={"symbol":"` + rd.Asset() + `"}&count=1`
	rtp := "GET"
	dt := cvt.StringToBytes("message=GoTrader bot&channelID=1")
	glt := api.ClientRobot(rtp, pth, dt)
	err := json.Unmarshal(glt, &ap)
	if err != nil {
		log.Println("Error to get pst:", err)
	}

	for _, v := range ap[:] {
		r = v.AvgEntryPrice
	}
	return r
}

// Price return the actual asset price
func Price() float64 {
	ast := rd.Asset()

	u := url.Values{}
	u.Set("symbol", ast)
	u.Add("count", "100")
	u.Add("reverse", "false")
	u.Add("columns", "lastPrice")

	p := ith + u.Encode()
	d := cvt.StringToBytes("message=GoTrader bot&channelID=1")
	g := api.ClientRobot("GET", p, d)

	return lastPrice(g)
}

// ClosePosition close all opened position
func ClosePosition(priceClose string) string {
	ast := rd.Asset()
	path := oph
	rtp := "POST"

	u := url.Values{}
	u.Set("symbol", ast)
	u.Add("execInst", "Close")
	u.Add("price", priceClose)
	u.Add("ordType", "Limit")

	data := cvt.StringToBytes(u.Encode())

	for {
		glt := api.ClientRobot(rtp, path, data)
		if waitCreateOrder() {
			return cvt.BytesToString(glt)
		}
	}
}

func setLeverge() {
	ast := rd.Asset()
	path := lth
	rtp := "POST"
	l := rd.Leverage()
	u := url.Values{}
	u.Set("symbol", ast)
	u.Add("leverage", l)
	data := cvt.StringToBytes(u.Encode())

	api.ClientRobot(rtp, path, data)
	log.Println(display.SetleverageMsg(rd.Asset(), l))
	api.MatrixSend(display.SetleverageMsg(rd.Asset(), l))

}

func statusOrder() bool {
	path := poh + `filter={"symbol":"` + rd.Asset() + `"}&count=1`
	data := cvt.StringToBytes("message=GoTrader bot&channelID=1")
	glt := api.ClientRobot("GET", path, data)
	return opening(glt)
}

func opening(data []byte) bool {
	ap := rd.APIArray()
	var r bool

	err := json.Unmarshal(data, &ap)
	if err != nil {
		log.Println("json open error:", err)
	}
	for _, v := range ap[:] {
		r = v.IsOpen
	}
	return r
}

// CreateOrder create the order on bitmex
func CreateOrder(typeOrder string) bool {
	setLeverge()
	orderTimeOut()
	makeOrder(typeOrder)

	for i := 0; i < 3; i++ {
		if waitCreateOrder() {
			log.Println(display.OrderCreatedMsg(rd.Asset(), typeOrder))
			api.MatrixSend(display.OrderCreatedMsg(rd.Asset(), typeOrder))
			return true
		}
		time.Sleep(time.Duration(1) * time.Minute)
	}
	return false
}

func waitCreateOrder() bool {
	if statusOrder() == true {
		log.Println(display.OrderDoneMsg(rd.Asset()))
		api.MatrixSend(display.OrderDoneMsg(rd.Asset()))
		return true
	}
	return false
}

func orderTimeOut() {
	poh := "/api/v1/order/cancelAllAfter?"
	data := cvt.StringToBytes("message=GoTrader bot&channelID=1")

	u := url.Values{}
	u.Set("timeout", "120000")

	p := poh + u.Encode()
	api.ClientRobot("POST", p, data)

}

// GetProfit waint to start a new trade round
func GetProfit() bool {
	log.Println(display.OrderWaintMsg(rd.Asset()))
	api.MatrixSend(display.OrderWaintMsg(rd.Asset()))

	for {
		if statusOrder() == false {
			log.Println(display.ProfitMsg(rd.Asset()))
			api.MatrixSend(display.ProfitMsg(rd.Asset()))
			time.Sleep(time.Duration(tlp) * time.Second)
			return true
		}
		time.Sleep(time.Duration(10) * time.Second)
	}
}
