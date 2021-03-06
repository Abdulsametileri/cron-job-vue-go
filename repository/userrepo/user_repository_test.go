package userrepo

import (
	"context"
	"github.com/Abdulsametileri/cron-job-vue-go/models"
	"github.com/Abdulsametileri/cron-job-vue-go/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func cleanCollection(t *testing.T, userCollection *mongo.Collection) {
	errDrop := userCollection.Drop(context.Background())
	require.NoError(t, errDrop)
}

func TestUserRepo_AddUser(t *testing.T) {
	client, errSetupDB := repository.SetupDB()
	require.NoError(t, errSetupDB)

	userCollection, errCollection := repository.SetupCollection(client, "users")
	require.NoError(t, errCollection)

	defer cleanCollection(t, userCollection)

	user := models.User{
		Token:               "123",
		TelegramId:          123,
		TelegramDisplayName: "Test",
	}
	userRepo := NewUserRepository(userCollection)
	errAddUser := userRepo.AddUser(user)
	require.NoError(t, errAddUser)
}

func TestUserRepo_GetUserByTelegramId(t *testing.T) {
	client, errSetupDB := repository.SetupDB()
	require.NoError(t, errSetupDB)

	userCollection, errCollection := repository.SetupCollection(client, "users")
	require.NoError(t, errCollection)

	defer cleanCollection(t, userCollection)

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
}

func TestUserRepo_GetUserByTelegramId_NotExistUser(t *testing.T) {
	client, errSetupDB := repository.SetupDB()
	require.NoError(t, errSetupDB)

	userCollection, errCollection := repository.SetupCollection(client, "users")
	require.NoError(t, errCollection)

	defer cleanCollection(t, userCollection)

	userRepo := NewUserRepository(userCollection)
	foundedUser, errGetUser := userRepo.GetUserByTelegramId(978)
	require.NoError(t, errGetUser)

	emptyUser := models.User{}
	assert.Equal(t, emptyUser.Token, foundedUser.Token)
	assert.Equal(t, emptyUser.TelegramId, foundedUser.TelegramId)
	assert.Equal(t, emptyUser.TelegramDisplayName, foundedUser.TelegramDisplayName)
}

func TestUserRepo_GetUserByToken(t *testing.T) {
	client, errSetupDB := repository.SetupDB()
	require.NoError(t, errSetupDB)

	userCollection, errCollection := repository.SetupCollection(client, "users")
	require.NoError(t, errCollection)

	defer cleanCollection(t, userCollection)

	user := models.User{
		Token:               "123",
		TelegramId:          123,
		TelegramDisplayName: "Test",
	}

	_, errAddUser := userCollection.InsertOne(context.Background(), user)
	require.NoError(t, errAddUser)

	userRepo := NewUserRepository(userCollection)
	foundedUser, errGetUser := userRepo.GetUserByToken(user.Token)
	require.NoError(t, errGetUser)

	assert.Equal(t, user.Token, foundedUser.Token)
	assert.Equal(t, user.TelegramId, foundedUser.TelegramId)
	assert.Equal(t, user.TelegramDisplayName, foundedUser.TelegramDisplayName)
}

func TestUserRepo_GetUserByToken_NotExistUser(t *testing.T) {
	client, errSetupDB := repository.SetupDB()
	require.NoError(t, errSetupDB)

	userCollection, errCollection := repository.SetupCollection(client, "users")
	require.NoError(t, errCollection)

	defer cleanCollection(t, userCollection)

	userRepo := NewUserRepository(userCollection)
	foundedUser, errGetUser := userRepo.GetUserByToken("notexistedtoken")
	require.NoError(t, errGetUser)

	emptyUser := models.User{}
	assert.Equal(t, emptyUser.Token, foundedUser.Token)
	assert.Equal(t, emptyUser.TelegramId, foundedUser.TelegramId)
	assert.Equal(t, emptyUser.TelegramDisplayName, foundedUser.TelegramDisplayName)
}
