package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/miekg/dns"

	"github.com/urlooker/agent/g"
)

func LookupIP(domain string, timeout int) ([]string, error) {
	var c dns.Client
	var err error
	ips := []string{}
	domain = strings.TrimRight(domain, ".") + "."
	c.DialTimeout = time.Duration(timeout) * time.Millisecond
	c.ReadTimeout = time.Duration(timeout) * time.Millisecond
	c.WriteTimeout = time.Duration(timeout) * time.Millisecond

	m := new(dns.Msg)
	m.SetQuestion(domain, dns.TypeA)

	dns := g.Config.Dns
	if dns == "" {
		dns = "114.114.114.114:53"
	}

	ret, _, err := c.Exchange(m, dns)
	if err != nil {
		domain = strings.TrimRight(domain, ".")
		e := fmt.Sprintf("lookup error: %s, %s", domain, err.Error())
		return ips, errors.New(e)
	}

	for _, i := range ret.Answer {
		result := strings.Split(i.String(), "\t")
		if result[3] == "A" && IsIP(result[4]) {
			ips = append(ips, result[4])
		}
	}
	return ips, err
}

func IsIP(ip string) bool {
	if ip != "" {
		isOk, _ := regexp.MatchString(`^(\d{1,3}\.){3}\d{1,3}$`, ip)
		if isOk {
			return isOk
		}
	}
	return false
}
