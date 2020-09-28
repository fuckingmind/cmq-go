package cmq_go

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
