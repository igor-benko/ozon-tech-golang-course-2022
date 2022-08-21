package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/igor.benko.1991/homework/pkg/api"
)

func TestService(t *testing.T) {
	t.Parallel()

	t.Run("Create", func(t *testing.T) {
		t.Run("success / all params are valid", func(t *testing.T) {
			// Arrange
			setUpPool()
			defer tearDownPool()

			person := api.CreatePersonRequest{
				FirstName: "Igor",
				LastName:  "Benko",
			}

			// Act
			result, err := personService.CreatePerson(context.Background(), &person)

			// Assert
			require.NoError(t, err)
			assert.NotEqual(t, uint64(0), result)
		})
	})
	t.Run("Update", func(t *testing.T) {
		t.Run("success / all params are valid", func(t *testing.T) {
			// Arrange
			setUpPool()
			defer tearDownPool()

			personCreate := api.CreatePersonRequest{
				FirstName: "Igor",
				LastName:  "Benko",
			}
			personUpdate := api.UpdatePersonRequest{
				FirstName: "Igor",
				LastName:  "Benko",
			}

			// Act
			result1, err := personService.CreatePerson(context.Background(), &personCreate)

			require.NoError(t, err)
			assert.NotEqual(t, uint64(0), result1)

			personUpdate.Id = result1.Id

			_, err = personService.UpdatePerson(context.Background(), &personUpdate)

			// Assert
			require.NoError(t, err)
		})
		t.Run("success / no person", func(t *testing.T) {
			// Arrange
			setUpPool()
			defer tearDownPool()

			personUpdate := api.UpdatePersonRequest{
				FirstName: "Igor",
				LastName:  "Benko",
			}

			// Act
			_, err := personService.UpdatePerson(context.Background(), &personUpdate)

			// Assert
			require.NoError(t, err)
		})
	})
	t.Run("Delete", func(t *testing.T) {
		t.Run("success / all params are valid", func(t *testing.T) {
			// Arrange
			setUpPool()
			defer tearDownPool()

			personCreate := api.CreatePersonRequest{
				FirstName: "Igor",
				LastName:  "Benko",
			}
			personDelete := api.DeletePersonRequest{}

			// Act
			result1, err := personService.CreatePerson(context.Background(), &personCreate)

			require.NoError(t, err)
			assert.NotEqual(t, uint64(0), result1)

			personDelete.Id = result1.Id

			_, err = personService.DeletePerson(context.Background(), &personDelete)

			// Assert
			require.NoError(t, err)
		})
		t.Run("success / no person", func(t *testing.T) {
			// Arrange
			setUpPool()
			defer tearDownPool()

			personDelete := api.DeletePersonRequest{
				Id: 0,
			}

			// Act
			_, err := personService.DeletePerson(context.Background(), &personDelete)

			// Assert
			require.NoError(t, err)
		})
	})
	t.Run("Get", func(t *testing.T) {
		t.Run("success / all params are valid", func(t *testing.T) {
			// Arrange
			setUpPool()
			defer tearDownPool()

			personCreate := api.CreatePersonRequest{
				FirstName: "Igor",
				LastName:  "Benko",
			}
			personGet := api.GetPersonRequest{}

			// Act
			result1, err := personService.CreatePerson(context.Background(), &personCreate)

			require.NoError(t, err)
			assert.NotEqual(t, uint64(0), result1)

			personGet.Id = result1.Id

			result2, err := personService.GetPerson(context.Background(), &personGet)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, personCreate.FirstName, result2.Person.FirstName)
			assert.Equal(t, personCreate.LastName, result2.Person.LastName)
			assert.Equal(t, result1.Id, result2.Person.Id)
		})
		t.Run("success / no person", func(t *testing.T) {
			// Arrange
			setUpPool()
			defer tearDownPool()

			person := api.GetPersonRequest{
				Id: 0,
			}

			// Act
			_, err := personService.GetPerson(context.Background(), &person)

			// Assert
			require.Error(t, err)
		})
	})
}
