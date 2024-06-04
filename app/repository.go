package store

// import (
// 	"bytes"
// 	"context"
// 	"database/sql"
// 	"encoding/base64"
// 	"encoding/gob"
// 	"errors"
// 	"fmt"
// 	"reflect"

// 	"github.com/drone/runner-go/logger"
// 	"github.com/jmoiron/sqlx"
// 	// "github.com/grafana/incident/api/backoff"
// 	// "github.com/grafana/incident/api/instrumentation"
// 	// "github.com/grafana/incident/api/logger"
// )

// // DataObject describes an object that has an organization ID.
// type DataObject interface {
// 	// GetOrgID gets the organization ID.

// }

// // Repo is an interface that describes a data repository.
// type Repo[T DataObject] interface {
// 	// FindByID returns a single item.
// 	FindOne(ctx context.Context, query string, args map[string]interface{}) (*T, error)
// 	// Insert adds a row.
// 	Insert(ctx context.Context, query string, item T) (sql.Result, error)
// 	// Update modifies a row.
// 	Update(ctx context.Context, query string, args ...any) error
// 	// UpdateWithOptimisticLocking modifies a row using optimistic locking.
// 	// UpdateWithOptimisticLocking(ctx context.Context, query string, getter UpdateWithVersionGetter[T], updater UpdateWithVersionUpdater[T]) (*T, error)
// 	// Delete removes a row.
// 	Delete(ctx context.Context, query string, args ...any) error
// 	// Query executes a raw SQL query, with the specified arguments.
// 	Query(ctx context.Context, query string, args map[string]interface{}) ([]T, error)
// 	// QuerySingle is a generic method that executes a SQL query with named parameters and scans the result into a single struct.
// 	QuerySingle(ctx context.Context, query string, args map[string]interface{}, dest interface{}) error
// 	// QueryWithCursor executes a raw SQL query, with the specified arguments and cursor.
// 	// QueryWithCursor(ctx context.Context, querySQL string, queryParams Query, cursor Cursor) ([]T, Cursor, error)
// 	// Exec executes a raw SQL query, with the specified arguments.
// 	// See sqlx.ExecContext for more information.
// 	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
// 	// Count returns the number of rows matching the query.
// 	Count(ctx context.Context, query string, args ...any) (int, error)
// 	// DecodeCursorValue decodes the Cursor string value.
// 	DecodeCursorValue(cursorValue string, dest any) (any, error)
// 	// EncodeCursorValue encodes the Cursor value into a string.
// 	EncodeCursorValue(field any) (string, error)
// 	// Unsafe returns a repository that does not check orgIDs on read.
// 	Unsafe() Repo[T]
// }

// type repository[T DataObject] struct {
// 	db        *sql.DB
// 	logger    *logger.Logger
// 	modelName string
// }

// var (
// 	// ErrNotFound indicates the record couldn't be found.
// 	ErrNotFound = errors.New("not found")
// 	// ErrVersionConflict when the row version to update is more recent that the last read
// 	ErrVersionConflict = errors.New("version conflict (did you do `obj.OldVersion = obj.Version` and `obj.Version++` in the updater function?)")
// 	// ErrSkipUpdate informs the UpdateWithOptimisticLocking function to skip the update
// 	ErrSkipUpdate = errors.New("skip update")
// 	// ErrNoRowsAffected indicates that no rows were affected by the sql statement
// 	ErrNoRowsAffected = errors.New("no rows affected")
// 	// ErrCount indicates that the count query failed
// 	ErrCount = errors.New("count error")
// 	// ErrAccessDenied indicates the current user has no access to this record
// 	ErrAccessDenied = errors.New("access denied")
// )

// // Unsafe returns a repository that does not check orgIDs on read.
// func (r repository[T]) Unsafe() Repo[T] {
// 	return &repository[T]{
// 		db:        r.db,
// 		logger:    r.logger,
// 		modelName: r.modelName,
// 	}
// }

// func NewRepo[T DataObject](logger *logger.Logger, db *sqlx.DB) *repository[T] {
// 	var o T
// 	return &repository[T]{
// 		db:        db,
// 		logger:    logger,
// 		modelName: reflect.TypeOf(o).Name(),
// 	}
// }

// func (r repository[T]) FindOne(ctx context.Context, query string, args map[string]interface{}) (*T, error) {
// 	method := fmt.Sprintf("%s.FindOne", r.modelName)
// 	// ctx, span, startedAt, hadError := instrumentation.StartSpanWithMetrics(ctx, instrumentation.InstrumentationServiceTypeDB, method)
// 	// defer instrumentation.FinishSpanWithMetrics(ctx, instrumentation.InstrumentationServiceTypeDB, method, span, startedAt, hadError)
// 	// span.SetTag("query", query)

// 	var ok bool

// 	var res T
// 	rows, err := r.db.NamedQueryContext(ctx, query, args)
// 	if err != nil {

// 		return nil, fmt.Errorf("%s: %w", method, err)
// 	}
// 	defer rows.Close()
// 	ok = rows.Next()
// 	if !ok {
// 		return nil, fmt.Errorf("%s: %w", method, ErrNotFound)
// 	}
// 	err = rows.StructScan(&res)
// 	if err != nil {

// 		return nil, fmt.Errorf("%s: %w", method, err)
// 	}

// 	return &res, nil
// }

// func (r repository[T]) Insert(ctx context.Context, query string, item T) (sql.Result, error) {
// 	method := fmt.Sprintf("%s.Insert", r.modelName)
// 	// ctx, span, startedAt, hadError := instrumentation.StartSpanWithMetrics(ctx, instrumentation.InstrumentationServiceTypeDB, method)
// 	// defer instrumentation.FinishSpanWithMetrics(ctx, instrumentation.InstrumentationServiceTypeDB, method, span, startedAt, hadError)
// 	// span.SetTag("query", query)

// 	res, err := r.db.NamedExecContext(ctx, query, item)
// 	if err != nil {
// 		// hadError.Err = err
// 		return nil, fmt.Errorf("%s: %w", method, err)
// 	}

// 	return res, nil
// }

// // Exect the update statement inside the transaction
// // In case of failure the transaction is rollbacked
// func (r repository[T]) Update(ctx context.Context, query string, args ...any) error {
// 	method := fmt.Sprintf("%s.Update", r.modelName)
// 	// ctx, span, startedAt, hadError := instrumentation.StartSpanWithMetrics(ctx, instrumentation.InstrumentationServiceTypeDB, method)
// 	// defer instrumentation.FinishSpanWithMetrics(ctx, instrumentation.InstrumentationServiceTypeDB, method, span, startedAt, hadError)
// 	// span.SetTag("query", query)

// 	tx, err := r.db.Beginx()
// 	if err != nil {
// 		// hadError.Err = err
// 		return fmt.Errorf("%s: %w", method, err)
// 	}
// 	// Defer a rollback in case anything fails.
// 	defer func() {
// 		err := tx.Rollback()
// 		if err != nil {
// 			r.logger.Warnf(ctx, "rollback failed: %v", err)
// 		}
// 	}()

// 	_, err = tx.NamedExecContext(ctx, query, args)
// 	if err != nil {
// 		// hadError.Err = err
// 		return fmt.Errorf("%s: %w", method, err)
// 	}

// 	// Commit the transaction.
// 	if err = tx.Commit(); err != nil {
// 		// hadError.Err = err
// 		return fmt.Errorf("%s: %w", method, err)
// 	}

// 	return nil
// }

// // Exect the update statement without transaction and we use optmistic locking strategy
// // func (r repository[T]) UpdateWithOptimisticLocking(ctx context.Context, query string, getter UpdateWithVersionGetter[T], updater UpdateWithVersionUpdater[T]) (*T, error) {
// // 	// method := fmt.Sprintf("%s.UpdateWithOptimisticLocking", r.modelName)
// // 	// ctx, span, startedAt, hadError := instrumentation.StartSpanWithMetrics(ctx, instrumentation.InstrumentationServiceTypeDB, method)
// // 	// defer instrumentation.FinishSpanWithMetrics(ctx, instrumentation.InstrumentationServiceTypeDB, method, span, startedAt, hadError)
// // 	// span.SetTag("query", query)

// // 	// attempts := 0
// // 	// boff := backoff.NewDoubleTimeBackoff(r.logger, 50*time.Millisecond, 2*time.Second, 5)
// // 	var updatedRow *T
// // 	// err := boff.Do(ctx, func() error {
// // 	// 	attempts++
// // 	// 	latestRow, err := getter(ctx)
// // 	// 	if err != nil {
// // 	// 		return errors.Join(fmt.Errorf("%s: %w", method, err), backoff.ErrAbort)
// // 	// 	}

// // 	// 	updatedRow, err = updater(ctx, *latestRow)
// // 	// 	if err != nil {
// // 	// 		return errors.Join(fmt.Errorf("%s: %w", method, err), backoff.ErrAbort)
// // 	// 	}

// // 	// 	res, err := r.db.NamedExecContext(ctx, query, updatedRow)
// // 	// 	if err != nil {
// // 	// 		return errors.Join(fmt.Errorf("%s: %w", method, err), backoff.ErrAbort)
// // 	// 	}
// // 	// 	rowsAffectd, err := res.RowsAffected()
// // 	// 	if err != nil {
// // 	// 		return errors.Join(fmt.Errorf("%s: %w", method, err), backoff.ErrAbort)
// // 	// 	}
// // 	// 	if rowsAffectd == 0 {
// // 	// 		// we only retry if the row was updated by someone else, the other errors are not retriable
// // 	// 		return ErrVersionConflict
// // 	// 	}
// // 	// 	return nil
// // 	// })

// // 	// span.LogFields(log.Int("attempts", attempts))
// // 	// if err != nil {
// // 	// 	if errors.Is(err, ErrVersionConflict) {
// // 	// 		err = fmt.Errorf("%s: updated failed after %d attempts", method, attempts)
// // 	// 	}
// // 	// 	if errors.Is(err, ErrSkipUpdate) {
// // 	// 		// not an error, just skip the update
// // 	// 		return updatedRow, nil
// // 	// 	}
// // 	// 	hadError.Err = err
// // 	// 	return nil, err
// // 	// }
// // 	return updatedRow, nil
// // // }

// func (r repository[T]) Delete(ctx context.Context, query string, args ...any) error {
// 	method := fmt.Sprintf("%s.Delete", r.modelName)
// 	// ctx, span, startedAt, hadError := instrumentation.StartSpanWithMetrics(ctx, instrumentation.InstrumentationServiceTypeDB, method)
// 	// defer instrumentation.FinishSpanWithMetrics(ctx, instrumentation.InstrumentationServiceTypeDB, method, span, startedAt, hadError)
// 	// span.SetTag("query", query)
// r.db.
// 	_, err := r.db.NamedExecContext(ctx, query, args)
// 	if err != nil {
// 		// hadError.Err = err
// 		return fmt.Errorf("%s: %w", method, err)
// 	}

// 	return nil
// }

// func (r repository[T]) Query(ctx context.Context, query string, args map[string]interface{}) ([]T, error) {
// 	method := fmt.Sprintf("%s.Query", r.modelName)
// 	// ctx, span, startedAt, hadError := instrumentation.StartSpanWithMetrics(ctx, instrumentation.InstrumentationServiceTypeDB, method)
// 	// defer instrumentation.FinishSpanWithMetrics(ctx, instrumentation.InstrumentationServiceTypeDB, method, span, startedAt, hadError)
// 	// span.SetTag("query", query)

// 	q, parsedArgs, err := sqlx.Named(query, interface{}(args))
// 	if err != nil {
// 		// hadError.Err = err
// 		return nil, fmt.Errorf("%s: %w", method, err)
// 	}
// 	q, parsedArgs, err = sqlx.In(q, parsedArgs...)
// 	if err != nil {
// 		// hadError.Err = err
// 		return nil, fmt.Errorf("%s: %w", method, err)
// 	}

// 	q = r.db.Rebind(q)
// 	rows, err := r.db.QueryxContext(ctx, q, parsedArgs...)
// 	if err != nil {
// 		// hadError.Err = err
// 		return nil, fmt.Errorf("%s: %w", method, err)
// 	}
// 	defer rows.Close()

// 	var items []T = make([]T, 0)
// 	for rows.Next() {
// 		var item T
// 		err := rows.StructScan(&item)
// 		if err != nil {
// 			// hadError.Err = err
// 			return nil, fmt.Errorf("%s: %w", method, err)
// 		}

// 		items = append(items, item)
// 	}

// 	return items, nil
// }

// // QuerySingle is a generic method that executes a SQL query with named parameters and scans the result into a single struct.
// // dest must be a pointer to the struct you want to fill with the query result. Useful for e.g., collecting multiple metrics from a single query.
// func (r repository[T]) QuerySingle(ctx context.Context, query string, args map[string]interface{}, dest interface{}) error {
// 	method := fmt.Sprintf("%s.QuerySingle", r.modelName)
// 	// ctx, span, startedAt, hadError := instrumentation.StartSpanWithMetrics(ctx, instrumentation.InstrumentationServiceTypeDB, method)
// 	// defer instrumentation.FinishSpanWithMetrics(ctx, instrumentation.InstrumentationServiceTypeDB, method, span, startedAt, hadError)
// 	// span.SetTag("query", query)

// 	q, parsedArgs, err := sqlx.Named(query, args)
// 	if err != nil {
// 		// hadError.Err = err
// 		return fmt.Errorf("%s: %w", method, err)
// 	}

// 	q, parsedArgs, err = sqlx.In(q, parsedArgs...)
// 	if err != nil {
// 		// hadError.Err = err
// 		return fmt.Errorf("%s: %w", method, err)
// 	}

// 	q = r.db.Rebind(q)

// 	err = r.db.GetContext(ctx, dest, q, parsedArgs...)
// 	if err != nil {
// 		// hadError.Err = err
// 		return fmt.Errorf("%s: %w", method, err)
// 	}

// 	return nil
// }

// func (r repository[T]) Count(ctx context.Context, query string, args ...any) (int, error) {
// 	method := fmt.Sprintf("%s.Count", r.modelName)
// 	// ctx, span, startedAt, hadError := instrumentation.StartSpanWithMetrics(ctx, instrumentation.InstrumentationServiceTypeDB, method)
// 	// defer instrumentation.FinishSpanWithMetrics(ctx, instrumentation.InstrumentationServiceTypeDB, method, span, startedAt, hadError)

// 	var count int
// 	res, err := r.db.NamedQueryContext(ctx, query, args)
// 	if err != nil {
// 		// hadError.Err = err
// 		return 0, fmt.Errorf("%s: %w", method, err)
// 	}
// 	defer res.Close()

// 	ok := res.Next()
// 	if !ok {
// 		return 0, fmt.Errorf("%s: %w", method, ErrCount)
// 	}

// 	err = res.Scan(&count)
// 	if err != nil {
// 		// hadError.Err = err
// 		return 0, fmt.Errorf("%s: %w", method, err)
// 	}

// 	return count, nil
// }

// // Exec a custom SQL that is not necessarily linked with the default model for the repo
// // and the execution is wrapped in a transaction
// // In case of failure the transaction is rollbacked
// func (r repository[T]) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
// 	method := fmt.Sprintf("%s.Exec", r.modelName)
// 	// ctx, span, startedAt, hadError := instrumentation.StartSpanWithMetrics(ctx, instrumentation.InstrumentationServiceTypeDB, method)
// 	// defer instrumentation.FinishSpanWithMetrics(ctx, instrumentation.InstrumentationServiceTypeDB, method, span, startedAt, hadError)
// 	// span.SetTag("query", query)

// 	tx, err := r.db.Beginx()
// 	if err != nil {
// 		// hadError.Err = err
// 		return nil, fmt.Errorf("%s: %w", method, err)
// 	}
// 	// Defer a rollback in case anything fails.
// 	defer func() {
// 		err := tx.Rollback()
// 		if err != nil {
// 			r.logger.Warnf(ctx, "rollback failed: %v", err)
// 		}
// 	}()

// 	res, err := tx.NamedExecContext(ctx, query, args)
// 	if err != nil {
// 		// hadError.Err = err
// 		return res, fmt.Errorf("%s: %w", method, err)
// 	}

// 	// Commit the transaction.
// 	if err = tx.Commit(); err != nil {
// 		// hadError.Err = err
// 		return res, fmt.Errorf("%s: %w", method, err)
// 	}

// 	return res, nil
// }

// // Encode the entity field value to cursor value
// func (r repository[T]) EncodeCursorValue(field any) (string, error) {
// 	var encoded bytes.Buffer
// 	if err := gob.NewEncoder(&encoded).Encode(field); err != nil {
// 		return "", fmt.Errorf("gob encode: %w", err)
// 	}
// 	base64Encoded := base64.StdEncoding.EncodeToString(encoded.Bytes())
// 	return base64Encoded, nil
// }

// // Decode the cursor value to entity field value
// func (r repository[T]) DecodeCursorValue(cursorValue string, dest any) (any, error) {
// 	if cursorValue == "" {
// 		return nil, nil
// 	}
// 	decoded, err := base64.StdEncoding.DecodeString(cursorValue)
// 	if err != nil {
// 		return nil, fmt.Errorf("base64 decode: %w", err)
// 	}

// 	var val = reflect.ValueOf(dest)
// 	if err := gob.NewDecoder(bytes.NewReader(decoded)).DecodeValue(val); err != nil {
// 		return nil, fmt.Errorf("gob decode: %w", err)
// 	}
// 	return val.Elem().Interface(), nil
// }

// // func (r repository[T]) QueryWithCursor(ctx context.Context, querySQL string, queryParams Query, cursor Cursor) ([]T, Cursor, error) {
// // 	var err error
// // 	var args map[string]interface{} = make(map[string]interface{})
// // 	var whereArgs map[string]interface{}
// // 	if queryParams.Limit > 0 {
// // 		queryParams.Limit++
// // 	}
// // 	if cursor.Value != "" {
// // 		var cursorValue int
// // 		cursorFieldValue, err := r.DecodeCursorValue(cursor.Value, &cursorValue)
// // 		if err != nil {
// // 			return nil, cursor, err
// // 		}
// // 		offset, ok := cursorFieldValue.(int)
// // 		if !ok {
// // 			return nil, cursor, errors.New("invalid cursor value")
// // 		}
// // 		queryParams.Offset = offset
// // 	}

// // 	querySQL, whereArgs, err = prepareQuery(querySQL, &queryParams)
// // 	if err != nil {
// // 		return nil, cursor, err
// // 	}
// // 	for k, v := range whereArgs {
// // 		args[k] = v
// // 	}

// // 	items, err := r.Query(ctx, querySQL, args)
// // 	if err != nil {
// // 		return nil, cursor, err
// // 	}
// // 	if len(items) >= queryParams.Limit && queryParams.Limit > 0 {
// // 		cursor.HasMore = true
// // 		items = items[:queryParams.Limit-1]
// // 		cursor.Value, err = r.EncodeCursorValue(len(items) + queryParams.Offset)
// // 		if err != nil {
// // 			return nil, cursor, err
// // 		}
// // 	} else {
// // 		cursor.HasMore = false
// // 		cursor.Value = ""
// // 	}

// // 	return items, cursor, nil
// // }
