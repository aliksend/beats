package simplerity

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	USER_CREDENTIALS_FILENAME  = "user_credentials.json"
	AGENT_CREDENTIALS_FILENAME = "agent_credentials.json"
)

type UserCredentials struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	AgentGuid    string `json:"agent_guid,omitempty"`
}

func loadUserCredentials() *UserCredentials {
	data, err := ioutil.ReadFile(USER_CREDENTIALS_FILENAME)
	if err != nil {
		fmt.Printf("Unable to read %s file. Using defaults\n", USER_CREDENTIALS_FILENAME)
		return &UserCredentials{}
	}
	credentials := UserCredentials{}
	err = json.Unmarshal(data, &credentials)
	if err != nil {
		fmt.Printf("Unable to unmarhal credentials. Using defaults\n")
		return &UserCredentials{}
	}
	return &credentials
}

func saveUserCredentials(credentials *UserCredentials) error {
	payload, err := json.Marshal(credentials)
	if err != nil {
		return fmt.Errorf("unable to marshal credentials: %v", err)
	}
	err = ioutil.WriteFile(USER_CREDENTIALS_FILENAME, payload, os.ModePerm)
	if err != nil {
		return fmt.Errorf("unable to write to file: %v", err)
	}
	return nil
}

type AgentCredentials struct {
	ApiSecret    string `json:"api_secret"`    // "Blowfish"
	ClientID     string `json:"client_id"`     // "1"
	ClientSecret string `json:"client_secret"` // abc123
	Scope        string `json:"scope"`         // "basic"
	PcName       string `json:"pc_name"`       // "demo"
	BuildVersion string `json:"build_version"` // "1.2.0.19"
}

func loadAgentCredentials() (*AgentCredentials, error) {
	data, err := ioutil.ReadFile(AGENT_CREDENTIALS_FILENAME)
	if err != nil {
		return nil, fmt.Errorf("unable to read %s file: %v", AGENT_CREDENTIALS_FILENAME, err)
	}
	credentials := AgentCredentials{}
	err = json.Unmarshal(data, &credentials)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarhal agent credentials: %v", err)
	}
	return &credentials, nil

}

func saveAgentCredentials(credentials *AgentCredentials) error {
	payload, err := json.Marshal(credentials)
	if err != nil {
		return fmt.Errorf("unable to marshal credentials: %v", err)
	}
	err = ioutil.WriteFile(AGENT_CREDENTIALS_FILENAME, payload, os.ModePerm)
	if err != nil {
		return fmt.Errorf("unable to write to file: %v", err)
	}
	return nil
}
