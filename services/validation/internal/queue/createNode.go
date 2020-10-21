package queue

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/MurmurationsNetwork/MurmurationsServices/common/logger"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/nats_utils"
	"github.com/MurmurationsNetwork/MurmurationsServices/common/randstr_util"
	"github.com/cenkalti/backoff"
	"github.com/nats-io/stan.go"
)

func Listen() {
	opts := []stan.Option{stan.NatsURL("http://nats-svc:4222")}

	var sc stan.Conn
	op := func() error {
		var err error
		sc, err = stan.Connect("murmurations", randstr_util.String(8), opts...)
		if err != nil {
			return err
		}
		return nil
	}

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 2 * time.Minute
	notify := func(err error, time time.Duration) {
		fmt.Printf("%+v \n", err)
		fmt.Printf("%+v \n", time)
	}
	err := backoff.RetryNotify(op, b, notify)
	if err != nil {
		logger.Panic("error when trying to connect nats", err)
	}

	nodeCreatedListener := nats_utils.NewListener(sc, "node:created", "indexer-svc-qgroup")
	err = nodeCreatedListener.Listen()
	if err != nil {
		sc.Close()
		log.Fatal(err)
	}

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			sc.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
