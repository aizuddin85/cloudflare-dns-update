package main

import (
        "bytes"
        "fmt"
        "io/ioutil"
        "net/http"
        "os"
        "strings"
)

const (
        zoneIdent   = xxxxxx
        recordIdent = yyyyyy
        authEmail   = "user@example.com"
        authKey     = "zzzzzz"
        recordName  = "www.example.com"
        recordType  = "A"
        recordTTL   = 1
        proxiedType = true
)

func main() {

        // Get public IP address
        pubIP, err := getPublicIP()
        if err != nil {
                fmt.Println("Error getting public IP:", err)
                return
        }

        // Read the stored IP from a file
        storedIP, err := readStoredIP()
        if err != nil {
                fmt.Println("Error reading stored IP:", err)
                return
        }

        // If the public IP hasn't changed, skip the update
        if pubIP == storedIP {
                fmt.Println("Public IP has not changed. Skipping update.")
                return
        }

        // Update the stored IP
        err = storeIP(pubIP)
        if err != nil {
                fmt.Println("Error storing IP:", err)
                return
        }

        // Cloudflare API request data
        recordData := fmt.Sprintf(`{"type":"%s","name":"%s","content":"%s","ttl":%d,"proxied":%v}`, recordType, recordName, pubIP, recordTTL,proxiedType)
        apiURL := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zoneIdent,recordIdent)
        // Create a PUT request
        req, err := http.NewRequest("PUT", apiURL, bytes.NewBuffer([]byte(recordData)))
        if err != nil {
                fmt.Println("Error creating request:", err)
                return
        }

        // Set headers
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("X-Auth-Email", authEmail)
        req.Header.Set("X-Auth-Key", authKey)

        // Perform the request
        client := http.Client{}
        resp, err := client.Do(req)
        if err != nil {
                fmt.Println("Error sending request:", err)
                return
        }
        defer resp.Body.Close()

        // Read and print the response
        respBody, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                fmt.Println("Error reading response:", err)
                return
        }
        fmt.Println("Response:", string(respBody))
}

func getPublicIP() (string, error) {
        resp, err := http.Get("https://ifconfig.co/ip")
        if err != nil {
                return "", err
        }
        defer resp.Body.Close()

        ipBytes, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                return "", err
        }

        return string(bytes.TrimSpace(ipBytes)), nil
}

func readStoredIP() (string, error) {
        // Check if the file exists
        if _, err := os.Stat("stored_ip.txt"); os.IsNotExist(err) {
                return "", nil // File doesn't exist, return empty IP
        }

        // Read the stored IP from the file
        ipBytes, err := ioutil.ReadFile("stored_ip.txt")
        if err != nil {
                return "", err
        }

        return strings.TrimSpace(string(ipBytes)), nil
}

func storeIP(ip string) error {
        // Write the IP to the file
        err := ioutil.WriteFile("stored_ip.txt", []byte(ip), 0644)
        if err != nil {
                return err
        }
        return nil
}
