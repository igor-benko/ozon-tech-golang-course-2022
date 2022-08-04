package postgres

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/entity"
)

type personRepo struct {
	pool *pgxpool.Pool
}

func NewPersonRepo(pool *pgxpool.Pool) *personRepo {
	return &personRepo{
		pool: pool,
	}
}

func (r *personRepo) Create(ctx context.Context, item entity.Person) (uint64, error) {
	query, args, err := sq.Insert("persons").
		SetMap(sq.Eq{
			"first_name": item.FirstName,
			"last_name":  item.LastName,
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

func (r *personRepo) Update(ctx context.Context, item entity.Person) error {
	query, args, err := sq.Update("persons").
		SetMap(sq.Eq{
			"first_name": item.FirstName,
			"last_name":  item.LastName,
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

func (r *personRepo) Delete(ctx context.Context, id uint64) error {
	query, args, err := sq.Delete("persons").Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return err
	}

	if _, err := r.pool.Exec(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *personRepo) Get(ctx context.Context, id uint64) (*entity.Person, error) {
	query, args, err := sq.Select("*").From("persons").Where(sq.Eq{"id": id}).Limit(1).
		PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return nil, err
	}

	person := &entity.Person{}
	if err := pgxscan.Get(ctx, r.pool, person, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrPersonNotFound
		}

		return nil, err
	}

	return person, nil
}

func (r *personRepo) List(ctx context.Context, filter entity.PersonFilter) (*entity.PersonPage, error) {
	baseQuery := sq.Select("*").From("persons")
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

	persons := []entity.Person{}
	if err := pgxscan.Select(ctx, r.pool, &persons, query, args...); err != nil {
		return nil, err
	}

	countQuery, args, err := sq.Select("count(*)").From("persons").PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var total uint64 = 0
	if err := pgxscan.Get(ctx, r.pool, &total, countQuery, args...); err != nil {
		return nil, err
	}

	return &entity.PersonPage{
		Persons: persons,
		Total:   total,
	}, nil
}
