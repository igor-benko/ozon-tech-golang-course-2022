package postgres

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/entity"
)

type vehicleRepo struct {
	pool PgxPool
}

func NewVehicleRepo(pool PgxPool) *vehicleRepo {
	return &vehicleRepo{
		pool: pool,
	}
}

func (r *vehicleRepo) Create(ctx context.Context, item entity.Vehicle) (uint64, error) {
	query, args, err := sq.Insert("vehicles").
		SetMap(sq.Eq{
			"brand":      item.Brand,
			"model":      item.Model,
			"reg_number": item.RegNumber,
			"person_id":  item.PersonID,
		}).Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return 0, err
	}

	var id uint64
	if err := r.pool.QueryRow(ctx, query, args...).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *vehicleRepo) Update(ctx context.Context, item entity.Vehicle) error {
	query, args, err := sq.Update("vehicles").
		SetMap(sq.Eq{
			"brand":      item.Brand,
			"model":      item.Model,
			"reg_number": item.RegNumber,
		}).
		Where(sq.Eq{"id": item.ID}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	if _, err := r.pool.Exec(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *vehicleRepo) Delete(ctx context.Context, id uint64) error {
	query, args, err := sq.Delete("vehicles").Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	if _, err := r.pool.Exec(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *vehicleRepo) Get(ctx context.Context, id uint64) (*entity.Vehicle, error) {
	query, args, err := sq.Select("*").From("vehicles").Where(sq.Eq{"id": id}).Limit(1).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	vehicle := &entity.Vehicle{}
	if err := pgxscan.Get(ctx, r.pool, vehicle, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrVehicleNotFound
		}

		return nil, err
	}

	return vehicle, nil
}

func (r *vehicleRepo) Exists(ctx context.Context, regNum string) (bool, error) {
	query, args, err := sq.Select("1").From("vehicles").Where(sq.Eq{"reg_number": regNum}).Limit(1).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return false, err
	}

	temp := 0
	if err := r.pool.QueryRow(ctx, query, args...).Scan(&temp); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *vehicleRepo) GetByPersonID(ctx context.Context, personID uint64) ([]entity.Vehicle, error) {
	query, args, err := sq.Select("*").From("vehicles").Where(sq.Eq{"person_id": personID}).Limit(1).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	vehicles := []entity.Vehicle{}
	if err := pgxscan.Select(ctx, r.pool, &vehicles, query, args...); err != nil {
		return nil, err
	}

	return vehicles, nil
}

func (r *vehicleRepo) List(ctx context.Context, filter entity.VehicleFilter) (*entity.VehiclePage, error) {
	baseQuery := sq.Select("*").From("vehicles")
	if len(filter.Order) != 0 {
		baseQuery = baseQuery.OrderBy(filter.Order)
	} else {
		baseQuery = baseQuery.OrderBy("id")
	}

	if filter.Offset > 0 {
		baseQuery = baseQuery.Offset(filter.Offset)
	}

	if filter.Limit > 0 {
		baseQuery = baseQuery.Limit(filter.Limit)
	}

	query, args, err := baseQuery.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	vehicles := []entity.Vehicle{}
	if err := pgxscan.Select(ctx, r.pool, &vehicles, query, args...); err != nil {
		return nil, err
	}

	countQuery, args, err := sq.Select("count(*)").From("vehicles").PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var total uint64 = 0
	if err := pgxscan.Get(ctx, r.pool, &total, countQuery, args...); err != nil {
		return nil, err
	}

	return &entity.VehiclePage{
		Vehicles: vehicles,
		Total:    total,
	}, nil
}
