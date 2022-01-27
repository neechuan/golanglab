package main

import (
	"fmt"
	"os/exec"
	"bytes"
	"log"
	appd "appdynamics"
	"time"
	"strconv"
"math/rand"
)


func main() {

	initAgent()

	maxBtCount := 500
	btCount := 0

	fmt.Print("Login ")
	for btCount < maxBtCount {
                s2 := strconv.Itoa(btCount)
                fmt.Printf("%v ", s2)
		btCount++

		// start the "Login" transaction
		btHandle := appd.StartBT("Login", "")
		time.Sleep(time.Duration(getRandomMilliseconds()) * time.Millisecond)
		setSnapshotAttributes(btHandle, "label", "blue")

		// end the transaction
		appd.EndBT(btHandle)
	}
	fmt.Print("\n")

	// Stop/Clean up the AppD SDK.
	appd.TerminateSDK()

}

func initAgent() {
        fmt.Println("Arch Check")

        // Exec ldd test
        ldd_cmd := exec.Command("ldd", "--version")

        var ldd_out bytes.Buffer
        ldd_cmd.Stdout = &ldd_out

        ldd_err := ldd_cmd.Run();
        if ldd_err != nil {
                log.Fatal(ldd_err);
        }

        fmt.Printf("ldd --version:\n")
        fmt.Printf("%s\n", ldd_out.String())

        // Exec uname test
        uname_cmd := exec.Command("uname", "-a")

        var uname_out bytes.Buffer
        uname_cmd.Stdout = &uname_out

        uname_err := uname_cmd.Run();
        if uname_err != nil {
                log.Fatal(uname_err);
        }

        fmt.Printf("uname -a:\n")
        fmt.Printf("%s\n", uname_out.String())

        // Configure AppD
        cfg := appd.Config{}

        // Controller
        //cfg.Controller.Host = "ec2-52-37-30-238.us-west-2.compute.amazonaws.com"
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
                appd.SetBTURL(bt, "/login")
        }
}

func getRandomMilliseconds() int {
    rand.Seed(time.Now().UnixNano())
    min := 500
    max := 3000
    return rand.Intn(max - min + 1) + min
}
