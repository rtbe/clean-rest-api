package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rtbe/clean-rest-api/domain/entity"
	"github.com/rtbe/clean-rest-api/internal/database"
	"github.com/rtbe/clean-rest-api/internal/logger"
	"golang.org/x/crypto/bcrypt"
)

// Postgre is an abstraction layer that manages user entities inside PostgreSQL DB.
type Postgre struct {
	db  *sqlx.DB
	log logger.Logger
}

// NewPostgreRepo creates a new PostgreSQL repository for User entity.
func NewPostgreRepo(db *sqlx.DB, l logger.Logger) *Postgre {
	return &Postgre{
		db:  db,
		log: l,
	}
}

// Create creates a new user in PostgreSQL.
func (r *Postgre) Create(ctx context.Context, nu entity.NewUser) (entity.User, error) {
	const q = `
	INSERT INTO users 
		(user_id, user_name, first_name, last_name, password, email, roles, date_created, date_updated) 
	VALUES
		(:user_id, :user_name, :first_name, :last_name, :password, :email, :roles, :date_created, :date_updated) `

	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return entity.User{}, errors.Wrap(err, "generating password hash")
	}

	u := entity.User{
		ID:          uuid.NewString(),
		UserName:    nu.UserName,
		FirstName:   nu.FirstName,
		LastName:    nu.LastName,
		Email:       nu.Email,
		Password:    hash,
		Roles:       nu.Roles,
		DateCreated: time.Now().UTC(),
		DateUpdated: time.Now().UTC(),
	}

	if _, err := r.db.NamedExecContext(ctx, q, u); err != nil {
		return entity.User{}, errors.Wrapf(err, "inserting a user")
	}

	return u, nil
}

// Query gets users from PostgreSQL.
// This query uses two provided values to implement pagination: last seen id and limit.
// Results of a query sorted by creation date of selected users.
func (r *Postgre) Query(ctx context.Context, lastSeenID, limit string) ([]entity.User, error) {
	query := `
	SELECT 
		* 
	FROM 
		users 
	WHERE 
		user_id <= :last_seen_id 
	ORDER BY 
		user_id DESC
	FETCH FIRST :limit ROWS ONLY`

	var users []entity.User

	data := struct {
		LastSeenID string `db:"last_seen_id"`
		Limit      string `db:"limit"`
	}{
		LastSeenID: lastSeenID,
		Limit:      limit,
	}

	err := database.QuerySlice(ctx, r.db, query, data, &users)
	if err != nil {
		return nil, errors.Wrap(err, "selecting users")
	}

	return users, nil
}

// QueryByID gets a user from PostgreSQL by given user id.
func (r *Postgre) QueryByID(ctx context.Context, userID string) (entity.User, error) {
	const q = `
	SELECT 
		* 
	FROM 
		users 
	WHERE 
		user_id = :user_id`

	data := struct {
		ID string `db:"user_id"`
	}{
		ID: userID,
	}

	var user entity.User

	err := database.QueryStruct(ctx, r.db, q, data, &user)
	if err != nil {
		return entity.User{}, errors.Wrapf(err, "getting a user with id %s", userID)
	}

	return user, nil
}

// Update updates a user inside PostgreSQL.
func (r *Postgre) Update(ctx context.Context, id string, uu entity.UpdateUser) error {
	u, err := r.QueryByID(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "error updating a user with id %s", id)
	}

	// You should not update not updated user fields.
	if uu.UserName != nil {
		u.UserName = *uu.UserName
	}
	if uu.FirstName != nil {
		u.FirstName = *uu.FirstName
	}
	if uu.LastName != nil {
		u.LastName = *uu.LastName
	}
	if uu.Email != nil {
		u.Email = *uu.Email
	}
	if uu.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*uu.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.Wrap(err, "generating password hash")
		}
		u.Password = hash
	}
	if uu.Roles != nil {
		u.Roles = uu.Roles
	}
	u.DateUpdated = time.Now().UTC()

	const q = `
	UPDATE 
		users 
	SET	
		"user_name" = :user_name, 
		"first_name" = :first_name, 
		"last_name" = :last_name, 
		"password" = :password,
		"email" = :email,
		"roles" = :roles,
		"date_updated" = :date_updated 
	WHERE 
		user_id = :user_id`

	_, err = r.db.NamedExecContext(ctx, q, u)
	if err != nil {
		return errors.Wrapf(err, "error updating a user with id %s", id)
	}

	return nil
}

// Delete deletes user from PostgreSQL by given user id.
func (r *Postgre) Delete(ctx context.Context, userID string) error {
	const q = `
	DELETE FROM 
		users 
	WHERE
		user_id = :user_id`

	data := struct {
		ID string `db:"user_id"`
	}{
		ID: userID,
	}

	if _, err := r.db.NamedExecContext(ctx, q, data); err != nil {
		return errors.Wrapf(err, "deleting a user with id %s", userID)
	}

	return nil
}

// DeleteByUserName deletes user from PostgreSQL by given user name.
func (r *Postgre) DeleteByUserName(ctx context.Context, userName string) error {
	const q = `
	DELETE FROM 
		users 
	WHERE 
		user_name = :user_name`

	data := struct {
		ID string `db:"user_name"`
	}{
		ID: userName,
	}

	if _, err := r.db.NamedExecContext(ctx, q, data); err != nil {
		return errors.Wrapf(err, "deleting a user with user_name %s", userName)
	}

	return nil
}
