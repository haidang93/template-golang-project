package emailservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/example/internal/config"
)

type MailTrapService struct {
	HttpClient *http.Client
	ApiToken   *string
	SendURL    *string
}

func CreateMailTrapService(env *config.Environment, httpClient *http.Client) *MailTrapService {
	return &MailTrapService{
		HttpClient: httpClient,
		ApiToken:   &env.MAILTRAP_API_TOKEN,
		SendURL:    &env.MAILTRAP_SEND_URL,
	}
}

func (s *MailTrapService) Send(data *SendData) error {

	body, err := json.Marshal(data.Data())
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, *s.SendURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+*s.ApiToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := s.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("delete failed: %s", string(bodyBytes))
	}

	fmt.Println("Delete succeeded")
	return nil
}
