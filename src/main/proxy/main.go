package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/panjf2000/ants"
	"golang.org/x/net/proxy"
	"h12.io/socks"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var mu sync.Mutex
var resultList []string
var successProxy int

// 不同验证节点，不同的截取信息
func responseOriginIp(proxyStr string, resp http.Response) interface{} {
	requestUrl := resp.Request.URL.String()
	switch requestUrl {
	case httpbinAddr:
		body, _ := ioutil.ReadAll(resp.Body)
		var resultMap map[string]string
		_ = json.Unmarshal(body, &resultMap)
		originProxy := resultMap["origin"]
		return originProxy
	case bilibiliAddr:
		if resp.StatusCode == 200 {
			return true
		}
	case ipApiAddr:
		body, _ := ioutil.ReadAll(resp.Body)
		var resultMap map[string]string
		_ = json.Unmarshal(body, &resultMap)
		country := resultMap["country"]
		logger.Println(proxyStr, string(body))
		return country == systemCountry
	case azenvAddr:
		body, _ := ioutil.ReadAll(resp.Body)
		bodyStr := string(body)
		remoteAddrRg := regexp.MustCompile("REMOTE_ADDR = ((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})(\\.((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})){3}")
		remoteAddrRgStr := remoteAddrRg.FindAllStringSubmatch(bodyStr, -1)
		if strings.Contains(bodyStr, "reached") || len(remoteAddrRgStr) == 0 || len(remoteAddrRgStr[0]) == 0 {
			return ""
		}
		remoteAddr := strings.Replace(remoteAddrRgStr[0][0], "REMOTE_ADDR = ", "", -1)
		return remoteAddr
	}
	if strings.Contains(requestUrl, ipApiAddr) {
		body, _ := ioutil.ReadAll(resp.Body)
		var resultMap map[string]string
		_ = json.Unmarshal(body, &resultMap)
		country := resultMap["country"]
		return country == systemCountry
	}

	return ""
}

// 发送http代理请求。
func checkHttpProtocol(protocol string, checkNode string, proxyList []string, coroutineNum int, outFile string, checkTimeout int) {
	wg := sync.WaitGroup{}

	pool, _ := ants.NewPoolWithFunc(coroutineNum, func(i interface{}) {
		defer func() {
			wg.Done()
		}()
		proxyStr := i.(string)
		urli := url.URL{}
		urlProxy, err := urli.Parse(protocol + "://" + proxyStr)
		// todo
		if err != nil {
			return
		}
		c := http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(urlProxy),
			},
			Timeout: time.Second * time.Duration(checkTimeout),
		}
		tmpCheckNode := checkNode
		if tmpCheckNode == ipApiAddr {
			ipAndPort := strings.Split(proxyStr, ":")
			tmpCheckNode += ipAndPort[0]
		}
		startRequestTime := time.Now()
		req, err := http.NewRequest("GET", tmpCheckNode, nil)
		// todo
		if err != nil {
			//log.Println(proxyStr, err)
			return
		}
		req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.115 Safari/537.36")
		req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")

		resp, err := c.Do(req)
		// todo
		if err != nil {
			//log.Println(proxyStr, err)
			return
		}
		defer resp.Body.Close()
		c.CloseIdleConnections()
		handleResult(protocol, startRequestTime, proxyStr, resp)
	})
	defer pool.Release()

	for _, proxy := range proxyList {
		wg.Add(1)
		pool.Invoke(string(proxy))
	}

	wg.Wait()
}

// 判断返回信息。
func handleResult(protocol string, startTime time.Time, proxyStr string, resp *http.Response) {
	var success bool
	pxyResp := responseOriginIp(proxyStr, *resp)
	requestAddr, ok := pxyResp.(string)
	if ok && requestAddr != "" {
		requestAddrs := strings.Split(requestAddr, ",")
		// 只要是ipv4/6 就表示代理成功
		if net.ParseIP(requestAddrs[0]) != nil {
			logStr := "Success: " + proxyStr + " [" + cutTime(time.Now().Sub(startTime).String()) + "]"
			logger.Printf("%s", logStr)
			mu.Lock()
			resultList = append(resultList, proxyStr)
			successProxy++
			mu.Unlock()
			success = true
		}
	}
	resb, ok := pxyResp.(bool)
	if ok {
		if resb {
			logStr := "Success: " + proxyStr + " [" + cutTime(time.Now().Sub(startTime).String()) + "]"
			logger.Printf("%s", logStr)
			mu.Lock()
			resultList = append(resultList, proxyStr)
			successProxy++
			mu.Unlock()
			success = true
		}
	}
	if !success {
		return
	}
	r, _ := timeoutClient.Get(fmt.Sprintf(`http://pool.88o.me/tools.php?type=%s&proxy=%s`, protocol, proxyStr))
	if r != nil {
		r.Body.Close()
	}
}

// 时间戳格式
func cutTime(s string) string {
	arr := strings.Split(s, ".")
	if len(arr) != 2 {
		return s
	}
	if strings.HasSuffix(arr[1], "ms") {
		x := strings.TrimSuffix(arr[1], "ms")
		if len(x) >= 6 {
			x = x[:6]
			return fmt.Sprintf("%s.%sms", arr[0], x)
		}
	}
	if strings.HasSuffix(arr[1], "s") {
		x := strings.TrimSuffix(arr[1], "s")
		if len(x) >= 6 {
			x = x[:6]
			return fmt.Sprintf("%s.%ss", arr[0], x)
		}
	}
	return s
}

var timeoutClient = &http.Client{
	Timeout: 5 * time.Second,
}

// 发送socks5代理请求
func checkSocketProtocol(checkNode string, proxyList []string, coroutineNum int, outFile string, checkTimeout int) {
	wg := sync.WaitGroup{}

	pool, _ := ants.NewPoolWithFunc(coroutineNum, func(i interface{}) {
		defer func() {
			wg.Done()
		}()
		proxyStr := i.(string)
		dialer, err := proxy.SOCKS5("tcp", proxyStr, nil, &net.Dialer{
			Timeout:   time.Second * time.Duration(checkTimeout),
			KeepAlive: time.Second * time.Duration(checkTimeout),
		})
		if err != nil {
			fmt.Printf("can't connect to the proxy: %s \n", err.Error())
		}
		httpTransport := &http.Transport{
			Proxy: nil,
			Dial:  dialer.Dial,
		}
		httpClient := &http.Client{
			Transport: httpTransport,
			Timeout:   time.Second * time.Duration(checkTimeout),
		}
		tmpCheckNode := checkNode
		if checkNode == ipApiAddr {
			ipAndPort := strings.Split(proxyStr, ":")
			tmpCheckNode += ipAndPort[0]
		}
		startRequestTime := time.Now()

		req, err := http.NewRequest("GET", tmpCheckNode, nil)
		// todo
		if err != nil {
			return
		}
		req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.115 Safari/537.36")
		req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")

		resp, err := httpClient.Do(req)
		// todo
		if err != nil {
			return
		}
		defer resp.Body.Close()
		httpClient.CloseIdleConnections()

		handleResult("socks5", startRequestTime, proxyStr, resp)
	})

	defer pool.Release()

	for _, proxy := range proxyList {
		wg.Add(1)
		pool.Invoke(string(proxy))
	}

	wg.Wait()
}

// 发送socks4代理请求
func checkSocket4Protocol(checkNode string, proxyList []string, coroutineNum int, outFile string, checkTimeout int) {
	wg := sync.WaitGroup{}

	pool, _ := ants.NewPoolWithFunc(coroutineNum, func(i interface{}) {
		defer func() {
			wg.Done()
		}()
		proxyStr := i.(string)
		dialSocksProxy := socks.Dial("socks4://" + proxyStr)
		httpTransport := &http.Transport{
			Proxy: nil,
			Dial:  dialSocksProxy,
		}

		httpClient := &http.Client{
			Transport: httpTransport,
			Timeout:   time.Second * time.Duration(checkTimeout),
		}
		tmpCheckNode := checkNode
		if checkNode == ipApiAddr {
			ipAndPort := strings.Split(proxyStr, ":")
			tmpCheckNode += ipAndPort[0]
		}
		startRequestTime := time.Now()
		req, err := http.NewRequest("GET", tmpCheckNode, nil)
		// todo
		if err != nil {
			return
		}
		req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.115 Safari/537.36")
		req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")

		resp, err := httpClient.Do(req)
		// todo
		if err != nil {
			return
		}
		httpClient.CloseIdleConnections()
		defer resp.Body.Close()
		handleResult("socks4", startRequestTime, proxyStr, resp)
	})

	defer pool.Release()

	for _, proxy := range proxyList {
		wg.Add(1)
		pool.Invoke(string(proxy))
	}

	wg.Wait()
}

// 读取文件并去除重复
func readInputFile(filePath string) ([]string, int64, int64, error) {
	totalNoRepeat := int64(0)
	buf, err := os.OpenFile(filePath, os.O_RDONLY, 0777)
	if err != nil {
		logger.Fatal("Input file does not exist...", err)
		return nil, 0, 0, err
	}
	bio := bufio.NewReader(buf)

	var resultProxyList []string
	var resultProxyMap = make(map[string]struct{})
	var repeat = 0
	for {
		data, _, err := bio.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			logger.Fatal(err)
			return nil, 0, 0, err
		}
		if strings.TrimSpace(string(data)) == "" {
			continue
		}
		totalNoRepeat++
		if _, ok := resultProxyMap[string(data)]; ok {
			repeat++
			continue
		}
		resultProxyMap[string(data)] = struct{}{}
		resultProxyList = append(resultProxyList, string(data))
	}

	return resultProxyList, totalNoRepeat, int64(repeat), nil

}

func RemoveRepeatedElementAndEmpty(arr []string) []string {
	newArr := make([]string, 0)
	for _, item := range arr {
		if "" == strings.TrimSpace(item) {
			continue
		}
		repeat := false
		if len(newArr) > 0 {
			for _, v := range newArr {
				if v == item {
					repeat = true
					break
				}
			}
		}
		if repeat {
			continue
		}
		newArr = append(newArr, item)
	}
	return newArr
}

// 验证节点对应的网站
const (
	bilibiliAddr = `http://139.162131.84/`
	httpbinAddr  = "http://httpbin.org/ip"
	azenvAddr    = "http://azenv.net/"
	ipApiAddr    = `http://ip-api.com/json/`
)

var (
	systemCountry = "World"
)

var logger = log.New(os.Stdout, "", log.Lmsgprefix)

func main() {
	if len(os.Args) < 7 {
		logger.Println("invalid number of args")
		return
	}
	startTimestamp := time.Now()
	//命令行参数
	node := os.Args[1]
	protocol := os.Args[2]
	input := os.Args[3]
	output := os.Args[4]
	timeout, err := strconv.Atoi(os.Args[5])
	if err != nil {
		logger.Println("timeout parameter is invalid...")
		return
	}
	coroutine, err := strconv.Atoi(os.Args[6])
	if err != nil {
		logger.Println("coroutine parameter is invalid...")
		return
	}
	if len(os.Args) >= 8 {
		countryParams := strings.Split(os.Args[7], "=")
		if len(countryParams) < 2 {
			logger.Println("systemCountry parameter is invalid...")
			return
		}

		systemCountry = countryParams[1]
	}

	file, err := os.OpenFile("sys.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Println("Faild to open error logger file:", err)
		return
	}
	log.SetOutput(file)

	//验证节点
	var checkNodeUrl string
	switch node {
	case "azenv":
		checkNodeUrl = azenvAddr
	case "httpbin":
		checkNodeUrl = httpbinAddr
	case "bilibili":
		checkNodeUrl = bilibiliAddr
	case "ipapi":
		checkNodeUrl = ipApiAddr
	default:
		fmt.Println("CheckNode parameter is invalid...")
		return
	}

	proxyList, totalNoRepeat, repeatProxyLen, err := readInputFile(input)
	totalProxyLen := len(proxyList)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	switch protocol {
	case "http", "https":
		checkHttpProtocol(protocol, checkNodeUrl, proxyList, coroutine, output, timeout)
	case "socks5":
		checkSocketProtocol(checkNodeUrl, proxyList, coroutine, output, timeout)
	case "socks4":
		checkSocket4Protocol(checkNodeUrl, proxyList, coroutine, output, timeout)
	default:
		fmt.Println("CheckProtocol parameter is invalid...")
		return
	}

	end := time.Now().Sub(startTimestamp)
	CurrentTime := time.Now().In(time.UTC).Add(8 * time.Hour).Format(`2006-01-02 15:04:05`)

	//打印结果信息。
	tailStr := fmt.Sprintf(
		"\u001B[1;31mProxy Check\u001B[0m\n"+
			"\u001B[1;31m=======================\u001B[0m\n"+
			"Check Addr: %s\n"+
			"Check Country: %s\n"+
			"Check Timeout: %d\n"+
			"Check Protocol: %s\n"+
			"Check Coroutine: %d\n"+
			"\u001B[32mTotal Proxy: %d\u001B[0m\n"+
			"\u001B[33mRepeat Proxy: %d\u001B[0m\n"+
			"\u001B[34mNotrepeat Proxy: %d\u001B[0m\n"+
			"\u001B[35mSuccess Proxy: %d\u001B[0m\n"+
			"\u001B[36mConsuming Time: %s\u001B[0m\n"+
			"\u001b[36mCurrent Time: %s\u001B[0m\n"+
			"\u001B[1;31m=======================\u001B[0m", node, systemCountry, timeout, protocol, coroutine, totalNoRepeat, repeatProxyLen, totalProxyLen, successProxy, cutTime(end.String()), CurrentTime)
	logger.Printf("%s", tailStr)
	resultStr := strings.Join(resultList, "\n")
	ioutil.WriteFile(output, []byte(resultStr), fs.ModePerm)
}
