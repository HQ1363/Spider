package main

import (
	"fmt"
	"runtime/debug"
	. "spider/colly"
	. "spider/selenium"
)

var (
	startCrawler = make(chan bool)
)

func runSeleniumCrawler() {
	// example 1
	// SetupWriter()
	// StartChrome()
	// StartCrawler()

	// example 2
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("some error has occurred, info: ", r)
			debug.PrintStack()
		}
	}()
	go StartLoopCrawler(startCrawler, "pc")
	for {
		select {
		case run := <- startCrawler:
			if run {
				fmt.Println("crawler pc run success")
			} else {
				fmt.Println("crawler pc run failure")
			}
			go StartLoopCrawler(startCrawler, "pc")
		}
	}
}

func runSeleniumCrawlerMobile() {
	// example 2
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("some error has occurred, info: ", r)
			debug.PrintStack()
		}
	}()
	go StartLoopCrawler(startCrawler, "mobile")
	for {
		select {
		case run := <- startCrawler:
			if run {
				fmt.Println("crawler mobile run success")
			} else {
				fmt.Println("crawler mobile run failure")
			}
			go StartLoopCrawler(startCrawler, "mobile")
		}
	}
}

func runCollyCrawler() {
	// example 1
	StartCrawlerHtmlPage()
}

func main()  {
	//runSeleniumCrawler()
	runSeleniumCrawlerMobile()
	//runCollyCrawler()
}
