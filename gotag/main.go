package main

import (
	"bufio"
	"fmt"
	"github.com/winkb/tcp1/gotag/chain"
	"os"
	"os/exec"
)

func doShell(str []string) {
	command := exec.Command(str[0], str[1:]...)

	output, err2 := command.CombinedOutput()
	if err2 != nil {
		fmt.Println(err2)
	}
	fmt.Println(string(output))
}

func main() {
	var inputCh = make(chan string)
	var collect = chain.NewCollect(inputCh,
		"1.0.0",
		chain.NewVersionCollect(),
	)

	go func() {
		collect.Collect()
		close(inputCh)
		tagCommand := collect.GetData()

		for _, v := range tagCommand {
			doShell(v)
		}

		fmt.Println("ok!")

		doShell([]string{"git", "tag"})
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		txt := scanner.Text()
		select {
		case <-inputCh:
			return
		default:
			inputCh <- txt
		}
	}
}
