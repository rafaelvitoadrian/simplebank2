package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/rafaelvitoadrian/simplebank2/utils"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	password := utils.RandomString(6)
	hashedPassword, err := utils.HashedPassword(password)
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       utils.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       utils.RandomOwner(),
		Email:          utils.RandomEmail(),
	}

	user, err := TestQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := TestQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := CreateRandomUser(t)

	newFullName := utils.RandomOwner()
	updateUser, err := TestQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.FullName, updateUser.FullName)
	require.Equal(t, newFullName, updateUser.FullName)
	require.Equal(t, oldUser.Username, updateUser.Username)
	require.Equal(t, oldUser.HashedPassword, updateUser.HashedPassword)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := CreateRandomUser(t)

	newEmail := utils.RandomEmail()
	updateUser, err := TestQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.Email, updateUser.Email)
	require.Equal(t, newEmail, updateUser.Email)
	require.Equal(t, oldUser.Username, updateUser.Username)
	require.Equal(t, oldUser.HashedPassword, updateUser.HashedPassword)
}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser := CreateRandomUser(t)

	newEmail := utils.RandomEmail()
	newFullName := utils.RandomOwner()
	newPassword := utils.RandomString(6)
	newHashedPasswrod, err := utils.HashedPassword(newPassword)
	require.NoError(t, err)

	updateUser, err := TestQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: newHashedPasswrod,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.HashedPassword, updateUser.HashedPassword)
	require.Equal(t, newHashedPasswrod, updateUser.HashedPassword)
	require.NotEqual(t, oldUser.FullName, updateUser.FullName)
	require.Equal(t, newFullName, updateUser.FullName)
	require.NotEqual(t, oldUser.Email, updateUser.Email)
	require.Equal(t, newEmail, updateUser.Email)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser := CreateRandomUser(t)

	newPassword := utils.RandomString(6)
	newHashedPasswrod, err := utils.HashedPassword(newPassword)
	require.NoError(t, err)

	updateUser, err := TestQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: newHashedPasswrod,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.HashedPassword, updateUser.HashedPassword)
	require.Equal(t, newHashedPasswrod, updateUser.HashedPassword)
	require.Equal(t, oldUser.Username, updateUser.Username)
	require.Equal(t, oldUser.FullName, updateUser.FullName)
}
