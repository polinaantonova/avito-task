package createTender

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"polina.com/m/internal/tender"
	"strings"
	"testing"
)

func TestTenderCreator_ServeHTTP(t *testing.T) {
	tenders := tender.NewTenderList()
	handler := NewTenderCreator(tenders)

	req := httptest.NewRequest(
		http.MethodPost,
		"/123123123",
		strings.NewReader(`{"name": "Доставка", "description": "Доставить товары из Казани в Москву", "serviceType": "Delivery", "creatorUsername": "user1"}`),
	)

	writer := httptest.NewRecorder()

	handler.ServeHTTP(writer, req)

	resp := writer.Result()
	require.Equal(t, 200, resp.StatusCode)
	require.Len(t, tenders.List(), 1)
}
