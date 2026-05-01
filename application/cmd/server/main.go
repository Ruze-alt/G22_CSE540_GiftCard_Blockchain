package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "path/filepath"

    "giftcard/application/api"
    "giftcard/application/api/fabric"
    "giftcard/application/api/handlers"
    adminrole "giftcard/application/roles/admin/service"
    customerrole "giftcard/application/roles/customer/service"
    issuerrole "giftcard/application/roles/issuer/service"
    retailerrole "giftcard/application/roles/retailer/service"
)

func getenv(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}

func main() {
    testNetworkPath := os.Getenv("TEST_NETWORK_PATH")
    if testNetworkPath == "" {
        log.Fatal("TEST_NETWORK_PATH environment variable is required")
    }
    channel := getenv("CHANNEL_NAME", "mychannel")
    chaincode := getenv("CHAINCODE_NAME", "giftCard")
    port := getenv("PORT", "8080")

    // Default to frontend/ relative to the working directory (where `go run` is invoked from).
    // Override with FRONTEND_DIR if needed (e.g. when running a compiled binary from another location).
    cwd, err := os.Getwd()
    if err != nil {
        log.Fatalf("get working directory: %v", err)
    }
    frontendDir := getenv("FRONTEND_DIR", filepath.Join(cwd, "frontend"))

    org1GW, err := fabric.NewFabricGateway(fabric.Org1Config(testNetworkPath, channel, chaincode))
    if err != nil {
        log.Fatalf("connect org1 gateway: %v", err)
    }
    defer org1GW.Close()

    org2GW, err := fabric.NewFabricGateway(fabric.Org2Config(testNetworkPath, channel, chaincode))
    if err != nil {
        log.Fatalf("connect org2 gateway: %v", err)
    }
    defer org2GW.Close()

    deps := handlers.HandlerDeps{
        Issuer:   issuerrole.New(org1GW),
        Retailer: retailerrole.New(org1GW),
        Admin:    adminrole.New(org1GW),
        Customer: customerrole.New(org2GW),
    }

    handler := api.NewServer(deps, frontendDir)

    addr := fmt.Sprintf(":%s", port)
    log.Printf("Starting server on http://localhost%s", addr)
    log.Printf("Serving frontend from: %s", frontendDir)
    if err := http.ListenAndServe(addr, handler); err != nil {
        log.Fatal(err)
    }
}
