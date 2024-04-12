package util

import (
	"fmt"
	"log"
	"sync"
)

func MyGoWg(wg *sync.WaitGroup, name string, f func()) {
	wg.Add(1)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Print(fmt.Sprintf("goroutine %s %s", name, err))

				MyGoWg(wg, name, f)

				wg.Done()
				return
			}

			fmt.Printf("goroutine %s defer\n", name)

			wg.Done()
		}()

		f()
	}()
}

func myGo(name string, f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Print(name, err)
				return
			}

			fmt.Printf("goroutine %s defer", name)
		}()

		f()
	}()
}
