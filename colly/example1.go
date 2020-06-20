package colly

import (
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	_ "github.com/gocolly/colly/v2/proxy"
	"github.com/jinzhu/gorm"
	"log"
	_ "log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"spider/utils"
	"time"
)

func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

type PageAddrStruct struct {
	gorm.Model
	Name         string  `gorm:"type:varchar(2048);column:name"`
	Address      string  `gorm:"index:addr;column:address"` // 给address字段创建名为addr的索引
	PageContent  string  `gorm:"type:mediumtext;column:page_content"`  // 存放整个网页信息
}

func (PageAddrStruct) TableName() string {
	return "WebHtmlPage"
}

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
		//{
		//	"city_abbr": "wx",
		//	"city_name": "无锡",
		//	"city_url":  "https://wuxi.newhouse.fang.com/house/s/",
		//},
		//{
		//	"city_abbr": "zhjg",
		//	"city_name": "张家港",
		//	"city_url":  "https://zjg.newhouse.fang.com/house/s/",
		//},
		//{
		//	"city_abbr": "chss",
		//	"city_name": "常熟",
		//	"city_url":  "https://changshu.newhouse.fang.com/house/s/",
		//},
		//{
		//	"city_abbr": "ks",
		//	"city_name": "昆山",
		//	"city_url":  "https://ks.newhouse.fang.com/house/s/",
		//},
		//{
		//	"city_abbr": "yix",
		//	"city_name": "宜兴",
		//	"city_url":  "https://yixing.newhouse.fang.com/house/s/",
		//},
		//{
		//	"city_abbr": "jiangy",
		//	"city_name": "江阴",
		//	"city_url":  "https://jy.newhouse.fang.com/house/s/",
		//},
		//{
		//	"city_abbr": "taic",
		//	"city_name": "太仓",
		//	"city_url":  "https://tc.newhouse.fang.com/house/s/",
		//},
		//{
		//	"city_abbr": "nt",
		//	"city_name": "南通",
		//	"city_url":  "https://nt.newhouse.fang.com/house/s/",
		//},
		//{
		//	"city_abbr": "hu",
		//	"city_name": "湖州",
		//	"city_url":  "https://huzhou.newhouse.fang.com/house/s/",
		//},
		//{
		//	"city_abbr": "zj",
		//	"city_name": "镇江",
		//	"city_url":  "https://zhenjiang.newhouse.fang.com/house/s/",
		//},
		//{
		//	"city_abbr": "jx",
		//	"city_name": "嘉兴",
		//	"city_url":  "https://jx.newhouse.fang.com/house/s/",
		//},
		//{
		//	"city_abbr": "yizh",
		//	"city_name": "仪征",
		//	"city_url":  "https://yizheng.newhouse.fang.com/house/s/",
		//},
		//{
		//	"city_abbr": "sh",
		//	"city_name": "上海",
		//	"city_url":  "https://sh.newhouse.fang.com/house/s/",
		//},
		//{
		//	"city_abbr": "tz",
		//	"city_name": "泰州",
		//	"city_url":  "https://taizhou.newhouse.fang.com/house/s/",
		//},
		//{
		//	"city_abbr": "km",
		//	"city_name": "昆明",
		//	"city_url":  "https://km.newhouse.fang.com/house/s/",
		//},
	}
	needDownPageList = []string{"楼盘详情", "楼盘动态", "楼盘点评", "房价走势"}
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

func StartCrawlerHtmlPage() {
	writer, err := os.OpenFile("../collector.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		//colly.AllowedDomains("*.newhouse.fang.com", "newhouse.fang.com", "fang.com"),
		colly.Async(false),
		colly.UserAgent(RandomString()),
		colly.Debugger(&debug.LogDebugger{Output: writer}),
		colly.MaxDepth(2),
		colly.AllowURLRevisit(),
	)

	// Create another collector to scrape course details
	detailCollector := c.Clone()

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
	//_ = c.Limit(&colly.LimitRule{
	//	DomainGlob:  "*anjuke.*",
	//	Parallelism: 5000,
	//})

	// 这里的是否访问过的逻辑重写了
	c.AllowURLRevisit = true

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

	// On every a element which has href attribute call callback
	c.OnHTML("div[class=nhouse_list]", func(e *colly.HTMLElement) {
		//fmt.Println("find house list dom")
		e.ForEach("div>ul>li", func(_ int, element *colly.HTMLElement) {
			houseDetailPageUrl := element.ChildAttr("div.clearfix > div:nth-child(2) > div:first-child > div:first-child a", "href")
			houseName := element.ChildText("div.clearfix > div:nth-child(2) > div:first-child > div:first-child a")
			fmt.Println("find house: ", ConvertToString(houseName, "gbk", "utf-8"), ", Url: ", e.Request.AbsoluteURL(houseDetailPageUrl))
			// Only those links are visited which are in AllowedDomains
			// Visit link found on page
			if detailErr := detailCollector.Visit(e.Request.AbsoluteURL(houseDetailPageUrl)); detailErr != nil {
				fmt.Println("visitor detail page failure, ", detailErr.Error())
			}
		})
	})

	// Extract details of the course
	detailCollector.OnHTML(`div[id=header-wrap]`, func(e *colly.HTMLElement) {
		log.Println("top nav found", e.Request.URL)
		e.ForEach("div[id=orginalNaviBox] a", func(_ int, el *colly.HTMLElement) {
			downUrl := e.Request.AbsoluteURL(el.Attr("href"))
			downName := ConvertToString(el.Text, "gbk", "utf-8")
			if utils.InList(downName, needDownPageList) {
				fmt.Println("start download page: ", downName, ", Url: ", downUrl)
			} else {
				fmt.Println("skip page: ", downName, ", url: ", downUrl)
			}
		})
	})

	// start scraping on each city
	for _, city := range cityList {
		fmt.Println("visit url: ", city["city_url"])
		time.Sleep(time.Second * 1)
		if err := c.Visit(city["city_url"]); err != nil {
			fmt.Println(err.Error())
		}
	}
}
