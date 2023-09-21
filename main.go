package main

import (
	"encoding/json"
	"fmt"
	"io"
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

	
				payload , err := ioutil.ReadAll(r.Body)
				if err != nil {
					fmt.Println(err)
					return
				}
		
				if value[0]["is_root"] == "true"{
					transaction.IsRoot = true
	
				}else{

					transaction.CausationId = value[0]["causation_id"]

					transaction.Correlation_id = value[0]["correlation_id"]
			 
				}

 
				transaction.Is_Command = value[0]["is_command"]
			    transaction.Is_Query = value[0]["is_query"]
				transaction.Payload = string(payload) 

				db.Insert(transaction)
				
				req, err := http.NewRequest("", redirectURL, nil)
 
				resp, err  = client.Do(req)
				if err != nil {
					fmt.Println(err)
					return
				}
			
				w.WriteHeader(resp.StatusCode)
				for key, values := range resp.Header {
					for _, value := range values {
						w.Header().Add(key, value)
				 
					}
				}
				_, err = io.Copy(w, resp.Body)
	
	 
			
			})

			
			hlr.ServeHTTP(r,w)
			break	
	 
		}else{
			fmt.Errorf("Not found")
		}

		
	 
   }
 
 
 
}
 
 
func main(){
	var handlers Handlers
  
 
 
	// Start the HTTP server on port 8080
	err := http.ListenAndServe(":8080", handlers)
	if err != nil {
		panic(err)
	}

}

 