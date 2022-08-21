package service

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/entity"
	"gitlab.ozon.dev/igor.benko.1991/homework/pkg/api"
)

func TestCreatePerson(t *testing.T) {
	t.Run("success create / all params are valid", func(t *testing.T) {
		// Arrange
		f := setUpPersonServiceFixture(t)

		fn, ln, id := "Igor", "Benko", uint64(1)
		apiPerson := api.CreatePersonRequest{
			FirstName: fn,
			LastName:  ln,
		}
		entityPerson := entity.Person{
			FirstName: fn,
			LastName:  ln,
		}

		f.personRepo.EXPECT().Create(gomock.Any(), entityPerson).Return(id, nil)

		// Act
		result, err := f.service.CreatePerson(f.Ctx, &apiPerson)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, id, result.Id)
	})
}

func TestUpdatePerson(t *testing.T) {
	t.Run("success update / all params are valid", func(t *testing.T) {
		// Arrange
		f := setUpPersonServiceFixture(t)

		fn, ln, id := "Igor", "Benko", uint64(1)
		apiPerson := api.UpdatePersonRequest{
			Id:        id,
			FirstName: fn,
			LastName:  ln,
		}
		entityPerson := entity.Person{
			ID:        id,
			FirstName: fn,
			LastName:  ln,
		}

		f.personRepo.EXPECT().Update(gomock.Any(), entityPerson).Return(nil)

		// Act
		_, err := f.service.UpdatePerson(f.Ctx, &apiPerson)

		// Assert
		require.NoError(t, err)
	})
}

func TestDeletePerson(t *testing.T) {
	t.Run("success delete / all params are valid", func(t *testing.T) {
		// Arrange
		f := setUpPersonServiceFixture(t)

		id := uint64(1)
		apiPerson := api.DeletePersonRequest{
			Id: id,
		}

		f.personRepo.EXPECT().Delete(gomock.Any(), id).Return(nil)

		// Act
		_, err := f.service.DeletePerson(f.Ctx, &apiPerson)

		// Assert
		require.NoError(t, err)
	})
}

func TestGetPerson(t *testing.T) {
	t.Run("success get / all params are valid", func(t *testing.T) {
		// Arrange
		f := setUpPersonServiceFixture(t)

		fn, ln, id := "Igor", "Benko", uint64(1)
		apiPerson := api.GetPersonRequest{
			Id: id,
		}
		entityPerson := &entity.Person{
			ID:        id,
			FirstName: fn,
			LastName:  ln,
		}
		entityVehicles := []entity.Vehicle{
			{ID: 2, Brand: "KIA", Model: "RIO", RegNumber: "P999PP39", PersonID: id},
		}

		f.personRepo.EXPECT().Get(gomock.Any(), id).Return(entityPerson, nil)
		f.vehicleRepo.EXPECT().GetByPersonID(gomock.Any(), id).Return(entityVehicles, nil)

		// Act
		result, err := f.service.GetPerson(f.Ctx, &apiPerson)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, id, result.Person.Id)
		assert.Equal(t, fn, result.Person.FirstName)
		assert.Equal(t, ln, result.Person.LastName)
		require.Equal(t, len(entityVehicles), len(result.Person.Vehicles))
		assert.Equal(t, entityVehicles[0].ID, result.Person.Vehicles[0].Id)
		assert.Equal(t, entityVehicles[0].Brand, result.Person.Vehicles[0].Brand)
		assert.Equal(t, entityVehicles[0].Model, result.Person.Vehicles[0].Model)
		assert.Equal(t, entityVehicles[0].RegNumber, result.Person.Vehicles[0].RegNumber)
		assert.Equal(t, entityVehicles[0].PersonID, result.Person.Vehicles[0].PersonId)
	})
}
