package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type page_cache struct {
	Url       string `json:"url"`
	FilePath  string `json:"file_path"`
	Timestamp int64  `json:"timestamp"`
}

type download_request struct {
	Url        string
	FilePath   string
	RetryLimit int
}

func log_errors(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func parse_json_msg(page page_cache) json.RawMessage {
	page_json, err := json.Marshal(page)
	log_errors(err)
	return page_json
}

func get_query_params(w http.ResponseWriter, r *http.Request, param_name string) string {
	value := r.URL.Query().Get(param_name)
	if value == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, param_name+" is required")
		return ""
	}

	return value
}

func send_json_response(w http.ResponseWriter, status int, page page_cache) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	page_json := parse_json_msg(page)
	w.Write(page_json)
}

var page_cache_list []page_cache
var download_queue chan download_request

func setup_worker_pool() {
	download_queue = make(chan download_request, 200)
	for w := 1; w <= 10; w++ {
		go download_worker(download_queue)
	}
}

func download_worker(download_queue <-chan download_request) {
	for download_job := range download_queue {
		download_file(download_job)
	}
}

func download_file(download_job download_request) {
	for i := 0; i < download_job.RetryLimit; i++ {
		file_instance, err := os.Create(download_job.FilePath)
		log_errors(err)
		defer file_instance.Close()
		response, err := http.Get(download_job.Url)
		log_errors(err)
		defer response.Body.Close()
		_, err = io.Copy(file_instance, response.Body)
		if err == nil {
			break
		}
		log_errors(err)
	}
}

func wait_for_download(file_path string) {
	for {
		if _, err := os.Stat(file_path); os.IsNotExist(err) {
			time.Sleep(1000000000)
		} else {
			break
		}
	}
}

func page_source_getter(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	url := get_query_params(w, r, "url")
	retryLimit := get_query_params(w, r, "retry_limit")
	retryLimitInt, err := strconv.Atoi(retryLimit)
	log_errors(err)
	timestamp := now.Unix()
	file_name := strconv.Itoa(len(page_cache_list)+1) + ".html"
	file_path := "downloads/" + file_name
	index := -1
	for ind, page := range page_cache_list {
		if page.Url == url {
			index = ind
		}
	}
	if index != -1 && timestamp-page_cache_list[index].Timestamp < 86400 {
		page_cache_list[index].Timestamp = timestamp
		file_path = page_cache_list[index].FilePath
	} else {
		download_queue <- download_request{Url: url, FilePath: file_path, RetryLimit: retryLimitInt}
		page_cache_list = append(page_cache_list, page_cache{Url: url, FilePath: file_path, Timestamp: timestamp})
		index = len(page_cache_list) - 1
	}
	wait_for_download(file_path)
	send_json_response(w, http.StatusOK, page_cache_list[index])
}

func main() {
	setup_worker_pool()
	http.HandleFunc("/pagesource", page_source_getter)
	log_errors(http.ListenAndServe("localhost:8000", nil))
}
