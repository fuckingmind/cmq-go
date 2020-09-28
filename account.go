package cmq_go

import (
	"fmt"
	"strconv"
)

type Account struct {
	client *CMQClient
}

func NewAccount(endpoint, secretId, secretKey string) *Account {
	return &Account{
		client: NewCMQClient(endpoint, "/v2/index.php", secretId, secretKey),
	}
}

func (this *Account) CreateQueue(queueName string, queueMeta QueueMeta) error {
	if queueName != "" {
		param := make(map[string]string)
		param["queueName"] = queueName
		if queueMeta.MaxMsgHeapNum > 0 {
			param["maxMsgHeapNum"] = strconv.Itoa(queueMeta.MaxMsgHeapNum)
		}
		if queueMeta.PollingWaitSeconds > 0 {
			param["pollingWaitSeconds"] = strconv.Itoa(queueMeta.PollingWaitSeconds)
		}
		if queueMeta.VisibilityTimeout > 0 {
			param["visibilityTimeout"] = strconv.Itoa(queueMeta.VisibilityTimeout)
		}
		if queueMeta.MaxMsgSize > 0 {
			param["maxMsgSize"] = strconv.Itoa(queueMeta.MaxMsgSize)
		}
		if queueMeta.MsgRetentionSeconds > 0 {
			param["msgRetentionSeconds"] = strconv.Itoa(queueMeta.MsgRetentionSeconds)
		}
		if queueMeta.RewindSeconds > 0 {
			param["rewindSeconds"] = strconv.Itoa(queueMeta.RewindSeconds)
		}

		return this.client.callWithoutResult("CreateQueue", param)
	}
	return nil
}

func (this *Account) DeleteQueue(queueName string) error {
	if queueName != "" {
		return this.client.callWithoutResult("DeleteQueue", map[string]string{"queueName": queueName})
	}
	return nil
}

func (this *Account) ListQueue(searchWord string, offset, limit int) (
	totalCount int, queueList []string, err error) {
	queueList = make([]string, 0)
	param := make(map[string]string)
	if searchWord != "" {
		param["searchWord"] = searchWord
	}
	if offset >= 0 {
		param["offset"] = strconv.Itoa(offset)
	}
	if limit > 0 {
		param["limit"] = strconv.Itoa(limit)
	}

	var resp struct {
		CommResp
		TotalCount int `json:"totalCount"`
		QueueList  []struct {
			QueueName string `json:"queueName"`
		} `json:"queueList"`
	}

	if err := this.client.call("ListQueue", param, &resp); err != nil {
		return 0, nil, err
	}

	if resp.Code != 0 {
		return 0, nil, &resp.CommResp
	}

	totalCount = resp.TotalCount
	for _, q := range resp.QueueList {
		queueList = append(queueList, q.QueueName)
	}

	return
}

func (this *Account) GetQueue(queueName string) (queue *Queue) {
	return NewQueue(queueName, this.client)
}

func (this *Account) GetTopic(topicName string) (topic *Topic) {
	return NewTopic(topicName, this.client)
}

func (this *Account) CreateTopic(topicName string, maxMsgSize int) (err error, code int) {
	err = _createTopic(this.client, topicName, maxMsgSize, 1)
	return
}

func _createTopic(client *CMQClient, topicName string, maxMsgSize, filterType int) (err error) {
	param := make(map[string]string)
	if topicName == "" {
		err = fmt.Errorf("createTopic failed: topicName is empty")
		return
	}
	param["topicName"] = topicName
	param["filterType"] = strconv.Itoa(filterType)
	if maxMsgSize < 1024 || maxMsgSize > 1048576 {
		err = fmt.Errorf("createTopic failed: Invalid parameter: maxMsgSize > 1024KB or maxMsgSize < 1KB")
		return
	}
	param["maxMsgSize"] = strconv.Itoa(maxMsgSize)

	return client.callWithoutResult("CreateTopic", param)

}

func (this *Account) DeleteTopic(topicName string) error {
	if topicName != "" {
		return this.client.callWithoutResult("DeleteTopic", map[string]string{"topicName": topicName})
	}
	return nil
}

func (this *Account) ListTopic(searchWord string, offset, limit int) (
	totalCount int, topicList []string, err error) {
	topicList = make([]string, 0)
	param := make(map[string]string)
	if searchWord != "" {
		param["searchWord"] = searchWord
	}
	if offset > 0 {
		param["offset"] = strconv.Itoa(offset)
	}
	if limit > 0 {
		param["limit"] = strconv.Itoa(limit)
	}

	var resp struct {
		CommResp
		TotalCount int `json:"totalCount"`
		TopicList  []struct {
			TopicName string `json:"topicName"`
		} `json:"topicList"`
	}

	if err := this.client.call("ListTopic", param, &resp); err != nil {
		return 0, nil, err
	}

	if resp.Code != 0 {
		return 0, nil, &resp.CommResp
	}

	totalCount = resp.TotalCount
	for _, q := range resp.TopicList {
		topicList = append(topicList, q.TopicName)
	}

	return
}

func (this *Account) CreateSubscribe(topicName, subscriptionName, endpoint, protocol, notifyContentFormat string) error {
	return _createSubscribe(this.client, topicName, subscriptionName, endpoint, protocol, nil, nil, "BACKOFF_RETRY", notifyContentFormat)
}

func _createSubscribe(client *CMQClient, topicName, subscriptionName, endpoint, protocol string, filterTag []string,
	bindingKey []string, notifyStrategy, notifyContentFormat string) (err error) {
	param := make(map[string]string)
	if topicName == "" {
		err = fmt.Errorf("createSubscribe failed: topicName is empty")
		return
	}
	param["topicName"] = topicName

	if subscriptionName == "" {
		err = fmt.Errorf("createSubscribe failed: subscriptionName is empty")
		return
	}
	param["subscriptionName"] = subscriptionName

	if endpoint == "" {
		err = fmt.Errorf("createSubscribe failed: endpoint is empty")
		return
	}
	param["endpoint"] = endpoint

	if protocol == "" {
		err = fmt.Errorf("createSubscribe failed: protocal is empty")
		return
	}
	param["protocol"] = protocol

	if notifyStrategy == "" {
		err = fmt.Errorf("createSubscribe failed: notifyStrategy is empty")
		return
	}
	param["notifyStrategy"] = notifyStrategy

	if notifyContentFormat == "" {
		err = fmt.Errorf("createSubscribe failed: notifyContentFormat is empty")
		return
	}
	param["notifyContentFormat"] = notifyContentFormat

	if filterTag != nil {
		for i, tag := range filterTag {
			param["filterTag."+strconv.Itoa(i+1)] = tag
		}
	}

	if bindingKey != nil {
		for i, key := range bindingKey {
			param["bindingKey."+strconv.Itoa(i+1)] = key
		}
	}

	return client.callWithoutResult("Subscribe", param)
}

func (this *Account) DeleteSubscribe(topicName, subscriptionName string) (err error) {
	param := make(map[string]string)
	if topicName == "" {
		err = fmt.Errorf("createSubscribe failed: topicName is empty")
		return
	}
	param["topicName"] = topicName

	if subscriptionName == "" {
		err = fmt.Errorf("createSubscribe failed: subscriptionName is empty")
		//log.Printf("%v", err.Error())
		return
	}
	param["subscriptionName"] = subscriptionName
	return this.client.callWithoutResult("Unsubscribe", param)
}

func (this *Account) GetSubscription(topicName, subscriptionName string) *Subscription {
	return NewSubscription(topicName, subscriptionName, this.client)
}
