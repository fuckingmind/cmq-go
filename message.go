package cmq_go

import "fmt"

type Message struct {
	/** 服务器返回的消息ID */
	MsgId string `json:"msgId"`
	/** 每次消费唯一的消息句柄，用于删除等操作 */
	ReceiptHandle string `json:"receiptHandle"`
	/** 消息体 */
	MsgBody string `json:"msgBody"`
	/** 消息发送到队列的时间，从 1970年1月1日 00:00:00 000 开始的毫秒数 */
	EnqueueTime int64 `json:"enqueueTime"`
	/** 消息下次可见的时间，从 1970年1月1日 00:00:00 000 开始的毫秒数 */
	NextVisibleTime int64 `json:"nextVisibleTime"`
	/** 消息第一次出队列的时间，从 1970年1月1日 00:00:00 000 开始的毫秒数 */
	FirstDequeueTime int64 `json:"firstDequeueTime"`
	/** 出队列次数 */
	DequeueCount int `json:"dequeueCount"`
	MsgTag       []string
}

// CommResp 通用返回
type CommResp struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"requestId"`
}

func (resp CommResp) Error() string {
	return fmt.Sprintf("request(%s) response=%d(%s)", resp.RequestID, resp.Code, resp.Message)
}
