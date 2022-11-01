package usersvc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ambowes87/betechtestv1.1/pkg/data"
	"github.com/ambowes87/betechtestv1.1/pkg/db"
	"github.com/ambowes87/betechtestv1.1/pkg/logger"
	"github.com/ambowes87/betechtestv1.1/pkg/notifications"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

const (
	// ServiceName the name of this service
	ServiceName = "UserService"
	// Version is the version of the running servicgo mod tidye
	Version = "1.0"

	notificationName = "user"
)

type UserService struct {
	address            string
	usersEndpoint      string
	port               int
	notificationBroker *notifications.Broker
	userStore          db.UserStore

	healthEndpoint string
}

func New(addr, ep string, port int, nb *notifications.Broker, us db.UserStore) *UserService {
	return &UserService{
		address:            addr,
		usersEndpoint:      ep,
		port:               port,
		notificationBroker: nb,
		userStore:          us,
		healthEndpoint:     "/health",
	}
}

// Start starts the service listening for requests, always returns a non-nil error on exit
func (s *UserService) Start() error {
	http.HandleFunc(s.usersEndpoint, s.handleUsersRequest)
	http.HandleFunc(s.healthEndpoint, s.handleHealthRequest)
	s.userStore.Open()
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", s.address, s.port), nil)
	if s.userStore != nil {
		s.userStore.Close()
	}
	return err
}

func (s *UserService) handleUsersRequest(w http.ResponseWriter, r *http.Request) {
	correlationID := uuid.New()
	logger.Log(fmt.Sprintf("%s | %s?%s | [%s]", r.Method, r.URL.EscapedPath(), r.URL.RawQuery, correlationID.String()))

	switch r.Method {
	case http.MethodPost:
		s.addUser(w, r, correlationID)
	case http.MethodGet:
		s.getUser(w, r, correlationID)
	case http.MethodPut:
		s.updateUser(w, r, correlationID)
	case http.MethodDelete:
		s.deleteUser(w, r, correlationID)
	default:
		handleError(http.StatusBadRequest, fmt.Sprintf("unsupported request type [%s]", r.Method), w, correlationID)
	}
}

func (s *UserService) handleHealthRequest(w http.ResponseWriter, r *http.Request) {
	correlationID := uuid.New()
	logger.Log(fmt.Sprintf("%s | %s?%s | [%s]", r.Method, r.URL.EscapedPath(), r.URL.RawQuery, correlationID.String()))

	w.Header().Add("Service-Version", Version)
	if s.userStore.Ping() != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (s *UserService) addUser(w http.ResponseWriter, r *http.Request, cid uuid.UUID) {
	userData, err := extractUserFromBody(r)
	if err != nil {
		handleError(http.StatusBadRequest, err.Error(), w, cid)
		return
	}
	err = s.userStore.Add(userData)
	if err != nil {
		handleError(http.StatusInternalServerError, err.Error(), w, cid)
		return
	}
	s.notificationBroker.Publish(notificationName, "add")
	writeSuccess(http.StatusOK, "Added user", w, cid)
}

func (s *UserService) getUser(w http.ResponseWriter, r *http.Request, cid uuid.UUID) {
	id, err := getUserID(r)
	if err != nil {
		handleError(http.StatusBadRequest, err.Error(), w, cid)
		return
	}
	user, err := s.userStore.Get(id)
	if err != nil {
		handleError(http.StatusInternalServerError, err.Error(), w, cid)
		return
	}
	buf, err := json.Marshal(user)
	if err != nil {
		handleError(http.StatusInternalServerError, err.Error(), w, cid)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	writeSuccessWithBody(http.StatusOK, buf, "Got user", w, cid)
}

func (s *UserService) updateUser(w http.ResponseWriter, r *http.Request, cid uuid.UUID) {
	userData, err := extractUserFromBody(r)
	if err != nil {
		handleError(http.StatusBadRequest, err.Error(), w, cid)
		return
	}
	err = s.userStore.Update(userData)
	if err != nil {
		handleError(http.StatusInternalServerError, err.Error(), w, cid)
		return
	}
	s.notificationBroker.Publish(notificationName, "update")
	writeSuccess(http.StatusOK, "Updated user", w, cid)
}

func (s *UserService) deleteUser(w http.ResponseWriter, r *http.Request, cid uuid.UUID) {
	id, err := getUserID(r)
	if err != nil {
		handleError(http.StatusBadRequest, err.Error(), w, cid)
		return
	}
	err = s.userStore.Delete(id)
	if err != nil {
		handleError(http.StatusInternalServerError, err.Error(), w, cid)
		return
	}
	s.notificationBroker.Publish(notificationName, "delete")
	writeSuccess(http.StatusOK, "Deleted user", w, cid)
}

func getUserID(r *http.Request) (string, error) {
	userID := r.FormValue("id")
	if len(userID) == 0 {
		return "", errors.New("No user ID provided on request")
	}
	return userID, nil
}

func extractUserFromBody(r *http.Request) (data.UserData, error) {
	buf, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return data.UserData{}, err
	}

	var userData data.UserData
	err = json.Unmarshal(buf, &userData)
	if err != nil {
		return data.UserData{}, err
	}
	return userData, nil
}

func handleError(responseCode int, errMsg string, w http.ResponseWriter, cid uuid.UUID) {
	logger.Log(errMsg + fmt.Sprintf(" [%s]", cid.String()))
	w.WriteHeader(responseCode)
}

func writeSuccess(responseCode int, successMsg string, w http.ResponseWriter, cid uuid.UUID) {
	logger.Log(successMsg + fmt.Sprintf(" [%s]", cid.String()))
	w.WriteHeader(responseCode)
}

func writeSuccessWithBody(responseCode int, body []byte, successMsg string, w http.ResponseWriter, cid uuid.UUID) {
	logger.Log(successMsg + fmt.Sprintf(" [%s]", cid.String()))
	w.Write(body)
}
