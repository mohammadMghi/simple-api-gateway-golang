package main

import (
	"encoding/json"
	"fmt"
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
					transaction.Is_root = true
				}else{

					transaction.RootId = value[0]["root_id"]

					transaction.Is_root = false
				}

 

				transaction.Payload = string(payload) 

				db.Insert(transaction)
				
	 
				http.Redirect(w, r, redirectURL, http.StatusFound)

				 
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

 