package cmq_go

import (
	"fmt"
	"strconv"
)

type QueueMeta struct {
	/** 最大堆积消息数 */
	MaxMsgHeapNum int `json:"maxMsgHeapNum"`
	/** 消息接收长轮询等待时间 */
	PollingWaitSeconds int `json:"pollingWaitSeconds"`
	/** 消息可见性超时 */
	VisibilityTimeout int `json:"visibilityTimeout"`
	/** 消息最大长度 */
	MaxMsgSize int `json:"maxMsgSize"`
	/** 消息保留周期 */
	MsgRetentionSeconds int `json:"msgRetentionSeconds"`
	/** 队列创建时间 */
	CreateTime int `json:"createTime"`
	/** 队列属性最后修改时间 */
	LastModifyTime int `json:"lastModifyTime"`
	/** 队列处于Active状态的消息总数 */
	ActiveMsgNum int `json:"activeMsgNum"`
	/** 队列处于Inactive状态的消息总数 */
	InactiveMsgNum int `json:"inactiveMsgNum"`
	/** 已删除的消息，但还在回溯保留时间内的消息数量 */
	RewindMsgNum int `json:"rewindMsgNum"`
	/** 消息最小未消费时间 */
	MinMsgTime int `json:"minMsgTime"`
	/** 延时消息数量 */
	DelayMsgNum int `json:"delayMsgNum"`
	/** 回溯时间 */
	RewindSeconds int `json:"rewindSeconds"`
}

type Queue struct {
	queueName string
	client    *CMQClient
}

func NewQueue(queueName string, client *CMQClient) (queue *Queue) {
	return &Queue{
		queueName: queueName,
		client:    client,
	}
}

func (this *Queue) SetQueueAttributes(queueMeta QueueMeta) (err error) {
	param := make(map[string]string)
	param["queueName"] = this.queueName

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

	return this.client.callWithoutResult("SetQueueAttributes", param)
}

func (this *Queue) GetQueueAttributes() (queueMeta QueueMeta, err error) {
	param := make(map[string]string)
	param["queueName"] = this.queueName

	var resp struct {
		CommResp
		QueueMeta
	}

	if err = this.client.call("GetQueueAttributes", param, &resp); err != nil {
		return
	}

	if resp.Code != 0 {
		err = &resp.CommResp
		return
	}

	queueMeta = resp.QueueMeta
	return
}

func (this *Queue) SendMessage(msgBody string) (string, error) {
	return _sendMessage(this.client, msgBody, this.queueName, 0)
}

func (this *Queue) SendDelayMessage(msgBody string, delaySeconds int) (string, error) {
	return _sendMessage(this.client, msgBody, this.queueName, delaySeconds)
}

func _sendMessage(client *CMQClient, msgBody, queueName string, delaySeconds int) (messageId string, err error) {
	param := make(map[string]string)
	param["queueName"] = queueName
	param["msgBody"] = msgBody
	param["delaySeconds"] = strconv.Itoa(delaySeconds)

	var resp struct {
		CommResp
		MsgID string `json:"msgId"`
	}

	if err = client.call("SendMessage", param, &resp); err != nil {
		return
	}

	if resp.Code != 0 {
		return "", &resp.CommResp
	}

	return resp.MsgID, nil
}

func (this *Queue) BatchSendMessage(msgBodys []string) ([]string, error) {
	return _batchSendMessage(this.client, msgBodys, this.queueName, 0)
}

func (this *Queue) BatchSendDelayMessage(msgBodys []string, delaySeconds int) ([]string, error) {
	return _batchSendMessage(this.client, msgBodys, this.queueName, delaySeconds)
}

func _batchSendMessage(client *CMQClient, msgBodys []string, queueName string, delaySeconds int) (messageIds []string, err error) {
	messageIds = make([]string, 0)

	if len(msgBodys) == 0 || len(msgBodys) > 16 {
		err = fmt.Errorf("message size is 0 or more than 16")
		return
	}

	param := make(map[string]string)
	param["queueName"] = queueName
	for i, msgBody := range msgBodys {
		param["msgBody."+strconv.Itoa(i+1)] = msgBody
	}
	param["delaySeconds"] = strconv.Itoa(delaySeconds)

	var resp struct {
		CommResp
		Msgs []struct {
			MsgID string `json:"msgId"`
		} `json:"msgList"`
	}

	if err = client.call("BatchSendMessage", param, &resp); err != nil {
		return
	}

	if resp.Code != 0 {
		return nil, &resp.CommResp
	}

	for _, msg := range resp.Msgs {
		messageIds = append(messageIds, msg.MsgID)
	}

	return messageIds, nil
}

func (this *Queue) ReceiveMessage(pollingWaitSeconds int) (Message, error) {
	param := make(map[string]string)
	param["queueName"] = this.queueName
	if pollingWaitSeconds >= 0 {
		param["UserpollingWaitSeconds"] = strconv.Itoa(pollingWaitSeconds * 1000)
		param["pollingWaitSeconds"] = strconv.Itoa(pollingWaitSeconds)
	} else {
		param["UserpollingWaitSeconds"] = strconv.Itoa(30000)
	}

	var resp struct {
		CommResp
		Message
	}

	if err := this.client.call("ReceiveMessage", param, &resp); err != nil {
		return resp.Message, err
	}

	if resp.Code != 0 {
		return resp.Message, &resp.CommResp
	}
	return resp.Message, nil
}

func (this *Queue) BatchReceiveMessage(numOfMsg, pollingWaitSeconds int) ([]Message, error) {
	param := make(map[string]string)
	param["queueName"] = this.queueName
	param["numOfMsg"] = strconv.Itoa(numOfMsg)
	if pollingWaitSeconds >= 0 {
		param["UserpollingWaitSeconds"] = strconv.Itoa(pollingWaitSeconds * 1000)
		param["pollingWaitSeconds"] = strconv.Itoa(pollingWaitSeconds)
	} else {
		param["UserpollingWaitSeconds"] = strconv.Itoa(30000)
	}

	var resp struct {
		CommResp
		Msgs []Message `json:"msgInfoList"`
	}

	if err := this.client.call("BatchReceiveMessage", param, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, &resp.CommResp
	}
	return resp.Msgs, nil
}

func (this *Queue) DeleteMessage(receiptHandle string) (err error) {
	param := make(map[string]string)
	param["queueName"] = this.queueName
	param["receiptHandle"] = receiptHandle

	return this.client.callWithoutResult("DeleteMessage", param)
}

func (this *Queue) BatchDeleteMessage(receiptHandles []string) (err error) {
	if len(receiptHandles) == 0 {
		return
	}
	param := make(map[string]string)
	param["queueName"] = this.queueName
	for i, receiptHandle := range receiptHandles {
		param["receiptHandle."+strconv.Itoa(i+1)] = receiptHandle
	}

	return this.client.callWithoutResult("BatchDeleteMessage", param)
}

func (this *Queue) RewindQueue(backTrackingTime int) (err error) {
	if backTrackingTime <= 0 {
		return
	}
	param := make(map[string]string)
	param["queueName"] = this.queueName
	param["startConsumeTime"] = strconv.Itoa(backTrackingTime)

	return this.client.callWithoutResult("RewindQueue", param)
}
