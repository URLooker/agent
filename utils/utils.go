package utils

import (
	"crypto/tls"
	//"html"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/httplib"
	webg "github.com/urlooker/web/g"

	"github.com/urlooker/agent/g"
)

const (
	NO_ERROR = 0
	REQ_TIMEOUT = 1
	INVALID_RESP_CODE = 2
	KEYWORD_UNMATCH = 3
	DNS_ERROR = 4
)

func CheckTargetStatus(item *webg.DetectedItem) {
	defer func() {
		<-g.WorkerChan
	}()

	checkResult := checkTargetStatus(item)
	g.CheckResultQueue.PushFront(checkResult)
}

func doCheckTargetStatus(item *webg.DetectedItem, req *httplib.BeegoHTTPRequest) (itemCheckResult *webg.CheckResult) {
	itemCheckResult = &webg.CheckResult{
		Sid:      item.Sid,
		Domain:   item.Domain,
		Creator:  item.Creator,
		Tag:      item.Tag,
		Target:   item.Target,
		Ip:       item.Ip,
		RespTime: item.Timeout,
		RespCode: "0",
	}

	reqStartTime := time.Now()

	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	req.SetTimeout(3 * time.Second, 10 * time.Second)
	req.Header("Content-Type", "application/json")
	req.SetHost(item.Domain)
	if len(item.PostData) > 0 && item.Method != "GET" {
		req.Body(item.PostData)
	}
	if item.Data != "" {
		req.Header("Cookie", item.Data)
	}

	resp, err := req.Response()
	itemCheckResult.PushTime = time.Now().Unix()

	if err != nil {
		log.Println("[ERROR]:", item.Sid, item.Domain, err)
		itemCheckResult.Status = REQ_TIMEOUT
		return
	}
	defer resp.Body.Close()

	respCode := strconv.Itoa(resp.StatusCode)
	itemCheckResult.RespCode = respCode

	respTime := int(time.Now().Sub(reqStartTime).Nanoseconds() / 1000000)
	itemCheckResult.RespTime = respTime

	log.Println("[req_status]:", respCode + "|" + item.Target + "|" + respTime + "|" + item.Timeout)
	if respTime > item.Timeout {
		itemCheckResult.Status = REQ_TIMEOUT
		return
	}

	if strings.Index(respCode, item.ExpectCode) == 0 || (len(item.ExpectCode) == 0 && respCode == "200") {
		if len(item.Keywords) > 0 {
			contents, _ := ioutil.ReadAll(resp.Body)
			contentStr := string(contents)
			if !strings.Contains(contentStr, item.Keywords) {
				log.Println("[result is not expected]: ", item.Keywords + "$$$$$$$$" + contentStr)
				itemCheckResult.Status = KEYWORD_UNMATCH
				return
			}
		}

		itemCheckResult.Status = NO_ERROR
		return

	} else {
		itemCheckResult.Status = INVALID_RESP_CODE
	}
	return
}

func checkTargetStatus(item *webg.DetectedItem) (itemCheckResult *webg.CheckResult) {
	method := item.Method
	switch method {
	case "GET":
		req := httplib.Get(item.Target)
		return doCheckTargetStatus(item, req)
	case "POST":
		req := httplib.Post(item.Target)
		return doCheckTargetStatus(item, req)
	case "PUT":
		req := httplib.Put(item.Target)
		return doCheckTargetStatus(item, req)
	case "DELETE":
		req := httplib.Delete(item.Target)
		return doCheckTargetStatus(item, req)
	}

	return
}
