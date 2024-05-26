package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/romanyakovlev/gophermart/internal/models"
)

type AccrualService struct {
	baseURL string
}

func (s *AccrualService) FetchOrderDetails(orderNumber string) (*models.OrderResponse, error) {
	url := fmt.Sprintf("%s/api/orders/%s", s.baseURL, orderNumber)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var orderResp models.OrderResponse
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bodyBytes, &orderResp); err != nil {
		return nil, err
	}

	return &orderResp, nil
}

func NewAccrualService(baseURL string) *AccrualService {
	return &AccrualService{baseURL: baseURL}
}
