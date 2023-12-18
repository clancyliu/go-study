package script

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	gomail "gopkg.in/mail.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var client = http.Client{
	Timeout: 5 * time.Second,
	//CheckRedirect: RedirectFunc,
}

// RedirectFunc 重定向禁止
func RedirectFunc(req *http.Request, via []*http.Request) error {
	fmt.Println(req.RequestURI)
	// 如果返回 非nil 则禁止向下重定向 返回nil 则 一直向下请求 10 次 重定向
	return http.ErrUseLastResponse
}

const (
	method      = "POST"
	submitUrl   = "https://www.jmhui.com/meeting/saveEnrollFormData.jspx"
	tokenUrl    = "https://www.jmhui.com/cas/login"
	contentType = "application/x-www-form-urlencoded"
	cookieFmt   = "%s; clientlanguage=zh_CN; _meeting_id_cookie=475"
)

var cookie = "JSESSIONID=925137739F9C048266F9E3907A805883; clientlanguage=zh_CN; _meeting_id_cookie=475"

var submitData = url.Values{
	"meetingid": {"475"},
	"rid":       {"590562770"},
	"formId":    {"496"},
	"pathId":    {"499"},
	"formData": {`{
			"attr_name": "姚婷",
			"attr_13": "女",
			"attr_age": "27",
			"attr_mobile": "15270813615",
			"attr_mail": "1075618950@qq.com",
			"attr_sfz": "360681199610111722",
			"attr__data_field_1951": "大学本科",
			"attr_14": "无锡爱尔眼科医院",
			"attr__data_field_2485": "病案科主管",
			"attr__data_field_2486": "3",
			"attr__data_field_1946": "江苏",
			"attr_zw": "病案主管",
			"attr__data_field_1949": "初级",
			"attr__data_field_1939": "否",
			"attr__data_field_1952": "是",
			"attr__data_field_1953": "是",
			"attr__data_field_1954": "姚婷",
			"attr__data_field_1955": "15270813615",
			"attr__data_field_1956": "江苏省无锡市滨湖区梁溪路11-1号  无锡爱尔眼科医院",
			"attr_accom_mode": "合住"
		}`},
}

var loginData = url.Values{
	"execution": {"c2c8ae93-6f5e-41e2-97e2-367a090aab84_ZXlKaGJHY2lPaUpJVXpVeE1pSjkucFYybkxhOVU3YkJyaEpLT0g2MGtqcjRhSlZHajlmN1doNUNGaDhQcGJ4RWpFc3VHbktBZVNpbzhxV0lMOG9hL3M3TWRVVDBDMlgwNFZINDI1Nm5PZGF6cy9aL3RiTWJ2MmxhMVhabVQwNWdQZ1k3OFZ4NUg3L0M0VlRJRGRWV0c5R3d1OUZZY01qV0NraVYrVmtDNTBlME5FcW9CdHhPTGMvOTRJTHRKRVBZY282QVBUdVc1a1NrUzJxaVBkR09NU1RHbnNVdVBCN2E0Rkpsdkx4TlI3Z2d1M2s1NzYvUExjZU8rVUFyMU1xODFzMWFhN1l1U250U1NmMElYRGZ1OWlHRXhyUE1nSFNxMGcxZ2tQZU52b01NYjhUVk80ZXpwWTg1RTZSTDNBQXg3Y2tmWE9HQmhlV2hFa0haU09vV1EzdkFxcjhQd21tR2RadFpDcXZvU1hBZzByMXBueXp1YnR1aDJoeEF3bjZqSE1STWtCMW5BN1kwRytCbCtySEFvUmV1b05ZMzJlMzZRdUM5Vnd6b1BNbDhwOXZCQ3duaUo3Nm14b3doaUc3czdFQk9vc1JzNWZwVFpibkhqMVZiNE53NUM1dkJqWkY3NGpna01XZ3lTblI3dDBJN3RTUnZCUkRDTmd0Vy9YcDNKRHVNbDZqazE0aGZBdkVsV1JaODVyTXBuVnRXdjRZeFFYUDNnT0RmOEhTWDkvaDNVVDFTQ2IxMFRxSXFwQ0hCR3NsZEVZYnRLUHp5WllRaUE1TEpQVTY3cDgvUEdSUU54YTU5eGd3ZE1jMFdaL1kzZTBlekhmSXlCcU9ObnlDTDN3MjlYcVZJTGdUbmdCWnZsdEtNYjdUczJIajNTdWp2UWhqQ3JKVm94aVJmVFhJdVdHZ1NtTUpNTTdEWFo1UmZlOUtERzVCMzRWSHFlMzZxMFlsQUNkcG5zWlk3aFN1dlk3ak1rNDA5VjJqUnBKVDg0b3RLa2p4bnFWODhBUUhuQnhGN3ovOEo0cHhUK1JFNjNiWGVZM0RsSHFINDdic1Nya3V1V2lYMWdnZW92bFZseS8xdElsWFJ3RE5XalhnTFM2dVk0cTZwQTdIbWg2cWtzYUN3Y2FpbjhDc0RRMXZMUnY0TW44YlNzdXBxM3FaYjRXOFJNdnZ5dnMrNncvL1ZqOTRVa3NyV0lxLy9zYVhyNUJWUEtPQmxaTk9LTk5QalNkK3Z1aVZSeVE3RVRxaDNKRnN5ZkUwR2VKMndidmdmVVVCM2s0Y0o1NzRxREgrRUxISXVuL3JvSTRJYUFsWllXbStySkl1MHhVVVlyMTJDV3hpZEE2ZzRpS0pwdHcxY3IzZlRZMDY5Zkl0SEdCV2JXTGVCK1R2WmxqYko4SWtyRnpJdWRFVmxjRFRKRGZ3T05PQkttQVVoVzlsbGF5Y2xtbWd3eUx0K3lqdFN5YzVVbnVLeUlod3drUW1HMXk5NTIyWHRWbDE3a2hvT2R4ak81TlZTRWkxZE1XUkJCaU45OWwwYUl1NnVlcC9kbUtQUjlWWlpSNGZxZXdBWFo3a1ZUTmNGR2tzOTNuZ1R4TDA1VlQ5aGsxY2xWaGpuTG93c215UWJpK2s5QklkZGpmRThWdGV2cklqRE9OMVI5eE9OOGNkUmQvVlJzNUFvTW5Nd2htUXV3U1kzaWtBSEFJS2ZPSDUvL0RscnZnKzJwclF5eFJRR01UNDNtYkhZV2Zhd0ZScUQrcnBFRmcyS2FCNFpEMjFRUG5CdjVQbWxqTTJUekFSY1hQT1pXYk5mdTY3N2VFRGpoblA5V3pPZGtFdW9xMUh6THp3RTQ2ZWRyY0tiQWtFYWdQVFhFVFlSSVVjcjZ3NXU5dFFiLzBCazdoWktUdmdwUUQ3Y3BPYkx5d2JSM005TXVOZWdGYUp1WlIzRGxwcVJvUy9UMHpxZXRlZnM4ZUhBQUgyYUFBLzVrVC9tWmpHUjZEZTUyMmZyR0pqME5YMk9jeHkxaHNnU1FDWkJvbTU0Tk5LKzRzSFRnMHFpcllNak5wYUhhdDY4cnphYW1ENy9LSnd1L05oRnR0bVN1TWJJUFBOa0R1aHNGMmRiZCtNdmJSaE1MZUZ6ZllaVGVPQzFpYjZWS2pqdEE1QnlDYUw1TlJ1ZVVOb3NNU20yMWZ6SFJHaW42R3I5SUIrWEU2TGpDUVFJOXIvM3c5bEdmZXlrR0MrU0JwcFo5c0RLQlQwVmg5ald6bXJjNVk5aGt1QUlFeDM4eFlOaXh5RFNzTHV6VTBzN3kyb3FUamFyZUZFZnptUnJwcnVwSHlDbU1JY2t5bXMzclA2cmhKSExkZ2VUMW9sZTBFYUtBTmowU0VadGpCWm8zSjR6NU91MC9uN0VoZk95Q1RYLzJPbWhSQ3p4RWNUcUJJSlhJUkxTM1BuZXl1M0cxSGNXc3RxNldzWFJGdEhia216QXFENVJSYWV6MkwxeU1hb00xMVRaNTA0SFRzc2pFUDNRLzRUWFA2RDUvMUFxWWE1YWdnMjZSNjl4ejF4WDhnamJlN2F2V2V6RkZkd3lmVi9GYk51a1Eza2xzK2Q1Z3lDREJMckpGWGc2QkZtY015SnpUNmVlM3g4aGJtTUdUUmI4cHpudzlDZWFpcjkycnZRa1pxVmN0dERwZnNUVmg3eTRWYTFPZEsvam1LYy9Cc2tZUDlJTEZMQzVSaXNsY2V6ZEc4K1Q5U2I0R05KdG1rVVRObElsblZMVndBMmEvNGJnaDY1TmhvbEMvRkZ6bVg4RlVoZVF4UGFKYlZkcW1TN3RESWxsQ3Z6Q0VnTXNMaGExb2lhY2hyK2hmVkJKcGpmUGl0OThlMlpJRGxHTVNnckNYSmhoVG54QSsxTHQyL2pvUkRFUE1uaWl1cWtWLzNNL2NoR2lodWozS1hXL0hSK240MjBlSy9jalFYY0pXUG8xUFZKRityWkhId1liVHQ1OFlUVVJMOWlkQTE3OUVYQ0VxU1NRQ0trb3NrZGl2SGhVU0hlbitHbGpZci9uWktDTUxyRmlCNStIVVkwZW01UjRKb216REs0WE9hNXB5KzNJaWxQcVd0QTk3MGJPck9GOTZGT3BXRFFnSUZ0bEw2YWtnRkR3RnRQMXU4ZkQ3SHNwWnBQUGZDSk1XN1MzUUcyZkN6ZVRIcjVzTUwrRFNveVhJMUE5T2RzMWw1a2lGNkRUMHl0NXB3N2lMaERoK1FlTE1xajlydVVYSXBZYU1zYVlmd3p3Q2hWQVczOStKeFVwNmFQNGNqQlhtTFR0V2MrMnNScXV3V1VLc2xLVmhTcnFGempTdENPK2g5VlRVcys2Mkc3VDJibzBBNUd0cWNuTmhWamtjUmRxa3VSY08vajE1bDNRMzJBaDZIK016cmIzcEJSSmJGY2ZkNTRXeDBsc3UvRlRRQnBFSU1MQ1ovdXRJU0xyVUdjVDFFOXBodlk0ejVsVlVCUXc0L3RNT0xTcnYwdVp1d05EVm5BQi92U3lhVzRocGZzV1llMVdRN1ZaTnMrOVRqYzZLV0Rwbm5Db3c1SDluTHJVK1ZlNXZIZnZrRHFTbys0Vm1IQzU1a0N3RlVya2t2NXU0bXFRSDZDNk13aWlUTzZ3NjdILy9Vdk5ZVHNtamttWk8xTlpnUUpLb1JoeUlheWg4Y1EzSFMvKzBWWEZxQzJJQ1hReFFDUEhxUGZkd2dEWUtBSmtPVk5HNmZHdFcvR2R3UW40RGxJNmZwRlp0R3pJb2dVb0RDUFFPeUtBbTJ2aUZ3UzMzRjZSNDVzQ1JMZFdmQkkyUXJWQnBjcVd1NUxJeXNqM0d6bjA4NFZleW5IdXE5K09qZlBQZGYrNkY0RmVta2oyY3ZHMnh0Z1RhOCtVbDBpQ3o0Z2lweThGQ0J1MzdYSkFnZXdCTnZ6UnpvTzFwdTEvVGxKKy9MKzMyWitieXpZNUhDOE12UmRUU1VvNWIrTkp4dTZCQ2Vmd05mQklHQnNsNUJOL1BTVU9aYlYwTFZrVUdmcU1MUVIrdFk4bWEuTUhnWml3Z3puSW1LOEJtVHNHbUlMNW53VExMbEgzdGJxUm9YQlJKWVI4TUFnb0xCT0o5ZmFOVlFTUjMwY3VwVHZPQjUyel9pb0tUNGxVS0dCRHNIckE="},
	"_eventId":  {"submit"},
	"username":  {"15270813615"},
	"password":  {"yt123456"},
	"type":      {"normal"},
}

var loginErr = errors.New("please login")

func main() {

	file, err := os.OpenFile("log.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.SetOutput(file)
	for {
		submit()
		time.Sleep(5 * time.Second)
	}
}

func doPostFormData(urlStr string, data url.Values, header http.Header) error {
	reqBody := data.Encode()
	request, err := http.NewRequest(method, urlStr, strings.NewReader(reqBody))

	if err != nil {
		log.Println(err)
		return err
	}

	request.Header = header

	response, err := client.Do(request)

	if err != nil {
		log.Println(err)
		return err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	respBodyStr := string(body)

	if strings.Contains(respBodyStr, "登录") {
		log.Println("cookie expire....")
		return loginErr
	}
	log.Println("respBody: ", respBodyStr)

	var respJson respStruct
	err = json.Unmarshal(body, &respJson)
	if err != nil {
		log.Println(err)
		return err
	}
	if respJson.Code != -1 {
		for i := 0; i < 10; i++ {
			if err := sendEmail(); err != nil {
				log.Println(err)
			}

		}
	}

	return nil
}

type respStruct struct {
	Code int
	Msg  string
}

func submit() {
	header := http.Header{
		"Content-Type": {contentType},
		"Cookie":       {cookie},
	}

	err := doPostFormData(submitUrl, submitData, header)
	if err != nil {
		log.Println(err)
		if errors.Is(err, loginErr) {
			for i := 0; i < 3; i++ {
				if login() == nil {
					break
				}
			}
		}
	}
}

func login() error {
	header := http.Header{
		"User-Agent":      {"PostmanRuntime/7.31.3"},
		"Accept":          {"*/*"},
		"Postman-Token":   {"8577c458-0dad-4f44-818c-7cc8ad21214b"},
		"Host":            {"www.jmhui.com"},
		"Accept-Encoding": {"gzip, deflate, br"},
		"Connection":      {"keep-alive"},
		//"Content-Type":    {"multipart/form-data; boundary=--------------------------862083641246614199800340"},
		"Content-Type": {contentType},
		//"Content-Length": {"4812"},
	}

	reqBody := loginData.Encode()
	request, err := http.NewRequest(method, tokenUrl, strings.NewReader(reqBody))

	if err != nil {
		log.Println(err)
		return err
	}

	request.Header = header

	tmpClient := http.Client{
		Timeout:       5 * time.Second,
		CheckRedirect: RedirectFunc,
	}

	response, err := tmpClient.Do(request)

	if err != nil {
		log.Println(err)
		return err
	}

	defer response.Body.Close()

	location := response.Header.Get("Location")

	resp, err := tmpClient.Get(location)
	if err != nil {
		log.Println(err)
		return err
	}

	setCookie := resp.Header.Get("Set-Cookie")

	split := strings.Split(setCookie, ";")
	cookie = fmt.Sprintf(cookieFmt, split[0])

	fmt.Println(cookie)

	return nil
}

func sendEmail() error {
	// Sender data.
	from := "clancy_liu@163.com"
	password := "IHCEGOGTYWIHNZCP"

	// Receiver email address.
	to := []string{"clancy_liu@qq.com", "996086041@qq.com"}

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", from)

	// Set E-Mail receivers
	m.SetHeader("To", to...)

	// Set E-Mail subject
	m.SetHeader("Subject", "报名成功！！！！！")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", "https://www.jmhui.com/event/202301-yxh")

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.163.com", 465, from, password)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
