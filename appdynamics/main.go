/*
Copyright 2014 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// outyet is a web server that announces whether or not a particular Go version
// has been tagged.
package main

import (
    "fmt"
	"net/http"
	appd "appdynamics"
)

func homePage(w http.ResponseWriter, r *http.Request){
	btHandle := appd.StartBT("/", "")
    fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
	appd.EndBT(btHandle)
}

func homePage1(w http.ResponseWriter, r *http.Request){
	btHandle := appd.StartBT("/one", "")
    fmt.Fprintf(w, "Welcome to the HomePage1!")
	fmt.Println("Endpoint Hit: homePage1")
	appd.EndBT(btHandle)
}

func handleRequests() {
    http.HandleFunc("/", homePage)
    http.HandleFunc("/one", homePage1)
    http.ListenAndServe(":8080", nil)
}

func main() {
    cfg := appd.Config{}
 	// Configure AppD
	// Controller
	//cfg.Controller.Host = "ec2-52-37-30-238.us-west-2.compute.amazonaws.com"
        cfg.Controller.Host = "ip-10-97-28-69.us-west-2.compute.internal"
	cfg.Controller.Port = 8090
	cfg.Controller.UseSSL = false
	cfg.Controller.Account = "customer1"
	cfg.Controller.AccessKey = "f7161983-e13b-4585-bd9d-abb885e852aa"

	// App Context
	cfg.AppName = "golanglab_main"
	cfg.TierName = "Golang"
	cfg.NodeName = "Golang1"

	// misc
	cfg.InitTimeoutMs = 1000
	cfg.Logging.MinimumLevel =appd.APPD_LOG_LEVEL_TRACE
	

	// init the SDK
	if err := appd.InitSDK(&cfg); err != nil {
		fmt.Printf("Error initializing the AppDynamics SDK\n")
	} else {
		fmt.Printf("Initialized AppDynamics SDK successfully\n")
	}
     
	 handleRequests()
}


