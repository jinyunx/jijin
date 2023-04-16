package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var code = "005827"
var fileName = code + ".txt"

func main() {
	var page string
	for i := 1; ; i++ {
		c, err := handleOnePage(code, i)
		if err != nil {
			break
		}
		page += c
	}
	_ = os.WriteFile(fileName, []byte(page), 0644)
}

func QueryJijin(code string, page int) (content string, err error) {

	host := "https://fundf10.eastmoney.com/F10DataApi.aspx"
	url := fmt.Sprintf("%s?type=lsjz&code=%s&per=20&page=%d", host, code, page)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		fmt.Println("query fail", err.Error())
		return
	}
	defer resp.Body.Close()

	txt, _ := ioutil.ReadAll(resp.Body)
	content = string(txt)
	if strings.Contains(content, "暂无数据") {
		return content, fmt.Errorf("end")
	}
	return content, nil
}

func handleHtml(n *html.Node) (result string) {
	// 第一层 html
	c := n.FirstChild
	for ; c != nil; c = c.NextSibling {
		//fmt.Println("->", c.Type, c.Data)
	}
	// 第二层 head/body
	c = n.FirstChild.FirstChild
	var body *html.Node
	for ; c != nil; c = c.NextSibling {
		//fmt.Println("->", c.Type, c.Data)
		if c.Data == "body" {
			body = c
		}
	}

	// body 底下一层
	c = body.FirstChild
	var table *html.Node
	for ; c != nil; c = c.NextSibling {
		//fmt.Println("->", c.Type, c.Data)
		if c.Data == "table" {
			table = c
		}
	}

	// table 底下一层
	c = table.FirstChild
	var tbody *html.Node
	for ; c != nil; c = c.NextSibling {
		//fmt.Println("->", c.Type, c.Data)
		if c.Data == "tbody" {
			tbody = c
		}
	}

	// tbody 底下一层
	c = tbody.FirstChild
	var trArr []*html.Node
	for ; c != nil; c = c.NextSibling {
		//fmt.Println("->", c.Type, c.Data)
		trArr = append(trArr, c)
	}

	for _, tr := range trArr {
		first := true
		for td := tr.FirstChild; td != nil; td = td.NextSibling {
			// td : <td>2023-03-15</td>
			// d : 2023-03-15
			d := td.FirstChild

			// add data to result
			if d != nil {
				space := " "
				if first {
					space = ""
					first = false
				}
				result += space + d.Data
			}
		}
		result += "\n"
	}
	return result
}

func handleOnePage(code string, page int) (result string, err error) {
	var content string
	if content, err = QueryJijin(code, page); err != nil {
		log.Println(err)
		return
	}
	content = strings.Split(content, "content:\"")[1]
	content = strings.Split(content, "\",records:")[0]
	// fmt.Println(content)

	var doc *html.Node
	if doc, err = html.Parse(strings.NewReader(content)); err != nil {
		log.Fatalln(err)
		return
	}
	result = handleHtml(doc)
	return
}

// var apidata={ content:"
//<table class='w782 comm lsjz'>
//    <thead>
//        <tr>
//            <th class='first'>净值日期</th>
//            <th>单位净值</th>
//            <th>累计净值</th>
//            <th>日增长率</th>
//            <th>申购状态</th>
//            <th>赎回状态</th>
//            <th class='tor last'>分红送配</th>
//        </tr>
//    </thead>
//    <tbody>
//        <tr>
//            <td>2023-04-12</td>
//            <td class='tor bold'>2.0934</td>
//            <td class='tor bold'>2.0934</td>
//            <td class='tor bold grn'>-1.90%</td>
//            <td>限制大额申购</td>
//            <td>开放赎回</td>
//            <td class='red unbold'></td>
//        </tr>
//        <tr>
//            <td>2023-04-11</td>
//            <td class='tor bold'>2.1339</td>
//            <td class='tor bold'>2.1339</td>
//            <td class='tor bold grn'>-0.66%</td>
//            <td>限制大额申购</td>
//            <td>开放赎回</td>
//            <td class='red unbold'></td>
//        </tr>
//        <tr>
//            <td>2023-04-10</td>
//            <td class='tor bold'>2.1480</td>
//            <td class='tor bold'>2.1480</td>
//            <td class='tor bold grn'>-0.56%</td>
//            <td>暂停申购</td>
//            <td>暂停赎回</td>
//            <td class='red unbold'></td>
//        </tr>
//        <tr>
//            <td>2023-04-07</td>
//            <td class='tor bold'>2.1601</td>
//            <td class='tor bold'>2.1601</td>
//            <td class='tor bold grn'>-0.09%</td>
//            <td>暂停申购</td>
//            <td>暂停赎回</td>
//            <td class='red unbold'></td>
//        </tr>
//        <tr>
//            <td>2023-04-06</td>
//            <td class='tor bold'>2.1620</td>
//            <td class='tor bold'>2.1620</td>
//            <td class='tor bold grn'>-0.82%</td>
//            <td>限制大额申购</td>
//            <td>开放赎回</td>
//            <td class='red unbold'></td>
//        </tr>
//        <tr>
//            <td>2023-04-04</td>
//            <td class='tor bold'>2.1799</td>
//            <td class='tor bold'>2.1799</td>
//            <td class='tor bold grn'>-0.27%</td>
//            <td>限制大额申购</td>
//            <td>开放赎回</td>
//            <td class='red unbold'></td>
//        </tr>
//        <tr>
//            <td>2023-04-03</td>
//            <td class='tor bold'>2.1859</td>
//            <td class='tor bold'>2.1859</td>
//            <td class='tor bold grn'>-0.23%</td>
//            <td>限制大额申购</td>
//            <td>开放赎回</td>
//            <td class='red unbold'></td>
//        </tr>
//        <tr>
//            <td>2023-03-31</td>
//            <td class='tor bold'>2.1910</td>
//            <td class='tor bold'>2.1910</td>
//            <td class='tor bold grn'>-0.24%</td>
//            <td>限制大额申购</td>
//            <td>开放赎回</td>
//            <td class='red unbold'></td>
//        </tr>
//        <tr>
//            <td>2023-03-30</td>
//            <td class='tor bold'>2.1963</td>
//            <td class='tor bold'>2.1963</td>
//            <td class='tor bold red'>0.59%</td>
//            <td>限制大额申购</td>
//            <td>开放赎回</td>
//            <td class='red unbold'></td>
//        </tr>
//        <tr>
//            <td>2023-03-29</td>
//            <td class='tor bold'>2.1834</td>
//            <td class='tor bold'>2.1834</td>
//            <td class='tor bold red'>1.11%</td>
//            <td>限制大额申购</td>
//            <td>开放赎回</td>
//            <td class='red unbold'></td>
//        </tr>
//        <tr>
//            <td>2023-03-28</td>
//            <td class='tor bold'>2.1595</td>
//            <td class='tor bold'>2.1595</td>
//            <td class='tor bold red'>1.56%</td>
//            <td>限制大额申购</td>
//            <td>开放赎回</td>
//            <td class='red unbold'></td>
//        </tr>
//        <tr>
//            <td>2023-03-27</td>
//            <td class='tor bold'>2.1263</td>
//            <td class='tor bold'>2.1263</td>
//            <td class='tor bold grn'>-1.43%</td>
//            <td>限制大额申购</td>
//            <td>开放赎回</td>
//            <td class='red unbold'></td>
//        </tr>
//        <tr>
//            <td>2023-03-24</td>
//            <td class='tor bold'>2.1572</td>
//            <td class='tor bold'>2.1572</td>
//            <td class='tor bold grn'>-0.69%</td>
//            <td>限制大额申购</td>
//            <td>开放赎回</td>
//            <td class='red unbold'></td>
//        </tr>
//        <tr>
//            <td>2023-03-23</td>
//            <td class='tor bold'>2.1721</td>
//            <td class='tor bold'>2.1721</td>
//            <td class='tor bold red'>2.17%</td>
//            <td>限制大额申购</td>
//            <td>开放赎回</td>
//            <td class='red unbold'></td>
//        </tr>
//        <tr>
//            <td>2023-03-22</td>
//            <td class='tor bold'>2.1259</td>
//            <td class='tor bold'>2.1259</td>
//            <td class='tor bold red'>0.44%</td>
//            <td>限制大额申购</td>
//            <td>开放赎回</td>
//            <td class='red unbold'></td>
//        </tr>
//        <tr>
//            <td>2023-03-21</td>
//            <td class='tor bold'>2.1165</td>
//            <td class='tor bold'>2.1165</td>
//            <td class='tor bold red'>2.51%</td>
//            <td>限制大额申购</td>
//            <td>开放赎回</td>
//            <td class='red unbold'></td>
//        </tr>
//        <tr>
//            <td>2023-03-20</td>
//            <td class='tor bold'>2.0647</td>
//            <td class='tor bold'>2.0647</td>
//            <td class='tor bold grn'>-1.94%</td>
//            <td>限制大额申购</td>
//            <td>开放赎回</td>
//            <td class='red unbold'></td>
//        </tr>
//        <tr>
//            <td>2023-03-17</td>
//            <td class='tor bold'>2.1055</td>
//            <td class='tor bold'>2.1055</td>
//            <td class='tor bold grn'>-0.48%</td>
//            <td>限制大额申购</td>
//            <td>开放赎回</td>
//            <td class='red unbold'></td>
//        </tr>
//        <tr>
//            <td>2023-03-16</td>
//            <td class='tor bold'>2.1157</td>
//            <td class='tor bold'>2.1157</td>
//            <td class='tor bold grn'>-0.76%</td>
//            <td>限制大额申购</td>
//            <td>开放赎回</td>
//            <td class='red unbold'></td>
//        </tr>
//        <tr>
//            <td>2023-03-15</td>
//            <td class='tor bold'>2.1320</td>
//            <td class='tor bold'>2.1320</td>
//            <td class='tor bold grn'>-0.20%</td>
//            <td>限制大额申购</td>
//            <td>开放赎回</td>
//            <td class='red unbold'></td>
//        </tr>
//    </tbody>
//</table>
//",records:1092,pages:55,curpage:1};
