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

	"github.com/google/uuid"
)

const (
	// ServiceName the name of this service
	ServiceName = "UserService"
)

// HandleRequest handles a call to add, get, update or delete a user
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	correlationID := uuid.New()
	logger.Log(fmt.Sprintf("%s | %s?%s | [%s]", r.Method, r.URL.EscapedPath(), r.URL.RawQuery, correlationID.String()))

	switch r.Method {
	case http.MethodPost:
		addUser(w, r, correlationID)
	case http.MethodGet:
		getUser(w, r, correlationID)
	case http.MethodPut:
		updateUser(w, r, correlationID)
	case http.MethodDelete:
		deleteUser(w, r, correlationID)
	default:
		handleError(http.StatusBadRequest, fmt.Sprintf("unsupported request type [%s]", r.Method), w, correlationID)
	}
}

func addUser(w http.ResponseWriter, r *http.Request, cid uuid.UUID) {
	userData, err := extractUserFromBody(r)
	if err != nil {
		handleError(http.StatusBadRequest, err.Error(), w, cid)
		return
	}
	err = db.AddUser(userData)
	if err != nil {
		handleError(http.StatusInternalServerError, err.Error(), w, cid)
		return
	}
	writeSuccess(http.StatusOK, "Added user", w, cid)
}

func getUser(w http.ResponseWriter, r *http.Request, cid uuid.UUID) {
	id, err := getUserID(r)
	if err != nil {
		handleError(http.StatusBadRequest, err.Error(), w, cid)
		return
	}
	user, err := db.GetUser(id)
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

func updateUser(w http.ResponseWriter, r *http.Request, cid uuid.UUID) {
	userData, err := extractUserFromBody(r)
	if err != nil {
		handleError(http.StatusBadRequest, err.Error(), w, cid)
		return
	}
	err = db.UpdateUser(userData)
	if err != nil {
		handleError(http.StatusInternalServerError, err.Error(), w, cid)
		return
	}
	writeSuccess(http.StatusOK, "Updated user", w, cid)
}

func deleteUser(w http.ResponseWriter, r *http.Request, cid uuid.UUID) {
	id, err := getUserID(r)
	if err != nil {
		handleError(http.StatusBadRequest, err.Error(), w, cid)
		return
	}
	err = db.DeleteUser(id)
	if err != nil {
		handleError(http.StatusInternalServerError, err.Error(), w, cid)
		return
	}
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
	logger.Log(errMsg + fmt.Sprintf("[%s]", cid.String()))
	w.WriteHeader(responseCode)
}

func writeSuccess(responseCode int, successMsg string, w http.ResponseWriter, cid uuid.UUID) {
	logger.Log(successMsg + fmt.Sprintf("[%s]", cid.String()))
	w.WriteHeader(responseCode)
}

func writeSuccessWithBody(responseCode int, body []byte, successMsg string, w http.ResponseWriter, cid uuid.UUID) {
	logger.Log(successMsg + fmt.Sprintf("[%s]", cid.String()))
	w.Write(body)
}
