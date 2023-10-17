package main

import (
	"bytes"
	"context"
	"strconv"

	"encoding/json"
	"fmt"
	"io"

	"io/ioutil"
	"log"

	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/go-redis/redis/v8"
	"github.com/mohammadMghi/apiGolangGateway/db"
	"github.com/mohammadMghi/apiGolangGateway/models"
)


type Handlers struct{
	http.Handler
 
}

 

type ResponseNode struct{
	CorrelationId int64 `json"correlation_id"`
	Causationid int64 `json"causation_id"`
	Status string `json"status"`
	EntitiyType string `json"entitiy_type"`
	Payload interface{} `json:"payload"`
 
}

 
type Config struct{
	Cfg []map[string][]Node `json:"config"`
} 
 
type Node struct{
	Balancing string `json:"balancing"`
	Sender string	`json:"sender"`
	AuthRequired bool `json:"auth_required"`
	Targets []map[string]string`json:"targets"`
	Next int `json:"next"`
}

var greenLog = color.New(color.FgGreen).PrintfFunc()
var redLog  = color.New(color.FgRed).PrintfFunc()
var config Config 
func (h  Handlers)ServeHTTP(r http.ResponseWriter, w  *http.Request){

 
	
 
	client := redis.NewClient(&redis.Options{
        Addr:	  "localhost:6379",
        Password: "", // no password set
        DB:		  0,  // use default DB
    })

    ctx := context.Background()


 
 
 
	var resp *http.Response
 
 
	file, err := os.Open("servers.json")
	if err != nil {
		redLog(err.Error())
		log.Fatal(err)
	}

	 
 
	jsonBytes, err := ioutil.ReadAll(file)


	
	if err != nil {
		redLog(err.Error())
		log.Fatal(err)
	}
	defer file.Close()

 
 
	err  = json.Unmarshal(jsonBytes, &config)

	if err != nil {
		redLog(err.Error())
		return
	}
 
 
	greenLog("\n \n",config)
	var redisSenderCount = 0
	for _ , nodes := range config.Cfg[0]{

		for _ , node := range nodes{
			 
	
		

			sender := node.Sender
 
 
	 

			if w.URL.Path == sender {

				tagetLen := len(node.Targets) -1
				val, err := client.Get(ctx, "sender").Result()
			
				if err == redis.Nil  {
					// if sender dosnt existed in redis this code will be execute
					err := client.Set(ctx, "sender", 0, 0).Err()
					if err != nil {
						panic(err)
					}
				} 
			 
		
		 

				if node.Balancing == "roundrobin"{
					greenLog("roundrobin")

					redisSenderCount, err = strconv.Atoi(val)
					if err != nil {
						redLog(err.Error())
						panic(err)
					}
 
					if  tagetLen < redisSenderCount{
						redisSenderCount = 0
				 
						err := client.Set(ctx, "sender", 0, 0).Err()
						if err != nil {
							redLog(err.Error())
							panic(err)
						}
					}
					 
				}else{
					redisSenderCount = 0
				}

				
		 
				targets := node.Targets

				authRequired := node.AuthRequired
	
	
			
	
				hlr :=http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
 
 
					redirectURL := targets[redisSenderCount  ]["url"]
					greenLog("\n \n","Target :: " + redirectURL)
		
					if node.Balancing == "roundrobin"{
						err := client.Set(ctx, "sender", redisSenderCount+1, 0).Err()
						if err != nil {
							redLog(err.Error())
							panic(err)
						}
					}
		
		 
				
	
					if node.AuthRequired == authRequired {
	
						authorizationHeader := r.Header.Get("Authorization")
					 
						if authorizationHeader == "" {
							http.Error(w, "Unauthorized", http.StatusUnauthorized)
							return
						}
				
					} 
					
					data := MakeResponse(r)
			
				
					m , e:= json.Marshal(data)
	
					if e != nil {
						redLog(err.Error())
						log.Fatal(e)
					}

 	
		 
					req, err := http.NewRequest("", redirectURL, bytes.NewBuffer(m))
	 
					if err != nil {
						redLog(err.Error())
						log.Fatal(err)
					}
	
	
					client := &http.Client{}
	
	
					resp, err  = client.Do(req)
					if err != nil {
						redLog(err.Error())
			 
						return
					}
	
	
					for key, values := range resp.Header {
						for _, value := range values {
							w.Header().Add(key, value)
					 
						}
					}
					
		
			 
					if err != nil {
						redLog(err.Error())
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
	
	
				
					InsertCorAndCuse(respNode)
	
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

 
			} 
 
			node.Next ++
	
		}
	} 
}

func InsertCorAndCuse(respNode map[string][]ResponseNode ){
	responseLen := len(respNode) -1
	for _ , value :=range respNode{
 

		body, _ := json.Marshal(value[responseLen].Payload)

	
		var tr models.Transaction
		db.Insert(&tr)
	 
		

		tr.Correlation_id = int64(value[responseLen].CorrelationId)
		tr.CausationId = int64(value[responseLen].Causationid)
		tr.EntitiyType = string(value[responseLen].EntitiyType)
		
		tr.Message =  string(body) 
		db.Update(&tr)
		responseLen++
	} 
}

func MakeResponse(r *http.Request) map[string]interface{}{
		 
	var transaction models.Transaction

	transaction , jsonString :=  InsertTransaction(r)


	return  map[string]interface{}{
		"correlation_id":    transaction.ID,
		"causation_id":   transaction.ID,
		"payload": map[string]interface{}{
			"Message": jsonString,
		},
	}
}
func InitSenderRedis(ctx context.Context,client redis.Client ) string{
	sender, err := client.Get(ctx, "sender").Result()
			
	if err == redis.Nil  {
		// if sender dosnt existed in redis this code will be execute
		err := client.Set(ctx, "sender", 0, 0).Err()
		if err != nil {
			panic(err)
		}
	} 

	return sender
}

func InsertTransaction(r *http.Request) (models.Transaction , string){
	var transaction models.Transaction
		 
	transaction.Status = "IN_PROCESS"


	body , err :=ioutil.ReadAll(r.Body)

	if err != nil {
		redLog(err.Error())
		log.Fatal(err)
	}
	jsonString := string(body)
	transaction.Message =jsonString

	db.Insert(&transaction)



	transaction.Correlation_id =  *transaction.ID 
	transaction.CausationId =    *transaction.ID 

	
	db.Update(&transaction)

	return transaction , jsonString
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
	greenLog := color.New(color.FgGreen).PrintfFunc()

	greenLog("Gateway is running ...")
	
 
	// Start the HTTP server on port 8080
	err := http.ListenAndServe(":8080", handlers)
	if err != nil {
		
		panic(err)
	}

}

 