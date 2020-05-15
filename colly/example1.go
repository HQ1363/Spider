package colly

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	_ "github.com/gocolly/colly/v2/proxy"
	_ "log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"
)

// colly + goquery 爬虫利器
// 考虑下分布式的问题
// blog: https://segmentfault.com/a/1190000019969473
// goquery: https://www.jianshu.com/p/ae172d60c431
// go系列教程: https://www.jianshu.com/nb/5264832   https://studygolang.com/subject/2

var (
	cityList = []map[string]string{
		{
			"city_abbr": "su",
			"city_name": "苏州",
			"city_url":  "https://suzhou.newhouse.fang.com/house/s/",
		},
		{
			"city_abbr": "wx",
			"city_name": "无锡",
			"city_url":  "https://wuxi.newhouse.fang.com/house/s/",
		},
		{
			"city_abbr": "zhjg",
			"city_name": "张家港",
			"city_url":  "https://zjg.newhouse.fang.com/house/s/",
		},
		{
			"city_abbr": "chss",
			"city_name": "常熟",
			"city_url":  "https://changshu.newhouse.fang.com/house/s/",
		},
		{
			"city_abbr": "ks",
			"city_name": "昆山",
			"city_url":  "https://ks.newhouse.fang.com/house/s/",
		},
		{
			"city_abbr": "yix",
			"city_name": "宜兴",
			"city_url":  "https://yixing.newhouse.fang.com/house/s/",
		},
		{
			"city_abbr": "jiangy",
			"city_name": "江阴",
			"city_url":  "https://jy.newhouse.fang.com/house/s/",
		},
		{
			"city_abbr": "taic",
			"city_name": "太仓",
			"city_url":  "https://tc.newhouse.fang.com/house/s/",
		},
		{
			"city_abbr": "nt",
			"city_name": "南通",
			"city_url":  "https://nt.newhouse.fang.com/house/s/",
		},
		{
			"city_abbr": "hu",
			"city_name": "湖州",
			"city_url":  "https://huzhou.newhouse.fang.com/house/s/",
		},
		{
			"city_abbr": "zj",
			"city_name": "镇江",
			"city_url":  "https://zhenjiang.newhouse.fang.com/house/s/",
		},
		{
			"city_abbr": "jx",
			"city_name": "嘉兴",
			"city_url":  "https://jx.newhouse.fang.com/house/s/",
		},
		{
			"city_abbr": "yizh",
			"city_name": "仪征",
			"city_url":  "https://yizheng.newhouse.fang.com/house/s/",
		},
		{
			"city_abbr": "sh",
			"city_name": "上海",
			"city_url":  "https://sh.newhouse.fang.com/house/s/",
		},
		{
			"city_abbr": "tz",
			"city_name": "泰州",
			"city_url":  "https://taizhou.newhouse.fang.com/house/s/",
		},
		{
			"city_abbr": "km",
			"city_name": "昆明",
			"city_url":  "https://km.newhouse.fang.com/house/s/",
		},
	}
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func SetupCollyProxy() {
	fmt.Println("配置colly代理")
}

func InitDbConnection() {
	fmt.Println("初始化MySQL连接")
}

func InitRedisConnection() {
	fmt.Println("初始化Redis连接")
}

func dealWithPagination() {
	fmt.Println("处理分页问题")
}

func main() {
	writer, err := os.OpenFile("../collector.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("newhouse.fang.com", "fang.com"),
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36"),
		colly.Debugger(&debug.LogDebugger{Output: writer}),
		colly.MaxDepth(2),
	)

	c.SetRedirectHandler(func(req *http.Request, via []*http.Request) error {
		// log.Println(len(via), 111, req.URL.String(), 222, via[0].URL.String())

		return http.ErrUseLastResponse
	})

	// * 取一个代理IP 并设置
	//subProxyStr = mproxy.GetSubProxy()
	//rp, err := proxy.RoundRobinProxySwitcher(subProxyStr)
	//if err != nil {
	//	log.Println("设置代理错误", subProxyStr, err)
	//}
	//c.SetProxyFunc(rp)

	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, // 超时时间
			KeepAlive: 30 * time.Second, // keepAlive 超时时间
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,              // 最大空闲连接数
		IdleConnTimeout:       30 * time.Second, // 空闲连接超时
		TLSHandshakeTimeout:   30 * time.Second, // TLS 握手超时
		ExpectContinueTimeout: 1 * time.Second,
	})

	c.WithTransport(&http.Transport{
		DisableKeepAlives: true,
	})

	// Limit the number of threads started by colly to two
	_ = c.Limit(&colly.LimitRule{
		DomainGlob:  "*anjuke.*",
		Parallelism: 5000,
		//Delay:      5 * time.Second,
	})

	// 这里的是否访问过的逻辑重写了
	c.AllowURLRevisit = true

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		_ = c.Visit(e.Request.AbsoluteURL(link))
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		// 动态更新user-agent信息
		//r.Headers.Set("User-Agent", RandomString())
		fmt.Println("Visiting", r.URL.String())
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		// 此处可能要考虑切换代理
	})

	// Create another collector to scrape course details
	//detailCollector := c.Clone()

	// Extract details of the course
	//detailCollector.OnHTML(`div[id=rendered-content]`, func(e *colly.HTMLElement) {
	//	log.Println("Course found", e.Request.URL)
	//	title := e.ChildText(".course-title")
	//	if title == "" {
	//		log.Println("No title found", e.Request.URL)
	//	}
	//	course := Course{
	//		Title:       title,
	//		URL:         e.Request.URL.String(),
	//		Description: e.ChildText("div.content"),
	//		Creator:     e.ChildText("div.creator-names > span"),
	//	}
	//	// Iterate over rows of the table which contains different information
	//	// about the course
	//	e.ForEach("table.basic-info-table tr", func(_ int, el *colly.HTMLElement) {
	//		switch el.ChildText("td:first-child") {
	//		case "Language":
	//			course.Language = el.ChildText("td:nth-child(2)")
	//		case "Level":
	//			course.Level = el.ChildText("td:nth-child(2)")
	//		case "Commitment":
	//			course.Commitment = el.ChildText("td:nth-child(2)")
	//		case "How To Pass":
	//			course.HowToPass = el.ChildText("td:nth-child(2)")
	//		case "User Ratings":
	//			course.Rating = el.ChildText("td:nth-child(2) div:nth-of-type(2)")
	//		}
	//	})
	//	courses = append(courses, course)
	//})

	// Start scraping on https://hackerspaces.org
	_ = c.Visit("https://hackerspaces.org/")
}