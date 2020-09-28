package cmq_go

import (
	"strconv"
)

type SubscriptionMeta struct {
	//Subscription 订阅的主题所有者的appId
	TopicOwner string `json:"topicOwner"`
	//订阅的终端地址
	Endpoint string `json:"endpoint"`
	//订阅的协议
	Protocal string `json:"protocol"`
	//推送消息出现错误时的重试策略
	NotifyStrategy string `json:"notifyStrategy"`
	//向 Endpoint 推送的消息内容格式
	NotifyContentFormat string `json:"notifyContentFormat"`
	//描述了该订阅中消息过滤的标签列表（仅标签一致的消息才会被推送）
	FilterTag []string `json:"filterTag"`
	//Subscription 的创建时间，从 1970-1-1 00:00:00 到现在的秒值
	CreateTime int `json:"createTime"`
	//修改 Subscription 属性信息最近时间，从 1970-1-1 00:00:00 到现在的秒值
	LastModifyTime int `json:"lastModifyTime"`
	//该订阅待投递的消息数
	MsgCount   int      `json:"msgCount"`
	BindingKey []string `json:"bindingKey"`
}

type Subscription struct {
	topicName        string
	subscriptionName string
	client           *CMQClient
}

func NewSubscription(topicName, subscriptionName string, client *CMQClient) *Subscription {
	return &Subscription{
		topicName:        topicName,
		subscriptionName: subscriptionName,
		client:           client,
	}
}

func (this *Subscription) ClearFilterTags() (err error) {
	param := make(map[string]string)
	param["topicName"] = this.topicName
	param["subscriptionName "] = this.subscriptionName

	return this.client.callWithoutResult("ClearSubscriptionFilterTags", param)
}

func (this *Subscription) SetSubscriptionAttributes(meta SubscriptionMeta) (err error) {
	param := make(map[string]string)
	param["topicName"] = this.topicName
	param["subscriptionName "] = this.subscriptionName
	if meta.NotifyStrategy != "" {
		param["notifyStrategy"] = meta.NotifyStrategy
	}
	if meta.NotifyContentFormat != "" {
		param["notifyContentFormat"] = meta.NotifyContentFormat
	}
	if meta.FilterTag != nil {
		for i, flag := range meta.FilterTag {
			param["filterTag."+strconv.Itoa(i+1)] = flag
		}
	}
	if meta.BindingKey != nil {
		for i, binding := range meta.BindingKey {
			param["bindingKey."+strconv.Itoa(i+1)] = binding
		}
	}

	return this.client.callWithoutResult("SetSubscriptionAttributes", param)
}

func (this *Subscription) GetSubscriptionAttributes() (*SubscriptionMeta, error) {
	param := make(map[string]string)
	param["topicName"] = this.topicName
	param["subscriptionName"] = this.subscriptionName

	var resp struct {
		CommResp
		SubscriptionMeta
	}

	if err := this.client.call("GetSubscriptionAttributes", param, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, &resp.CommResp
	}

	return &resp.SubscriptionMeta, nil
}
