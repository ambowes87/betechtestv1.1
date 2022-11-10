package usersvc

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ambowes87/betechtestv1.1/pkg/data"
	"github.com/ambowes87/betechtestv1.1/pkg/db"
	"github.com/google/uuid"
)

var (
	service   *UserService
	addr      = "localhost"
	ep        = "/test"
	port      = 9090
	cid       uuid.UUID
	userStore = db.NewUserMemStore()
)

func TestMain(m *testing.M) {
	service = New(addr, ep, port, nil, userStore)
	cid = uuid.New()
	service.Start()

	exitCode := m.Run()
	service = nil
	os.Exit(exitCode)
}

func Test_Add(t *testing.T) {
	buf, err := json.Marshal(data.UserData{
		ID: "testID",
	})
	if err != nil {
		t.FailNow()
	}
	req := httptest.NewRequest(http.MethodPost, ep, bytes.NewReader(buf))
	w := httptest.NewRecorder()

	service.addUser(w, req, cid)

	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Errorf("Expected status 200, got [%d]", res.StatusCode)
	}
}
