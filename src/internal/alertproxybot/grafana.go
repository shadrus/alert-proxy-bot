package alertproxybot

import (
	"AlertProxyBot/src/internal/senders"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"net/http"
)

/*
{
  "dashboardId": 1,
  "evalMatches": [
    {
      "value": 1,
      "metric": "Count",
      "tags": {}
    }
  ],
  "imageUrl": "https://grafana.com/assets/img/blog/mixed_styles.png",
  "message": "Notification Message",
  "orgId": 1,
  "panelId": 2,
  "ruleId": 1,
  "ruleName": "Panel Title alert",
  "ruleUrl": "http://localhost:3000/d/hZ7BuVbWz/test-dashboard?fullscreen\u0026edit\u0026tab=alert\u0026panelId=2\u0026orgId=1",
  "state": "alerting",
  "tags": {
    "tag name": "tag value"
  },
  "title": "[Alerting] Panel Title alert"
}
*/

type grafanaRequest struct {
	DashboardId int `json:"dashboardId"`
	EvalMatches []struct {
		Value  int    `json:"value"`
		Metric string `json:"metric"`
		Tags   struct {
		} `json:"tags"`
	} `json:"evalMatches"`
	ImageUrl string `json:"imageUrl"`
	Message  string `json:"message"`
	OrgId    int    `json:"orgId"`
	PanelId  int    `json:"panelId"`
	RuleId   int    `json:"ruleId"`
	RuleName string `json:"ruleName"`
	RuleUrl  string `json:"ruleUrl"`
	State    string `json:"state"`
	Tags     struct {
		TagName string `json:"tag name"`
	} `json:"tags"`
	Title string `json:"title"`
}

func (request grafanaRequest) String() string {
	return fmt.Sprintf("Title: %s\r\nState: %s\r\nMessage: %s\r\n", request.Title, request.State, request.Message)
}

func serveGrafanaAlert(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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
		var grafana grafanaRequest
		err := json.NewDecoder(request.Body).Decode(&grafana)
		if err != nil {
			log.WithField("source", "jsonDecoder").Error(err)
			writeErrorResponse(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if grafana.Title == ""{
			log.WithField("source", "jsonDecoder").Error(err)
			writeErrorResponse(writer, "Not valid request body", http.StatusBadRequest)
		}
		log.Debug(grafana)
		err = sender.Send(grafana.String(), chatId)
		if err != nil {
			log.WithField("source", "telegram").Error(err)
		}
		writer.WriteHeader(200)
	}
}
