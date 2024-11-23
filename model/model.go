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

// grafana告警结构体
/*
{
  "receiver": "My Super Webhook",
  "status": "firing",
  "orgId": 1,
  "alerts": [
    {
      "status": "firing",
      "labels": {
        "alertname": "High memory usage",
        "team": "blue",
        "zone": "us-1"
      },
      "annotations": {
        "description": "The system has high memory usage",
        "runbook_url": "https://myrunbook.com/runbook/1234",
        "summary": "This alert was triggered for zone us-1"
      },
      "startsAt": "2021-10-12T09:51:03.157076+02:00",
      "endsAt": "0001-01-01T00:00:00Z",
      "generatorURL": "https://play.grafana.org/alerting/1afz29v7z/edit",
      "fingerprint": "c6eadffa33fcdf37",
      "silenceURL": "https://play.grafana.org/alerting/silence/new?alertmanager=grafana&matchers=alertname%3DT2%2Cteam%3Dblue%2Czone%3Dus-1",
      "dashboardURL": "",
      "panelURL": "",
      "values": {
        "B": 44.23943737541908,
        "C": 1
      }
    },
    {
      "status": "firing",
      "labels": {
        "alertname": "High CPU usage",
        "team": "blue",
        "zone": "eu-1"
      },
      "annotations": {
        "description": "The system has high CPU usage",
        "runbook_url": "https://myrunbook.com/runbook/1234",
        "summary": "This alert was triggered for zone eu-1"
      },
      "startsAt": "2021-10-12T09:56:03.157076+02:00",
      "endsAt": "0001-01-01T00:00:00Z",
      "generatorURL": "https://play.grafana.org/alerting/d1rdpdv7k/edit",
      "fingerprint": "bc97ff14869b13e3",
      "silenceURL": "https://play.grafana.org/alerting/silence/new?alertmanager=grafana&matchers=alertname%3DT1%2Cteam%3Dblue%2Czone%3Deu-1",
      "dashboardURL": "",
      "panelURL": "",
      "values": {
        "B": 44.23943737541908,
        "C": 1
      }
    }
  ],
  "groupLabels": {},
  "commonLabels": {
    "team": "blue"
  },
  "commonAnnotations": {},
  "externalURL": "https://play.grafana.org/",
  "version": "1",
  "groupKey": "{}:{}",
  "truncatedAlerts": 0,
  "title": "[FIRING:2]  (blue)",
  "state": "alerting",
  "message": "**Firing**\n\nLabels:\n - alertname = T2\n - team = blue\n - zone = us-1\nAnnotations:\n - description = This is the alert rule checking the second system\n - runbook_url = https://myrunbook.com\n - summary = This is my summary\nSource: https://play.grafana.org/alerting/1afz29v7z/edit\nSilence: https://play.grafana.org/alerting/silence/new?alertmanager=grafana&matchers=alertname%3DT2%2Cteam%3Dblue%2Czone%3Dus-1\n\nLabels:\n - alertname = T1\n - team = blue\n - zone = eu-1\nAnnotations:\nSource: https://play.grafana.org/alerting/d1rdpdv7k/edit\nSilence: https://play.grafana.org/alerting/silence/new?alertmanager=grafana&matchers=alertname%3DT1%2Cteam%3Dblue%2Czone%3Deu-1\n"
}
*/

// Body grafana webhook request body. 详细见：https://grafana.com/docs/grafana/latest/alerting/configure-notifications/manage-contact-points/integrations/webhook-notifier/
type Body struct {
	Receiver        string  `json:"receiver"`
	Status          string  `json:"status"`
	OrgId           int     `json:"orgId"`
	Alerts          []Alert `json:"alerts"`
	ExternalURL     string  `json:"externalURL"`
	Version         string  `json:"version"`
	GroupKey        string  `json:"groupKey"`
	TruncatedAlerts int     `json:"truncatedAlerts"`
	Title           string  `json:"title"`
	State           string  `json:"state"`
	Message         string  `json:"message"`
}

/*{
"alerts": [
{
"annotations": {
"description": "aviation带宽持续满载",
"summary": "aviation带宽持续满载，可能有恶意流量，请手动排查。"
},
"dashboardURL": "https://grafana.redrock.team/d/cdn592h0ghzwgf?orgId=1",
"endsAt": "0001-01-01T00:00:00Z",
"fingerprint": "1894ec986fe5d383",
"generatorURL": "https://grafana.redrock.team/alerting/grafana/adn3239c7vawwc/view?orgId=1",
"labels": {
"alertname": "aviation telecom带宽持续满载",
"device": "telecom",
"grafana_folder": "mirror",
"instance": "198.18.114.2:9100",
"job": "node-exporter",
"namespace": "mirror"
},
"panelURL": "https://grafana.redrock.team/d/cdn592h0ghzwgf?orgId=1&viewPanel=9",
"silenceURL": "https://grafana.redrock.team/alerting/silence/new?alertmanager=grafana&matcher=alertname%3Daviation+telecom%E5%B8%A6%E5%AE%BD%E6%8C%81%E7%BB%AD%E6%BB%A1%E8%BD%BD&matcher=device%3Dtelecom&matcher=grafana_folder%3Dmirror&matcher=instance%3D198.18.114.2%3A9100&matcher=job%3Dnode-exporter&matcher=namespace%3Dmirror&orgId=1",
"startsAt": "2024-05-30T10:25:30Z",
"status": "firing",
"valueString": "[ var='A' labels={device=telecom, instance=198.18.114.2:9100, job=node-exporter} value=5.091481347669635 ], [ var='C' labels={device=telecom, instance=198.18.114.2:9100, job=node-exporter} value=1 ]",
"values": {
"A": 5.091481347669635,
"C": 1
}
}
],
"commonAnnotations": {
"description": "aviation带宽持续满载",
"summary": "aviation带宽持续满载，可能有恶意流量，请手动排查。"
},
"commonLabels": {
"alertname": "aviation telecom带宽持续满载",
"device": "telecom",
"grafana_folder": "mirror",
"instance": "198.18.114.2:9100",
"job": "node-exporter",
"namespace": "mirror"
},
"externalURL": "https://grafana.redrock.team/",
"groupKey": "{}/{namespace=\"mirror\"}:{alertname=\"aviation telecom带宽持续满载\", grafana_folder=\"mirror\"}",
"groupLabels": {
"alertname": "aviation telecom带宽持续满载",
"grafana_folder": "mirror"
},
"message": "Firing\n\nValue: A=5.091481347669635, C=1\nLabels:\n - alertname = aviation telecom带宽持续满载\n - device = telecom\n - grafana_folder = mirror\n - instance = 198.18.114.2:9100\n - job = node-exporter\n - namespace = mirror\nAnnotations:\n - description = aviation带宽持续满载\n - summary = aviation带宽持续满载，可能有恶意流量，请手动排查。\nSource: https://grafana.redrock.team/alerting/grafana/adn3239c7vawwc/view?orgId=1\nSilence: https://grafana.redrock.team/alerting/silence/new?alertmanager=grafana&matcher=alertname%3Daviation+telecom%E5%B8%A6%E5%AE%BD%E6%8C%81%E7%BB%AD%E6%BB%A1%E8%BD%BD&matcher=device%3Dtelecom&matcher=grafana_folder%3Dmirror&matcher=instance%3D198.18.114.2%3A9100&matcher=job%3Dnode-exporter&matcher=namespace%3Dmirror&orgId=1\nDashboard: https://grafana.redrock.team/d/cdn592h0ghzwgf?orgId=1\nPanel: https://grafana.redrock.team/d/cdn592h0ghzwgf?orgId=1&viewPanel=9\n",
"orgId": 1,
"receiver": "mirror-test",
"state": "alerting",
"status": "firing",
"title": "[FIRING:1] aviation telecom带宽持续满载 mirror (telecom 198.18.114.2:9100 node-exporter mirror)",
"truncatedAlerts": 0,
"version": "1"
}
*/

type Alert struct {
	Annotations  map[string]string `json:"annotations"`
	DashboardURL string            `json:"dashboardURL,omitempty"`
	EndsAt       string            `json:"endsAt"`
	Fingerprint  string            `json:"fingerprint"`
	GeneratorURL string            `json:"generatorURL,omitempty"`
	Labels       struct {
		AlertName string `json:"alertname"`
	} `json:"labels,omitempty"`
	PanelURL    string             `json:"panelURL,omitempty"`
	SilenceURL  string             `json:"silenceURL"`
	ImageURL    string             `json:"imageURL,omitempty"`
	StartsAt    string             `json:"startsAt"`
	Status      string             `json:"status"`
	ValueString string             `json:"valueString,omitempty"`
	Values      map[string]float64 `json:"values,omitempty"`
}

// 查询历史消息
type Message struct {
	MessageID  string `json:"message_id"`
	CreateTime string `json:"create_time"`
	Sender     struct {
		ID string `json:"id"`
	} `json:"sender"`
}

type MessageListResponse struct {
	Data struct {
		Items []Message `json:"items"`
	} `json:"data"`
}
