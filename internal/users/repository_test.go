package users

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mramsden/go-rooms/internal/testutils"
	"github.com/testcontainers/testcontainers-go"

	_ "github.com/lib/pq"
)

func setupDatabase(ctx context.Context) (testcontainers.Container, *sql.DB) {
	container, db, err := testutils.CreateTestContainer(ctx, "testdb")
	if err != nil {
		log.Fatal(err)
	}

	mig, err := testutils.NewPgMigrator(db)
	if err != nil {
		log.Fatal(err)
	}

	err = mig.Up()
	if err != nil {
		log.Fatal(err)
	}

	return container, db
}

func TestUserRepository(t *testing.T) {
	ctx := context.Background()
	container, db := setupDatabase(ctx)
	defer db.Close()
	defer container.Terminate(ctx)

	r := &UsersRepository{db}
	created, err := r.CreateUser(ctx, "username", "password")
	if err != nil {
		t.Fatalf("failed to create user: %s", err)
	}

	// Can retrieve user by id
	retrievedById, err := r.GetUser(ctx, created.Id)
	if err != nil {
		t.Fatalf("failed to retrieve user by id: %s", err)
	}
	if created.Id != retrievedById.Id {
		t.Errorf("created.Id (%d) != retrievedById.Id (%d)", created.Id, retrievedById.Id)
	}
	if created.Username != retrievedById.Username {
		t.Errorf("created.Username (%s) != retrievedById.Username (%s)", created.Username, retrievedById.Username)
	}

	// Can retrieve user by username
	retrievedByUsername, err := r.GetUserByUsername(ctx, created.Username)
	if err != nil {
		t.Fatalf("failed to retrieve user by username: %s", err)
	}
	if created.Id != retrievedByUsername.Id {
		t.Errorf("created.Id (%d) != retrievedByUsername.Id (%d)", created.Id, retrievedByUsername.Id)
	}
	if created.Username != retrievedByUsername.Username {
		t.Errorf("created.Username (%s) != retrievedByUsername.Username (%s)", created.Username, retrievedByUsername.Username)
	}

	// Username must be unique
	_, err = r.CreateUser(ctx, "username", "password")
	if err == nil {
		t.Errorf("expected creating user with a duplicate username to fail")
	}

	// Unknown user id returns an error
	retrieved, err := r.GetUser(ctx, 2)
	if retrieved != nil {
		t.Errorf("expected no user to be retrieved for incorrect id")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("unexpected error for incorrect user id: %s", err)
	}

	// Unknown username returns an error
	retrieved, err = r.GetUserByUsername(ctx, "unknown")
	if retrieved != nil {
		t.Errorf("expected no user to be retrieved for incorrect username")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("unexpected error for incorrect username: %s", err)
	}
}
