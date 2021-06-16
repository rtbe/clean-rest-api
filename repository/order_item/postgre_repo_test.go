package orderitem

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
	"github.com/rtbe/clean-rest-api/repository/order"
	"github.com/rtbe/clean-rest-api/repository/product"
	"github.com/rtbe/clean-rest-api/repository/user"
)

var pgUserRepo *user.Postgre
var pgProductRepo *product.Postgre
var pgOrderRepo *order.Postgre
var pgOrderItemRepo *Postgre
var validUser entity.User
var validProduct entity.Product
var validOrder entity.Order
var validOrderItem entity.OrderItem

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
		pgUserRepo = user.NewPostgreRepo(db, nil)
		pgProductRepo = product.NewPostgreRepo(db, nil)
		pgOrderRepo = order.NewPostgreRepo(db, nil)
		pgOrderItemRepo = NewPostgreRepo(db, nil)

		ctx := context.Background()
		newUser := entity.NewUser{
			UserName:        "AlanKay",
			FirstName:       "Alan",
			LastName:        "Kay",
			Password:        "OOP_is_about_messages",
			PasswordConfirm: "OOP_is_about_messages",
			Email:           "AlanKay@zeroxparc.com",
			Roles:           []string{"admin"},
		}
		validUser, err = pgUserRepo.Create(ctx, newUser)
		if err != nil {
			return err
		}

		newProduct := entity.NewProduct{
			Title:       "Apple Juice",
			Description: "Just an apple juice",
			Price:       1.01,
			Stock:       10,
		}
		validProduct, err = pgProductRepo.Create(ctx, newProduct)
		if err != nil {
			return err
		}

		newOrder := entity.NewOrder{
			UserID: validUser.ID,
			Status: "test",
		}
		validOrder, err = pgOrderRepo.Create(ctx, newOrder)
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
	t.Run("Given the need to create a order item inside PostgreSQL", func(t *testing.T) {
		tt := []struct {
			testName  string
			orderID   string
			productID string
			quantity  int
		}{
			{testName: "Create an order item", orderID: validOrder.ID, productID: validProduct.ID, quantity: 10},
		}

		for testID, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				newOrderItem := entity.NewOrderItem{
					OrderID:   tc.orderID,
					ProductID: tc.productID,
					Quantity:  tc.quantity,
				}

				savedOrderItem, err := pgOrderItemRepo.Create(context.Background(), newOrderItem)
				if err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to create a order item. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to create a order item.", tests.Success, testID)

				if newOrderItem.OrderID != savedOrderItem.OrderID {
					t.Fatalf("\t%s\tTest %d:\tWant order id: %s, got: %s", tests.Failed, testID, newOrderItem.OrderID, savedOrderItem.OrderID)
				}
				t.Logf("\t%s\tTest %d:\tWant order id: %s, got: %s", tests.Success, testID, newOrderItem.OrderID, savedOrderItem.OrderID)

				if newOrderItem.ProductID != savedOrderItem.ProductID {
					t.Fatalf("\t%s\tTest %d:\tWant product id: %s, got: %s", tests.Failed, testID, newOrderItem.ProductID, savedOrderItem.ProductID)
				}
				t.Logf("\t%s\tTest %d:\tWant product id: %s, got: %s", tests.Success, testID, newOrderItem.ProductID, savedOrderItem.ProductID)

				// Save product as valid product for further testing
				validOrderItem = savedOrderItem
			})
		}
	})

	t.Run("Given the need to query order items from PostgreSQL", func(t *testing.T) {
		tt := []struct {
			testName   string
			lastSeenID string
			limit      string
		}{
			{testName: "list existing order items", lastSeenID: validOrderItem.ID, limit: "1"},
		}
		for testID, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				orderItems, err := pgOrderItemRepo.Query(context.Background(), tc.lastSeenID, tc.limit)
				if err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to get a list of order items. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to get a list of order items.", tests.Success, testID)

				for _, orderItem := range orderItems {
					if tc.lastSeenID != orderItem.ID {
						t.Fatalf("\t%s\tTest %d:\tWant id: %s, got: %s", tests.Failed, testID, tc.lastSeenID, orderItem.ID)
					}
					t.Logf("\t%s\tTest %d:\tWant id: %s, got: %s", tests.Success, testID, tc.lastSeenID, orderItem.ID)
				}
			})
		}
	})

	t.Run("Given the need to get an order item by it`s id from PostgreSQL", func(t *testing.T) {
		tt := []struct {
			testName string
			id       string
			valid    bool
		}{
			{testName: "Order item with valid id", id: validOrderItem.ID, valid: true},
			{testName: "Order item with invalid id", id: "Not valid UUID id", valid: false},
		}
		for testID, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				retrievedOrderItem, err := pgOrderItemRepo.QueryByID(context.Background(), tc.id)
				if err != nil && tc.valid {
					t.Fatalf("\t%s\tTest %d:\tShould be able to get an order by it`s user id. Error: %s", tests.Failed, testID, err)
				}
				// Skip test if user with invalid id produce an error (wanted behavior)
				if err != nil && !tc.valid {
					t.SkipNow()
				}
				if err == nil && !tc.valid {
					t.Fatalf("\t%s\tTest %d:\tShould not be able to get a order item by it`s id.", tests.Failed, testID)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to get a order item by it`s id.", tests.Success, testID)

				if tc.id != retrievedOrderItem.ID {
					t.Fatalf("\t%s\tTest %d:\tWant id: %s, got: %s", tests.Failed, testID, tc.id, retrievedOrderItem.ID)
				}
				t.Logf("\t%s\tTest %d:\tWant id: %s, got: %s", tests.Success, testID, tc.id, retrievedOrderItem.ID)
			})
		}
	})

	t.Run("Given the need to get order items by their`s order id from PostgreSQL", func(t *testing.T) {
		tt := []struct {
			testName string
			orderID  string
			valid    bool
		}{
			{testName: "Order items with valid order id", orderID: validOrder.ID, valid: true},
			{testName: "Order items with invalid order id", orderID: "Not valid UUID id", valid: false},
		}
		for testID, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				orderItems, err := pgOrderItemRepo.QueryByOrderID(context.Background(), tc.orderID)
				if err != nil && tc.valid {
					t.Fatalf("\t%s\tTest %d:\tShould be able to get orders by their`s user id. Error: %s", tests.Failed, testID, err)
				}
				// Skip test if user with invalid id produce an error (wanted behavior)
				if err != nil && !tc.valid {
					t.SkipNow()
				}
				if err == nil && !tc.valid {
					t.Fatalf("\t%s\tTest %d:\tShould not be able to get order items by their`s id.", tests.Failed, testID)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to get order items by their`s id.", tests.Success, testID)

				for _, orderItem := range orderItems {
					if tc.orderID != orderItem.OrderID {
						t.Fatalf("\t%s\tTest %d:\tWant order id: %s, got: %s", tests.Failed, testID, tc.orderID, orderItem.OrderID)
					}
					t.Logf("\t%s\tTest %d:\tWant order id: %s, got: %s", tests.Success, testID, tc.orderID, orderItem.OrderID)
				}
			})
		}
	})

	t.Run("Given the need to update an order item inside PostgreSQL", func(t *testing.T) {
		tt := []struct {
			testName        string
			id              string
			updateOrderItem entity.UpdateOrderItem
		}{
			{
				testName: "Update a order item",
				id:       validOrderItem.ID,
				updateOrderItem: entity.UpdateOrderItem{
					Quantity: tests.IntPtr(2),
				}},
		}
		for testID, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				ctx := context.Background()

				err := pgOrderItemRepo.Update(ctx, tc.id, tc.updateOrderItem)
				if err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to update a order item. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to update an order item.", tests.Success, testID)

				retrievedOrderItem, err := pgOrderItemRepo.QueryByID(ctx, tc.id)
				if err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to get a order item by it`s id. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to get a order item by it`s id.", tests.Success, testID)

				if *tc.updateOrderItem.Quantity != retrievedOrderItem.Quantity {
					t.Fatalf("\t%s\tTest %d:\tWant quantity: %d, got: %d", tests.Failed, testID, *tc.updateOrderItem.Quantity, retrievedOrderItem.Quantity)
				}
				t.Logf("\t%s\tTest %d:\tWant quantity: %d, got: %d", tests.Success, testID, *tc.updateOrderItem.Quantity, retrievedOrderItem.Quantity)
			})
		}
	})

	t.Run("Given the need to delete a order item from PostgreSQL by it`s id", func(t *testing.T) {
		tt := []struct {
			testName     string
			newOrderItem entity.NewOrderItem
		}{
			{
				testName: "Delete an existing order item by it`s id",
				newOrderItem: entity.NewOrderItem{
					OrderID:   validOrder.ID,
					ProductID: validProduct.ID,
					Quantity:  1,
				}},
		}

		for testID, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				ctx := context.Background()

				savedOrderItem, err := pgOrderItemRepo.Create(ctx, tc.newOrderItem)
				if err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to create a order item. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to create a order item.", tests.Success, testID)

				if err := pgOrderItemRepo.Delete(ctx, savedOrderItem.ID); err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to delete a order item by it`s id. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to delete a order item by it`s id.", tests.Success, testID)
			})
		}
	})

	t.Run("Given the need to delete order items from PostgreSQL by their`s order id", func(t *testing.T) {
		tt := []struct {
			testName     string
			newOrderItem entity.NewOrderItem
		}{
			{
				testName: "Delete an existing order items by their`s order id",
				newOrderItem: entity.NewOrderItem{
					OrderID:   validOrder.ID,
					ProductID: validProduct.ID,
					Quantity:  1,
				}},
		}

		for testID, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				ctx := context.Background()

				savedOrderItem, err := pgOrderItemRepo.Create(ctx, tc.newOrderItem)
				if err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to create a order item. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to create a order item.", tests.Success, testID)

				if err := pgOrderItemRepo.DeleteByOrderID(ctx, savedOrderItem.OrderID); err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to delete a order item by their`s order id. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to delete order items by their`s order id.", tests.Success, testID)
			})
		}
	})
}
