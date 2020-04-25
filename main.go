package main

import (
	. "spider/selenium"
)

func main()  {
	SetupWriter()
	StartChrome()
	StartCrawler()
}
