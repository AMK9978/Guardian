package services

import (
	"bytes"
	"fmt"
	"net/http"

	"guardian/internal/models"
)

type ProcessingTask interface {
	Process(req models.SendRequest) (bool, error)
}

//
//type FewShotTask struct {
//	Template string
//	ApiUrl string
//}

//func (f *FewShotTask) Process(req models.SendRequest) (bool, error) {
//	resp, err := httpAPICall(f.ApiUrl, fmt.Sprintf("%s\n%s", f.Template, req.Prompt))
//	if err != nil {
//		return false, err
//	}
//	return , nil
//}

type ExternalHttpServiceTask struct {
	ApiUrl string
}

func (e *ExternalHttpServiceTask) Process(req models.SendRequest) (bool, error) {
	resp, err := httpAPICall(e.ApiUrl, req.Prompt)
	if err != nil {
		return false, err
	}
	return resp, nil
}

func httpAPICall(url string, prompt string) (bool, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(prompt)))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return true, nil
}
