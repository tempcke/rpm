package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tempcke/rpm/api/rest"
	"github.com/tempcke/rpm/repository"
)

type jsonMap map[string]interface{}

var propRepo = repository.NewInMemoryRepo()
var server = rest.NewServer(propRepo)

func getJsonMapFromResponseBody(t *testing.T, r *httptest.ResponseRecorder) jsonMap {
	t.Helper()
	var m jsonMap
	err := json.Unmarshal(r.Body.Bytes(), &m)
	assert.Nil(t, err)
	return m
}
func TestPostNewProperty(t *testing.T) {
	t.Run("POST and GET Property", func(t *testing.T) {
		postResponse := httptestPost("/property", `{
			"street": "123 N Fake st.",
			"city": "Dallas",
			"state": "TX",
			"zip": "75001"
		}`)

		m := getJsonMapFromResponseBody(t, postResponse)
		propertyID := m["id"].(string)

		t.Run("verify post response", func(t *testing.T) {
			assertEqual(t, http.StatusCreated, postResponse.Code)
			assertEqual(t, "123 N Fake st.", m["street"])
			assertEqual(t, "Dallas", m["city"])
			assertEqual(t, "TX", m["state"])
			assertEqual(t, "75001", m["zip"])
			assert.NotEmpty(t, m["id"])
			assert.NotEmpty(t, m["createdAt"])
			assertValidAndRecentDateString(t, m["createdAt"].(string))
			propertyID = m["id"].(string)
		})

		t.Run("Get property", func(t *testing.T) {
			getResponse := httptestGet("/property/" + propertyID)
			assertEqual(t, http.StatusOK, getResponse.Code)

			// getResponse data should match postResponse data
			m2 := getJsonMapFromResponseBody(t, getResponse)
			assert.Equal(t, m, m2)
		})
	})

	t.Run("post invalid json, expect 400", func(t *testing.T) {
		postResponse := httptestPost("/property", `{
			"street": "123 N Fake st.",
			"city": "Dallas",
			"state": "TX",
			"zip": "75001",
		}`) // notice the last ',' shouldn't be there
		assertEqual(t, http.StatusBadRequest, postResponse.Code)
	})

	t.Run("post invalid property, expect 400", func(t *testing.T) {
		postResponse := httptestPost("/property", `{
			"city": "Dallas",
			"state": "TX",
			"zip": "75001"
		}`) // street is missing
		assertEqual(t, http.StatusBadRequest, postResponse.Code)
	})

	t.Run("unknown property, expect 404", func(t *testing.T) {
		getResponse := httptestGet("/property/does-Not-Exist")
		assertEqual(t, http.StatusNotFound, getResponse.Code)
	})
}

// Helper functions
func httptestPost(uri, jsonStr string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(http.MethodPost, uri, jsonReader(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	return execReq(req)
}

func httptestGet(uri string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(http.MethodGet, uri, nil)
	return execReq(req)
}

func execReq(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)
	return rr
}

func jsonReader(jsonStr string) *bytes.Buffer {
	return bytes.NewBuffer([]byte(jsonStr))
}

// Custom Assertions
func assertEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if expected != actual {
		t.Errorf(
			"Values not equal \nWant: %v \t%T\nGot:  %v \t%T\n",
			expected, expected, actual, actual,
		)
	}
}

func assertValidAndRecentDateString(t *testing.T, timeStr string) {
	t.Helper()

	// validate string as RFC3339
	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	assert.Nil(t, err)

	// ensure recent
	secondsAgo := time.Now().Sub(parsedTime).Seconds()
	lowerBound, upperBound := 0.0, 5.0
	assert.GreaterOrEqual(t, secondsAgo, lowerBound)
	assert.Less(t, secondsAgo, upperBound)
}
