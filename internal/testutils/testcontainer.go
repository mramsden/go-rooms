package testutils

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func CreateTestContainer(ctx context.Context, dbname string) (testcontainers.Container, *sql.DB, error) {
	env := map[string]string{
		"POSTGRES_PASSWORD": "password",
		"POSTGRES_USER":     "postgres",
		"POSTGRES_DB":       dbname,
	}
	port := "5432/tcp"
	dbURL := func(port nat.Port) string {
		return fmt.Sprintf("postgres://postgres:password@localhost:%s/%s?sslmode=disable", port.Port(), dbname)
	}

	// Mute logging from testcontainers, remove this to see messages
	testcontainers.Logger = log.New(ioutil.Discard, "", 0)

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:12",
			ExposedPorts: []string{port},
			Cmd:          []string{"postgres", "-c", "fsync=off"},
			Env:          env,
			WaitingFor:   wait.ForSQL(nat.Port(port), "postgres", dbURL).Timeout(time.Second * 5),
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, nil, fmt.Errorf("failed to start container: %s", err)
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(port))
	if err != nil {
		return container, nil, fmt.Errorf("failed to get container external port: %s", err)
	}

	url := fmt.Sprintf("postgres://postgres:password@localhost:%s/%s?sslmode=disable", mappedPort.Port(), dbname)
	db, err := sql.Open("postgres", url)
	if err != nil {
		return container, db, fmt.Errorf("failed to establish database connection: %s", err)
	}

	return container, db, nil
}
