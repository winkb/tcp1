package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	std := os.Stdin
	scanner := bufio.NewScanner(std)

	go func() {
		tk := time.NewTimer(time.Second)
		select {
		case <-tk.C:
			err := std.Close()
			if err != nil {
				panic(err)
			}
		}
	}()

	for scanner.Scan() { // 无法停止,只能等待输入exit
		text := scanner.Text()
		if text == "exit;" {
			break
		}
		fmt.Println("input:", text)
	}
}
