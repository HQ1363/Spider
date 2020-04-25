package selenium

import (
	"encoding/csv"
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"strings"
	"time"
)

var (
	urlBeijing = "https://www.che168.com/beijing/list/#pvareaid=104646"
	webDriver  selenium.WebDriver
	dateTime   string
	writer     *csv.Writer
	service    *selenium.Service
	csvFile    *os.File
)

func sendResult(fileName string) {
	email := gomail.NewMessage()
	email.SetAddressHeader("From", "re**ng@163.com", "张**")
	email.SetHeader("To", email.FormatAddress("li**yang@163.com", "李**"))
	email.SetHeader("Cc", email.FormatAddress("zhang**tao@163.net", "张**"))
	email.SetHeader("Subject", "二手车之家-北京-二手车信息")
	email.SetBody("text/plain;charset=UTF-8", "本周抓取到的二手车信息数据，请注意查收！\n")
	email.Attach(fmt.Sprintf("data/%s.csv", fileName))

	dialer := &gomail.Dialer{
		Host:     "smtp.163.com",
		Port:     25,
		Username: "${your_email}",    // 替换自己的邮箱地址
		Password: "${smtp_password}", // 自定义smtp服务器密码
		SSL:      false,
	}
	if err := dialer.DialAndSend(email); err != nil {
		log.Println("邮件发送失败！err: ", err)
		return
	}
	log.Println("邮件发送成功！")
}

// StartCrawler 开始爬取数据
func StartCrawler() {
	log.Println("Start Crawling at ", time.Now().Format("2006-01-02 15:04:05"))
	pageIndex := 0

	for {
		listContainer, err := webDriver.FindElement(selenium.ByXPATH, "//*[@class=\"viewlist_ul\"]")
		if err != nil {
			panic(err)
		}
		lists, err := listContainer.FindElements(selenium.ByClassName, "carinfo")
		if err != nil {
			panic(err)
		}
		log.Println("数据量：", len(lists))
		pageIndex++
		log.Printf("正在抓取第%d页数据...\n", pageIndex)
		for i := 0; i < len(lists); i++ {
			var urlElem selenium.WebElement
			if pageIndex == 1 {
				urlElem, err = webDriver.FindElement(selenium.ByXPATH, fmt.Sprintf("//*[@class='viewlist_ul']/li[%d]/a", i+13))
			} else {
				urlElem, err = webDriver.FindElement(selenium.ByXPATH, fmt.Sprintf("//*[@class='viewlist_ul']/li[%d]/a", i+1))
			}
			if err != nil {
				break
			}
			url, err := urlElem.GetAttribute("href")
			if err != nil {
				break
			}
			if err := webDriver.Get(url); err != nil {
				log.Println(err.Error())
			}
			title, _ := webDriver.Title()
			log.Printf("当前页面标题：%s\n", title)

			modelElem, err := webDriver.FindElement(selenium.ByXPATH, "/html/body/div[5]/div[2]/div[1]/h3")
			var model string
			if err != nil {
				log.Println(err)
				model = "暂无"
			} else {
				model, _ = modelElem.Text()
			}
			log.Printf("model=[%s]\n", model)

			priceElem, err := webDriver.FindElement(selenium.ByXPATH, "/html/body/div[5]/div[2]/div[2]/span")
			var price string
			if err != nil {
				log.Println(err)
				price = "暂无"
			} else {
				price, _ = priceElem.Text()
				price = fmt.Sprintf("%s万", price)
			}
			log.Printf("price=[%s]\n", price)

			milesElem, err := webDriver.FindElement(selenium.ByXPATH, "/html/body/div[5]/div[2]/ul/li[1]/h4")
			var miles string
			if err != nil {
				log.Println(err)
				milesElem, err := webDriver.FindElement(selenium.ByXPATH, "/html/body/div[5]/div[2]/ul/li[1]/h4")
				if err != nil {
					log.Println(err)
					miles = "暂无"
				} else {
					miles, _ = milesElem.Text()
				}
			} else {
				miles, _ = milesElem.Text()
			}
			log.Printf("miles=[%s]\n", miles)

			timeElem, err := webDriver.FindElement(selenium.ByXPATH, "/html/body/div[5]/div[2]/ul/li[2]/h4")
			var date string
			if err != nil {
				log.Println(err)
				timeElem, err := webDriver.FindElement(selenium.ByXPATH, "/html/body/div[5]/div[2]/ul/li[2]/h4")
				if err != nil {
					log.Println(err)
					date = "暂无"
				} else {
					date, _ = timeElem.Text()
				}
			} else {
				date, _ = timeElem.Text()
			}
			log.Printf("time=[%s]\n", date)

			positionElem, err := webDriver.FindElement(selenium.ByXPATH, "/html/body/div[5]/div[2]/ul/li[4]/h4")
			var position string
			if err != nil {
				log.Println(err)
				positionElem, err := webDriver.FindElement(selenium.ByXPATH, "/html/body/div[5]/div[2]/ul/li[4]/h4")
				if err != nil {
					log.Println(err)
					position = "暂无"
				} else {
					position, _ = positionElem.Text()
				}
			} else {
				position, _ = positionElem.Text()
			}
			log.Printf("position=[%s]\n", position)

			storeElem, err := webDriver.FindElement(selenium.ByXPATH, "/html/body/div[5]/div[2]/div[1]/div/div/div")
			var store string
			if err != nil {
				log.Println(err)
				store = "暂无"
			} else {
				store, _ = storeElem.Text()
				store = strings.Replace(store, "商家|", "", -1)
				if strings.Contains(store, "金牌店铺") {
					store = strings.Replace(store, "金牌店铺", "", -1)
				}
			}
			log.Printf("store=[%s]\n", store)
			_ = writer.Write([]string{model, miles, date, price, position, store})
			writer.Flush()
			_ = webDriver.Back()
		}
		log.Printf("第%d页数据已经抓取完毕，开始下一页...\n", pageIndex)
		nextButton, err := webDriver.FindElement(selenium.ByClassName, "page-item-next")
		if err != nil {
			log.Println("所有数据抓取完毕！")
			break
		}
		_ = nextButton.Click()
	}
	log.Println("Crawling Finished at ", time.Now().Format("2006-01-02 15:04:05"))
	sendResult(dateTime)

	defer service.Stop()    // 停止chromedriver
	defer webDriver.Quit()  // 关闭浏览器
	defer csvFile.Close()   // 关闭文件流
}

// SetupWriter 初始化CSV
func SetupWriter() {
	dateTime = time.Now().Format("2006-01-02") // 格式字符串是固定的，据说是go语言诞生时间，谷歌的恶趣味...
	_ = os.Mkdir("data", os.ModePerm)
	csvFile, err := os.Create(fmt.Sprintf("data/%s.csv", dateTime))
	if err != nil {
		panic(err)
	}
	_, _ = csvFile.WriteString("\xEF\xBB\xBF")
	writer = csv.NewWriter(csvFile)
	_ = writer.Write([]string{"车型", "行驶里程", "首次上牌", "价格", "所在地", "门店"})
}

// StartChrome 启动谷歌浏览器headless模式
func StartChrome() {
	var opts []selenium.ServiceOption
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}

	// 禁止加载图片，加快渲染速度
	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}

	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Path:  "",
		Args: []string{
			//"--headless", // 设置Chrome无头模式
			"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/604.4.7 (KHTML, like Gecko) Version/11.0.2 Safari/604.4.7", // 模拟user-agent，防反爬
		},
	}
	caps.AddChrome(chromeCaps)
	// 启动chromedriver，端口号可自定义
	_, err := selenium.NewChromeDriverService("C:\\Program Files (x86)\\Google\\Chrome\\Application\\chromedriver", 9515, opts...)
	if err != nil {
		log.Printf("Error starting the ChromeDriver server: %v", err)
	}
	// 调起chrome浏览器
	webDriver, err = selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", 9515))
	if err != nil {
		panic(err)
	}
	// 这是目标网站留下的坑，不加这个在linux系统中会显示手机网页，每个网站的策略不一样，需要区别处理。
	_ = webDriver.AddCookie(&selenium.Cookie{
		Name:  "defaultJumpDomain",
		Value: "www",
	})
	// 导航到目标网站
	err = webDriver.Get(urlBeijing)
	if err != nil {
		panic(fmt.Sprintf("Failed to load page: %s\n", err))
	}
	log.Println(webDriver.Title())
}
