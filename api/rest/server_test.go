package rest_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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

var (
	street = "123 N Fake st."
	city   = "Dallas"
	state  = "TX"
	zip    = "75001"
)
var propJsonTemplate = `{"street": "%v", "city": "%v", "state": "%v", "zip": "%v"}`

var propertyJson = fmt.Sprintf(
	propJsonTemplate,
	street,
	city,
	state,
	zip,
)

func TestPostNewProperty(t *testing.T) {
	t.Run("POST and GET Property", func(t *testing.T) {

		// post new property
		postResponse := httptestPost("/property", propertyJson)

		// parse response
		pr := getJsonMapFromResponseBody(t, postResponse)
		assertEqual(t, http.StatusCreated, postResponse.Code)
		assert.NotEmpty(t, pr["id"])
		propertyID := pr["id"].(string)

		// check property stored in repo
		p, err := propRepo.RetrieveProperty(propertyID)
		assert.Nil(t, err)
		assertEqual(t, street, p.Street)
		assertEqual(t, city, p.City)
		assertEqual(t, state, p.StateCode)
		assertEqual(t, zip, p.Zip)

		// check response data structure
		assertEqual(t, street, pr["street"])
		assertEqual(t, city, pr["city"])
		assertEqual(t, state, pr["state"])
		assertEqual(t, zip, pr["zip"])
		assertValidAndRecentDateString(t, pr["createdAt"].(string))
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

	t.Run("post with missing field, expect 400", func(t *testing.T) {
		postResponse := httptestPost("/property", `{
			"city": "Dallas",
			"state": "TX",
			"zip": "75001"
		}`) // street is missing
		assertEqual(t, http.StatusBadRequest, postResponse.Code)
	})

}

func TestListProperties(t *testing.T) {
	t.Run("expect empty set when no properties exist", func(t *testing.T) {
		// reset repo and server
		propRepo = repository.NewInMemoryRepo()
		server = rest.NewServer(propRepo)

		getResponse := httptestGet("/property")
		assertEqual(t, http.StatusOK, getResponse.Code)

		var propList struct {
			Items []map[string]interface{} `json:"items"`
		}
		err := json.Unmarshal(getResponse.Body.Bytes(), &propList)
		assert.Nil(t, err)
		assert.Len(t, propList.Items, 0)
	})

	t.Run("list two properties", func(t *testing.T) {
		// create two properties in propRepo
		p1 := propRepo.NewProperty("100 Main st", city, state, zip)
		p2 := propRepo.NewProperty("102 Main st", city, state, zip)
		propRepo.StoreProperty(p1)
		propRepo.StoreProperty(p2)

		// list via API
		getResponse := httptestGet("/property")
		assertEqual(t, http.StatusOK, getResponse.Code)

		var propList struct {
			Items []map[string]interface{} `json:"items"`
		}
		err := json.Unmarshal(getResponse.Body.Bytes(), &propList)
		assert.Nil(t, err)

		// ensure two results
		assert.Len(t, propList.Items, 2)

		// ensure those results are among the properties added
		for _, p := range propList.Items {
			id := p["id"]
			if id != p1.ID && id != p2.ID {
				t.Errorf("property id %v was not added?", id)
				break
			}
		}

		// make sure the same property wasn't just listed twice...
		assert.NotEqual(t, propList.Items[0]["id"], propList.Items[1]["id"])
	})
}

func TestGetProperty(t *testing.T) {
	t.Run("Get property", func(t *testing.T) {
		// create property in propRepo
		p := propRepo.NewProperty(street, city, state, zip)
		err := propRepo.StoreProperty(p)
		assert.Nil(t, err)

		// get property via API
		getResponse := httptestGet("/property/" + p.ID)
		assertEqual(t, http.StatusOK, getResponse.Code)

		// check response data structure
		m := getJsonMapFromResponseBody(t, getResponse)
		assertEqual(t, p.ID, m["id"])
		assertEqual(t, street, m["street"])
		assertEqual(t, city, m["city"])
		assertEqual(t, state, m["state"])
		assertEqual(t, zip, m["zip"])
		assertEqual(t, p.CreatedAt.Format(time.RFC3339), m["createdAt"])
	})

	t.Run("get unknown property, expect 404", func(t *testing.T) {
		getResponse := httptestGet("/property/doesNotExist")
		assertEqual(t, http.StatusNotFound, getResponse.Code)
	})
}

func TestDeleteProperty(t *testing.T) {
	t.Run("post then delete", func(t *testing.T) {
		// create property in propRepo
		p := propRepo.NewProperty(street, city, state, zip)
		err := propRepo.StoreProperty(p)
		assert.Nil(t, err)

		// delete via API
		delResponse := httptestDelete("/property/" + p.ID)
		assertEqual(t, http.StatusNoContent, delResponse.Code)

		// property should not be retrievable by repo anymore
		_, err = propRepo.RetrieveProperty(p.ID)
		assert.Error(t, err)
	})

	// should a restful DELETE on a resource that does not exist
	// result in a 404 or not?
	// https://stackoverflow.com/a/16632048/2683059
	// a lot of conflicting answers on this one, I'm going to chose no
	// for now because I can't think of a reason why the client should care
	t.Run("unknown property, expect 204", func(t *testing.T) {
		delResponse := httptestDelete("/property/doesNotExist")
		assertEqual(t, http.StatusNoContent, delResponse.Code)
	})
}

// http request helper functions
func httptestPost(uri, jsonStr string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(http.MethodPost, uri, jsonReader(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	return execReq(req)
}

func httptestGet(uri string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(http.MethodGet, uri, nil)
	return execReq(req)
}

func httptestDelete(uri string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(http.MethodDelete, uri, nil)
	return execReq(req)
}

func execReq(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)
	return rr
}

// json helper functions
func jsonReader(jsonStr string) *bytes.Buffer {
	return bytes.NewBuffer([]byte(jsonStr))
}

func getJsonMapFromResponseBody(t *testing.T, r *httptest.ResponseRecorder) jsonMap {
	t.Helper()
	var m jsonMap
	err := json.Unmarshal(r.Body.Bytes(), &m)
	assert.Nil(t, err)
	return m
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
	secondsAgo := time.Since(parsedTime).Seconds()
	lowerBound, upperBound := 0.0, 5.0
	assert.GreaterOrEqual(t, secondsAgo, lowerBound)
	assert.Less(t, secondsAgo, upperBound)
}
