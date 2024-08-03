package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type publicIp struct {
	Ip string `json:"ip"`
}

type newCloudflareDNSRecord struct {
	Content string `json:"content"`
}

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func (c publicIp) String() string {
	return fmt.Sprintf("%s", c.Ip)
}

func NewConfig() (publicIp, error) {
	var cf publicIp
	confFile, err := os.Open("./config.json")

	if err != nil {
		if os.IsNotExist(err) {
			// File does not exist, create it with default values
			defaultConfig := publicIp{
				Ip: "",
			}
			configFile, err := os.Create("./config.json")
			if err != nil {
				return cf, err
			}
			defer configFile.Close()

			encoder := json.NewEncoder(configFile)
			if err := encoder.Encode(&defaultConfig); err != nil {
				return cf, err
			}

			// Reopen the newly created file
			confFile, err = os.Open("./config.json")
			if err != nil {
				return cf, err
			}
		} else {
			// Some other error occurred
			log.Fatal(err)
		}
	}

	defer confFile.Close()

	jsonParser := json.NewDecoder(confFile)
	if err := jsonParser.Decode(&cf); err != nil {
		return cf, err
	}

	return cf, nil
}

func updateConfigFile(ipAddress string) {
	var cfg, err = NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	cfg.Ip = ipAddress

	file, err := json.MarshalIndent(cfg, "", " ")

	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("./config.json", file, 0644)

	if err != nil {
		log.Fatal(err)
	}
}

func getPublicIp() string {
	req, err := http.Get("https://api.ipify.org?format=json")

	if err != nil {
		log.Fatal(err)
		return err.Error()
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)

	log.Println(string(body))

	if err != nil {
		log.Fatal(err)
		return err.Error()
	}

	var publicIp publicIp
	json.Unmarshal(body, &publicIp)
	log.Println("Public IP: " + publicIp.Ip)

	return publicIp.Ip
}

func updateDnsRecord(apikey string, zoneId string, recordId string, ipAddress string) {

	var cfg, err = NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Ip == ipAddress {
		log.Println("IP Address is the same. Exiting...")
		return
	}

	newRecord := &newCloudflareDNSRecord{
		Content: ipAddress,
	}

	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(newRecord)

	client := &http.Client{}
	req, err := http.NewRequest("PATCH", "https://api.cloudflare.com/client/v4/zones/"+zoneId+"/dns_records/"+recordId, body)

	log.Println("https://api.cloudflare.com/client/v4/zones/" + zoneId + "/dns_records/" + recordId)
	log.Println("IP Address: " + ipAddress)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+apikey)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)

	updateConfigFile(ipAddress)

	log.Println(res)

	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	apikey := goDotEnvVariable("API_KEY")
	zoneId := goDotEnvVariable("ZONE_ID")
	recordId := goDotEnvVariable("DNS_RECORD_ID")

	updateDnsRecord(apikey, zoneId, recordId, getPublicIp())

}
