package confunc

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

type EventData struct {
	Str string
}

func Test_f(t *testing.T) {
	fn0 := func(event interface{}) {
		data := event.(*EventData)
		time.Sleep(time.Second)
		fmt.Println("fn0: ", data.Str)
	}
	fn1 := func(event interface{}) {
		data := event.(*EventData)
		time.Sleep(time.Second)
		fmt.Println("fn1: ", data.Str)
	}
	con, err := NewConfunc(10, 2, fn0, fn1)
	if err != nil {
		fmt.Println(err)
		return
	}
	con.Start()
	go func() {
		for {
			con.Handle(&EventData{Str: "test"})
			fmt.Println("inserted.")
			time.Sleep(time.Millisecond * 100)
			fmt.Println("goroutine: ", runtime.NumGoroutine())
		}
	}()
	select {}
}
