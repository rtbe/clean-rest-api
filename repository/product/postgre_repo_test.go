package product

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

var pgProductRepo *Postgre
var pgUserRepo *user.Postgre
var validUser entity.User
var validProduct entity.Product

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
		pgProductRepo = NewPostgreRepo(db, nil)
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
	t.Run("Given the need to create a product inside PostgreSQL", func(t *testing.T) {
		tt := []struct {
			testName    string
			title       string
			description string
			price       float32
			stock       int
		}{
			{testName: "Create a single product", title: "Apple Juice", description: "Just an apple juice", price: 10.2, stock: 10},
		}

		for testID, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				newProduct := entity.NewProduct{
					Title:       tc.title,
					Description: tc.description,
					Price:       tc.price,
					Stock:       tc.stock,
				}

				savedProduct, err := pgProductRepo.Create(context.Background(), newProduct)
				if err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to create a product. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to create a product.", tests.Success, testID)

				if newProduct.Title != savedProduct.Title {
					t.Fatalf("\t%s\tTest %d:\tWant title: %s, got: %s", tests.Failed, testID, newProduct.Title, savedProduct.Title)
				}
				t.Logf("\t%s\tTest %d:\tWant title: %s, got: %s", tests.Success, testID, newProduct.Title, savedProduct.Title)

				if newProduct.Description != savedProduct.Description {
					t.Fatalf("\t%s\tTest %d:\tWant description: %s, got: %s", tests.Failed, testID, newProduct.Description, savedProduct.Description)
				}
				t.Logf("\t%s\tTest %d:\tWant description: %s, got: %s", tests.Success, testID, newProduct.Description, savedProduct.Description)

				if newProduct.Price != savedProduct.Price {
					t.Fatalf("\t%s\tTest %d:\tWant price: %.2f, got: %.2f", tests.Failed, testID, newProduct.Price, savedProduct.Price)
				}
				t.Logf("\t%s\tTest %d:\tWant price: %.2f, got: %.2f", tests.Success, testID, newProduct.Price, savedProduct.Price)

				if newProduct.Stock != savedProduct.Stock {
					t.Fatalf("\t%s\tTest %d:\tWant stock: %d, got: %d", tests.Failed, testID, newProduct.Stock, savedProduct.Stock)
				}
				t.Logf("\t%s\tTest %d:\tWant stock: %d, got: %d", tests.Success, testID, newProduct.Stock, savedProduct.Stock)

				// Save product as valid product for further testing
				validProduct = savedProduct
			})
		}
	})

	t.Run("Given the need to get a product by it`s id from PostgreSQL", func(t *testing.T) {
		tt := []struct {
			testName string
			id       string
			valid    bool
		}{
			{testName: "Product with valid id", id: validProduct.ID, valid: true},
			{testName: "Product with invalid id", id: "Not valid UUID id", valid: false},
		}
		for testID, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				retrievedProduct, err := pgProductRepo.QueryByID(context.Background(), tc.id)
				if err != nil && tc.valid {
					t.Fatalf("\t%s\tTest %d:\tShould be able to get a product by it`s id. Error: %s", tests.Failed, testID, err)
				}
				// Skip test if user with invalid id produce an error (wanted behavior)
				if err != nil && !tc.valid {
					t.SkipNow()
				}
				if err == nil && !tc.valid {
					t.Fatalf("\t%s\tTest %d:\tShould not be able to get a product by it`s id.", tests.Failed, testID)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to get a product by it`s id.", tests.Success, testID)

				if tc.id != retrievedProduct.ID {
					t.Fatalf("\t%s\tTest %d:\tWant id: %s, got: %s", tests.Failed, testID, tc.id, retrievedProduct.ID)
				}
				t.Logf("\t%s\tTest %d:\tWant id: %s, got: %s", tests.Success, testID, tc.id, retrievedProduct.ID)
			})
		}
	})

	t.Run("Given the need to query products from PostgreSQL", func(t *testing.T) {
		tt := []struct {
			testName   string
			lastSeenID string
			limit      string
		}{
			{testName: "list existing products", lastSeenID: validProduct.ID, limit: "1"},
		}
		for testID, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				products, err := pgProductRepo.Query(context.Background(), tc.lastSeenID, tc.limit)
				if err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to get a list of products. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to get a list of products.", tests.Success, testID)

				for _, product := range products {
					if tc.lastSeenID != product.ID {
						t.Fatalf("\t%s\tTest %d:\tWant id: %s, got: %s", tests.Failed, testID, tc.lastSeenID, product.ID)
					}
					t.Logf("\t%s\tTest %d:\tWant id: %s, got: %s", tests.Success, testID, tc.lastSeenID, product.ID)
				}
			})
		}
	})

	t.Run("Given the need to update a product inside PostgreSQL", func(t *testing.T) {
		tt := []struct {
			testName      string
			id            string
			updateProduct entity.UpdateProduct
		}{
			{
				testName: "Update a product",
				id:       validProduct.ID,
				updateProduct: entity.UpdateProduct{
					Title:       tests.StrPtr("Orange juice"),
					Description: tests.StrPtr("Just an orange juice"),
					Price:       tests.Float32Ptr(2.02),
					Stock:       tests.IntPtr(33),
				}},
		}
		for testID, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				ctx := context.Background()

				err := pgProductRepo.Update(ctx, tc.id, tc.updateProduct)
				if err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to update a product. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to update a product.", tests.Success, testID)

				retrievedProduct, err := pgProductRepo.QueryByID(ctx, tc.id)
				if err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to get a product by it`s id. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to get a product by it`s id.", tests.Success, testID)

				if *tc.updateProduct.Title != retrievedProduct.Title {
					t.Fatalf("\t%s\tTest %d:\tWant title: %s, got: %s", tests.Failed, testID, *tc.updateProduct.Title, retrievedProduct.Title)
				}
				t.Logf("\t%s\tTest %d:\tWant title: %s, got: %s", tests.Success, testID, *tc.updateProduct.Title, retrievedProduct.Title)

				if *tc.updateProduct.Description != retrievedProduct.Description {
					t.Fatalf("\t%s\tTest %d:\tWant description: %s, got: %s", tests.Failed, testID, *tc.updateProduct.Description, retrievedProduct.Description)
				}
				t.Logf("\t%s\tTest %d:\tWant description: %s, got: %s", tests.Success, testID, *tc.updateProduct.Description, retrievedProduct.Description)

				if *tc.updateProduct.Price != retrievedProduct.Price {
					t.Fatalf("\t%s\tTest %d:\tWant price: %.2f, got: %.2f", tests.Failed, testID, *tc.updateProduct.Price, retrievedProduct.Price)
				}
				t.Logf("\t%s\tTest %d:\tWant price: %.2f, got: %.2f", tests.Success, testID, *tc.updateProduct.Price, retrievedProduct.Price)

				if *tc.updateProduct.Stock != retrievedProduct.Stock {
					t.Fatalf("\t%s\tTest %d:\tWant stock: %d, got: %d", tests.Failed, testID, *tc.updateProduct.Stock, retrievedProduct.Stock)
				}
				t.Logf("\t%s\tTest %d:\tWant price: %d, got: %d", tests.Success, testID, *tc.updateProduct.Stock, retrievedProduct.Stock)

			})
		}
	})

	t.Run("Given the need to delete a product from PostgreSQL by it`s id", func(t *testing.T) {
		tt := []struct {
			testName   string
			newProduct entity.NewProduct
		}{
			{
				testName: "Delete a existing product by it`s id",
				newProduct: entity.NewProduct{
					Title:       "test",
					Description: "test",
					Price:       1.01,
					Stock:       1,
				}},
		}

		for testID, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				ctx := context.Background()

				savedProduct, err := pgProductRepo.Create(ctx, tc.newProduct)
				if err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to create a product. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to create a product.", tests.Success, testID)

				if err := pgProductRepo.Delete(ctx, savedProduct.ID); err != nil {
					t.Fatalf("\t%s\tTest %d:\tShould be able to delete a product by it`s id. Error: %s", tests.Failed, testID, err)
				}
				t.Logf("\t%s\tTest %d:\tShould be able to delete a product by it`s id.", tests.Success, testID)
			})
		}
	})
}
