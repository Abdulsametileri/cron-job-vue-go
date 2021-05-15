package userrepo

import (
	"context"
	"github.com/Abdulsametileri/cron-job-vue-go/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func setupDB() (*mongo.Client, error) {
	client, err := mongo.
		Connect(context.Background(),
			options.Client().
				ApplyURI("mongodb://localhost:27017/testReminder"))

	return client, err
}

func setupCollection(client *mongo.Client) (*mongo.Collection, error) {
	userCollection := client.Database("testReminder").Collection("users")
	errDrop := userCollection.Drop(context.Background())
	return userCollection, errDrop
}

func cleanCollection(t *testing.T, userCollection *mongo.Collection) {
	errDrop := userCollection.Drop(context.Background())
	require.NoError(t, errDrop)
}

func TestUserRepo_AddUser(t *testing.T) {
	client, errSetupDB := setupDB()
	require.NoError(t, errSetupDB)

	userCollection, errCollection := setupCollection(client)
	require.NoError(t, errCollection)

	user := models.User{
		Token:               "123",
		TelegramId:          123,
		TelegramDisplayName: "Test",
	}
	userRepo := NewUserRepository(userCollection)
	errAddUser := userRepo.AddUser(user)
	require.NoError(t, errAddUser)

	defer cleanCollection(t, userCollection)
}

func TestUserRepo_GetUserByTelegramId(t *testing.T) {
	client, errSetupDB := setupDB()
	require.NoError(t, errSetupDB)

	userCollection, errCollection := setupCollection(client)
	require.NoError(t, errCollection)

	user := models.User{
		Token:               "123",
		TelegramId:          123,
		TelegramDisplayName: "Test",
	}

	_, errAddUser := userCollection.InsertOne(context.Background(), user)
	require.NoError(t, errAddUser)

	userRepo := NewUserRepository(userCollection)
	foundedUser, errGetUser := userRepo.GetUserByTelegramId(user.TelegramId)
	require.NoError(t, errGetUser)

	assert.Equal(t, user.Token, foundedUser.Token)
	assert.Equal(t, user.TelegramId, foundedUser.TelegramId)
	assert.Equal(t, user.TelegramDisplayName, foundedUser.TelegramDisplayName)

	defer cleanCollection(t, userCollection)
}

func TestUserRepo_GetUserByTelegramId_NotExistUser(t *testing.T) {
	client, errSetupDB := setupDB()
	require.NoError(t, errSetupDB)

	userCollection, errCollection := setupCollection(client)
	require.NoError(t, errCollection)

	user := models.User{
		Token:               "123",
		TelegramId:          123,
		TelegramDisplayName: "Test",
	}

	_, errAddUser := userCollection.InsertOne(context.Background(), user)
	require.NoError(t, errAddUser)

	userRepo := NewUserRepository(userCollection)
	foundedUser, errGetUser := userRepo.GetUserByTelegramId(978)
	require.NoError(t, errGetUser)

	emptyUser := models.User{}
	assert.Equal(t, emptyUser.Token, foundedUser.Token)
	assert.Equal(t, emptyUser.TelegramId, foundedUser.TelegramId)
	assert.Equal(t, emptyUser.TelegramDisplayName, foundedUser.TelegramDisplayName)

	defer cleanCollection(t, userCollection)
}
