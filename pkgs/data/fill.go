package data

import (
	"EasyBanner/model"
	"bytes"
	"text/template"
)

// 填充 ban 模板数据，生成动态 JSON
func GenerateTemplate(ipDataList []model.IPData, showButton bool, resultText string) (string, error) {
	const tpl = `{
	"config": {
		"update_multi": true
	},
	"i18n_elements": {
		"zh_cn": [
			{
				"tag": "table",
				"columns": [
					{
						"data_type": "text",
						"name": "ip",
						"display_name": "IP",
						"horizontal_align": "left",
						"width": "auto"
					},
					{
						"data_type": "options",
						"name": "count",
						"display_name": "次数",
						"horizontal_align": "left",
						"width": "auto"
					}
				],
				"rows": [
					{{- $lastIndex := sub (len .IPDataList) 1}}
					{{- range $index, $item := .IPDataList}}
					{
						"ip": "{{.IP}}",
						"count": [
							{
								"text": "{{.Count}}",
								"color": "{{if gt .Count 300}}blue{{else}}green{{end}}"
							}
						]
					}{{if lt $index $lastIndex}},{{end}}
					{{- end}}
				],
				"row_height": "low",
				"header_style": {
					"background_style": "none",
					"bold": true,
					"lines": 1
				},
				"page_size": 5
			},
			{{if .ShowButton}}
			{
				"tag": "action",
				"actions": [
					{
						"tag": "button",
						"text": {
							"tag": "plain_text",
							"content": "BAN"
						},
						"type": "danger",
						"complex_interaction": true,
						"width": "default",
						"size": "medium",
						"behaviors": [
							{
								"type": "callback",
								"value": {
									"action": "ban_ip"
								}
							}
						]
					}
				]
			}
			{{else}}
			{
				"tag": "div",
				"text": {
					"tag": "plain_text",
					"content": "{{.ResultText}}"
				}
			}
			{{end}}
		]
	},
	"i18n_header": {
		"zh_cn": {
			"title": {
				"tag": "plain_text",
				"content": "恶意IP列表"
			},
			"subtitle": {
				"tag": "plain_text",
				"content": "本小时内连续访问次数超过250次的IP"
			},
			"template": "red",
			"ud_icon": {
				"tag": "standard_icon",
				"token": "warning_outlined"
			}
		}
	}
}`

	// 创建模板
	tmpl, err := template.New("jsonTemplate").Funcs(template.FuncMap{
		"sub": func(a, b int) int { return a - b },
	}).Parse(tpl)
	if err != nil {
		return "", err
	}

	// 使用模板生成JSON
	var result bytes.Buffer
	data := struct {
		IPDataList []model.IPData
		ShowButton bool
		ResultText string
	}{
		IPDataList: ipDataList,
		ShowButton: showButton,
		ResultText: resultText,
	}
	if err := tmpl.Execute(&result, data); err != nil {
		return "", err
	}

	return result.String(), nil
}
