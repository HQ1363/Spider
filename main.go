package main

import (
	"fmt"
	"runtime/debug"
	. "spider/selenium"
)

var (
	startCrawler = make(chan bool)
)

func main()  {
	// example 1
	//SetupWriter()
	//StartChrome()
	//StartCrawler()

	// example2
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("some error has occurred, info: ", r)
			debug.PrintStack()
		}
	}()
	go StartLoopCrawler(startCrawler)
	for {
		select {
		 	case run := <- startCrawler:
		 		if run {
					fmt.Println("crawler run success")
				} else {
					fmt.Println("crawler run failure")
				}
				go StartLoopCrawler(startCrawler)
		}
	}
}
