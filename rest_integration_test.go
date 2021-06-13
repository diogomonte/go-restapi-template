package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

const ApiKey string = "123"
const MqttUrl string = "root:root@tcp(localhost:%s)/device_registry?charset=utf8&parseTime=True&loc=Local"

var app App

func TestMain(m *testing.M) {
	ctx := context.Background()

	mysqlContainer, _ := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "wangxian/alpine-mysql",
			ExposedPorts: []string{"3306/tcp", "33060/tcp"},
			Env: map[string]string{
				"MYSQL_ROOT_USER": 		"root",
				"MYSQL_ROOT_PASSWORD": 	"root",
				"MYSQL_PASSWORD": 		"root",
				"MYSQL_ROOT_USERNAME": 	"root",
				"MYSQL_DATABASE":      	"device_registry",
			},
			WaitingFor: wait.ForListeningPort("3306").WithStartupTimeout(45 * time.Second),
		},
		Started:          true,
	})
	defer mysqlContainer.Terminate(ctx)

	port, _ := mysqlContainer.MappedPort(ctx, "3306")

	app.Initialize(ApiKey, fmt.Sprintf(MqttUrl, port.Port()))
	code := m.Run()
	os.Exit(code)
}

func TestHttpSaveDevice(t *testing.T) {

	req, _ := http.NewRequest("POST", "/helloworld", bytes.NewBuffer([]byte("")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "ApiKey 123")

	response := httptest.NewRecorder()
	app.Router.ServeHTTP(response, req)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if response.Code != 200 {
		t.Errorf("expected response conde '200' but. Got '%d'", response.Code)
	}
}