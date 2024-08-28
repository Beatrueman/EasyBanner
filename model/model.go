package model

// 解析飞书的验证请求
type ChallengeRequest struct {
	Challenge string `json:"challenge"`
}

// 解析token
var ResponseData struct {
	TenantAccessToken string `json:"tenant_access_token"`
}

// MessageBody 定义了飞书消息事件的结构体
type MessageBody struct {
	Header struct {
		EventType string `json:"event_type"`
		EventID   string `json:"event_id"`
	} `json:"header"`
	Event struct {
		Message struct {
			Content   string `json:"content"`
			MessageID string `json:"message_id"`
			ChatID    string `json:"chat_id"`
			Mentions  []struct {
				Key string `json:"key"`
				ID  struct {
					UnionID string `json:"union_id"`
					UserID  string `json:"user_id"`
					OpenID  string `json:"open_id"`
				} `json:"id"`
				Name      string `json:"name"`
				TenantKey string `json:"tenant_key"`
			} `json:"mentions"`
		} `json:"message"`
		Sender struct {
			SenderID struct {
				UserID string `json:"user_id"`
			} `json:"sender_id"`
		} `json:"sender"`
	} `json:"event"`
}

// 请求体的数据
type CardCallback struct {
	Schema string `json:"schema"`
	Header struct {
		EventID    string `json:"event_id"`
		Token      string `json:"token"`
		CreateTime string `json:"create_time"`
		EventType  string `json:"event_type"`
		TenantKey  string `json:"tenant_key"`
		AppID      string `json:"app_id"`
	} `json:"header"`
	Event struct {
		Operator struct {
			TenantKey string `json:"tenant_key"`
			UserID    string `json:"user_id"`
			OpenID    string `json:"open_id"`
		} `json:"operator"`
		Token  string `json:"token"`
		Action struct {
			Value     map[string]interface{} `json:"value"`
			Tag       string                 `json:"tag"`
			Timezone  string                 `json:"timezone"`
			FormValue map[string]interface{} `json:"form_value"`
			Name      string                 `json:"name"`
		} `json:"action"`
		Host         string `json:"host"`
		DeliveryType string `json:"delivery_type"`
		Context      struct {
			URL           string `json:"url"`
			PreviewToken  string `json:"preview_token"`
			OpenMessageID string `json:"open_message_id"`
			OpenChatID    string `json:"open_chat_id"`
		} `json:"context"`
	} `json:"event"`
}

// IP和count
type IPData struct {
	Count int    `json:"count"`
	IP    string `json:"ip"`
}

// TemplateData 代表用于模板的数据结构
type TemplateData struct {
	NeedBanIPs []IPData `json:"needBanIPs"`
}

// BanIPRequest 用于发送至 ban API 的请求结构体
type BanIPRequest struct {
	IPList []string `json:"ip_list"`
}
