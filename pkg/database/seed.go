package database

import (
	"car-rental/pkg/models"
	"log"
)

// SeedData initializes the database with default memberships and booking types
func SeedData() {
	log.Println("Checking if database needs seeding...")

	// Check if memberships table is empty
	var membershipCount int64
	DB.Model(&models.Membership{}).Count(&membershipCount)

	if membershipCount == 0 {
		log.Println("Seeding memberships...")
		// Seed memberships
		memberships := []models.Membership{
			{
				MembershipName: "Bronze",
				Discount:       4.0, // 4%
			},
			{
				MembershipName: "Silver",
				Discount:       7.0, // 7%
			},
			{
				MembershipName: "Gold",
				Discount:       15.0, // 15%
			},
		}

		for _, membership := range memberships {
			if err := DB.Create(&membership).Error; err != nil {
				log.Printf("Error creating membership %s: %v", membership.MembershipName, err)
			} else {
				log.Printf("Created membership: %s (%v%% discount)", membership.MembershipName, membership.Discount)
			}
		}
	} else {
		log.Printf("Memberships table already has %d records, skipping seed", membershipCount)
	}

	// Check if booking types table is empty
	var bookingTypeCount int64
	DB.Model(&models.BookingType{}).Count(&bookingTypeCount)

	if bookingTypeCount == 0 {
		log.Println("Seeding booking types...")
		// Seed booking types
		bookingTypes := []models.BookingType{
			{
				BookingType: "Car Only",
				Description: "Rent Car only",
			},
			{
				BookingType: "Car & Driver",
				Description: "Rent Car and a Driver",
			},
		}

		for _, bookingType := range bookingTypes {
			if err := DB.Create(&bookingType).Error; err != nil {
				log.Printf("Error creating booking type %s: %v", bookingType.BookingType, err)
			} else {
				log.Printf("Created booking type: %s", bookingType.BookingType)
			}
		}
	} else {
		log.Printf("Booking types table already has %d records, skipping seed", bookingTypeCount)
	}

	log.Println("Database seeding completed")
}
