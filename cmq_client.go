package cmq_go

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	CURRENT_VERSION = "SDK_GO_1.3"
)

type CMQClient struct {
	Endpoint  string
	Path      string
	SecretId  string
	SecretKey string
	conn      *http.Client
}

func NewCMQClient(endpoint, path, secretId, secretKey string) *CMQClient {
	return &CMQClient{
		Endpoint:  endpoint,
		Path:      path,
		SecretId:  secretId,
		SecretKey: secretKey,
		conn: &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				MaxIdleConns:          500,
				MaxIdleConnsPerHost:   100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		},
	}
}

func (this *CMQClient) call(action string, param map[string]string) (resp string, err error) {
	param["Action"] = action
	param["Nonce"] = strconv.Itoa(rand.Int())
	param["SecretId"] = this.SecretId
	param["Timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	param["RequestClient"] = CURRENT_VERSION
	param["SignatureMethod"] = "HmacSHA1"
	sortedParamKeys := make([]string, 0)
	for k := range param {
		sortedParamKeys = append(sortedParamKeys, k)
	}
	sort.Strings(sortedParamKeys)

	host := ""
	if strings.HasPrefix(this.Endpoint, "https") {
		host = this.Endpoint[8:]
	} else {
		host = this.Endpoint[7:]
	}

	src := http.MethodPost + host + this.Path + "?"
	flag := false
	for _, key := range sortedParamKeys {
		if flag {
			src += "&"
		}
		src += key + "=" + param[key]
		flag = true
	}
	mac := hmac.New(sha1.New, []byte(this.SecretKey))
	mac.Write([]byte(src))
	param["Signature"] = base64.StdEncoding.EncodeToString(mac.Sum(nil))

	reqStr := ""
	urlStr := this.Endpoint + this.Path
	flag = false
	for k, v := range param {
		if flag {
			reqStr += "&"
		}
		reqStr += k + "=" + url.QueryEscape(v)
		flag = true
	}

	userTimeout := 0
	if UserpollingWaitSeconds, found := param["UserpollingWaitSeconds"]; found {
		userTimeout, err = strconv.Atoi(UserpollingWaitSeconds)
		if err != nil {
			return "", fmt.Errorf("strconv failed: %v", err.Error())
		}
	}

	resp, err = this.doPost(urlStr, reqStr, userTimeout)
	return
}

func (this *CMQClient) doPost(urlStr, reqStr string, userTimeout int) (result string, err error) {
	timeout := 3000
	if userTimeout >= 0 {
		timeout += userTimeout
	}
	this.conn.Timeout = time.Duration(timeout) * time.Millisecond

	req, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewReader([]byte(reqStr)))
	if err != nil {
		return "", fmt.Errorf("make http req error %v", err)
	}
	resp, err := this.conn.Do(req)
	if err != nil {
		return "", fmt.Errorf("http error  %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http error code %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read http resp body error %v", err)
	}
	result = string(body)
	return
}
