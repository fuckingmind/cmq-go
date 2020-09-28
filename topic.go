package cmq_go

import (
	"fmt"
	"strconv"
)

type Topic struct {
	topicName string
	client    *CMQClient
}

func NewTopic(topicName string, client *CMQClient) (queue *Topic) {
	return &Topic{
		topicName: topicName,
		client:    client,
	}
}

func (this *Topic) SetTopicAttributes(maxMsgSize int) (err error) {
	if maxMsgSize < 1024 || maxMsgSize > 1048576 {
		err = fmt.Errorf("Invalid parameter maxMsgSize < 1KB or maxMsgSize > 1024KB")
		return
	}
	param := make(map[string]string)
	param["topicName"] = this.topicName
	param["maxMsgSize"] = strconv.Itoa(maxMsgSize)

	return this.client.callWithoutResult("SetTopicAttributes", param)
}

type TopicMeta struct {
	// 当前该主题的消息堆积数
	MsgCount int `json:"msgCount"`
	// 消息最大长度，取值范围1024-1048576 Byte（即1-1024K），默认1048576
	MaxMsgSize int `json:"maxMsgSize"`
	//消息在主题中最长存活时间，从发送到该主题开始经过此参数指定的时间后，
	//不论消息是否被成功推送给用户都将被删除，单位为秒。固定为一天，该属性不能修改。
	MsgRetentionSeconds int `json:"msgRetentionSeconds"`
	//创建时间
	CreateTime int `json:"createTime"`
	//修改属性信息最近时间
	LastModifyTime int `json:"lastModifyTime"`
	LoggingEnabled int `json:"loggingEnabled"`
	FilterType     int `json:"filterType"`
}

func (this *Topic) GetTopicAttributes() (TopicMeta, error) {
	param := make(map[string]string)
	param["topicName"] = this.topicName

	var resp struct {
		CommResp
		TopicMeta
	}

	if err := this.client.call("GetTopicAttributes", param, &resp); err != nil {
		return resp.TopicMeta, err
	}

	if resp.Code != 0 {
		return resp.TopicMeta, &resp.CommResp
	}

	return resp.TopicMeta, nil
}

func (this *Topic) PublishMessage(message string, tagList []string) (msgId string, err error) {
	msgId, err = _publishMessage(this.client, this.topicName, message, tagList, "")
	return
}

func _publishMessage(client *CMQClient, topicName, msg string, tagList []string, routingKey string) (string, error) {
	param := make(map[string]string)
	param["topicName"] = topicName
	param["msgBody"] = msg
	if routingKey != "" {
		param["routingKey"] = routingKey
	}
	if tagList != nil {
		for i, tag := range tagList {
			param["msgTag."+strconv.Itoa(i+1)] = tag
		}
	}

	var resp struct {
		CommResp
		MsgID string `json:"msgId"`
	}

	if err := client.call("PublishMessage", param, &resp); err != nil {
		return "", err
	}
	if resp.Code != 0 {
		return "", &resp.CommResp
	}
	return resp.MsgID, nil
}

func (this *Topic) BatchPublishMessage(msgList []string) (msgIds []string, err error) {
	msgIds, err = _batchPublishMessage(this.client, this.topicName, msgList, nil, "")
	return
}

func _batchPublishMessage(client *CMQClient, topicName string, msgList, tagList []string, routingKey string) (msgIds []string, err error) {
	param := make(map[string]string)
	param["topicName"] = topicName
	if routingKey != "" {
		param["routingKey"] = routingKey
	}
	if msgList != nil {
		for i, msg := range msgList {
			param["msgBody."+strconv.Itoa(i+1)] = msg
		}
	}
	if tagList != nil {
		for i, tag := range tagList {
			param["msgTag."+strconv.Itoa(i+1)] = tag
		}
	}

	var resp struct {
		CommResp
		MsgList []struct {
			MsgID string `json:"msgId"`
		} `json:"msgList"`
	}

	if err := client.call("BatchPublishMessage", param, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, &resp.CommResp
	}

	for _, msg := range resp.MsgList {
		msgIds = append(msgIds, msg.MsgID)
	}

	return
}

func (this *Topic) ListSubscription(offset, limit int, searchWord string) (totalCount int, subscriptionList []string, err error) {
	param := make(map[string]string)
	param["topicName"] = this.topicName
	if searchWord != "" {
		param["searchWord "] = searchWord
	}
	if offset >= 0 {
		param["offset "] = strconv.Itoa(offset)
	}
	if limit > 0 {
		param["limit "] = strconv.Itoa(limit)
	}

	var resp struct {
		CommResp
		TotalCount       int `json:"totalCount"`
		SubscriptionList []struct {
			SubscriptionName string `json:"subscriptionName"`
		} `json:"subscriptionList"`
	}

	if err := this.client.call("ListSubscriptionByTopic", param, &resp); err != nil {
		return 0, nil, err
	}

	if resp.Code != 0 {
		return 0, nil, &resp.CommResp
	}

	totalCount = resp.TotalCount
	for _, sub := range resp.SubscriptionList {
		subscriptionList = append(subscriptionList, sub.SubscriptionName)
	}
	return
}
