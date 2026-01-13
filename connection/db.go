package connection

import (
	"log"
	"os"

	"attendance-system/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// Load local .env if exists (for local testing)
	godotenv.Load()

	// Prefer DATABASE_URL (Neon or Render)
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Println("DATABASE_URL not found, falling back to local config")

		// Fallback to local DB (optional, for local testing only)
		localHost := os.Getenv("DB_HOST")
		localPort := os.Getenv("DB_PORT")
		localUser := os.Getenv("DB_USER")
		localPass := os.Getenv("DB_PASSWORD")
		localDB := os.Getenv("DB_NAME")

		dsn = "host=" + localHost +
			" user=" + localUser +
			" password=" + localPass +
			" dbname=" + localDB +
			" port=" + localPort +
			" sslmode=disable"
	}

	// Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	// Optional: Auto migrate tables (uncomment if starting fresh)
	// db.AutoMigrate(
	// 	&models.User{},
	// 	&models.PendingUser{},
	// 	&models.PasswordReset{},
	// 	&models.Event{},
	// 	&models.Attendance{},
	// )

	// Create audit_logs table if it doesn't exist
	if !db.Migrator().HasTable("audit_logs") {
		db.Exec(`
			CREATE TABLE audit_logs (
				id SERIAL PRIMARY KEY,
				action VARCHAR(255) NOT NULL,
				actor_id VARCHAR(255) NOT NULL,
				target_id VARCHAR(255),
				details TEXT,
				ip_address VARCHAR(255),
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			)
		`)
		db.Exec("CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action)")
		db.Exec("CREATE INDEX IF NOT EXISTS idx_audit_logs_actor_id ON audit_logs(actor_id)")
		db.Exec("CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at)")
		db.Exec("CREATE INDEX IF NOT EXISTS idx_audit_logs_target_id ON audit_logs(target_id)")
		log.Println("audit_logs table created successfully!")
	}

	// Ensure `tagged_courses` column exists in events table
	if !db.Migrator().HasColumn(&models.Event{}, "TaggedCoursesCSV") {
		if err := db.Migrator().AddColumn(&models.Event{}, "TaggedCoursesCSV"); err != nil {
			log.Printf("Failed to add column TaggedCoursesCSV: %v", err)
			if execErr := db.Exec("ALTER TABLE events ADD COLUMN IF NOT EXISTS tagged_courses text").Error; execErr != nil {
				log.Printf("Fallback ALTER TABLE failed: %v", execErr)
			}
		}
	}

	DB = db
	log.Println("Database connected successfully!")
}
