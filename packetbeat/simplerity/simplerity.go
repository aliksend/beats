package simplerity

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
)

func Login(login string, password string) {
	agentCredentials, err := loadAgentCredentials()
	if err != nil {
		log.Fatalf("unable to load agent credentials: %v", err)
	}
	credentials := loadUserCredentials()

	var accessTokenResp struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Endpoints    []struct {
			Id    string `json:"id"`
			Title string `json:"title"`
		} `json:"endpoints"`
	}
	err = makeRequest("POST", "https://api.simplerity.com/access_token", "", map[string]interface{}{
		"grant_type":    "password",
		"username":      login,
		"password":      password,
		"client_id":     agentCredentials.ClientID,
		"client_secret": agentCredentials.ClientSecret,
		"scope":         agentCredentials.Scope,
		"pcName":        agentCredentials.PcName,
	}, agentCredentials.ApiSecret, &accessTokenResp)
	if err != nil {
		log.Fatalf("unable to request access_token: %v", err)
	}

	credentials.AccessToken = accessTokenResp.AccessToken
	credentials.ExpiresIn = accessTokenResp.ExpiresIn
	credentials.RefreshToken = accessTokenResp.RefreshToken

	err = saveUserCredentials(credentials)
	if err != nil {
		log.Fatalf("unable to save user credentials: %v", err)
	}

	fmt.Println("Authorized. Now you should select agent to use:")
	for _, endpoint := range accessTokenResp.Endpoints {
		fmt.Printf("%s: %s\n", endpoint.Id, endpoint.Title)
	}
	fmt.Println("Use `packetbeat simplerity select --agent AGENT_ID` to select agent")
}

func Select(agentID string) {
	agentCredentials, err := loadAgentCredentials()
	if err != nil {
		log.Fatalf("unable to load agent credentials: %v", err)
	}
	credentials := loadUserCredentials()
	if credentials.AccessToken == "" {
		log.Fatal("You should run `packetbeat simplerity login --login YOUR_LOGIN --password YOUR_PASSWORD` first")
	}

	var selectAgentResp struct {
		Ok   bool   `json:"ok"`
		Guid string `json:"guid"`
	}
	err = makeRequest("POST", "https://api.simplerity.com/api/v2/agents/select", credentials.AccessToken, map[string]interface{}{
		"agentId":    agentID,
		"group_name": "ua,office-10",
	}, agentCredentials.ApiSecret, &selectAgentResp)
	if err != nil {
		log.Fatalf("unable to select agent: %v", err)
	}
	if !selectAgentResp.Ok {
		log.Fatalf("unable to select agent: %v", selectAgentResp)
	}

	credentials.AgentGuid = selectAgentResp.Guid

	err = saveUserCredentials(credentials)
	if err != nil {
		log.Fatalf("unable to save user credentials: %v", err)
	}

	fmt.Println("Agent selected. Now you can download configuration using `packetbeat simplerity load`")
}

func LoadConfig(configFileName string) {
	agentCredentials, err := loadAgentCredentials()
	if err != nil {
		log.Fatalf("unable to load agent credentials: %v", err)
	}
	credentials := loadUserCredentials()
	if credentials.AccessToken == "" {
		log.Fatal("You should run `packetbeat simplerity login --login YOUR_LOGIN --password YOUR_PASSWORD` first")
	}
	if credentials.AgentGuid == "" {
		log.Fatal("You should run `packetbeat simplerity select --agent AGENT_ID` first")
	}

	var packetBeatConfigResp struct {
		Data string `json:"data"`
	}
	err = makeRequest("POST", "https://api.simplerity.com/api/v2/"+credentials.AgentGuid+"/packetbeat/config", credentials.AccessToken, map[string]interface{}{
		"build_version": agentCredentials.BuildVersion,
		"config":        "mid",
	}, agentCredentials.ApiSecret, &packetBeatConfigResp)
	if err != nil {
		log.Fatalf("unable to request packetbeat config: %v", err)
	}
	err = ioutil.WriteFile(configFileName, []byte(packetBeatConfigResp.Data), os.ModePerm)
	if err != nil {
		log.Fatalf("unable to save packetbeat config: %v", err)
	}

	fmt.Printf("Config downloaded. Now you can launch packetbeat using `packetbeat --c %s run`\n", configFileName)
}

func generatePayloadSignature(payload map[string]interface{}, apiSecret string) string {
	var keys []string
	for k := range payload {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	stringToEncode := ""
	for _, k := range keys {
		stringToEncode += fmt.Sprintf("%s%v", k, payload[k])
	}

	h := hmac.New(sha256.New, []byte(apiSecret))
	_, err := h.Write([]byte(stringToEncode))
	if err != nil {
		log.Fatal(err)
	}
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

func makeRequest(method string, uri string, authHeader string, payload map[string]interface{}, apiSecret string, responsePtr interface{}) error {
	payload["signature"] = generatePayloadSignature(payload, apiSecret)
	data := url.Values{}
	for k, v := range payload {
		data.Set(k, fmt.Sprintf("%v", v))
	}

	req, err := http.NewRequest(method, uri, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("unable to create http request: %v", err)
	}
	req.Header.Add("Authorization", authHeader)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("unable to make request to api: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read body of response: %v", err)
	}
	// log.Printf("API RESPONSE %s: %s", uri, body)
	err = json.Unmarshal(body, responsePtr)
	if err != nil {
		return fmt.Errorf("unable to parse response: %v", err)
	}
	return nil
}
