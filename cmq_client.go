package cmq_go

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

const (
	CURRENT_VERSION = "SDK_GO_1.3"
)

type CMQClient struct {
	uri       *url.URL
	SecretId  string
	SecretKey string
	conn      *http.Client
}

func NewCMQClient(endpoint, path, secretId, secretKey string) *CMQClient {
	client := &CMQClient{
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

	client.uri, _ = url.Parse(endpoint + path)
	return client
}

func (this *CMQClient) callWithoutResult(action string, param map[string]string) error {
	res := &CommResp{}
	if err := this.call(action, param, res); err != nil {
		return err
	}
	if res.Code != 0 {
		return res
	}
	return nil
}

func (this *CMQClient) call(action string, param map[string]string, ires interface{}) error {
	uriParams := make(url.Values)
	for k, v := range param {
		uriParams.Set(k, v)
	}
	uriParams.Set("Action", action)
	uriParams.Set("Nonce", strconv.Itoa(rand.Int()))
	uriParams.Set("SecretId", this.SecretId)
	uriParams.Set("Timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	uriParams.Set("RequestClient", CURRENT_VERSION)
	uriParams.Set("SignatureMethod", "HmacSHA1")

	sortedParamKeys := make([]string, 0)
	for k := range param {
		sortedParamKeys = append(sortedParamKeys, k)
	}
	sort.Strings(sortedParamKeys)

	paramStr := uriParams.Encode()
	mac := hmac.New(sha1.New, []byte(this.SecretKey))
	mac.Write([]byte(http.MethodPost + this.uri.Host + this.uri.Path + "?" + paramStr))
	paramStr += "&Signature=" + base64.StdEncoding.EncodeToString(mac.Sum(nil))

	userTimeout := 0
	if UserpollingWaitSeconds, found := param["UserpollingWaitSeconds"]; found {
		userTimeout, _ = strconv.Atoi(UserpollingWaitSeconds)
	}

	this.conn.Timeout = time.Duration(3000+userTimeout) * time.Millisecond

	req, err := http.NewRequest(http.MethodPost, this.uri.String(), bytes.NewReader([]byte(paramStr)))
	if err != nil {
		return err
	}
	resp, err := this.conn.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http error code %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, ires)
}
