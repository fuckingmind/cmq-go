package cmq_go

import (
	"strconv"
)

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
