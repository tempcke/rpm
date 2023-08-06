package rest_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tempcke/rpm/actions"
	"github.com/tempcke/rpm/api/rest"
	"github.com/tempcke/rpm/entity"
	"github.com/tempcke/rpm/entity/fake"
	"github.com/tempcke/rpm/internal"
	"github.com/tempcke/rpm/repository"
)

var ctx = context.Background()

func TestPutProperty(t *testing.T) {
	var (
		repo    = repository.NewInMemoryRepo()
		s       = rest.NewServer(actions.NewActions(repo)).WithConfig(noAuthConf(t))
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

		res := handleReq(t, s, postReq(t, route, body, headers))
		assertResCode(t, res, http.StatusCreated)
		assertApplicationJson(t, res.Header)
		var created rest.PropertyModel
		require.NoError(t, json.NewDecoder(res.Body).Decode(&created))

		// check response body for the echoed entity
		assert.NotEmpty(t, created.ID)
		assertEqual(t, p1.Street, created.Street)
		assertEqual(t, p1.City, created.City)
		assertEqual(t, p1.StateCode, created.State)
		assertEqual(t, p1.Zip, created.Zip)
		resCreatedAt, err := time.Parse(time.RFC3339, created.CreatedAt)
		require.NoError(t, err)
		assertTimeRecent(t, resCreatedAt)

		// check location header
		// we should be able to use the location header to fetch the resource
		loc := res.Header.Get("Location")
		require.NotEmpty(t, loc)

		// get property from route in location header
		res = handleReq(t, s, getReq(t, loc, headers))
		require.Equal(t, http.StatusOK, res.StatusCode)
		assertApplicationJson(t, res.Header)

		var fetched rest.PropertyModel
		assert.NoError(t, json.NewDecoder(res.Body).Decode(&fetched))
		assertEqual(t, created.ID, fetched.ID)
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

		res := handleReq(t, s, putReq(t, route, body, headers))
		require.Equal(t, http.StatusCreated, res.StatusCode)
		var created rest.PropertyModel
		require.NoError(t, json.NewDecoder(res.Body).Decode(&created))

		// check response body for the echoed entity
		assertEqual(t, p1.ID, created.ID)
		assertEqual(t, p1.Street, created.Street)
		assertEqual(t, p1.City, created.City)
		assertEqual(t, p1.StateCode, created.State)
		assertEqual(t, p1.Zip, created.Zip)
		resCreatedAt, err := time.Parse(time.RFC3339, created.CreatedAt)
		require.NoError(t, err)
		require.False(t, resCreatedAt.IsZero())
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
		res := handleReq(t, s, putReq(t, route, body, headers))
		require.Equal(t, http.StatusCreated, res.StatusCode)

		// update property fixing typo mistake
		body["street"] = p1.Street
		res = handleReq(t, s, putReq(t, route, body, headers))
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var updated rest.PropertyModel
		require.NoError(t, json.NewDecoder(res.Body).Decode(&updated))

		// check response body for the echoed entity
		assertEqual(t, p1.ID, updated.ID)
		assertEqual(t, p1.Street, updated.Street)
		assertEqual(t, p1.City, updated.City)
		assertEqual(t, p1.StateCode, updated.State)
		assertEqual(t, p1.Zip, updated.Zip)
	})
}
func TestAddProperty_badRequest(t *testing.T) {
	var (
		headers map[string]string
		route   = "/property"
		repo    = repository.NewInMemoryRepo()
		s       = rest.NewServer(actions.NewActions(repo)).WithConfig(noAuthConf(t))
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
			res := handleReq(t, s, postReq(t, route, body, headers))
			assertEqual(t, http.StatusBadRequest, res.StatusCode)
		})
	}
}
func TestListProperties(t *testing.T) {
	var (
		route   = "/property"
		headers map[string]string
	)
	t.Run("200 expect empty set when no properties exist", func(t *testing.T) {
		var s = newServer(t)
		res := handleReq(t, s, getReq(t, route, headers))
		assertResCode(t, res, http.StatusOK)

		var propList rest.PropertyList
		require.NoError(t, json.NewDecoder(res.Body).Decode(&propList))
		assert.Len(t, propList.Items, 0)
	})
	t.Run("200 list two properties", func(t *testing.T) {
		var (
			repo = repository.NewInMemoryRepo()
			s    = rest.NewServer(actions.NewActions(repo)).
				WithConfig(noAuthConf(t))
		)

		// create two properties in propRepo
		p1 := fake.Property()
		p2 := fake.Property()
		require.NoError(t, repo.StoreProperty(ctx, p1))
		require.NoError(t, repo.StoreProperty(ctx, p2))

		// list via API
		res := handleReq(t, s, getReq(t, route, headers))
		require.Equal(t, http.StatusOK, res.StatusCode)

		var propList rest.PropertyList
		require.NoError(t, json.NewDecoder(res.Body).Decode(&propList))
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
		s         = rest.NewServer(actions.NewActions(repo)).
				WithConfig(noAuthConf(t))
	)
	t.Run("200 get property", func(t *testing.T) {
		var (
			p1    = fake.Property()
			route = routeBase + p1.ID
		)
		// create property in propRepo
		require.NoError(t, repo.StoreProperty(ctx, p1))

		// get property via API
		res := handleReq(t, s, getReq(t, route, headers))
		require.Equal(t, http.StatusOK, res.StatusCode)

		// check response data structure
		var resData rest.PropertyModel
		require.NoError(t, json.NewDecoder(res.Body).Decode(&resData))
		assertEqual(t, p1.ID, resData.ID)
		assertEqual(t, p1.Street, resData.Street)
		assertEqual(t, p1.City, resData.City)
		assertEqual(t, p1.StateCode, resData.State)
		assertEqual(t, p1.Zip, resData.Zip)
	})
	t.Run("404 get unknown property", func(t *testing.T) {
		var (
			p1    = fake.Property()
			route = routeBase + p1.ID
		)

		res := handleReq(t, s, getReq(t, route, headers))
		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}
func TestDeleteProperty(t *testing.T) {
	var (
		routeBase = "/property/"
		headers   map[string]string
		repo      = repository.NewInMemoryRepo()
		s         = rest.NewServer(actions.NewActions(repo)).
				WithConfig(noAuthConf(t))
	)
	t.Run("204 delete property", func(t *testing.T) {
		var (
			p1    = fake.Property()
			route = routeBase + p1.ID
		)
		// seed property
		require.NoError(t, repo.StoreProperty(ctx, p1))

		// del property via API
		res := handleReq(t, s, delReq(t, route, headers))
		require.Equal(t, http.StatusNoContent, res.StatusCode)

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
		res := handleReq(t, s, req)
		require.Equal(t, http.StatusNoContent, res.StatusCode)

		// del property again via API - idempotent check
		res = handleReq(t, s, req)
		require.Equal(t, http.StatusNoContent, res.StatusCode)
	})
}

func TestAccessViaAPIKeyAndSecret(t *testing.T) {
	var (
		repo = repository.NewInMemoryRepo()
		acts = actions.NewActions(repo)
	)

	tests := map[string]struct {
		k, s   string // key and secret from env
		rk, rs string // key and secret used in request header
		code   int    // expected response code
	}{
		// 401 when set and no match
		"secret1": {"", "c", "", "", 401},
		"secret2": {"", "c", "", "d", 401},
		"key1":    {"e", "", "", "", 401},
		"key2":    {"e", "", "f", "", 401},
		"both1":   {"g", "h", "", "", 401},
		"both2":   {"g", "h", "g", "", 401},
		"both3":   {"g", "h", "g", "a", 401},
		"both4":   {"g", "h", "", "h", 401},
		"both5":   {"g", "h", "a", "h", 401},
		"both6":   {"g", "h", "a", "b", 401},

		// 201 when unset or match
		"unset1":  {"", "", "", "", 201},
		"unset2":  {"", "", "a", "b", 201},
		"secret3": {"", "c", "", "c", 201},
		"key3":    {"e", "", "e", "", 201},
		"both7":   {"g", "h", "g", "h", 201},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				p1        = fake.Property()
				route     = "/property/" + p1.ID
				key       = tc.k
				secret    = tc.s
				reqKey    = tc.rk
				reqSecret = tc.rs
			)

			s := rest.NewServer(acts).WithConfig(authConf(t, key, secret))

			body := map[string]string{
				"street": p1.Street,
				"city":   p1.City,
				"state":  p1.StateCode,
				"zip":    p1.Zip,
			}

			res := handleReq(t, s, putReq(t, route, body, map[string]string{
				rest.HeaderAPIKey:    reqKey,
				rest.HeaderAPISecret: reqSecret,
			}))
			assertResCode(t, res, tc.code, "%s:%s %s:%s", tc.k, tc.s, tc.rk, tc.rs)
		})
	}
}
func TestHealth(t *testing.T) {
	var (
		headers map[string]string
		s       = newServer(t)
	)
	t.Run("health", func(t *testing.T) {
		res := handleReq(t, s, getReq(t, "/health", headers))
		require.Equal(t, http.StatusOK, res.StatusCode)
	})
	t.Run("ready", func(t *testing.T) {
		res := handleReq(t, s, getReq(t, "/health/ready", headers))
		require.Equal(t, http.StatusOK, res.StatusCode)
	})
	t.Run("live", func(t *testing.T) {
		res := handleReq(t, s, getReq(t, "/health/live", headers))
		require.Equal(t, http.StatusOK, res.StatusCode)
	})
}

func assertResCode(t testing.TB, res *http.Response, code int, msgAndArgs ...any) {
	t.Helper()
	if res.StatusCode != code {
		msg := fmtMsgAndArgs(msgAndArgs...)
		t.Fatalf(
			"unexpected response code\ngot  %d\nwant %d\nres  %s\nmsg  %v",
			res.StatusCode, code, internal.JSONString(t, res.Body), msg)
	}
}
func fmtMsgAndArgs(v ...any) string {
	switch len(v) {
	case 0:
		return ""
	case 1:
		return v[0].(string)
	default:
		var msg = v[0].(string)
		var args = v[1:]
		return fmt.Sprintf(msg, args...)
	}
}
func newServer(t testing.TB) *rest.Server {
	var repo = repository.NewInMemoryRepo()
	return rest.NewServer(actions.NewActions(repo)).WithConfig(noAuthConf(t))
}
