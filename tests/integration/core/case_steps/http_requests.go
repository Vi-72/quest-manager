package casesteps

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
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
	req.Headers = withAuthHeader(req.Headers)
	var body io.Reader

	// Подготавливаем тело запроса
	if req.Body != nil {
		switch v := req.Body.(type) {
		case string:
			// Raw string body (for malformed JSON)
			body = bytes.NewReader([]byte(v))
		case json.RawMessage:
			// Raw JSON message
			body = bytes.NewReader([]byte(v))
		default:
			// Regular objects - marshal to JSON
			bodyBytes, err := json.Marshal(req.Body)
			if err != nil {
				return nil, err
			}
			body = bytes.NewReader(bodyBytes)
		}
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
		Headers:     withAuthHeader(nil),
		ContentType: "application/json",
	}
}

// AssignQuestHTTPRequestWithStringID создает HTTP запрос с строковым ID (для тестирования невалидных UUID)
func AssignQuestHTTPRequestWithStringID(questID string) HTTPRequest {
	return HTTPRequest{
		Method:      "POST",
		URL:         "/api/v1/quests/" + questID + "/assign",
		Headers:     withAuthHeader(nil),
		ContentType: "application/json",
	}
}

// AssignQuestHTTPRequest создает HTTP запрос для назначения квеста
func AssignQuestHTTPRequest(questID uuid.UUID) HTTPRequest {
	return HTTPRequest{
		Method:      "POST",
		URL:         "/api/v1/quests/" + questID.String() + "/assign",
		Headers:     withAuthHeader(nil),
		ContentType: "application/json",
	}
}

// GetQuestHTTPRequest создает HTTP запрос для получения квеста
func GetQuestHTTPRequest(questID uuid.UUID) HTTPRequest {
	return HTTPRequest{
		Method:  "GET",
		URL:     "/api/v1/quests/" + questID.String(),
		Headers: withAuthHeader(nil),
	}
}

// GetQuestHTTPRequestWithStringID создает HTTP запрос с строковым ID (для тестирования невалидных UUID)
func GetQuestHTTPRequestWithStringID(questID string) HTTPRequest {
	return HTTPRequest{
		Method:  "GET",
		URL:     "/api/v1/quests/" + questID,
		Headers: withAuthHeader(nil),
	}
}

// ListQuestsHTTPRequest создает HTTP запрос для получения списка квестов
func ListQuestsHTTPRequest(status string) HTTPRequest {
	url := "/api/v1/quests"
	if status != "" {
		url += "?status=" + status
	}

	return HTTPRequest{
		Method:  "GET",
		URL:     url,
		Headers: withAuthHeader(nil),
	}
}

// ListAssignedQuestsHTTPRequest создает HTTP запрос для получения квестов назначенных аутентифицированному пользователю
// User ID теперь берется из JWT токена, поэтому не передается в query параметрах
func ListAssignedQuestsHTTPRequest() HTTPRequest {
	return HTTPRequest{
		Method:  "GET",
		URL:     "/api/v1/quests/assigned",
		Headers: withAuthHeader(nil),
	}
}

// SearchQuestsByRadiusHTTPRequest создает HTTP запрос для поиска квестов по радиусу
func SearchQuestsByRadiusHTTPRequest(lat, lon, radiusKm float32) HTTPRequest {
	url := fmt.Sprintf("/api/v1/quests/search-radius?lat=%f&lon=%f&radius_km=%f", lat, lon, radiusKm)
	return HTTPRequest{
		Method:  "GET",
		URL:     url,
		Headers: withAuthHeader(nil),
	}
}

// ChangeQuestStatusHTTPRequest создает HTTP запрос для изменения статуса квеста
func ChangeQuestStatusHTTPRequest(questID uuid.UUID, statusRequest interface{}) HTTPRequest {
	return HTTPRequest{
		Method:      "PATCH",
		URL:         "/api/v1/quests/" + questID.String() + "/status",
		Body:        statusRequest,
		Headers:     withAuthHeader(nil),
		ContentType: "application/json",
	}
}

// ChangeQuestStatusHTTPRequestWithStringID создает HTTP запрос с строковым ID (для тестирования невалидных UUID)
func ChangeQuestStatusHTTPRequestWithStringID(questID string, statusRequest interface{}) HTTPRequest {
	return HTTPRequest{
		Method:      "PATCH",
		URL:         "/api/v1/quests/" + questID + "/status",
		Body:        statusRequest,
		Headers:     withAuthHeader(nil),
		ContentType: "application/json",
	}
}

// CreateMalformedJSONRequest создает HTTP запрос с невалидным JSON
func CreateMalformedJSONRequest(method, url string) HTTPRequest {
	return HTTPRequest{
		Method:      method,
		URL:         url,
		Body:        `{"status": invalid-json}`, // Malformed JSON
		Headers:     withAuthHeader(nil),
		ContentType: "application/json",
	}
}

func withAuthHeader(headers map[string]string) map[string]string {
	if headers == nil {
		headers = make(map[string]string)
	}

	result := make(map[string]string, len(headers)+1)
	for k, v := range headers {
		result[k] = v
	}
	// Only add Authorization if not explicitly provided (even if empty)
	if _, exists := result["Authorization"]; !exists {
		result["Authorization"] = "Bearer test-token"
	}
	return result
}
