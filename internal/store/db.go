package store

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/tanvir0188/vcita-ai-agent/internal/crypto"
)

// DB wraps the GORM database and the PHI encryptor.
// All store methods live on this struct so encryption is always available.
type DB struct {
	gorm *gorm.DB
	enc  *crypto.Encryptor
}

func (db *DB) GetGorm() *gorm.DB {
	return db.gorm
}

// New opens the SQLite database via GORM, auto-migrates the schema,
// and returns a ready-to-use DB handle.
func New(dsn string, enc *crypto.Encryptor, log *zap.Logger) (*DB, error) {
	// GORM logger: log slow queries only; never log values (may contain PHI)
	gormLog := logger.New(
		newGormZapWriter(log),
		logger.Config{
			SlowThreshold:             200,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true, // CRITICAL: prevents values from appearing in logs
		},
	)

	gormDB, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger:                 gormLog,
		SkipDefaultTransaction: true, // better performance for single-writer SQLite
	})
	if err != nil {
		return nil, fmt.Errorf("store: failed to open database: %w", err)
	}

	// Configure the underlying connection pool via the sql.DB handle
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("store: get sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)

	// AutoMigrate creates/updates tables to match the model structs.
	// Safe to run on every startup — it only adds columns and indexes, never drops.
	if err := gormDB.AutoMigrate(
		&Conversation{},
		&User{},
		&AuditLog{},
	); err != nil {
		return nil, fmt.Errorf("store: auto-migration failed: %w", err)
	}

	return &DB{gorm: gormDB, enc: enc}, nil
}

// Close closes the underlying database connection.
func (d *DB) Close() error {
	sqlDB, err := d.gorm.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// gormZapWriter adapts zap.Logger to GORM's logger.Writer interface.
type gormZapWriter struct{ log *zap.Logger }

func newGormZapWriter(log *zap.Logger) *gormZapWriter { return &gormZapWriter{log} }

func (w *gormZapWriter) Printf(format string, args ...interface{}) {
	w.log.Sugar().Debugf("[gorm] "+format, args...)
}
