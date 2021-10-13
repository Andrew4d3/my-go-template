// Sample Code. Delete it once you understand how it works

package users

import (
	"go.mongodb.org/mongo-driver/bson"
)

// User ....
type User struct {
	repo UserRepository
}

// UserDocument ....
type UserDocument struct {
	ID   string `bson:"_id" json:"_id"`
	Name string `bson:"name" json:"name"`
}

// UserRepository ....
type UserRepository interface {
	FindOne(filter bson.M, result interface{}) error
}

// New ....
func New(repo UserRepository) *User {
	return &User{repo: repo}
}

// FindByName ....
func (u User) FindByName(name string) (UserDocument, error) {
	result := UserDocument{}

	if err := u.repo.FindOne(bson.M{"name": name}, &result); err != nil {
		return UserDocument{}, err
	}

	return result, nil
}
