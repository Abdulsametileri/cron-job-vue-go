package userrepo

import (
	"context"
	"github.com/Abdulsametileri/cron-job-vue-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Repo interface {
	AddUser(user models.User) error
	GetUserByTelegramId(telegramId int) (models.User, error)
}

type userRepo struct {
	collection *mongo.Collection
}

func NewUserRepository(col *mongo.Collection) Repo {
	return &userRepo{
		collection: col,
	}
}

func (u userRepo) AddUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := u.collection.InsertOne(ctx, user)
	return err
}

func (u userRepo) GetUserByTelegramId(telegramId int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"telegramId": telegramId,
	}

	var user models.User
	err := u.collection.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return models.User{}, nil
	}

	return user, err
}
