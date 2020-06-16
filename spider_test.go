package main

import (
	. "spider/selenium"
	"testing"
)

func TestExample2(t *testing.T) {
	// 爬房源信息
	var startCrawler chan bool
	StartLoopCrawler(startCrawler, "pc")
}
