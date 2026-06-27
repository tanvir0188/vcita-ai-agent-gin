package store

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/tanvir0188/vcita-ai-agent/internal/ai"
	"github.com/tanvir0188/vcita-ai-agent/internal/crypto"
)

// Store is the interface that appointment (and other) packages use.
// Using an interface keeps those packages decoupled from the concrete *DB
// and makes unit-testing with mocks straightforward.
type Store interface {
	UpsertAppointment(appt *Appointment) error
	SetAutoReplyOff(conversationID string) error
	SetAutoReplyOn(conversationID string) error
	CreateOrUpdateReminder(matterUID string, resp *ai.MedicationRefillResponse) error
}

// DB wraps the GORM database and the PHI encryptor.
// All store methods live on this struct so encryption is always available.
type DB struct {
	gorm *gorm.DB
	enc  *crypto.Encryptor
}

// Ensure *DB satisfies the Store interface at compile time.
var _ Store = (*DB)(nil)

// UpsertAppointment inserts or updates an appointment record keyed on
// (ConversationID, ServiceID). A record is created if none exists;
// otherwise Status, StartTime, EndTime and VcitaAppointmentID are updated.
func (db *DB) UpsertAppointment(appt *Appointment) error {
	return db.gorm.
		Where(Appointment{
			ConversationID: appt.ConversationID,
			ServiceID:      appt.ServiceID,
		}).
		Assign(Appointment{
			ContactID:          appt.ContactID,
			ServiceName:        appt.ServiceName,
			Status:             appt.Status,
			StartTime:          appt.StartTime,
			EndTime:            appt.EndTime,
			VcitaAppointmentID: appt.VcitaAppointmentID,
			VcitaClientID:      appt.VcitaClientID,
		}).
		FirstOrCreate(appt).
		Error
}

// SetAutoReplyOff disables the AI auto-reply for a conversation while
// appointment scheduling is in progress.
func (db *DB) SetAutoReplyOff(conversationID string) error {
	return db.gorm.
		Model(&Conversation{}).
		Where("conversation_id = ?", conversationID).
		Update("auto_reply_off", true).
		Error
}

// SetAutoReplyOn re-enables the AI auto-reply after scheduling completes.
func (db *DB) SetAutoReplyOn(conversationID string) error {
	return db.gorm.
		Model(&Conversation{}).
		Where("conversation_id = ?", conversationID).
		Update("auto_reply_off", false).
		Error
}

func (db *DB) CreateOrUpdateReminder(matterUID string, resp *ai.MedicationRefillResponse) error {
	// Find existing sync state
	var syncState ClientSyncState
	err := db.gorm.Where("matter_uid = ?", matterUID).First(&syncState).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	now := time.Now()
	// Prepare fields based on response
	remind := resp.ReminderNeeded
	partial := resp.PartialReminder
	hasMeds := len(resp.Medications) > 0
	if err == gorm.ErrRecordNotFound {
		// Create new record
		syncState = ClientSyncState{
			MatterUID:             matterUID,
			RemindMedication:      remind,
			HasMedications:        hasMeds,
			PartialReminderNeeded: partial,
			LastAICheckAt:         &now,
		}
		return db.gorm.Create(&syncState).Error
	}
	// Update existing record
	syncState.RemindMedication = remind
	syncState.HasMedications = hasMeds
	syncState.PartialReminderNeeded = partial
	syncState.LastAICheckAt = &now
	return db.gorm.Save(&syncState).Error
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
		&Appointment{},
		&ClientSyncState{},
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
