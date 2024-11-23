package grafana

import (
	"EasyBanner/model"
	"errors"
	"log"
)

func GetAlertName(body *model.Body) (string, error) {
	if len(body.Alerts) == 0 {
		log.Println("Alerts is empty!")
		return "", errors.New("Alerts is empty!")
	}
	for _, alert := range body.Alerts {
		if alert.Labels.AlertName == "" {
			log.Println("alertname is empty!")
			return "", errors.New("alertname is empty!")
		}
		return alert.Labels.AlertName, nil
	}
	return "", errors.New("No valid alertname found in Alerts!")
}
