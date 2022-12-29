package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"ot-recorder/app/response"
	"ot-recorder/app/server"
	"ot-recorder/infrastructure/config"
	"ot-recorder/infrastructure/db"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/suite"
)

type e2eTestSuite struct {
	suite.Suite
	db           *sql.DB
	dbMigration  *migrate.Migrate
	apiBaseURL   string
	hooksBaseURL string
}

const (
	username = "dev"
	device   = "phone"
	clientIP = "127.0.0.1"
)

var (
	epoch      = time.Now().Unix()
	pingReqStr = fmt.Sprintf(`{"_type":"location","acc":13,"alt":-42,"batt":40,"bs":1, "created_at":%d,
"lat":23.0000000,"lon":90.0000000,"m":1,"tid":"p1","topic":"owntracks/dev/phone","tst":%d,"vac":1,"vel":0}`,
		epoch, epoch)
	wifiInfo = `,"BSSID":"c0:00:00:00:00:00","SSID":"dev-test"}`
)

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

func (s *e2eTestSuite) SetupSuite() {
	s.Require().NoError(config.Load("./config.yml"))

	cfg := config.Get()
	connStr := fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/%s?multiStatements=true",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	var err error

	s.dbMigration, err = migrate.New("file://../../infrastructure/db/migrations/mysql", connStr)
	s.Require().NoError(err)

	if err := s.dbMigration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		s.Require().NoError(err)
	}

	serverReady := make(chan bool)
	httpServer := server.Server{ServerReady: serverReady}

	go httpServer.Serve()

	// wait until api server is start
	<-serverReady

	s.db = db.GetClient()
	s.apiBaseURL = fmt.Sprintf("http://localhost:%d/api/v1", cfg.App.Port)
	s.hooksBaseURL = fmt.Sprintf("http://localhost:%d/hooks", cfg.App.Port)
}

func (s *e2eTestSuite) TearDownSuite() {
	p, _ := os.FindProcess(syscall.Getpid())
	_ = p.Signal(syscall.SIGINT)
}

func (s *e2eTestSuite) SetupTest() {
	if err := s.dbMigration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		s.Require().NoError(err)
	}
}

func (s *e2eTestSuite) TearDownTest() {
	s.NoError(s.dbMigration.Down())
}

func (s *e2eTestSuite) Test_EndToEnd_Ping_WIFI() {
	reqStr := pingReqStr[0:len(pingReqStr)-1] + wifiInfo
	body := postPing(s, reqStr)

	s.Equal(`[]`, strings.Trim(string(body), "\n"))
}

func (s *e2eTestSuite) Test_EndToEnd_Ping_Mobile() {
	body := postPing(s, pingReqStr)

	s.Equal(`[]`, strings.Trim(string(body), "\n"))
}

func (s *e2eTestSuite) Test_EndToEnd_Ping_Waypoint() {
	reqStr := fmt.Sprintf(`{"_type":"waypoint","desc":"home","lat":23.0000000,"lon":90.0000000,"rad":50,
"topic":"owntracks/dev/phone/waypoints","tst":%d}`, epoch)
	body := postPing(s, reqStr)

	s.Equal(`[]`, strings.Trim(string(body), "\n"))
}

func (s *e2eTestSuite) Test_EndToEnd_Last_Location() {
	reqStr := pingReqStr[0:len(pingReqStr)-1] + wifiInfo
	_ = postPing(s, reqStr)

	req, err := http.NewRequestWithContext(
		context.Background(),
		echo.GET, s.apiBaseURL+"/last-location?username="+username,
		nil,
	)
	s.NoError(err)

	client := http.Client{}
	res, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusOK, res.StatusCode)

	byteBody, err := io.ReadAll(res.Body)

	s.NoError(err)

	_ = res.Body.Close()

	var r response.Response

	s.NoError(json.Unmarshal(byteBody, &r))

	resultsMap := r.Data.(map[string]interface{})
	s.Equal(device, resultsMap["device"])
	s.Equal("dev-test", resultsMap["wifi_name"])
}

func (s *e2eTestSuite) Test_EndToEnd_Last_Location_Not_Found() {
	req, err := http.NewRequestWithContext(
		context.Background(),
		echo.GET, s.apiBaseURL+"/last-location?username="+username,
		nil,
	)
	s.NoError(err)

	client := http.Client{}
	res, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusNotFound, res.StatusCode)

	_ = res.Body.Close()
}

func (s *e2eTestSuite) Test_EndToEnd_Telegram_Hook() {
	reqStr := pingReqStr[0:len(pingReqStr)-1] + wifiInfo
	_ = postPing(s, reqStr)

	reqStr = fmt.Sprintf(`{"update_id":1,"message":{"message_id":1,"message_thread_id":1,"date":%d,
"text":"/loc dev","from":{"is_bot":true,"id":1,"first_name":"test","last_name":"test","username":"test"},
"chat":{"id":-123,"type":"private","title":"test_group","username":"test","first_name":"test","last_name":"test"}}}`,
		epoch)
	body := postTelegramHook(s, reqStr, "secret")

	bodystr := string(body)
	s.Contains(bodystr, "Username: *dev*")
	s.Contains(bodystr, "openstreetmap.org")
}

func (s *e2eTestSuite) Test_EndToEnd_Telegram_Hook_Help() {
	reqStr := fmt.Sprintf(`{"update_id":1,"message":{"message_id":1,"message_thread_id":1,"date":%d,
"text":"/help","from":{"is_bot":true,"id":1,"first_name":"test","last_name":"test","username":"test"},
"chat":{"id":-123,"type":"private","title":"test_group","username":"test","first_name":"test","last_name":"test"}}}`,
		epoch)
	body := postTelegramHook(s, reqStr, "secret")

	bodystr := string(body)
	s.Contains(bodystr, "/help")
}

func postPing(s *e2eTestSuite, payload string) []byte {
	req, err := http.NewRequestWithContext(
		context.Background(),
		echo.POST,
		s.apiBaseURL+"/ping",
		strings.NewReader(payload),
	)

	s.NoError(err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("x-limit-u", username)
	req.Header.Set("x-limit-d", device)
	req.Header.Set("X-Real-IP", clientIP)

	client := http.Client{}
	res, err := client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusOK, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	s.NoError(err)

	_ = res.Body.Close()

	return body
}

func postTelegramHook(s *e2eTestSuite, payload, secret string) []byte {
	req, err := http.NewRequestWithContext(
		context.Background(),
		echo.POST,
		s.hooksBaseURL+"/telegram",
		strings.NewReader(payload),
	)

	s.NoError(err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("X-Telegram-Bot-Api-Secret-Token", secret)

	client := http.Client{}
	res, err := client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusOK, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	s.NoError(err)

	_ = res.Body.Close()

	return body
}
