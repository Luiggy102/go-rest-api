package repository

import (
	"context"

	"github.com/Luiggy102/go-rest-ws/models"
)

// methods for Repository
type Repository interface {
	InsertUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, id string) (*models.User, error)
	UserEmailExits(ctx context.Context, email string) (bool, error)
	InsertPost(ctx context.Context, post *models.Post) error            // create
	GetPostById(ctx context.Context, id string) (*models.Post, error)   // read
	UpdatePost(ctx context.Context, post *models.Post) error            // update
	DeletePost(ctx context.Context, id string, userId string) error     // delete
	ListPosts(ctx context.Context, page uint64) ([]*models.Post, error) // pagination
	Close() error
}

var implementation Repository

// UserRepo is equal at any implementation of the interface
// different databases for example
func SetRepo(repository Repository) { // postgres implementation, mongodb implementation, etc
	// dependency injection
	implementation = repository
}

func InsertUser(ctx context.Context, user *models.User) error {
	return implementation.InsertUser(ctx, user) // implementation
}

func GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return implementation.GetUserByID(ctx, id) // implementation
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return implementation.GetUserByEmail(ctx, email) // implementation
}

func UserEmailExists(ctx context.Context, email string) (bool, error) {
	return implementation.UserEmailExits(ctx, email)
}

func InsertPost(ctx context.Context, post *models.Post) error {
	return implementation.InsertPost(ctx, post)
}

func GetPostById(ctx context.Context, id string) (*models.Post, error) {
	return implementation.GetPostById(ctx, id)
}
func UpdatePost(ctx context.Context, post *models.Post) error {
	return implementation.UpdatePost(ctx, post)
}
func DeletePost(ctx context.Context, id string, userId string) error {
	return implementation.DeletePost(ctx, id, userId)
}
func ListPosts(ctx context.Context, page uint64) ([]*models.Post, error) {
	return implementation.ListPosts(ctx, page)
}

func Close() error {
	return implementation.Close()
}
