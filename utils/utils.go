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
	NO_ERROR          = 0
	REQ_TIMEOUT       = 1
	INVALID_RESP_CODE = 2
	KEYWORD_UNMATCH   = 3
	DNS_ERROR         = 4
)

func CheckTargetStatus(item *webg.DetectedItem) {
	defer func() {
		<-g.WorkerChan
	}()

	checkResult := checkTargetStatus(item)
	g.CheckResultQueue.PushFront(checkResult)
}

func checkTargetStatus(item *webg.DetectedItem) (itemCheckResult *webg.CheckResult) {
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
	req := httplib.Get(item.Target)
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	req.SetTimeout(3*time.Second, 10*time.Second)
	req.Header("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.SetHost(item.Domain)
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

	if respTime > item.Timeout {
		itemCheckResult.Status = REQ_TIMEOUT
		return
	}

	if strings.Index(respCode, item.ExpectCode) == 0 || (len(item.ExpectCode) == 0 && respCode == "200") {
		if len(item.Keywords) > 0 {
			contents, _ := ioutil.ReadAll(resp.Body)
			if !strings.Contains(string(contents), item.Keywords) {
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
