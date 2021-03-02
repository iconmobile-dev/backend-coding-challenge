package user

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/iconmobile-dev/go-coding-challenge/lib/storage"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	serverTest      *Server
	failingDBServer *Server
	ts              *httptest.Server
	failingDBTs     *httptest.Server
)

func TestMain(m *testing.M) {
	// setup before tests
	var err error

	// bootstrap logger and config
	SetupLoggerAndConfig("engagement", true)

	// database
	db, err := storage.NewDB(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.SSLMode)
	if err != nil {
		log.Error(err, "database postgres new")
		os.Exit(1)
	}
	log.Verbose("connected to Postgres at", cfg.DB.Host)
	err = db.Reset()
	if err != nil {
		log.Error(err, "database reset")
		os.Exit(1)
	}

	// cache
	cache, err := storage.NewCache(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password)
	if err != nil {
		log.Error(err, "cache database redis new")
		os.Exit(1)
	}
	log.Verbose("connected to redis at", cfg.DB.Host)

	err = cache.Reset()
	if err != nil {
		log.Error(err, "cache reset")
		os.Exit(1)
	}

	// init server for test
	serverTest = New(db, cache)

	ts = httptest.NewServer(serverTest)

	// init server with failing DB for test
	{
		failingDB := &storage.DB{}
		db, err := sqlx.Open("postgres", "")
		if err != nil {
			log.Error(err, "failed to open database")
			os.Exit(1)
		}
		failingDB.DB = db
		failingDBServer = New(failingDB, cache)
		failingDBTs = httptest.NewServer(failingDBServer)
	}

	// run tests
	code := m.Run()

	// shutdown after tests
	db.Close()
	cache.Close()
	ts.Close()

	os.Exit(code)
}

// PostRequest sends a POST request with user ID, Role required headers and object as JSON
// required headers are dumy values
func PostRequest(myURL string, data interface{}) (*http.Response, error) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrapf(err, "POST request to %s, JSON Marshal error", myURL)
	}
	client := &http.Client{}
	req, _ := http.NewRequest("POST", myURL, bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Language", "en")

	return client.Do(req)
}

// utility func to load response
func loadFromResponse(resp *http.Response, obj interface{}) error {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, obj)
}

// mustPostRequest wraps PostRequest and fails the test if an error or an unexpected http status code is returned
func mustPostRequest(t *testing.T, myURL string, data interface{}, expectedStatusCode int) *http.Response {
	resp, err := PostRequest(myURL, data)
	if !assert.NoError(t, err) || !assert.Equal(t, expectedStatusCode, resp.StatusCode) {
		t.FailNow()
	}

	return resp
}

// mustLoadFromResponse wraps loadFromResponse and fails the test if an error is returned
func mustLoadFromResponse(t *testing.T, resp *http.Response, obj interface{}) {
	err := loadFromResponse(resp, obj)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
}
