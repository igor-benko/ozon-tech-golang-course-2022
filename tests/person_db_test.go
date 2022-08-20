package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/entity"
)

func TestDB(t *testing.T) {
	t.Run("Create", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			// Arrange
			setUpPool()
			defer tearDownPool()

			person := entity.Person{
				FirstName: "Igor",
				LastName:  "Benko",
			}

			// Act
			result, err := personRepo.Create(context.Background(), person)

			// Assert
			require.NoError(t, err)
			assert.NotEqual(t, uint64(0), result)
		})
	})
	t.Run("Update", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			// Arrange
			setUpPool()
			defer tearDownPool()

			person := entity.Person{
				FirstName: "Igor",
				LastName:  "Benko",
			}

			// Act
			result1, err := personRepo.Create(context.Background(), person)

			require.NoError(t, err)
			assert.NotEqual(t, uint64(0), result1)

			person.ID = result1

			err = personRepo.Update(context.Background(), person)

			// Assert
			require.NoError(t, err)
		})
		t.Run("success/no person", func(t *testing.T) {
			// Arrange
			setUpPool()
			defer tearDownPool()

			person := entity.Person{
				ID:        0,
				FirstName: "Igor",
				LastName:  "Benko",
			}

			// Act
			err := personRepo.Update(context.Background(), person)

			// Assert
			require.NoError(t, err)
		})
	})
	t.Run("Delete", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			// Arrange
			setUpPool()
			defer tearDownPool()

			person := entity.Person{
				FirstName: "Igor",
				LastName:  "Benko",
			}

			// Act
			result1, err := personRepo.Create(context.Background(), person)

			require.NoError(t, err)
			assert.NotEqual(t, uint64(0), result1)

			err = personRepo.Delete(context.Background(), result1)

			// Assert
			require.NoError(t, err)
		})
		t.Run("success/no person", func(t *testing.T) {
			// Arrange
			setUpPool()
			defer tearDownPool()

			person := entity.Person{
				ID:        0,
				FirstName: "Igor",
				LastName:  "Benko",
			}

			// Act
			err := personRepo.Delete(context.Background(), person.ID)

			// Assert
			require.NoError(t, err)
		})
	})
	t.Run("Get", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			// Arrange
			setUpPool()
			defer tearDownPool()

			person := entity.Person{
				FirstName: "Igor",
				LastName:  "Benko",
			}

			// Act
			result1, err := personRepo.Create(context.Background(), person)

			require.NoError(t, err)
			assert.NotEqual(t, uint64(0), result1)

			result2, err := personRepo.Get(context.Background(), result1)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, person.FirstName, result2.FirstName)
			assert.Equal(t, person.LastName, result2.LastName)
			assert.Equal(t, result1, result2.ID)
		})
		t.Run("success/no person", func(t *testing.T) {
			// Arrange
			setUpPool()
			defer tearDownPool()

			person := entity.Person{
				ID:        0,
				FirstName: "Igor",
				LastName:  "Benko",
			}

			// Act
			_, err := personRepo.Get(context.Background(), person.ID)

			// Assert
			require.Error(t, err)
		})
	})
	t.Run("List", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			// Arrange
			setUpPool()
			defer tearDownPool()

			person := entity.Person{
				FirstName: "Igor",
				LastName:  "Benko",
			}
			filter := entity.PersonFilter{
				Limit:  1,
				Offset: 0,
				Order:  "id",
			}

			// Act
			result1, err := personRepo.Create(context.Background(), person)

			require.NoError(t, err)
			assert.NotEqual(t, uint64(0), result1)

			result2, err := personRepo.List(context.Background(), filter)

			// Assert
			require.NoError(t, err)
			require.Equal(t, 1, len(result2.Persons))
			assert.Equal(t, person.FirstName, result2.Persons[0].FirstName)
			assert.Equal(t, person.LastName, result2.Persons[0].LastName)
		})
		t.Run("success/no person", func(t *testing.T) {
			// Arrange
			setUpPool()
			defer tearDownPool()

			filter := entity.PersonFilter{
				Limit:  1,
				Offset: 0,
				Order:  "id",
			}

			// Act

			result, err := personRepo.List(context.Background(), filter)

			// Assert
			require.NoError(t, err)
			require.Equal(t, 0, len(result.Persons))
		})
	})
}
