package casesteps

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

// HTTPRequest представляет HTTP запрос для тестирования
type HTTPRequest struct {
	Method      string
	URL         string
	Body        interface{}
	Headers     map[string]string
	ContentType string
}

// HTTPResponse представляет HTTP ответ
type HTTPResponse struct {
	StatusCode int
	Body       string
	Headers    http.Header
}

// ExecuteHTTPRequest выполняет HTTP запрос через тестовый сервер
func ExecuteHTTPRequest(ctx context.Context, handler http.Handler, req HTTPRequest) (*HTTPResponse, error) {
	var body io.Reader

	// Подготавливаем тело запроса
	if req.Body != nil {
		bodyBytes, err := json.Marshal(req.Body)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(bodyBytes)
	}

	// Создаем HTTP запрос
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.URL, body)
	if err != nil {
		return nil, err
	}

	// Устанавливаем заголовки
	if req.ContentType != "" {
		httpReq.Header.Set("Content-Type", req.ContentType)
	} else if req.Body != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Выполняем запрос
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httpReq)

	// Возвращаем ответ
	return &HTTPResponse{
		StatusCode: recorder.Code,
		Body:       recorder.Body.String(),
		Headers:    recorder.Header(),
	}, nil
}

// CreateQuestHTTPRequest создает HTTP запрос для создания квеста
func CreateQuestHTTPRequest(questData interface{}) HTTPRequest {
	return HTTPRequest{
		Method:      "POST",
		URL:         "/api/v1/quests",
		Body:        questData,
		ContentType: "application/json",
	}
}

// AssignQuestHTTPRequest создает HTTP запрос для назначения квеста
func AssignQuestHTTPRequest(questID string, userID string) HTTPRequest {
	return HTTPRequest{
		Method:      "POST",
		URL:         "/api/v1/quests/" + questID + "/assign",
		Body:        map[string]string{"user_id": userID},
		ContentType: "application/json",
	}
}

// ChangeQuestStatusHTTPRequest создает HTTP запрос для изменения статуса квеста
func ChangeQuestStatusHTTPRequest(questID string, status string) HTTPRequest {
	return HTTPRequest{
		Method:      "PATCH",
		URL:         "/api/v1/quests/" + questID + "/status",
		Body:        map[string]string{"status": status},
		ContentType: "application/json",
	}
}

// GetQuestHTTPRequest создает HTTP запрос для получения квеста
func GetQuestHTTPRequest(questID string) HTTPRequest {
	return HTTPRequest{
		Method: "GET",
		URL:    "/api/v1/quests/" + questID,
	}
}

// ListQuestsHTTPRequest создает HTTP запрос для получения списка квестов
func ListQuestsHTTPRequest(status string) HTTPRequest {
	url := "/api/v1/quests"
	if status != "" {
		url += "?status=" + status
	}

	return HTTPRequest{
		Method: "GET",
		URL:    url,
	}
}

// ListAssignedQuestsHTTPRequest создает HTTP запрос для получения квестов назначенных пользователю
func ListAssignedQuestsHTTPRequest(userID string) HTTPRequest {
	return HTTPRequest{
		Method: "GET",
		URL:    "/api/v1/quests/assigned?user_id=" + userID,
	}
}

// SearchQuestsByRadiusHTTPRequest создает HTTP запрос для поиска квестов по радиусу
func SearchQuestsByRadiusHTTPRequest(lat, lon, radiusKm float64) HTTPRequest {
	url := fmt.Sprintf("/api/v1/quests/search-radius?lat=%f&lon=%f&radius_km=%f", lat, lon, radiusKm)
	return HTTPRequest{
		Method: "GET",
		URL:    url,
	}
}
