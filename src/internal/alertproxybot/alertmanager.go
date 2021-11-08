package alertproxybot

import (
	"AlertProxyBot/src/internal/senders"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type alert struct {
	Status string `json:"status"`
	Labels struct {
		Alertname  string `json:"alertname"`
		Prometheus string `json:"prometheus"`
		Severity   string `json:"severity"`
	} `json:"labels"`
	Annotations struct {
		Description string `json:"description"`
		RunbookUrl  string `json:"runbook_url"`
		Summary     string `json:"summary"`
	} `json:"annotations"`
	StartsAt     time.Time `json:"startsAt"`
	EndsAt       time.Time `json:"endsAt"`
	GeneratorURL string    `json:"generatorURL"`
	Fingerprint  string    `json:"fingerprint"`
}

func (request alert) String() string {
	return fmt.Sprintf("Title: %s\r\nState: %s\r\nMessage: %s\r\n", request.Labels.Alertname, request.Status, request.Annotations.Description)
}

type alertManagerRequest struct {
	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Alerts   []alert `json:"alerts"`
	GroupLabels struct {
	} `json:"groupLabels"`
	CommonLabels struct {
		Prometheus string `json:"prometheus"`
	} `json:"commonLabels"`
	CommonAnnotations struct {
	} `json:"commonAnnotations"`
	ExternalURL     string `json:"externalURL"`
	Version         string `json:"version"`
	GroupKey        string `json:"groupKey"`
	TruncatedAlerts int    `json:"truncatedAlerts"`
}



func serveAlertManagerAlert(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	chatId := request.URL.Query().Get("chat_id")
	log.Debugf("chat_id, %s!", request.URL.Query().Get("chat_id"))
	if chatId == "" {
		writeErrorResponse(writer, "param chat_id was not found in this request", http.StatusInternalServerError)
		return
	}
	target := request.URL.Query().Get("target")
	if target == "" {
		writeErrorResponse(writer, "param target was not found in this request", http.StatusInternalServerError)
		return
	}
	sender, err := senders.NewAlertSender(request.URL.Query().Get("target"))
	if err != nil {
		log.Error(err)
		writeErrorResponse(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	} else {
		var alertManager alertManagerRequest
		err := json.NewDecoder(request.Body).Decode(&alertManager)
		if err != nil {
			log.WithField("source", "jsonDecoder").Error(err)
			writeErrorResponse(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if alertManager.Status == ""{
			log.WithField("source", "jsonDecoder").Error(err)
			writeErrorResponse(writer, "Not valid request body", http.StatusBadRequest)
			return
		}
		log.Debug(alertManager)
		for _, a := range alertManager.Alerts {
			err = sender.Send(a.String(), chatId)
			if err != nil {
				log.WithField("source", "telegram").Error(err)
			}
		}
		writer.WriteHeader(200)
	}
}
