package network

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"time"
)

var (
	HostReg = regexp.MustCompile(`^[a-zA-Z0-9][-a-zA-Z0-9_]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9_]{0,62})+$`)
	UrlReg  = regexp.MustCompile(`((http|ftp|https)://)(([a-zA-Z0-9\._-]+\.[a-zA-Z]{2,6})|([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}))(:[0-9]{1,4})*(/[a-zA-Z0-9\&%_\./-~-]*)?`)
	AddrReg = regexp.MustCompile(`^[a-zA-Z0-9][-a-zA-Z0-9_]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9_]{0,62})+:(\d{1,5})$`)
)

func Cmd(ctx context.Context, cmdName string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, cmdName, args...)
	outer := bytes.NewBuffer([]byte{})
	cmd.Stdout = outer
	cmd.Stderr = outer
	err := cmd.Start()
	if err != nil {
		return "", err
	}
	_ = cmd.Wait()

	reader := bufio.NewReader(outer)
	result := ""
	for {
		l, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		result += string(l) + "\n"
	}

	return result, nil
}

func Ping(ctx context.Context, host string) (string, error) {
	if !HostReg.MatchString(host) {
		return "", errors.New("invalid host")
	}
	return Cmd(ctx, "ping", "-c", "4", host)
}

func Curl(ctx context.Context, addr string, timeout time.Duration) (string, error) {
	if !UrlReg.MatchString(addr) {
		return "", errors.New("invalid url")
	}

	reqUrl, err := url.Parse(addr)
	if err != nil {
		return "", errors.WithMessage(err, "parse url")
	}

	// 协议过滤
	if reqUrl.Scheme != "http" && reqUrl.Scheme != "https" {
		return "", errors.Errorf("unSupport %s scheme", reqUrl.Scheme)
	}

	// ip过滤
	ips, err := net.LookupIP(reqUrl.Hostname())
	if err != nil {
		return "", errors.WithMessage(err, "lookup ip")
	}

	for _, ip := range ips {
		if ip.IsLoopback() {
			return "", errors.New("can not access inner ip address")
		}
	}

	req, err := http.NewRequestWithContext(ctx, "GET", addr, nil)
	if err != nil {
		return "", err
	}
	client := http.Client{Timeout: timeout}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func Telnet(host string, port int, timeout time.Duration) (bool, time.Duration, error) {

	timeStart := time.Now()
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
	if nil != err {
		return false, 0, err
	}

	_ = conn.Close()

	return true, time.Now().Sub(timeStart), nil
}
