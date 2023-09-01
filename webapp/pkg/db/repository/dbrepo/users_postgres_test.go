package dbrepo

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/felicia/testing_course/webapp/pkg/data"
	"github.com/felicia/testing_course/webapp/pkg/db/repository"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "postgres"
	dbName   = "users_test"
	port     = "5435"
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5"
)

var resource *dockertest.Resource
var pool *dockertest.Pool
var testDB *sql.DB
var testRepo repository.DatabaseRepo

func TestMain(m *testing.M) {
	//connect to docker
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker; is it running?%s", err)
	}
	pool = p
	fmt.Println("Initializing Docker pool...")
	//set up docker option, specifying the image and so forth

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14.5",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	//gets a resource (docker image)
	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		//_ = pool.Purge(resource)
		log.Fatalf("could not connect to resource: %s", err)
	}
	fmt.Println("Docker pool connected successfully.")
	//start the image and wait until its ready
	if err := pool.Retry(func() error {
		var err error
		testDB, err = sql.Open("pgx", fmt.Sprintf(dsn, host, port, user, password, dbName))
		if err != nil {
			log.Println("Error:", err)
			return err
		}
		defer testDB.Close()
		return testDB.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not connect to database: %s", err)
	}
	fmt.Println("Connection database successfully.")
	//populate database with empty tables
	err = createTables()
	if err != nil {
		log.Fatalf("error creating tables: %s", err)
	}
	fmt.Println("Creating table successfully.")
	testRepo = &PostgresDBRepo{DB: testDB}
	//run tests
	code := m.Run()
	//clean up
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}

	os.Exit(code)
}

func createTables() error {
	tableSQL, err := os.ReadFile("./testdata/users.sql")
	if err != nil {
		fmt.Print(err)
		return err
	}
	_, err = testDB.Exec(string(tableSQL))

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func Test_pingDB(t *testing.T) {
	err := testDB.Ping()
	if err != nil {
		t.Error("can't ping database")
	}
}

func Test_RepoInsertUser(t *testing.T) {
	testUser := data.User{
		FirstName: "Admin",
		LastName:  "User",
		Email:     "admin@example.com",
		Password:  "secret",
		IsAdmin:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	id, err := testRepo.InsertUser(testUser)
	if err != nil {
		t.Errorf("Insert user returned error: %s", err)
	}
	if id != 1 {
		t.Errorf("Insert user returned wrong id; expected got 1 but got %d", id)
	}
}

func Test_DBRepoAllUsers(t *testing.T) {
	users, err := testRepo.AllUsers()
	if err != nil {
		t.Errorf("all users report an error: %s", err)
	}
	if len(users) != 1 {
		t.Errorf("all users report wrong size; expected got 1 but got %d", len(users))
	}

	testUser := data.User{
		FirstName: "Jack",
		LastName:  "Smith",
		Email:     "jack@smith.com",
		Password:  "secret",
		IsAdmin:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, _ = testRepo.InsertUser(testUser)

	users, err = testRepo.AllUsers()

	if len(users) != 2 {
		t.Errorf("all users report wrong size after insert user; expected got 2 but got %d", len(users))
	}
}

func Test_DBRepoGetUsers(t *testing.T) { //get user by id
	user, err := testRepo.GetUser(1)
	if err != nil {
		t.Errorf("error getting users by id: %s", err)
	}

	if user.Email != "admin@example.com" {
		t.Errorf("wrong email returned by GetUser; expected admin@example.com but got %s", user.Email)
	}

	_, _ = testRepo.GetUser(3)
	if err == nil {
		t.Error("no error reported when getting non-existent user by getUserid")
	}
}

func Test_DBRepoGetUserByEmail(t *testing.T) { //get user by email
	user, err := testRepo.GetUserByEmail("jack@smith.com")
	if err != nil {
		t.Errorf("error getting users by email:%s", err)
	}
	if user.ID != 2 {
		t.Errorf("wrong id returned by GetUserByEmail; expected 2 but got %d", user.ID)
	}
}

func Test_DBRepoUpdateUser(t *testing.T) {
	user, err := testRepo.GetUser(2)
	user.FirstName = "Jane"
	user.Email = "jane@smith.com"
	err = testRepo.UpdateUser(*user)
	if err != nil {
		t.Errorf("error updating user %d : %s", user.ID, err)
	}
	user, _ = testRepo.GetUser(2)
	if user.FirstName != "Jane" || user.Email != "jane@smith.com" {
		t.Errorf("expected got Jane and email jane@smith.com, but got %s %s", user.FirstName, user.Email)
	}
}

func Test_DBRepoDeleteUser(t *testing.T) {
	err := testRepo.DeleteUser(2)
	if err != nil {
		t.Errorf("error delete user 2: %s", err)
	}
	_, err = testRepo.GetUser(2)
	if err == nil {
		t.Error("retrieved user id 2, which should have been deleted")
	}
}