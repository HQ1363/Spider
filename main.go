package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/abiosoft/ishell"
	"runtime/debug"
	. "spider/colly"
	. "spider/selenium"
	"strings"
	"time"
)

var (
	startCrawler = make(chan bool)
	Shell = ishell.New()
)

var BoldGreen = color.New(color.FgGreen, color.Bold).SprintFunc()
var BoldMagenta = color.New(color.FgMagenta, color.Bold).SprintFunc()
var BoldBlack = color.New(color.FgBlack, color.Bold).SprintFunc()
var BoldHiBlack = color.New(color.FgHiBlack, color.Bold).SprintFunc()

func SetCurrentPrompt(prompt string) {
	Shell.SetPrompt(fmt.Sprintf("%s%s%s$ ", BoldHiBlack("["), prompt, BoldHiBlack("]")))
}

func GreetWord() {
	greet := ""
	extend := ""
	username := fmt.Sprintf("%s", BoldBlack("HQ"))
	switch hour := time.Now().Hour(); {
	case hour > 0 && hour <= 6:
		greet = fmt.Sprintf("%s", BoldMagenta("凌晨好"))
		extend = fmt.Sprintf("%s", BoldGreen("熬夜对身体不好, 要早点休息哦 ╮(￣▽￣)╭"))
	case hour > 6 && hour <= 11:
		greet = fmt.Sprintf("%s", BoldMagenta("早上好"))
		extend = fmt.Sprintf("%s", BoldGreen("一日之计在于晨, 希望今天又是棒棒的一天 (*^▽^*)"))
	case hour > 11 && hour <= 13:
		greet = fmt.Sprintf("%s", BoldMagenta("中午好"))
		extend = fmt.Sprintf("%s", BoldGreen("现在是午饭午觉时间呢, 你可真勤奋 ヾ(◍°∇°◍)ﾉﾞ"))
	case hour > 13 && hour <= 19:
		greet = fmt.Sprintf("%s", BoldMagenta("下午好"))
		extend = fmt.Sprintf("%s", BoldGreen("快点干完, 早点下班 (＾－＾)V"))
	case hour > 19 && hour <= 23:
		greet = fmt.Sprintf("%s", BoldMagenta("晚上好"))
		extend = fmt.Sprintf("%s", BoldGreen("大晚上的还在在忙活, 你就是家庭的顶梁柱 ╮(￣▽￣)╭"))
	}
	Shell.Println(strings.Repeat("#", len(greet+username+extend)-45))
	Shell.Println(greet, username, extend)
	Shell.Println(strings.Repeat("#", len(greet+username+extend)-45))
}

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
		case run := <-startCrawler:
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
		case run := <-startCrawler:
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

func main() {
	//runSeleniumCrawler()
	//runCollyCrawler()
	GreetWord()
	Shell.AddCmd(&ishell.Cmd{
		Name: "crawlerMobile",
		Help: "run selenium crawler on mobile",
		Func: func(c *ishell.Context) {
			runSeleniumCrawlerMobile()
		},
	})
	Shell.Run()
}
