package ssq

import (
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	State    int      `json:"state"`
	Message  string   `json:"message"`
	Total    int      `json:"total"`
	PageNum  int      `json:"pageNum"`
	PageNo   int      `json:"pageNo"`
	PageSize int      `json:"pageSize"`
	Tflag    int      `json:"Tflag"`
	Result   []Detail `json:"result"`
}
type Detail struct {
	Name        string        `json:"name"`
	Code        string        `json:"code"`
	DetailsLink string        `json:"detailsLink"`
	VideoLink   string        `json:"videoLink"`
	Date        string        `json:"date"`
	Week        string        `json:"week"`
	Red         string        `json:"red"`
	Blue        string        `json:"blue"`
	Blue2       string        `json:"blue_2"`
	Sales       string        `json:"sales"`
	Poolmoney   string        `json:"poolmoney"`
	Content     string        `json:"content"`
	Addmoney    string        `json:"addmoney"`
	Addmoney2   string        `json:"addmoney2"`
	Msg         string        `json:"msg"`
	Z2add       string        `json:"z2add"`
	M2add       string        `json:"m2add"`
	Prizegrades []Prizegrades `json:"prizegrades"`
}
type Prizegrades struct {
	Type      int    `json:"type"`
	Typenum   string `json:"typenum"`
	Typemoney string `json:"typemoney"`
}
type Ssq struct {
	Id          string
	Date        string
	Week        string
	Red_numbers string
	Red_number1 string
	Red_number2 string
	Red_number3 string
	Red_number4 string
	Red_number5 string
	Red_number6 string
	Blue        string
	Sales       string
	Pool_amount string
	Prizegrades []Prize
}
type Prize struct {
	code          string
	number        int
	people_number string
	money         string
}

func Crawler() {
	baseUrl := "https://www.cwl.gov.cn/cwl_admin/front/cwlkj/search/kjxx/findDrawNotice?name=ssq&issueCount=&issueStart=&issueEnd=&dayStart=&dayEnd=&pageNo=@@&pageSize=30&week=&systemType=PC"
	for i := 1; i < 58; i++ {
		url := strings.Replace(baseUrl, "@@", strconv.Itoa(i), 1)
		ssqarray := getAndAnalysis(url)
		log.Println("解析第", i, "页数据 ", len(ssqarray), "条")
		InsertDatas(ssqarray)
		time.Sleep(time.Duration(rand.Int63n(10)) * time.Millisecond)
	}
	CloseDB()

}
func getAndAnalysis(url string) []Ssq {
	log.Println(url)
	client := http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	request.Header.Add("Cookie", "HMF_CI=2eed9bd84ac2226580165f747d6a2306cabd4d428cd359076b2b5ec52e39d1345b03fc20a8fd808c1c552e7953105026b20da8ffb4febfb8a23474551af4048e6b; 21_vq=5")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	request.Header.Add("Host", "www.cwl.gov.cn")
	request.Header.Add("Sec-Ch-Ua-Platform", "Windows")
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	if resp.StatusCode == 200 {
		// 正确
		return parserSsqJson(resp.Body)
	} else {
		//异常
		log.Fatalln("response code", resp.StatusCode)
	}
	return nil
}

func parserSsqJson(body io.ReadCloser) []Ssq {
	defer body.Close()
	bytes, _ := io.ReadAll(body)
	response := Response{}
	err := json.Unmarshal(bytes, &response)
	if err != nil {
		log.Fatalln(err)
	}
	var ssqs []Ssq
	if "查询成功" == response.Message {
		details := response.Result
		for _, v := range details {
			ssq := Ssq{}
			ssq.Id = v.Code
			ssq.Week = v.Week
			index := strings.Index(v.Date, "(")
			ssq.Date = v.Date[0:index]
			ssq.Red_numbers = v.Red
			ssq.Blue = v.Blue
			split := strings.Split(v.Red, ",")
			ssq.Red_number1 = split[0]
			ssq.Red_number2 = split[1]
			ssq.Red_number3 = split[2]
			ssq.Red_number4 = split[3]
			ssq.Red_number5 = split[4]
			ssq.Red_number6 = split[5]
			ssq.Sales = v.Sales
			ssq.Pool_amount = v.Poolmoney
			prizes := []Prize{}
			for _, p := range v.Prizegrades {
				if p.Typenum != "" {
					prize := Prize{code: v.Code, number: p.Type, people_number: p.Typenum, money: p.Typemoney}
					prizes = append(prizes, prize)
				}

			}
			ssq.Prizegrades = prizes
			ssqs = append(ssqs, ssq)
		}
	}
	return ssqs
}
