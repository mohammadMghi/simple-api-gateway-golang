package main

import (
	"bytes"
 
	"encoding/json"
	"fmt"
	"io"
	"time"

	"io/ioutil"
	"log"

	"net/http"
	"os"
 
	"github.com/mohammadMghi/apiGolangGateway/db"
	"github.com/mohammadMghi/apiGolangGateway/models"
)


type Handlers struct{
	http.Handler
 
}

 

type ResponseNode struct{
	Correlation_id int64 `json"correlation_id"`
	Causation_id int64 `json"causation_id"`
	Status string `json"status"`
	Payload interface{} `json:"payload"`
 
}

 
type Config struct{
	Cfg []map[string][]Node `json:"config"`
} 
 
type Node struct{
	Balancing string `json:"balancing"`
	Sender string	`json:"sender"`
	AuthRequired bool `json:"auth_required"`
	Targets []map[string]interface{} `json:"targets"`
}


func (h  Handlers)ServeHTTP(r http.ResponseWriter, w  *http.Request){

 
	var resp *http.Response
 
 
	file, err := os.Open("servers.json")
	if err != nil {
		log.Fatal(err)
	}

	 
 
	jsonBytes, err := ioutil.ReadAll(file)


	
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

 
 
	var config Config
	

 

	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
		select {
			case <- ticker.C:
				 
				


			case <- quit:
				ticker.Stop()
				return
			}
		}
	}()





	services  := make(map[string][]map[string][]map[string]string) 
 
	err  = json.Unmarshal(jsonBytes, &config)

 
	fmt.Println(config)

	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	var previousRequestPath = ""

	for serviceName, serviceList := range services  {
		fmt.Println("Service Name:", serviceName)
		for index, services := range serviceList {
			// Loop through the slice of maps
 
 

	 

			for _, pods := range services {

					println( index)
					
					var pod = pods[index]

					sender := pod["sender"]

					sender = previousRequestPath

					receiver := pod["receiver"]

					authRequired := pod["authRequired"]
		 




					//new request is like pre request

					//check in redis and latest request if same then count 1 and send request to the next
 

	 


















					if w.URL.Path == sender {
				
			
						var transaction models.Transaction
			
						hlr :=http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					 
							fmt.Println(sender)
			
							redirectURL := receiver
			
						
			
							if authRequired== "true" {
			
								authorizationHeader := r.Header.Get("Authorization")
							 
								if authorizationHeader == "" {
									http.Error(w, "Unauthorized", http.StatusUnauthorized)
									return
								}
						
							} 
			
				 
			 
				 
							transaction.Status = "IN_PROCESS"
					 
				
							
						
			
							body , err :=ioutil.ReadAll(r.Body)
					 
							if err != nil {
								log.Fatal(err)
							}
							jsonString := string(body)
							transaction.Message =jsonString
			
							db.Insert(&transaction)
							transaction.Correlation_id =  *transaction.ID 
							transaction.CausationId =    *transaction.ID 
			
							
							db.Update(&transaction)
							data := map[string]interface{}{
								"correlation_id":    transaction.ID,
								"causation_id":   transaction.ID,
								"payload": map[string]interface{}{
									"Message": jsonString,
								},
							}
						
							m , e:= json.Marshal(data)
			
							if e != nil {
								log.Fatal(e)
							}
				 
							req, err := http.NewRequest("", redirectURL, bytes.NewBuffer(m))
			 
							if err != nil {
								log.Fatal(err)
							}
			
			
							client := &http.Client{}
			
			
							resp, err  = client.Do(req)
							if err != nil {
								fmt.Println(err)
								return
							}
			
			
							for key, values := range resp.Header {
								for _, value := range values {
									w.Header().Add(key, value)
							 
								}
							}
			
							
					 
							if err != nil {
								log.Fatal(err)
							}
			
							var respNode = make(map[string][]ResponseNode)
							
							bodyResponse , err :=ioutil.ReadAll(resp.Body)
			
							fmt.Printf(string(bodyResponse))
							err =json.Unmarshal(bodyResponse, &respNode)
			
			
							fmt.Print( (respNode))
							if err != nil {
								fmt.Println(err)
								return
							}
			
			
							responseLen := len(respNode) -1
							for _ , value :=range respNode{
						 
			  
								b, _ := json.Marshal(value[responseLen].Payload)
			
							
								var tr models.Transaction
								db.Insert(&tr)
			
			
								tr.Correlation_id = int64(value[responseLen].Correlation_id)
								tr.CausationId = int64(value[responseLen].Causation_id)
								tr.Message =  string(b) 
								db.Update(&tr)
								responseLen++
							} 
					 
					 
						 
								 
						 
			
							w.WriteHeader(resp.StatusCode)
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
						 fmt.Println("Error : url not found")
					}
			
					
				 
			   
				}
 
			}
		}
	
 


 
 
}

func normalizationDatabase(){

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

 