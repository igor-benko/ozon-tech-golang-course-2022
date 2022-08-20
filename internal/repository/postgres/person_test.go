package postgres

import (
	"context"
	"testing"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/entity"
)

func TestCreatePerson(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		f := setUpPersonRepoFixture(t)
		defer f.tearDownPersonRepoFixture()

		arg1, arg2 := "Igor", "Benko"
		id := uint64(1)
		query := "INSERT INTO persons (first_name,last_name) VALUES ($1,$2) RETURNING id"
		rows := pgxpoolmock.NewRows([]string{"id"}).AddRow(id).ToPgxRows()

		f.pool.EXPECT().Query(context.Background(), query, arg1, arg2).Return(rows, nil)

		// Act
		result, err := f.repo.Create(context.Background(), entity.Person{
			FirstName: arg1,
			LastName:  arg2,
		})

		// Assert
		require.NoError(t, err)
		assert.Equal(t, id, result)
	})
}

func TestUpdatePerson(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		f := setUpPersonRepoFixture(t)
		defer f.tearDownPersonRepoFixture()

		arg1, arg2, arg3 := "Igor", "Benko", uint64(1)
		query := "UPDATE persons SET first_name = $1, last_name = $2 WHERE id = $3"

		f.pool.EXPECT().Exec(context.Background(), query, arg1, arg2, arg3)

		// Act
		err := f.repo.Update(context.Background(), entity.Person{
			ID:        arg3,
			FirstName: arg1,
			LastName:  arg2,
		})

		// Assert
		require.NoError(t, err)
	})
}

func TestDeletePerson(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		f := setUpPersonRepoFixture(t)
		defer f.tearDownPersonRepoFixture()

		id := uint64(1)
		query := "DELETE FROM persons WHERE id = $1"

		f.pool.EXPECT().Exec(context.Background(), query, id).Return(nil, nil)

		// Act
		err := f.repo.Delete(context.Background(), 1)

		// Assert
		require.NoError(t, err)
	})
}

func TestGetPerson(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		f := setUpPersonRepoFixture(t)
		defer f.tearDownPersonRepoFixture()

		arg1, arg2 := "Igor", "Benko"
		id := uint64(1)
		query := "SELECT * FROM persons WHERE id = $1 LIMIT 1"
		rows := pgxpoolmock.NewRows([]string{"id", "last_name", "first_name"}).AddRow(id, arg2, arg1).ToPgxRows()

		f.pool.EXPECT().Query(context.Background(), query, id).Return(rows, nil)

		// Act
		result, err := f.repo.Get(context.Background(), 1)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, id, result.ID)
		assert.Equal(t, arg1, result.FirstName)
		assert.Equal(t, arg2, result.LastName)
	})
}

func TestListPerson(t *testing.T) {
	t.Run("success/limit < count", func(t *testing.T) {
		// Arrange
		f := setUpPersonRepoFixture(t)
		defer f.tearDownPersonRepoFixture()

		filter := entity.PersonFilter{
			Limit:  0,
			Offset: 0,
			Order:  "",
		}
		arg1, arg2, id, count := "Igor", "Benko", uint64(1), uint64(1)

		query1 := "SELECT * FROM persons ORDER BY id"
		rows1 := pgxpoolmock.NewRows([]string{"id", "last_name", "first_name"}).AddRow(id, arg2, arg1).ToPgxRows()

		f.pool.
			EXPECT().
			Query(context.Background(), query1).
			Return(rows1, nil)

		query2 := "SELECT count(*) FROM persons"
		rows2 := pgxpoolmock.NewRows([]string{"count"}).AddRow(count).ToPgxRows()

		f.pool.
			EXPECT().
			Query(context.Background(), query2).
			Return(rows2, nil)

		// Act
		result, err := f.repo.List(context.Background(), filter)

		// Assert
		require.NoError(t, err)
		require.Equal(t, count, result.Total)
		assert.Equal(t, id, result.Persons[0].ID)
		assert.Equal(t, arg1, result.Persons[0].FirstName)
		assert.Equal(t, arg2, result.Persons[0].LastName)
	})
	t.Run("success/limit > count", func(t *testing.T) {
		// Arrange
		f := setUpPersonRepoFixture(t)
		defer f.tearDownPersonRepoFixture()

		filter := entity.PersonFilter{
			Limit:  2,
			Offset: 0,
			Order:  "",
		}
		arg1, arg2, id, count := "Igor", "Benko", uint64(1), uint64(1)

		query1 := "SELECT * FROM persons ORDER BY id LIMIT 2"
		rows1 := pgxpoolmock.NewRows([]string{"id", "last_name", "first_name"}).AddRow(id, arg2, arg1).ToPgxRows()

		f.pool.
			EXPECT().
			Query(context.Background(), query1).
			Return(rows1, nil)

		query2 := "SELECT count(*) FROM persons"
		rows2 := pgxpoolmock.NewRows([]string{"count"}).AddRow(count).ToPgxRows()

		f.pool.
			EXPECT().
			Query(context.Background(), query2).
			Return(rows2, nil)

		// Act
		result, err := f.repo.List(context.Background(), filter)

		// Assert
		require.NoError(t, err)
		require.Equal(t, count, result.Total)
		assert.Equal(t, id, result.Persons[0].ID)
		assert.Equal(t, arg1, result.Persons[0].FirstName)
		assert.Equal(t, arg2, result.Persons[0].LastName)
	})
	t.Run("success/offset > count", func(t *testing.T) {
		// Arrange
		f := setUpPersonRepoFixture(t)
		defer f.tearDownPersonRepoFixture()

		filter := entity.PersonFilter{
			Limit:  0,
			Offset: 1,
			Order:  "",
		}
		count := uint64(0)

		query1 := "SELECT * FROM persons ORDER BY id OFFSET 1"
		rows1 := pgxpoolmock.NewRows([]string{}).AddRow().ToPgxRows()

		f.pool.
			EXPECT().
			Query(context.Background(), query1).
			Return(rows1, nil)

		query2 := "SELECT count(*) FROM persons"
		rows2 := pgxpoolmock.NewRows([]string{"count"}).AddRow(count).ToPgxRows()

		f.pool.
			EXPECT().
			Query(context.Background(), query2).
			Return(rows2, nil)

		// Act
		result, err := f.repo.List(context.Background(), filter)

		// Assert
		require.NoError(t, err)
		require.Equal(t, count, result.Total)
	})
	t.Run("order is empty", func(t *testing.T) {
		// Arrange
		f := setUpPersonRepoFixture(t)
		defer f.tearDownPersonRepoFixture()

		filter := entity.PersonFilter{
			Limit:  0,
			Offset: 0,
			Order:  "",
		}
		arg1, arg2, id, count := "Igor", "Benko", uint64(1), uint64(1)

		query1 := "SELECT * FROM persons ORDER BY id"
		rows1 := pgxpoolmock.NewRows([]string{"id", "last_name", "first_name"}).AddRow(id, arg2, arg1).ToPgxRows()

		f.pool.
			EXPECT().
			Query(context.Background(), query1).
			Return(rows1, nil)

		query2 := "SELECT count(*) FROM persons"
		rows2 := pgxpoolmock.NewRows([]string{"count"}).AddRow(count).ToPgxRows()

		f.pool.
			EXPECT().
			Query(context.Background(), query2).
			Return(rows2, nil)

		// Act
		result, err := f.repo.List(context.Background(), filter)

		// Assert
		require.NoError(t, err)
		require.Equal(t, count, result.Total)
		assert.Equal(t, id, result.Persons[0].ID)
		assert.Equal(t, arg1, result.Persons[0].FirstName)
		assert.Equal(t, arg2, result.Persons[0].LastName)
	})
}
