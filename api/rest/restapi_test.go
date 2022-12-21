package rest_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tempcke/rpm/api/rest"
	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/entity/fake"
	"github.com/tempcke/rpm/repository"
)

var ctx = context.Background()

func TestAddProperty(t *testing.T) {
	var (
		repo = repository.NewInMemoryRepo()
		s    = rest.NewServer(repo)

		headers map[string]string
	)
	t.Run("201 post creates a new property with server generated id", func(t *testing.T) {
		var (
			route = "/property"
			p1    = fake.Property()
		)
		body := map[string]string{
			"street": p1.Street,
			"city":   p1.City,
			"state":  p1.StateCode,
			"zip":    p1.Zip,
		}

		rr := httptest.NewRecorder()
		s.ServeHTTP(rr, postReq(t, route, body, headers))

		require.Equal(t, http.StatusCreated, rr.Code)
		var res rest.PropertyModel
		require.NoError(t, json.NewDecoder(rr.Body).Decode(&res))

		// check response body for the echoed entity
		assert.NotEmpty(t, res.ID)
		assertEqual(t, p1.Street, res.Street)
		assertEqual(t, p1.City, res.City)
		assertEqual(t, p1.StateCode, res.State)
		assertEqual(t, p1.Zip, res.Zip)
		resCreatedAt, err := time.Parse(time.RFC3339, res.CreatedAt)
		require.NoError(t, err)
		assertTimeRecent(t, resCreatedAt)

		// check location header
		// we should be able to use the location header to fetch the resource
		loc := rr.Header().Get("Location")
		require.NotEmpty(t, loc)

		// get property from route in location header
		rr = httptest.NewRecorder()
		s.ServeHTTP(rr, getReq(t, loc, headers))
		require.Equal(t, http.StatusOK, rr.Code)

		var fetched rest.PropertyModel
		assert.NoError(t, json.Unmarshal(rr.Body.Bytes(), &fetched))
		require.NoError(t, err)
		assertEqual(t, res.ID, fetched.ID)
		assertEqual(t, p1.Street, fetched.Street)
		assertEqual(t, p1.City, fetched.City)
		assertEqual(t, p1.StateCode, fetched.State)
		assertEqual(t, p1.Zip, fetched.Zip)
	})
	t.Run("201 put new property", func(t *testing.T) {
		var (
			p1    = fake.Property()
			route = "/property/" + p1.ID
		)
		body := map[string]string{
			"street": p1.Street,
			"city":   p1.City,
			"state":  p1.StateCode,
			"zip":    p1.Zip,
		}

		rr := httptest.NewRecorder()
		s.ServeHTTP(rr, putReq(t, route, body, headers))
		require.Equal(t, http.StatusCreated, rr.Code)
		var res rest.PropertyModel
		assert.NoError(t, json.Unmarshal(rr.Body.Bytes(), &res))

		// check response body for the echoed entity
		assertEqual(t, p1.ID, res.ID)
		assertEqual(t, p1.Street, res.Street)
		assertEqual(t, p1.City, res.City)
		assertEqual(t, p1.StateCode, res.State)
		assertEqual(t, p1.Zip, res.Zip)
		resCreatedAt, err := time.Parse(time.RFC3339, res.CreatedAt)
		require.NoError(t, err)
		assertTimeRecent(t, resCreatedAt)
	})
	t.Run("200 put updated property", func(t *testing.T) {
		var (
			p1    = fake.Property()
			route = "/property/" + p1.ID
		)

		// create property with a typo mistake of some kind
		body := map[string]string{
			"street": p1.Street + "typo",
			"city":   p1.City,
			"state":  p1.StateCode,
			"zip":    p1.Zip,
		}
		rr := httptest.NewRecorder()
		s.ServeHTTP(rr, putReq(t, route, body, headers))
		require.Equal(t, http.StatusCreated, rr.Code)

		// update property fixing typo mistake
		body["street"] = p1.Street
		rr = httptest.NewRecorder()
		s.ServeHTTP(rr, putReq(t, route, body, headers))
		assert.Equal(t, http.StatusOK, rr.Code)

		var res rest.PropertyModel
		assert.NoError(t, json.Unmarshal(rr.Body.Bytes(), &res))

		// check response body for the echoed entity
		assertEqual(t, p1.ID, res.ID)
		assertEqual(t, p1.Street, res.Street)
		assertEqual(t, p1.City, res.City)
		assertEqual(t, p1.StateCode, res.State)
		assertEqual(t, p1.Zip, res.Zip)
	})
}
func TestAddProperty_badRequest(t *testing.T) {
	var (
		headers map[string]string
		route   = "/property"
		repo    = repository.NewInMemoryRepo()
		s       = rest.NewServer(repo)
	)
	tests := map[string]struct {
		body string
	}{
		"no street": {`{"city": "b", "state": "TX", "zip": "12345"}`},
		"no city":   {`{"street": "a", "state": "TX", "zip": "12345"}`},
		"no state":  {`{"street": "a", "city": "b", "zip": "12345"}`},
		"no zip":    {`{"street": "a", "city": "b", "state": "TX"}`},
		"invalid json": {
			`{
			"street": "a", 
			"city": "b", 
			"state": "TX", 
			"zip": "12345",
		}`,
		}, // notice the last ',' shouldn't be there
	}
	for name, body := range tests {
		t.Run(name, func(t *testing.T) {
			req := postReq(t, route, body, headers)
			rr := httptest.NewRecorder()
			s.ServeHTTP(rr, req)
			assertEqual(t, http.StatusBadRequest, rr.Code)
		})
	}
}

func TestListProperties(t *testing.T) {
	var (
		route   = "/property"
		headers map[string]string
	)
	t.Run("200 expect empty set when no properties exist", func(t *testing.T) {
		var (
			repo = repository.NewInMemoryRepo()
			s    = rest.NewServer(repo)
		)

		req := getReq(t, route, headers)
		rr := httptest.NewRecorder()
		s.ServeHTTP(rr, req)
		require.Equal(t, http.StatusOK, rr.Code)

		var propList rest.PropertyList
		require.NoError(t, json.NewDecoder(rr.Body).Decode(&propList))
		assert.Len(t, propList.Items, 0)
	})
	t.Run("200 list two properties", func(t *testing.T) {
		var (
			repo = repository.NewInMemoryRepo()
			s    = rest.NewServer(repo)
		)

		// create two properties in propRepo
		p1 := fake.Property()
		p2 := fake.Property()
		require.NoError(t, repo.StoreProperty(ctx, p1))
		require.NoError(t, repo.StoreProperty(ctx, p2))

		// list via API
		req := getReq(t, route, headers)
		rr := httptest.NewRecorder()
		s.ServeHTTP(rr, req)
		require.Equal(t, http.StatusOK, rr.Code)

		var propList rest.PropertyList
		assert.NoError(t, json.Unmarshal(rr.Body.Bytes(), &propList))
		assert.Len(t, propList.Items, 2)

		// ensure those results are among the properties added
		propMap := map[string]entity.Property{
			p1.ID: p1,
			p2.ID: p2,
		}
		for _, item := range propList.Items {
			p, ok := propMap[item.ID]
			require.True(t, ok)
			assertEqual(t, p.Street, item.Street)
			assertEqual(t, p.City, item.City)
			assertEqual(t, p.StateCode, item.State)
			assertEqual(t, p.Zip, item.Zip)
		}

		// make sure the same property wasn't just listed twice...
		assert.NotEqual(t, propList.Items[0].ID, propList.Items[1].ID)
	})
}

func TestGetProperty(t *testing.T) {
	var (
		routeBase = "/property/"
		headers   map[string]string
		repo      = repository.NewInMemoryRepo()
		s         = rest.NewServer(repo)
	)
	t.Run("200 get property", func(t *testing.T) {
		var (
			p1    = fake.Property()
			route = routeBase + p1.ID
		)
		// create property in propRepo
		require.NoError(t, repo.StoreProperty(ctx, p1))

		// get property via API
		req := getReq(t, route, headers)
		rr := httptest.NewRecorder()
		s.ServeHTTP(rr, req)
		require.Equal(t, http.StatusOK, rr.Code)

		// check response data structure
		var res rest.PropertyModel
		require.NoError(t, json.NewDecoder(rr.Body).Decode(&res))
		assertEqual(t, p1.ID, res.ID)
		assertEqual(t, p1.Street, res.Street)
		assertEqual(t, p1.City, res.City)
		assertEqual(t, p1.StateCode, res.State)
		assertEqual(t, p1.Zip, res.Zip)
	})
	t.Run("404 get unknown property", func(t *testing.T) {
		var (
			p1    = fake.Property()
			route = routeBase + p1.ID
		)

		req := getReq(t, route, headers)
		rr := httptest.NewRecorder()
		s.ServeHTTP(rr, req)
		require.Equal(t, http.StatusNotFound, rr.Code)
	})
}

func TestDeleteProperty(t *testing.T) {
	var (
		routeBase = "/property/"
		headers   map[string]string
		repo      = repository.NewInMemoryRepo()
		s         = rest.NewServer(repo)
	)
	t.Run("204 delete property", func(t *testing.T) {
		var (
			p1    = fake.Property()
			route = routeBase + p1.ID
		)
		// seed property
		require.NoError(t, repo.StoreProperty(ctx, p1))

		// del property via API
		req := delReq(t, route, headers)
		rr := httptest.NewRecorder()
		s.ServeHTTP(rr, req)
		require.Equal(t, http.StatusNoContent, rr.Code)

		// property should not be retrievable by repo anymore
		_, err := repo.RetrieveProperty(ctx, p1.ID)
		assert.Error(t, err)
	})
	t.Run("204 delete is idempotent", func(t *testing.T) {
		// should a restful DELETE on a resource that does not exist result in a 404 or not?
		// https://stackoverflow.com/a/16632048/2683059
		// a lot of conflicting answers on this one, I'm going to choose no
		// for now because I can't think of a reason why the client should care

		var (
			p1    = fake.Property()
			route = routeBase + p1.ID
		)
		// seed property
		require.NoError(t, repo.StoreProperty(ctx, p1))

		// del property via API
		req := delReq(t, route, headers)
		rr := httptest.NewRecorder()
		s.ServeHTTP(rr, req)
		require.Equal(t, http.StatusNoContent, rr.Code)

		// del property again via API - idempotent check
		rr = httptest.NewRecorder()
		s.ServeHTTP(rr, req)
		require.Equal(t, http.StatusNoContent, rr.Code)
	})
}
