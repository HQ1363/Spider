package main

import (
	. "spider/selenium"
)

func main()  {
	// example 1
	//SetupWriter()
	//StartChrome()
	//StartCrawler()

	// example2
	//if port, err := PickUnusedPort(); err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	fmt.Println(port)
	//}
	StartLoopCrawler()
}
