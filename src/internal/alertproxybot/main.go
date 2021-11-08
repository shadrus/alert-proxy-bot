package alertproxybot

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type ResponseDetail struct {
	Detail string `json:"detail"`
}

func writeErrorResponse (writer http.ResponseWriter, message string, status int){
	js, err := json.Marshal(ResponseDetail{Detail: message})
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	writer.Write(js)
}

func StartServer() {
	router := httprouter.New()
	router.POST("/alert/grafana/*params", serveGrafanaAlert)
	router.POST("/alert/alertmanager/*params", serveAlertManagerAlert)
	log.Fatal(http.ListenAndServe(":8080", router))
}
