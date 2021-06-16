package order

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
	"github.com/rtbe/clean-rest-api/repository/user"
)

var pgOrderRepo *Postgre
var pgUserRepo *user.Postgre
var validUser entity.User
var validOrder entity.Order

func TestMain(m *testing.M) {
	var err error
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
		var err error
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
		// Init global package dependencies after successfull connection to a database
		pgOrderRepo = NewPostgreRepo(db, nil)
		pgUserRepo = user.NewPostgreRepo(db, nil)

		newUser := entity.NewUser{
			UserName:        "AlanKay",
			FirstName:       "Alan",
			LastName:        "Kay",
			Password:        "OOP_is_about_messages",
			PasswordConfirm: "OOP_is_about_messages",
			Email:           "AlanKay@zeroxparc.com",
			Roles:           []string{"admin"},
		}
		validUser, err = pgUserRepo.Create(context.Background(), newUser)
		if err != nil {
			return err
		}

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
	t.Run("Given the need to create a order inside PostgreSQL", func(t *testing.T) {
		tt := []struct {
			testName string
			userID   string
			status   string
		}{
			{testName: "Create an single order", userID: validUser.ID, status: "test"},
		}

		for testID, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				newOrder := entity.NewOrder{
					UserID: tc.userID,
					Status: tc.status,
				}

				savedOrder, err := pgOrderRepo.Create(context.Background(), newOrder)
				if err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to create a product. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to create a product.", tests.Success, testID)

				if newOrder.UserID != savedOrder.UserID {
					t.Fatalf("\t%s\tTest %d:\tWant user id: %s, got: %s", tests.Failed, testID, newOrder.UserID, savedOrder.UserID)
				}
				t.Logf("\t%s\tTest %d:\tWant user id: %s, got: %s", tests.Success, testID, newOrder.UserID, savedOrder.UserID)

				if newOrder.Status != savedOrder.Status {
					t.Fatalf("\t%s\tTest %d:\tWant status: %s, got: %s", tests.Failed, testID, newOrder.Status, savedOrder.Status)
				}
				t.Logf("\t%s\tTest %d:\tWant status: %s, got: %s", tests.Success, testID, newOrder.Status, savedOrder.Status)

				// Save product as valid product for further testing
				validOrder = savedOrder
			})
		}
	})

	t.Run("Given the need to get an order by it`s id from PostgreSQL", func(t *testing.T) {
		tt := []struct {
			testName string
			id       string
			valid    bool
		}{
			{testName: "Order with valid id", id: validOrder.ID, valid: true},
			{testName: "Order with invalid id", id: "Not valid UUID id", valid: false},
		}
		for testID, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				retrievedOrder, err := pgOrderRepo.QueryByID(context.Background(), tc.id)
				if err != nil && tc.valid {
					t.Fatalf("\t%s\tTest %d:\tShould be able to get an order by it`s id. Error: %s", tests.Failed, testID, err)
				}
				// Skip test if user with invalid id produce an error (wanted behavior)
				if err != nil && !tc.valid {
					t.SkipNow()
				}
				if err == nil && !tc.valid {
					t.Fatalf("\t%s\tTest %d:\tShould not be able to get a product by it`s id.", tests.Failed, testID)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to get a product by it`s id.", tests.Success, testID)

				if tc.id != retrievedOrder.ID {
					t.Fatalf("\t%s\tTest %d:\tWant id: %s, got: %s", tests.Failed, testID, tc.id, retrievedOrder.ID)
				}
				t.Logf("\t%s\tTest %d:\tWant id: %s, got: %s", tests.Success, testID, tc.id, retrievedOrder.ID)
			})
		}
	})

	t.Run("Given the need to get an order by it`s user id from PostgreSQL", func(t *testing.T) {
		tt := []struct {
			testName string
			userID   string
			valid    bool
		}{
			{testName: "Order with valid id", userID: validOrder.UserID, valid: true},
			{testName: "Order with invalid id", userID: "Not valid UUID id", valid: false},
		}
		for testID, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				retrievedOrders, err := pgOrderRepo.QueryByUserID(context.Background(), tc.userID)
				if err != nil && tc.valid {
					t.Fatalf("\t%s\tTest %d:\tShould be able to get an order by it`s user id. Error: %s", tests.Failed, testID, err)
				}
				// Skip test if user with invalid id produce an error (wanted behavior)
				if err != nil && !tc.valid {
					t.SkipNow()
				}
				if err == nil && !tc.valid {
					t.Fatalf("\t%s\tTest %d:\tShould not be able to get a product by it`s id.", tests.Failed, testID)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to get a product by it`s id.", tests.Success, testID)

				for _, order := range retrievedOrders {
					if tc.userID != order.UserID {
						t.Fatalf("\t%s\tTest %d:\tWant id: %s, got: %s", tests.Failed, testID, tc.userID, order.UserID)
					}
					t.Logf("\t%s\tTest %d:\tWant id: %s, got: %s", tests.Success, testID, tc.userID, order.UserID)
				}
			})
		}
	})

	t.Run("Given the need to query orders from PostgreSQL", func(t *testing.T) {
		tt := []struct {
			testName   string
			lastSeenID string
			limit      string
		}{
			{testName: "List existing orders", lastSeenID: validOrder.ID, limit: "1"},
		}
		for testID, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				orders, err := pgOrderRepo.Query(context.Background(), tc.lastSeenID, tc.limit)
				if err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to get a list of orders. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to get a list of orders.", tests.Success, testID)

				for _, order := range orders {
					if tc.lastSeenID != order.ID {
						t.Fatalf("\t%s\tTest %d:\tWant id: %s, got: %s", tests.Failed, testID, tc.lastSeenID, order.ID)
					}
					t.Logf("\t%s\tTest %d:\tWant id: %s, got: %s", tests.Success, testID, tc.lastSeenID, order.ID)
				}
			})
		}
	})

	t.Run("Given the need to update an order inside PostgreSQL", func(t *testing.T) {
		tt := []struct {
			testName    string
			id          string
			updateOrder entity.UpdateOrder
		}{
			{
				testName: "Update an order",
				id:       validOrder.ID,
				updateOrder: entity.UpdateOrder{
					Status: tests.StrPtr("test"),
				}},
		}
		for testID, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				ctx := context.Background()

				err := pgOrderRepo.Update(ctx, tc.id, tc.updateOrder)
				if err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to update an order. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to update an order.", tests.Success, testID)

				retrievedOrder, err := pgOrderRepo.QueryByID(ctx, tc.id)
				if err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to get an order by it`s id. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to get a order by it`s id.", tests.Success, testID)

				if *tc.updateOrder.Status != retrievedOrder.Status {
					t.Fatalf("\t%s\tTest %d:\tWant title: %s, got: %s", tests.Failed, testID, *tc.updateOrder.Status, retrievedOrder.Status)
				}
				t.Logf("\t%s\tTest %d:\tWant title: %s, got: %s", tests.Success, testID, *tc.updateOrder.Status, retrievedOrder.Status)
			})
		}
	})

	t.Run("Given the need to delete an order from PostgreSQL by it`s id", func(t *testing.T) {
		tt := []struct {
			testName string
			newOrder entity.NewOrder
		}{
			{
				testName: "Delete an existing order by it`s id",
				newOrder: entity.NewOrder{
					UserID: validUser.ID,
					Status: "test",
				}},
		}

		for testID, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				ctx := context.Background()

				savedOrder, err := pgOrderRepo.Create(ctx, tc.newOrder)
				if err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to create an order. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to create an order.", tests.Success, testID)

				if err := pgOrderRepo.Delete(ctx, savedOrder.ID); err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to delete an order by it`s id. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to delete an order by it`s id.", tests.Success, testID)
			})
		}
	})

	t.Run("Given the need to delete orders from PostgreSQL by their`s user id", func(t *testing.T) {
		tt := []struct {
			testName string
			newOrder entity.NewOrder
		}{
			{
				testName: "Delete existing products by their`s user id",
				newOrder: entity.NewOrder{
					UserID: validUser.ID,
					Status: "test",
				}},
		}

		for testID, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				ctx := context.Background()

				savedOrder, err := pgOrderRepo.Create(ctx, tc.newOrder)
				if err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to create an order. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to create an order.", tests.Success, testID)

				if err := pgOrderRepo.DeleteByUserID(ctx, savedOrder.ID); err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to delete orders by their`s id. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to delete orders by their`s id.", tests.Success, testID)
			})
		}
	})
}
