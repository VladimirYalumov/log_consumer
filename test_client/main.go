package main

import (
	"log_consumer/logger"
	"log_consumer/rabbit"
	"log_consumer/test_client/request"
	"log_consumer/test_client/response"
	"log_consumer/test_client/test_handler"
	"encoding/json"
	"log"
	"net/http"
)

const numbersSumRoute = "numbers_sum"
const numbersDifferenceRoute = "numbers_difference"

func setClient()  {
	var client rabbit.Client
	err := client.SetConnection("guest", "guest", "127.0.0.1:5673")
	if err != nil {
		panic(err)
	}
	logInstance := logger.GetInstance()
	logInstance.InitPusher(client)
}

func main() {
	setClient()
	http.HandleFunc("/" +numbersSumRoute, func(w http.ResponseWriter, r *http.Request) {
		logger.GetInstance().DefinePid(r.Host, numbersSumRoute)
		datas := make(map[string]interface{}) // it is only example, don't be critical))
		datas["msg"] = "start Action: numbers"
		datas["user_ip"] = r.Host
		datas["test_array"] = []int{1,2,3,4}
		delete(datas, "user_ip")

		var request request.NumRequest
		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			datas["msg"] = "Invalid params: " + err.Error()
			_ = logger.GetInstance().Error(datas)

			datas["msg"] = "finish Action: numbers"
			_ = logger.GetInstance().Info(datas)

			w.WriteHeader(500)
			response := response.NumResponse{ErrorMsg: err.Error()}
			_ = json.NewEncoder(w).Encode(&response)
			return
		}

		result := test_handler.AddNumbers(request.FirstNumber, request.SecondNumber)

		w.WriteHeader(200)
		response := response.NumResponse{Result: result, ErrorMsg: ""}
		_ = json.NewEncoder(w).Encode(&response)

		datas["msg"] = "finish Action: numbers"
		_ = logger.GetInstance().Info(datas)

		return
	})

	http.HandleFunc("/" + numbersDifferenceRoute, func(w http.ResponseWriter, r *http.Request) {
		logger.GetInstance().DefinePid(r.Host, numbersDifferenceRoute)
		datas := make(map[string]interface{}) // it is only example, don't be critical)
		datas["msg"] = "start Action: numbers"
		datas["user_ip"] = r.Host
		_ = logger.GetInstance().Info(datas)
		delete(datas, "user_ip")

		var request request.NumRequest
		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			datas["msg"] = "Invalid params: " + err.Error()
			_ = logger.GetInstance().Error(datas)

			datas["msg"] = "finish Action: numbers"
			_ = logger.GetInstance().Info(datas)

			w.WriteHeader(500)
			response := response.NumResponse{ErrorMsg: err.Error()}
			_ = json.NewEncoder(w).Encode(&response)
			return
		}

		result := test_handler.SubtractNumbers(request.FirstNumber, request.SecondNumber)

		w.WriteHeader(200)
		response := response.NumResponse{Result: result, ErrorMsg: ""}
		_ = json.NewEncoder(w).Encode(&response)

		datas["msg"] = "finish Action: numbers"
		_ = logger.GetInstance().Info(datas)

		return
	})

	log.Fatal(http.ListenAndServe(":8089", nil))
}
