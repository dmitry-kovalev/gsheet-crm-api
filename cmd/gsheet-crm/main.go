package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/dmitry-kovalev/internal/gsheet-crm-api/googleclient"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger
	gc  *googleclient.GoogleClient
)

func processQuery(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	log.Infof("GET with params %v", values)
	w.Header().Add("Content-Type", "text/html; charset=UTF-8")

	spreadsheetID := values.Get("spreadsheetID")
	if spreadsheetID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`Missed spreadsheetID`))
		return
	}

	cellsRange := values.Get("cellsRange")
	if cellsRange == "" {
		cellsRange = "A2:B"
	}
	matched, err := regexp.MatchString(`^(.{1,}\!){0,1}([A-Z]+[0-9]+:[A-Z]+)$`, cellsRange)
	if err != nil || !matched {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`Invalid cells range`))
		return
	}

	phone := values.Get("phone")
	if phone == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`Missed phone number`))
		return
	}

	nameCol_ := values.Get("nameCol")
	nameCol, err := strconv.Atoi(nameCol_)
	if err != nil || nameCol_ == "" {
		nameCol = 0
	}

	phoneCol_ := values.Get("phoneCol")
	phoneCol, err := strconv.Atoi(phoneCol_)
	if err != nil || phoneCol_ == "" {
		phoneCol = 1
	}

	phoneMask_ := values.Get("phoneMask")
	phoneMask, err := strconv.Atoi(phoneMask_)
	if err != nil || phoneMask_ == "" {
		phoneMask = 10
	}
	if phoneMask > len(phone) {
		phoneMask = len(phone)
	}
	phone = phone[len(phone)-phoneMask:]

	data, err := gc.Query(spreadsheetID, cellsRange)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`Empty data`))
		return
	}
	for _, row := range data {
		if len(row) <= phoneCol {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`Invalid number of phone column`))
			return
		}
		if len(row) <= nameCol {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`Invalid number of name column`))
			return
		}
		val := fmt.Sprintf("%v", row[phoneCol])
		log.Infof("%v", val)
		if strings.HasSuffix(val, phone) {
			log.Infof("Found customer %v with phone %v", data, val)
			w.WriteHeader(http.StatusOK)
			res := fmt.Sprintf("Hello %v", row[nameCol])
			w.Write([]byte(res))
			return
		}
	}
	log.Infof("Didn't find any customer with phone %v", phone)
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`Not found`))
}

func main() {
	var err error
	log = logrus.New()
	log.SetOutput(os.Stdout)

	log.Info("Starting the app...")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Port is not set")
	}

	gc, err = googleclient.Init(log)
	if err != nil {
		log.Fatal("Can't connect to Google API because %v", err)
		return
	}

	router := mux.NewRouter()
	router.HandleFunc("/", processQuery).Methods("GET")

	serv := http.Server{
		Addr:    net.JoinHostPort("", port),
		Handler: router,
	}
	go serv.ListenAndServe()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt

	log.Info("Stopping app...")

	timeout, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	err = serv.Shutdown(timeout)
	if err != nil {
		log.Error("Error when shutdown app: %v", err)
	}

	log.Info("The app stopped")
}
