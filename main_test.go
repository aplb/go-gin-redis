package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "encoding/base64"
    "encoding/json"
    "bytes"
)

func basicAuth(username string, pwd string) string {
    auth := username + ":" + pwd
    return base64.StdEncoding.EncodeToString([]byte(auth))
}

func TestPingRoute(t *testing.T) {
    router := setupRouter()
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/ping", nil)
    router.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Errorf("ping not ok")
    }

    if w.Body.String() != "pong" {
        t.Errorf("doesn't respond with pong")
    }
}

func TestUnauthorized(t *testing.T) {
    router := setupRouter()
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/admin", nil)
    req.Header.Add("Content-Type", "application/json")
    router.ServeHTTP(w, req)

    if w.Code != http.StatusUnauthorized {
        t.Errorf("user should be unauthorized")
    }
}

func TestMissingRequiredValue(t *testing.T) {
    router := setupRouter()
    w := httptest.NewRecorder()
    body, _ := json.Marshal(map[string]string{"some": "asd"})
    req, _ := http.NewRequest("POST", "/admin", bytes.NewBuffer(body))
    req.Header.Add("Authorization", "Basic " + basicAuth("foo", "123"))
    req.Header.Add("Content-Type", "application/json")
    router.ServeHTTP(w, req)

    if w.Code != http.StatusBadRequest {
        t.Errorf("should be Bad Request")
    }
}

func TestWriteUserValue(t *testing.T) {
    setupRedis()
    router := setupRouter()
    w := httptest.NewRecorder()
    body, _ := json.Marshal(map[string]string{"value": "asd"})
    req, _ := http.NewRequest("POST", "/admin", bytes.NewBuffer(body))
    req.Header.Add("Authorization", "Basic " + basicAuth("foo", "123"))
    req.Header.Add("Content-Type", "application/json")
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("user write not ok")
    }
}

