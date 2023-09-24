package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"

	"net/http"
	"os"

	"github.com/mohammadMghi/apiGolangGateway/db"
	"github.com/mohammadMghi/apiGolangGateway/models"
)


type Handlers struct{
	http.Handler
 
}

 


func (h  Handlers)ServeHTTP(r http.ResponseWriter, w  *http.Request){

 
	var resp *http.Response
 
	file, err := os.Open("servers.json")
	if err != nil {
		log.Fatal(err)
	}
	var nodes map[string][]map[string]string
	 
 
	jsonBytes, err := ioutil.ReadAll(file)


	
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

 
	json.Unmarshal(jsonBytes, &nodes)
	
	client := &http.Client{}
 

	for _ , value := range nodes {
	
		if w.URL.Path == value[0]["sender"] {
	
			var transaction models.Transaction

			hlr :=http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		 
				fmt.Println(value[0]["sender"])

				redirectURL := value[0]["receiver"] 

			

				if value[0]["auth_required"]== "true" {

					authorizationHeader := r.Header.Get("Authorization")
				 
					if authorizationHeader == "" {
						http.Error(w, "Unauthorized", http.StatusUnauthorized)
						return
					}
			
				} 

	 
 
				transaction.Correlation_id = 1
				transaction.CausationId = 1
	 
				transaction.Status = "IN_PROCESS"
		 
	
				
				req, err := http.NewRequest("", redirectURL, nil)
 

				body , err :=ioutil.ReadAll(r.Body)
		 
				if err != nil {
					log.Fatal(err)
				}
				jsonString := string(body)

				data := map[string]interface{}{
					"correlation_id":    transaction.Correlation_id,
					"causation_id":   transaction.CausationId,
					"payload": map[string]interface{}{
						"Message": jsonString,
					},
				}
			
				transaction.Message =jsonString
				db.Insert(transaction)

				resp, err  = client.Do(req)
				if err != nil {
					fmt.Println(err)
					return
				}


				var response map[string][]map[string]string
			 
			
				w.WriteHeader(resp.StatusCode)
				for key, values := range resp.Header {
					for _, value := range values {
						w.Header().Add(key, value)
				 
					}
				}

				bodyResponse , err :=ioutil.ReadAll(resp.Body)
		 
				if err != nil {
					log.Fatal(err)
				}

				err =json.Unmarshal(bodyResponse, response)
				if err != nil {
					fmt.Println(err)
					return
				}


				for _ , value := range response{
			 

					correlation_id, err := strconv.ParseInt(value[0]["correlation_id"], 10, 64)

					if err != nil {
						fmt.Println("Error:", err)
						return
					}
					causation_id, err := strconv.ParseInt(value[0]["causation_id"], 10, 64)

					if err != nil {
						fmt.Println("Error:", err)
						return
					}
			 
					transaction.Message =  value[0]["payload"]



					var tr models.Transaction
					tr.Correlation_id =correlation_id
					tr.CausationId = causation_id
					db.Insert(transaction)
				} 

				_, err = io.Copy(w, resp.Body)
				jsonResp, err := json.Marshal(data)
				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				}
				w.Write(jsonResp)
			
			})

			
			hlr.ServeHTTP(r,w)

			
			break	
	 
		}else{
			fmt.Errorf("Not found")
		}

		
	 
   }
 
 
 
}
 
func ReadUserIP(r *http.Request) string {
    IPAddress := r.Header.Get("X-Real-Ip")
    if IPAddress == "" {
        IPAddress = r.Header.Get("X-Forwarded-For")
    }
    if IPAddress == "" {
        IPAddress = r.RemoteAddr
    }
    return IPAddress
}
func main(){
	var handlers Handlers
  
 
 
	// Start the HTTP server on port 8080
	err := http.ListenAndServe(":8080", handlers)
	if err != nil {
		panic(err)
	}

}

 