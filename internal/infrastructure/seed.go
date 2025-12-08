package infrastructure

import (
	"log"

	"github.com/zercle/gofiber-skelton/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedDefaultAdmin creates default admin if not exists
func SeedDefaultAdmin(db *gorm.DB) error {
	var count int64
	db.Model(&models.Staff{}).Where("email = ?", "admin@system.com").Count(&count)

	if count > 0 {
		log.Println("✅ Default admin already exists")
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte("Admin@123456"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Printf("❌ Failed to hash admin password: %v", err)
		return err
	}

	admin := models.Staff{
		Prefix:      "System",
		FirstName:   "Admin",
		LastName:    "Default",
		Email:       "admin@system.com",
		Password:    string(hashedPassword),
		Role:        models.RoleAdmin,
		Status:      models.StatusStaffActive,
		PhoneNumber: "0000000000",
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Printf("❌ Failed to create admin: %v", err)
		return err
	}

	log.Println("🎉 ========================================")
	log.Println("   ✅ Default Admin Created Successfully!")
	log.Println("   📧 Email:    admin@system.com")
	log.Println("   🔑 Password: Admin@123456")
	log.Println("   ⚠️  Please change password after first login")
	log.Println("   ========================================")

	return nil
}

// SeedTestData seeds test data for development
func SeedTestData(db *gorm.DB) error {
	log.Println("🌱 Starting test data seeding...")

	// ========================================
	// 1. Seed Test Staff (Driver & Collector)
	// ========================================

	var driver, collector models.Staff

	// ✅ Check และสร้าง Driver แยกกัน
	err := db.Where("email = ?", "driver@test.com").First(&driver).Error
	if err == gorm.ErrRecordNotFound {
		log.Println("   Creating test driver...")

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("driver123"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("❌ Failed to hash driver password: %v", err)
			return err
		}

		driver = models.Staff{
			Prefix:      "นาย",
			FirstName:   "สมชาย",
			LastName:    "คนขับรถ",
			Email:       "driver@test.com",
			Password:    string(hashedPassword),
			Role:        models.RoleDriver,
			Status:      models.StatusStaffActive,
			PhoneNumber: "0812345678",
		}

		if err := db.Create(&driver).Error; err != nil {
			log.Printf("❌ Failed to create driver: %v", err)
			return err
		}
		log.Printf("   ✅ Driver created (ID: %d)", driver.ID)
	} else if err != nil {
		// ถ้ามี error อื่นๆ ที่ไม่ใช่ record not found
		log.Printf("❌ Error loading driver: %v", err)
		return err
	} else {
		log.Printf("   ⏭️  Driver already exists (ID: %d)", driver.ID)
	}

	// ✅ Check และสร้าง Collector แยกกัน
	err = db.Where("email = ?", "collector@test.com").First(&collector).Error
	if err == gorm.ErrRecordNotFound {
		log.Println("   Creating test collector...")

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("collector123"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("❌ Failed to hash collector password: %v", err)
			return err
		}

		collector = models.Staff{
			Prefix:      "นาย",
			FirstName:   "สมหมาย",
			LastName:    "คนเก็บขยะ",
			Email:       "collector@test.com",
			Password:    string(hashedPassword),
			Role:        models.RoleCollector,
			Status:      models.StatusStaffActive,
			PhoneNumber: "0823456789",
		}

		if err := db.Create(&collector).Error; err != nil {
			log.Printf("❌ Failed to create collector: %v", err)
			return err
		}
		log.Printf("   ✅ Collector created (ID: %d)", collector.ID)
	} else if err != nil {
		log.Printf("❌ Error loading collector: %v", err)
		return err
	} else {
		log.Printf("   ⏭️  Collector already exists (ID: %d)", collector.ID)
	}

	// ========================================
	// 2. Seed Test Vehicles
	// ========================================

	// ✅ Check และสร้าง Vehicle 1 แยกกัน
	var vehicle1 models.Vehicle
	err = db.Where("registration_number = ?", "กข-1234").First(&vehicle1).Error
	if err == gorm.ErrRecordNotFound {
		log.Println("   Creating vehicle 1 (กข-1234)...")

		vehicle1 = models.Vehicle{
			RegistrationNumber:        "กข-1234",
			VehicleType:               "รถบรรทุก 6 ล้อ",
			Status:                    models.StatusActive,
			RegularWasteCapacityKg:    5000,
			RecyclableWasteCapacityKg: 2000,
			CurrentDriverID:           &driver.ID,
			FuelType:                  "Diesel",
			DepreciationValuePerYear:  50000,
		}

		if err := db.Create(&vehicle1).Error; err != nil {
			log.Printf("❌ Failed to create vehicle 1: %v", err)
			return err
		}
		log.Printf("   ✅ Vehicle 1 created (ID: %d, Driver: %d)", vehicle1.ID, *vehicle1.CurrentDriverID)
	} else if err != nil {
		log.Printf("❌ Error loading vehicle 1: %v", err)
		return err
	} else {
		log.Printf("   ⏭️  Vehicle 1 already exists (ID: %d)", vehicle1.ID)
	}

	// ✅ Check และสร้าง Vehicle 2 แยกกัน
	var vehicle2 models.Vehicle
	err = db.Where("registration_number = ?", "คง-5678").First(&vehicle2).Error
	if err == gorm.ErrRecordNotFound {
		log.Println("   Creating vehicle 2 (คง-5678)...")

		vehicle2 = models.Vehicle{
			RegistrationNumber:        "คง-5678",
			VehicleType:               "รถกระบะ",
			Status:                    models.StatusActive,
			RegularWasteCapacityKg:    1000,
			RecyclableWasteCapacityKg: 500,
			FuelType:                  "Gasoline",
			DepreciationValuePerYear:  30000,
		}

		if err := db.Create(&vehicle2).Error; err != nil {
			log.Printf("❌ Failed to create vehicle 2: %v", err)
			return err
		}
		log.Printf("   ✅ Vehicle 2 created (ID: %d)", vehicle2.ID)
	} else if err != nil {
		log.Printf("❌ Error loading vehicle 2: %v", err)
		return err
	} else {
		log.Printf("   ⏭️  Vehicle 2 already exists (ID: %d)", vehicle2.ID)
	}

	log.Println("✅ ========================================")
	log.Println("   Test data seeding completed!")
	log.Println("   👤 Driver:    driver@test.com / driver123")
	log.Println("   👤 Collector: collector@test.com / collector123")
	log.Println("   🚛 Vehicle 1: กข-1234 (Driver assigned)")
	log.Println("   🚛 Vehicle 2: คง-5678 (No driver)")
	log.Println("   ========================================")

	return nil
}
