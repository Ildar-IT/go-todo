package userRepo_test

import (
	"fmt"
	"testing"
	"todo/internal/database/pg"
	"todo/internal/entity"
	userRepo "todo/internal/repository/user"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUserCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err.Error())
	}

	defer db.Close()

	r := userRepo.NewUserRepository(db)

	type mockBehavior func(user *entity.User)

	testTable := []struct {
		name         string
		user         *entity.User
		id           int
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "OK",
			user: &entity.User{
				Username:      "test",
				Email:         "test@mail.ru",
				Password_hash: "passsss_hasshhh",
			},
			id: 2,
			mockBehavior: func(user *entity.User) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(2)
				mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", pg.UsersTable)).WithArgs("test", "test@mail.ru", "passsss_hasshhh").WillReturnRows(rows)
			},
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.user)
			id, err := r.Create(test.user)

			if test.wantErr {
				assert.Error(t, err)
			}

			assert.NoError(t, err)
			assert.Equal(t, test.id, id)

		})
	}

}
