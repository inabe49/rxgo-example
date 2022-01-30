package main

import (
	"context"
	"fmt"
	"github.com/reactivex/rxgo/v2"
	"time"
)

func ToCancelable(observable rxgo.Observable) (rxgo.Observable, func()) {
	canceled := false
	cancel := func() {
		canceled = true
	}

	return observable.Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		if canceled {
			return nil, fmt.Errorf("canceled")
		}

		return i, nil
	}), cancel
}

func main() {
	ch := make(chan rxgo.Item)
	observable := rxgo.FromEventSource(ch, rxgo.WithBackPressureStrategy(rxgo.Drop))

	// ctx, cancel := context.WithCancel(context.Background())

	ch <- rxgo.Of(1)
	ch <- rxgo.Of(2)
	ch <- rxgo.Of(3)

	go func() {
		for i := 4; i < 10; i++ {
			ch <- rxgo.Of(i)
			<-time.After(1 * time.Second)
		}

		close(ch)
	}()

	go func() {
		// Observer 1
		observable.ForEach(func(i interface{}) {
			println("Observer 1 : value : ", i.(int))
		}, func(err error) {
			println("Observer 1 : error : ", err.Error())
		}, func() {
			println("Observer 1 : finish")
		})
	}()

	go func() {
		// Observer 2
		observable.ForEach(func(i interface{}) {
			println("Observer 2 : value : ", i.(int))
		}, func(err error) {
			println("Observer 2 : error : ", err.Error())
		}, func() {
			println("Observer 2 : finish")
		})
	}()

	obs, cancel := ToCancelable(observable)

	go func() {
		<-time.After(1 * time.Second)

		cancel()
	}()

	go func() {
		// Observer 3
		obs.ForEach(func(i interface{}) {
			println("Observer 3 : value : ", i.(int))
		}, func(err error) {
			println("Observer 3 : error : ", err.Error())
		}, func() {
			println("Observer 3 : finish")
		})
	}()

	/*
		for item := range observable.Map(func(ctx context.Context, i interface{}) (interface{}, error) {
			return nil, fmt.Errorf("error")
		}).Observe() {
			println(item.E)
		}
		observable.Observe()
	*/

	/*
		for item := range observable.Observe() {
			println(item.V.(int))

			<-time.After(2 * time.Second)
		}
	*/

	<-time.After(10 * time.Second)
}
