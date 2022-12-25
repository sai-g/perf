package main

import (
	// Standard libraries
	"flag"
	"log"
	"os"
	"time"

	"github.com/quickfixgo/quickfix/field"

	// Custom libraries
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/enum"
	"github.com/quickfixgo/quickfix/fix42/newordersingle"
)

var (
	fixconfig  = flag.String("fixconfig", "outbound.cfg", "FIX config file")
	sampleSize = flag.Int("samplesize", 1000, "Expected sample size")
)

var (
	SessionID quickfix.SessionID
	start     = make(chan interface{})
	app       = &OutboundRig{}
)

func main() {
	flag.Parse()

	cfg, err := os.Open(*fixconfig)
	if err != nil {
		log.Fatal(err)
	}

	appSettings, err := quickfix.ParseSettings(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// logFactory, err := quickfix.NewFileLogFactory(appSettings)
	logFactory := quickfix.NewNullLogFactory()
	if err != nil {
		log.Fatal(err)
	}

	storeFactory := quickfix.NewMemoryStoreFactory()

	initiator, err := quickfix.NewInitiator(app, storeFactory, appSettings, logFactory)
	if err != nil {
		log.Fatal(err)
	}
	if err = initiator.Start(); err != nil {
		log.Fatal(err)
	}

	<-start

	for i := 0; i < *sampleSize; i++ {
		order := newordersingle.New(
			field.NewClOrdID("100"),
			field.NewHandlInst(enum.HandlInst_AUTOMATED_EXECUTION_ORDER_PRIVATE_NO_BROKER_INTERVENTION),
			field.NewSymbol("TSLA"),
			field.NewSide(enum.Side_BUY),
			field.NewTransactTime(time.Now()),
			field.NewOrdType(enum.OrdType_MARKET),
		)
		quickfix.SendToTarget(order, SessionID)
		// time.Sleep(1 * time.Millisecond)
	}

	time.Sleep(2 * time.Second)
}

type OutboundRig struct {
	quickfix.SessionID
}

func (e OutboundRig) OnCreate(sessionID quickfix.SessionID) {}
func (e *OutboundRig) OnLogon(sessionID quickfix.SessionID) {
	SessionID = sessionID
	start <- "START"
}
func (e OutboundRig) OnLogout(sessionID quickfix.SessionID)                              {}
func (e OutboundRig) ToAdmin(msgBuilder *quickfix.Message, sessionID quickfix.SessionID) {}
func (e OutboundRig) ToApp(msgBuilder *quickfix.Message, sessionID quickfix.SessionID) (err error) {
	return
}

func (e OutboundRig) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	return
}

func (e OutboundRig) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	return
}
