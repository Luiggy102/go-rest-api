package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Luiggy102/go-rest-ws/models"
	_ "github.com/lib/pq"
)

// create a PostgresRepo for implement a userRepo
type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(url string) (*PostgresRepo, error) { // constructor
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepo{db: db}, nil
}

// implement the user repository interface
// (GetUserById - InsertUser - Close)

func (prepo *PostgresRepo) InsertUser(ctx context.Context, user *models.User) error {
	// ExecContext execute the query without returning any rows
	_, err := prepo.db.ExecContext(ctx,
		"insert into users (id, email, password) values ($1, $2, $3);",
		user.Id,
		user.Email,
		user.Password,
	)
	return err
}

func (prepo *PostgresRepo) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	// QueryContext returns the rows of the query result
	rows, err := prepo.db.QueryContext(ctx, "select id, email from users where id = $1;", id)
	if err != nil {
		return nil, err
	}
	// IMPORTANT
	// always close the rows reader
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// transform the query result into a user struct
	u := models.User{}
	for rows.Next() {
		err = rows.Scan(&u.Id, &u.Email)
		if err != nil {
			return nil, err
		}
	}
	return &u, nil
}

func (prepo *PostgresRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	rows, err := prepo.db.QueryContext(ctx, "select id, email, password from users where email = $1;", email)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	u := models.User{}
	for rows.Next() {
		err = rows.Scan(&u.Id, &u.Email, &u.Password)
		if err != nil {
			return nil, err
		}
	}
	return &u, nil
}

func (prepo *PostgresRepo) UserEmailExits(ctx context.Context, email string) (bool, error) {
	q := fmt.Sprintf(" select exists(select email from users where email='%s');", email)
	row := prepo.db.QueryRowContext(ctx, q)
	var dbResult string
	err := row.Scan(&dbResult)
	if err != nil {
		return false, err
	}
	if dbResult == "true" {
		return true, nil
	} else {
		return false, nil
	}
}

// insert post in the db
// create
func (prepo *PostgresRepo) InsertPost(ctx context.Context, post *models.Post) error {
	_, err := prepo.db.ExecContext(ctx,
		"insert into posts (id, post_content, user_id) values ($1, $2, $3);",
		post.Id,
		post.PostContent,
		post.UserId,
	)
	return err
}

// read
func (prepo *PostgresRepo) GetPostById(ctx context.Context, id string) (*models.Post, error) {
	rows, err := prepo.db.QueryContext(ctx,
		"select id, post_content, created_at, user_id from posts where id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	p := models.Post{}
	for rows.Next() {
		err = rows.Scan(&p.Id, &p.PostContent, &p.CreatedAt, &p.UserId)
		if err != nil {
			return nil, err
		}
	}
	return &p, nil
}

// update
func (prepo *PostgresRepo) UpdatePost(ctx context.Context, post *models.Post) error {
	_, err := prepo.db.ExecContext(ctx,
		"update posts set post_content = $1 where id = $2 and user_id = $3",
		post.PostContent, post.Id, post.UserId,
	)
	return err
}

// delete
func (prepo *PostgresRepo) DeletePost(ctx context.Context, id string, userId string) error {
	_, err := prepo.db.ExecContext(ctx,
		"delete from posts where id = $1 and user_id = $2",
		id, userId,
	)
	return err
}

// close the PostgresRepo db
func (prepo *PostgresRepo) Close() error {
	err := prepo.db.Close()
	if err != nil {
		return err
	}
	return nil
}
