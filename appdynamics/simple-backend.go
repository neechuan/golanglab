package main

import (
	appd "appdynamics"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
"strconv"
"time"
"math/rand"
"os"
)

func main() {
    initAgent()

    maxBtCount := 500
    btCount := 0

    for btCount < maxBtCount {
        s2 := strconv.Itoa(btCount)
        fmt.Printf("%v ", s2)
        time.Sleep(time.Duration(getRandomMilliseconds()) * time.Millisecond)
        btCount++

        // Create Backend connection details
        backendName := "EUMAdrum"
        backendType := "HTTP"
        backendProperties := map[string]string{"URL": "https://cdn.appdynamics.com/adrum/adrum-latest.js"}
        resolveBackend := false
        appd.AddBackend(backendName, backendType, backendProperties, resolveBackend)

        // Start AppD BT
        btHandle2 := appd.StartBT("CheckAdrum", "")
        ecHandle := appd.StartExitcall(btHandle2, backendName)
        hdr := appd.GetExitcallCorrelationHeader(ecHandle)
        response, err := http.Get("https://cdn.appdynamics.com/adrum/adrum-latest.js")

        // Show the Exit call correlation header
        fmt.Printf("===== Show exit call correlation header =======\n")
        fmt.Print("hdr: ")
	fmt.Println(hdr)

        // Add Error to BT
        if response.StatusCode != 200 {
		fmt.Println("AddBTError")
                appd.AddBTError(btHandle2, appd.APPD_LEVEL_ERROR, response.Status, true)
                appd.AddUserDataToBT(btHandle2, "Error", "EUM adrum is unreachable")
                appd.SetBTURL(btHandle2, "/checkadrum")
        }
        if err != nil {
		fmt.Print("error: ")
                fmt.Println(err.Error())
                os.Exit(1)
        }
        responseData, err := ioutil.ReadAll(response.Body)
        if err != nil {
                fmt.Print("error2: ")
		fmt.Println(err.Error())
        }
        fmt.Printf("\n=====================================\n\n")
	fmt.Println("response data: ")
        fmt.Println(string(responseData))
        appd.EndExitcall(ecHandle)
	setSnapshotAttributes(btHandle2, "label", "blue")
        appd.EndBT(btHandle2)
    }
    // Stop/Clean up the AppD SDK.
    appd.TerminateSDK()
}

func initAgent() {
        fmt.Println("Arch Check")

        // Exec ldd test
        ldd_cmd := exec.Command("ldd", "--version")

        var ldd_out bytes.Buffer
        ldd_cmd.Stdout = &ldd_out

        ldd_err := ldd_cmd.Run()
        if ldd_err != nil {
                log.Fatal(ldd_err)
        }

        fmt.Printf("ldd --version:\n")
        fmt.Printf("%s\n", ldd_out.String())

        // Exec uname test
        uname_cmd := exec.Command("uname", "-a")

        var uname_out bytes.Buffer
        uname_cmd.Stdout = &uname_out

        uname_err := uname_cmd.Run()
        if uname_err != nil {
                log.Fatal(uname_err)
        }

        fmt.Printf("uname -a:\n")
        fmt.Printf("%s\n", uname_out.String())

        // Configure AppD
        cfg := appd.Config{}

        // Controller
        cfg.Controller.Host = "ip-10-97-28-69.us-west-2.compute.internal"
        cfg.Controller.Port = 8090
        cfg.Controller.UseSSL = false
        cfg.Controller.Account = "customer1"
        cfg.Controller.AccessKey = "f7161983-e13b-4585-bd9d-abb885e852aa"

        // App Context
        cfg.AppName = "golandlab_simple"
        cfg.TierName = "GolangTier1"
        cfg.NodeName = "GolangNode1"

        // misc
        cfg.InitTimeoutMs = 1000

        // init the SDK
        if err := appd.InitSDK(&cfg); err != nil {
                fmt.Printf("Error initializing the AppDynamics SDK\n")
        } else {
                fmt.Printf("Initialized AppDynamics SDK successfully\n")
        }
}

func setSnapshotAttributes(bt appd.BtHandle, key, value string) {
        if appd.IsBTSnapshotting(bt) {
                appd.AddUserDataToBT(bt, key, value)
                appd.SetBTURL(bt, "/checkadrum")
        }
}

func getRandomMilliseconds() int {
    rand.Seed(time.Now().UnixNano())
    min := 10
    max := 3000
    return rand.Intn(max - min + 1) + min
}
