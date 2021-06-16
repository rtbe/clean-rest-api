package user

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/rtbe/clean-rest-api/domain/entity"
	"github.com/rtbe/clean-rest-api/internal/tests"
	"golang.org/x/crypto/bcrypt"
)

var pgUserRepo *Postgre
var validUser entity.User

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	absFilepath, _ := filepath.Abs("../../internal/tests")
	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12.3",
		Env: []string{
			"POSTGRES_USER=" + tests.PgUser,
			"POSTGRES_PASSWORD=" + tests.PgPassword,
			"POSTGRES_DB=" + tests.PgDB,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: tests.PgPort},
			},
		},
		Mounts: []string{absFilepath + ":/docker-entrypoint-initdb.d/"},
	}

	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		log.Fatalf("could not start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		db, err := sqlx.Connect("postgres", fmt.Sprintf(
			"postgres://%s:%s@localhost:%s/%s?sslmode=disable",
			tests.PgUser,
			tests.PgPassword,
			resource.GetPort("5432/tcp"),
			tests.PgDB,
		))
		if err != nil {
			return err
		}

		// Init global package dependencies after
		// successfull connection to a database
		pgUserRepo = NewPostgreRepo(db, nil)

		return db.Ping()
	}); err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	code := m.Run()

	// When you're done, kill and remove the container
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestPostgre(t *testing.T) {
	t.Run("Given the need to create a user inside PostgreSQL", func(t *testing.T) {
		tt := []struct {
			testName        string
			userName        string
			firstName       string
			lastName        string
			password        string
			passwordConfirm string
			email           string
			roles           []string
		}{
			{testName: "Create a single user", userName: "AlanKay", firstName: "Alan", lastName: "Kay", password: "OOP_is_about_messages", passwordConfirm: "OOP_is_about_messages", email: "alanKay@rocks.com", roles: []string{"admin"}},
		}
		for testID, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				newUser := entity.NewUser{
					UserName:        tc.userName,
					FirstName:       tc.firstName,
					LastName:        tc.lastName,
					Password:        tc.password,
					PasswordConfirm: tc.passwordConfirm,
					Email:           tc.email,
					Roles:           tc.roles,
				}

				savedUser, err := pgUserRepo.Create(context.Background(), newUser)
				if err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to create a user. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to create user.", tests.Success, testID)

				if newUser.UserName != savedUser.UserName {
					t.Fatalf("\t%s\tTest %d:\tWant user name: %s, got: %s", tests.Failed, testID, newUser.UserName, savedUser.UserName)
				}
				t.Logf("\t%s\tTest %d:\tWant user name: %s, got: %s", tests.Success, testID, newUser.UserName, savedUser.UserName)

				if newUser.FirstName != savedUser.FirstName {
					t.Fatalf("\t%s\tTest %d:\tWant first name: %s, got: %s", tests.Failed, testID, newUser.FirstName, savedUser.FirstName)
				}
				t.Logf("\t%s\tTest %d:\tWant first name: %s, got: %s", tests.Success, testID, newUser.FirstName, savedUser.FirstName)

				if newUser.LastName != savedUser.LastName {
					t.Fatalf("\t%s\tTest %d:\tWant last name: %s, got: %s", tests.Failed, testID, newUser.LastName, savedUser.LastName)
				}
				t.Logf("\t%s\tTest %d:\tWant last name: %s, got: %s", tests.Success, testID, newUser.LastName, savedUser.LastName)

				if newUser.Email != savedUser.Email {
					t.Fatalf("\t%s\tTest %d:\tWant email : %s, got: %s", tests.Failed, testID, newUser.Email, savedUser.Email)
				}
				t.Logf("\t%s\tTest %d:\tWant email : %s, got: %s", tests.Success, testID, newUser.Email, savedUser.Email)

				if err := bcrypt.CompareHashAndPassword(savedUser.Password, []byte(newUser.Password)); err != nil {
					t.Fatalf("\t%s\tTest %d:\tWant password to be equal to saved password", tests.Failed, testID)
				}
				t.Logf("\t%s\tTest %d:\tWant password to be equal to saved password", tests.Success, testID)

				for i, role := range newUser.Roles {
					if role != savedUser.Roles[i] {
						t.Fatalf("\t%s\tTest %d:\tWant role: %s, got: %s", tests.Failed, testID, role, savedUser.Roles[i])
					}
				}
				t.Logf("\t%s\tTest %d:\tWant equal roles", tests.Success, testID)

				// Save user as valid user for further testing
				validUser = savedUser
			})
		}
	})

	t.Run("Given the need to get a user by his id from PostgreSQL", func(t *testing.T) {
		tt := []struct {
			testName string
			id       string
			valid    bool
		}{
			{testName: "User with valid id", id: validUser.ID, valid: true},
			{testName: "User with invalid id", id: "Not valid UUID id", valid: false},
		}
		for testID, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				retrievedUser, err := pgUserRepo.QueryByID(context.Background(), tc.id)
				if err != nil && tc.valid {
					t.Fatalf("\t%s\tTest %d:\tShould be able to get a user by his id. Error: %s", tests.Failed, testID, err)
				}
				// Skip test if user with invalid id produce an error (wanted behavior)
				if err != nil && !tc.valid {
					t.SkipNow()
				}
				if err == nil && !tc.valid {
					t.Fatalf("\t%s\tTest %d:\tShould not be able to get a user by his id.", tests.Failed, testID)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to get a user by his id.", tests.Success, testID)

				if tc.id != retrievedUser.ID {
					t.Fatalf("\t%s\tTest %d:\tWant id: %s, got: %s", tests.Failed, testID, tc.id, retrievedUser.ID)
				}
				t.Logf("\t%s\tTest %d:\tWant id: %s, got: %s", tests.Success, testID, tc.id, retrievedUser.ID)
			})
		}
	})

	t.Run("Given the need to query users from PostgreSQL", func(t *testing.T) {
		tt := []struct {
			testName   string
			lastSeenID string
			limit      string
		}{
			{testName: "list existing user", lastSeenID: validUser.ID, limit: "1"},
		}
		for testID, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				users, err := pgUserRepo.Query(context.Background(), tc.lastSeenID, tc.limit)
				if err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to get a list of users. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to get a list of users.", tests.Success, testID)

				for _, user := range users {
					if tc.lastSeenID != user.ID {
						t.Fatalf("\t%s\tTest %d:\tWant id: %s, got: %s", tests.Failed, testID, tc.lastSeenID, user.ID)
					}
					t.Logf("\t%s\tTest %d:\tWant id: %s, got: %s", tests.Success, testID, tc.lastSeenID, user.ID)
				}
			})
		}
	})

	t.Run("Given the need to update a user inside PostgreSQL", func(t *testing.T) {
		tt := []struct {
			testName   string
			id         string
			updateUser entity.UpdateUser
		}{
			{
				testName: "Update a user",
				id:       validUser.ID,
				updateUser: entity.UpdateUser{
					UserName:        tests.StrPtr("KenThompson"),
					FirstName:       tests.StrPtr("Ken"),
					LastName:        tests.StrPtr("Thompson"),
					Password:        tests.StrPtr("CAndUnix"),
					PasswordConfirm: tests.StrPtr("CAndUnix"),
					Email:           tests.StrPtr("KenThompson@bellabs.com"),
					Roles:           []string{"admin", "user"},
				}},
		}
		for testID, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				ctx := context.Background()

				err := pgUserRepo.Update(ctx, tc.id, tc.updateUser)
				if err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to update a user. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to update a user.", tests.Success, testID)

				retrievedUser, err := pgUserRepo.QueryByID(ctx, tc.id)
				if err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to get a user by his id. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to get a user by his id.", tests.Success, testID)

				if *tc.updateUser.UserName != retrievedUser.UserName {
					t.Fatalf("\t%s\tTest %d:\tWant user name: %s, got: %s", tests.Failed, testID, *tc.updateUser.UserName, retrievedUser.UserName)
				}
				t.Logf("\t%s\tTest %d:\tWant user name: %s, got: %s", tests.Success, testID, *tc.updateUser.UserName, retrievedUser.UserName)

				if *tc.updateUser.FirstName != retrievedUser.FirstName {
					t.Fatalf("\t%s\tTest %d:\tWant first name: %s, got: %s", tests.Failed, testID, *tc.updateUser.FirstName, retrievedUser.FirstName)
				}
				t.Logf("\t%s\tTest %d:\tWant first name: %s, got: %s", tests.Success, testID, *tc.updateUser.FirstName, retrievedUser.FirstName)

				if *tc.updateUser.LastName != retrievedUser.LastName {
					t.Fatalf("\t%s\tTest %d:\tWant last name: %s, got: %s", tests.Failed, testID, *tc.updateUser.LastName, retrievedUser.LastName)
				}
				t.Logf("\t%s\tTest %d:\tWant last name: %s, got: %s", tests.Success, testID, *tc.updateUser.LastName, retrievedUser.LastName)

				if err := bcrypt.CompareHashAndPassword(retrievedUser.Password, []byte(*tc.updateUser.Password)); err != nil {
					t.Fatalf("\t%s\tTest %d:\tWant password to be equal to retrieved password", tests.Failed, testID)
				}
				t.Logf("\t%s\tTest %d:\tWant password to be equal to retrieved password", tests.Success, testID)

				if err := bcrypt.CompareHashAndPassword(retrievedUser.Password, []byte(*tc.updateUser.PasswordConfirm)); err != nil {
					t.Fatalf("\t%s\tTest %d:\tWant confirm password to be equal to retrieved password", tests.Failed, testID)
				}
				t.Logf("\t%s\tTest %d:\tWant confirm password to be equal to retrieved password", tests.Success, testID)

				if *tc.updateUser.Email != retrievedUser.Email {
					t.Fatalf("\t%s\tTest %d:\tWant email : %s, got: %s", tests.Failed, testID, *tc.updateUser.Email, retrievedUser.Email)
				}
				t.Logf("\t%s\tTest %d:\tWant email : %s, got: %s", tests.Success, testID, *tc.updateUser.Email, retrievedUser.Email)

				for i, role := range tc.updateUser.Roles {
					if role != retrievedUser.Roles[i] {
						t.Fatalf("\t%s\tTest %d:\tWant role: %s, got: %s", tests.Failed, testID, role, tc.updateUser.Roles[i])
					}
				}
				t.Logf("\t%s\tTest %d:\tWant equal roles", tests.Success, testID)
			})
		}
	})

	t.Run("Given the need to delete a user from PostgreSQL by his id", func(t *testing.T) {
		tt := []struct {
			testName string
			newUser  entity.NewUser
		}{
			{
				testName: "Delete a existing user by his id",
				newUser: entity.NewUser{
					UserName:        "Test",
					FirstName:       "Test",
					LastName:        "Test",
					Password:        "testtest",
					PasswordConfirm: "testtest",
					Email:           "test@gmail.com",
					Roles:           []string{"test"},
				}},
		}

		for testID, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				ctx := context.Background()

				savedUser, err := pgUserRepo.Create(ctx, tc.newUser)
				if err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to create a user. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to create a user.", tests.Success, testID)

				if err := pgUserRepo.Delete(ctx, savedUser.ID); err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to delete a user by his id. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to delete a user by his id.", tests.Success, testID)
			})
		}
	})

	t.Run("Given the need to delete a user from PostgreSQL by his user name", func(t *testing.T) {
		tt := []struct {
			testName string
			newUser  entity.NewUser
		}{
			{
				testName: "Delete a existing user by his user name",
				newUser: entity.NewUser{
					UserName:        "Test",
					FirstName:       "Test",
					LastName:        "Test",
					Password:        "testtest",
					PasswordConfirm: "testtest",
					Email:           "test@gmail.com",
					Roles:           []string{"test"},
				}},
		}

		for testID, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				ctx := context.Background()

				savedUser, err := pgUserRepo.Create(ctx, tc.newUser)
				if err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to create a user. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to create a user.", tests.Success, testID)

				if err := pgUserRepo.DeleteByUserName(ctx, savedUser.UserName); err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to delete a user by his user name. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to delete a user by his user name.", tests.Success, testID)
			})
		}
	})
}
