package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUsers(t *testing.T) {
	testUserDB := NewUserDatabase(getTestInMemoryDBClient())

	users := Users{
		&User{Username: "user1", PasswordHash: "hash1"},
		&User{Username: "user2", PasswordHash: "hash2"},
	}
	assert.NoError(t, testUserDB.CreateUsers(users))

	user1, err := testUserDB.GetUserByUsername("user1")
	assert.NoError(t, err)
	assert.Equal(t, users[0].Username, user1.Username)
	assert.Equal(t, users[0].PasswordHash, user1.PasswordHash)

	user2, err := testUserDB.GetUserByUsername("user2")
	assert.NoError(t, err)
	assert.Equal(t, users[1].Username, user2.Username)
	assert.Equal(t, users[1].PasswordHash, user2.PasswordHash)
}

func TestGetUserByUsername(t *testing.T) {
	testUserDB := NewUserDatabase(getTestInMemoryDBClient())

	_, err := testUserDB.GetUserByUsername("non existing user")
	assert.Error(t, err)
}

func TestGetUsersByUsernames(t *testing.T) {
	testUserDB := NewUserDatabase(getTestInMemoryDBClient())

	users := Users{
		&User{Username: "user1", PasswordHash: "hash1"},
		&User{Username: "user2", PasswordHash: "hash2"},
	}
	assert.NoError(t, testUserDB.CreateUsers(users))

	usersFromDB, err := testUserDB.GetUsersByUsernames([]string{users[0].Username, users[1].Username})
	assert.NoError(t, err)

	assert.Len(t, usersFromDB, len(users))
	for i, userFromDB := range usersFromDB {
		assert.Equal(t, users[i].Username, userFromDB.Username)
		assert.Equal(t, users[i].PasswordHash, userFromDB.PasswordHash)
	}
}
