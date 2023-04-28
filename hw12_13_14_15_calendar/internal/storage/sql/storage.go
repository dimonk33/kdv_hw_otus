package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/dimonk33/kdv_hw_otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db           *sqlx.DB
	log          *logger.Logger
	connectParam string
}

func New(param string, logger *logger.Logger) *Storage {
	return &Storage{connectParam: param, log: logger}
}

func (s *Storage) Connect(ctx context.Context) (err error) {
	s.db, err = sqlx.ConnectContext(ctx, "postgres", s.connectParam)
	return
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) Create(ctx context.Context, data *storage.Event) (int64, error) {
	query := `
		INSERT INTO 
		    events
		    (
				id,
				title,
				start_time,
				end_time,
				description,
				own_user_id
			)
		OVERRIDING USER VALUE
		VALUES
		    ($1, $2, $3, $4, $5, $6)`
	res, err := s.db.ExecContext(
		ctx,
		query,
		data.ID,
		data.Title,
		data.StartTime.Unix(),
		data.EndTime.Unix(),
		data.Description,
		data.OwnUserID,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Storage) Update(ctx context.Context, data *storage.Event) error {
	query := `
		UPDATE 
		    events 
		SET 
		    title = $2,
		    start_time = $3,
		    end_time = $4,
		    description = $5,
		    own_user_id = $6 
		WHERE
		    id = $1`
	_, err := s.db.ExecContext(
		ctx,
		query,
		data.ID,
		data.Title,
		data.StartTime.Unix(),
		data.EndTime.Unix(),
		data.Description,
		data.OwnUserID,
	)
	return err
}

func (s *Storage) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM 
			events
		WHERE
		    id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	return err
}

func (s *Storage) ListOnDate(ctx context.Context, year int, month int, day int) ([]storage.Event, error) {
	var out []storage.Event
	query := `
		SELECT
		    * 
		FROM 
		    events 
		WHERE
		    date(start_time) = $1 or date(end_time) = $1`
	rows, err := s.db.QueryContext(ctx, query, fmt.Sprintf("%d-%d-%d", year, month, day))
	if err != nil {
		return out, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			s.log.Warning(err.Error())
		}
	}()
	return s.getSliceResult(ctx, rows)
}

func (s *Storage) ListOnWeek(ctx context.Context, year int, week int) ([]storage.Event, error) {
	var out []storage.Event
	query := `
		SELECT
		    * 
		FROM 
		    events 
		WHERE 
		    (date_part('year', start_time) = $1 and date_part('week', start_time) = $2)
		   or (date_part('year', end_time) = $1 and date_part('week', end_time) = $2)`
	rows, err := s.db.QueryContext(ctx, query, year, week)
	if err != nil {
		return out, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			s.log.Warning(err.Error())
		}
	}()
	return s.getSliceResult(ctx, rows)
}

func (s *Storage) ListOnMonth(ctx context.Context, year int, month int) ([]storage.Event, error) {
	var out []storage.Event
	query := `
		SELECT
		    * 
		FROM 
		    events 
		WHERE 
		    (date_part('year', start_time) = $1 and date_part('month', start_time) = $2)
		   or (date_part('year', end_time) = $1 and date_part('month', end_time) = $2)`
	rows, err := s.db.QueryContext(ctx, query, year, month)
	if err != nil {
		return out, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			s.log.Warning(err.Error())
		}
	}()
	return s.getSliceResult(ctx, rows)
}

func (s *Storage) ListLessDate(ctx context.Context, year, month, day int) ([]storage.Event, error) {
	var out []storage.Event
	query := `
		SELECT
		    * 
		FROM 
		    events 
		WHERE
		    date(end_time) < $1`
	rows, err := s.db.QueryContext(ctx, query, fmt.Sprintf("%d-%d-%d", year, month, day))
	if err != nil {
		return out, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			s.log.Warning(err.Error())
		}
	}()
	return s.getSliceResult(ctx, rows)
}

func (s *Storage) getSliceResult(ctx context.Context, rows *sql.Rows) ([]storage.Event, error) {
	var out []storage.Event
	for rows.Next() {
		select {
		case <-ctx.Done():
			return out, ctx.Err()
		default:
			var row storage.Event
			err := rows.Scan(&row)
			if err != nil {
				return out, err
			}
			out = append(out, row)
		}
	}
	return out, nil
}
