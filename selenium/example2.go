package selenium

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
	"net"
	"runtime/debug"
	"strings"
	"time"
)

var (
	urlBaidu = "https://www.baidu.com/s?wd=%s&rsv_spt=1&rsv_iqid=0xd67f9b330001219c&issp=1&f=8&rsv_bp=1&rsv_idx=2&ie=utf-8&rqlang=cn&tn=baiduhome_pg&rsv_enter=0&rsv_dl=tb&rsv_t=c721ztPqoWrMLlh87vzqI58VrneEgNDDUV42nx9LrTE9gk9OAhCn%2Baq9GIdlzUOhLKSp&oq=%25E6%2581%2592%25E5%25A4%25A7%25E6%259E%2597%25E6%25BA%25AA%25E9%2583%25A1&rsv_btype=t&rsv_pq=d7ed47ae0002dc5f"
	wd       selenium.WebDriver
	sv       *selenium.Service
	keywords = []string{"龙湖大境天成", "龙湖首开湖西星辰", "九龙仓翠樾庭", "中骏云景台", "碧桂园伴山澜湾", "阳光城檀苑", "苏州恒大悦珑湾", "鲁能公馆", "中骏天荟", "中交路劲璞玥风华", "蔚蓝四季花园", "天房美瑜兰庭", "苏悦湾", "明发江湾新城", "南山楠", "九龙仓天曦", "中铁诺德姑苏上府", "中粮天悦悦茏雅苑", "姑苏金茂悦", "恒大珺睿庭", "阳光城平江悦", "高铁新城朗诗蔚蓝广场", "金科仁恒浅棠平江", "金辉姑苏铭著", "江南沄著", "保利天樾人家", "当代蘇洲府", "东原千浔", "九龙仓天灏", "中南滨江铂郡", "中旅运河名著", "城仕高尔夫", "融信海月平江", "首开金茂熙悦", "洛克公园", "绿景苏州公馆", "新城十里锦绣", "万科公园大道", "苏州万和悦花园", "宽泰铂园", "正荣香山麓院", "泉山39", "和昌紫竹云山墅", "万泽太湖庄园", "柳岸晓风", "弘阳上熙名苑", "上海浦西玫瑰园", "绿都苏和雅集", "银城原溪"}
	//keywords = []string{"恒大林溪郡"}
)

// 获取本地未使用的端口号
func PickUnusedPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	port := l.Addr().(*net.TCPAddr).Port
	if err := l.Close(); err != nil {
		return 0, err
	}
	return port, nil
}

func StartLoopCrawler(runFlag chan bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("some error has occurred, info: ", r)
			debug.PrintStack()
			if val, ok := <- runFlag; ok {
				runFlag <- val
			} else {
				runFlag <- false
			}
		}
	}()
	StartChromeBrowser()
	log.Println("Start Crawling at ", time.Now().Format("2006-01-02 15:04:05"))
	for _, keyword := range keywords {
		StartCrawlerBDAds(keyword)
	}
	log.Println("Crawling Finished at ", time.Now().Format("2006-01-02 15:04:05"))
	defer sv.Stop() // 停止chromedriver
	defer wd.Quit() // 关闭浏览器
	runFlag <- true
}

// StartCrawler 开始爬取数据
func StartCrawlerBDAds(keyword string) {
	// 导航到目标网站
	err := wd.Get(fmt.Sprintf(urlBaidu, keyword))
	if err != nil {
		panic(fmt.Sprintf("Failed to load page: %s\n", err))
	}
	log.Println(wd.Title())
	leftContent, err := wd.FindElement(selenium.ByXPATH, "//*[@id=\"content_left\"]")
	if err != nil {
		panic(err)
	}
	lists, err := leftContent.FindElements(selenium.ByXPATH, "//div[@id>3000]")
	if err != nil {
		panic(err)
	}
	log.Println("推广总数：", len(lists))
	for i := 0; i < len(lists); i++ {
		title, err := lists[i].Text()
		log.Printf("正在浏览第%d个房源数据(%s)...\n", i+1, title)

		urlElem, err := lists[i].FindElement(selenium.ByXPATH, "div[1]/h3/a")
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		//url, err := urlElem.GetAttribute("href")
		//if err != nil {
		//	break
		//}
		//fmt.Println("url:", url)

		if domainUrlElem, err := lists[i].FindElement(selenium.ByXPATH, "div[2]/div/div[2]/div[2]/a|div[3]/a"); domainUrlElem != nil && err == nil {
			if domainText, domainErr := domainUrlElem.Text(); domainErr == nil {
				fmt.Println("domainText: ", domainText)
				if strings.Contains(domainText, "yuanhaowang") || strings.Contains(domainText, "xunkyz") || strings.Contains(domainText, "源浩网") || strings.Contains(domainText, "讯客驿站") {
					continue
				}
			}
		}

		if err = urlElem.Click(); err != nil {
			fmt.Println(err.Error())
		}

		time.Sleep(time.Second * 3)
		winHandlers, err := wd.WindowHandles()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("winHandlers: ", winHandlers)
		currentUrl, err := wd.CurrentURL()
		fmt.Printf("before currnetUrl: %s\n", currentUrl)
		if len(winHandlers) > 1 {
			// 切换到新窗口
			_ = wd.SwitchWindow(winHandlers[1])

			currentUrl, err = wd.CurrentURL()
			fmt.Printf("after currnetUrl: %s\n", currentUrl)
			// 窗口2上下滑动
			if err = wd.KeyDown(selenium.EndKey); err != nil {
				fmt.Println(err.Error())
			}
			for j := 0; j < 4; j++ {
				time.Sleep(time.Second * 1)
				if err = wd.KeyDown(selenium.PageUpKey); err != nil {
					fmt.Println(err.Error())
				}
			}

			// 关闭窗口2，回到窗口1
			if err = wd.CloseWindow(winHandlers[1]); err != nil {
				fmt.Println(err.Error())
			}
			_ = wd.SwitchWindow(winHandlers[0])
			for j := 0; j < 4; j++ {
				time.Sleep(time.Second * 1)
				if err = wd.KeyDown(selenium.PageUpKey); err != nil {
					fmt.Println(err.Error())
				}
			}
		} else {
			fmt.Println("未打开窗口2")
		}
	}
}

// StartChrome 启动谷歌浏览器headless模式
func StartChromeBrowser() {
	//selenium.SetDebug(true)
	var opts []selenium.ServiceOption
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}

	// 配置chrome浏览器基础能力
	chromeCaps := chrome.Capabilities{
		Path: "",
		Args: []string{
			//"--headless", // 设置Chrome无头模式
			"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/604.4.7 (KHTML, like Gecko) Version/11.0.2 Safari/604.4.7", // 模拟user-agent，防反爬
		},
	}
	caps.AddChrome(chromeCaps)

	// 添加代理
	chromeProxy := selenium.Proxy{
		Type: selenium.Manual,
		HTTP: "45.159.75.71",
		HTTPPort: 8080,
	}
	caps.AddProxy(chromeProxy)

	// 启动chromedriver，端口号可自定义
	var err error
	sv, err = selenium.NewChromeDriverService("C:\\Program Files (x86)\\Google\\Chrome\\Application\\chromedriver", 9515, opts...)
	if err != nil {
		log.Printf("Error starting the ChromeDriver server: %v", err)
	}
	// 调起chrome浏览器
	wd, err = selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", 9515))
	if err != nil {
		panic(err)
	}
}
