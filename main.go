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

type Nodes struct{
 	data	map[string]string
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
	

	//get width of request (size , content , check )


	//monitor

 

	for _ , value := range nodes {
	
		if w.URL.Path == value[0]["sender"] {
		 
			var transaction models.Transaction
			
			http.HandleFunc(value[0]["sender"], func(w http.ResponseWriter, r *http.Request) {
	   
				redirectURL := value[0]["receiver"] 

				 
				payload , err := ioutil.ReadAll(r.Body)
				if err != nil {
					fmt.Println(err)
					return
				}
 

				transaction.Payload = string(payload) 

				db.Insert(transaction)
				
	 
				http.Redirect(w, r, redirectURL, http.StatusFound)
			})
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

 