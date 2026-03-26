package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aceextension/core/config"
	"github.com/aceextension/core/db"
	"github.com/aceextension/core/logger"
	"github.com/aceextension/identity/models"
	"github.com/aceextension/identity/repository"
	"github.com/aceextension/identity/service"
)

func main() {
	// 1. Define Flags
	name := flag.String("name", "System Admin", "User's full name")
	email := flag.String("email", "", "User's email (required)")
	phone := flag.String("phone", "", "User's phone number")
	password := flag.String("password", "", "User's password (required)")
	role := flag.String("role", "superadmin", "User's role")

	flag.Parse()

	if *email == "" || *password == "" {
		fmt.Println("Error: email and password are required")
		flag.Usage()
		os.Exit(1)
	}

	// 2. Load Configuration & Init DB
	cfg := config.Load()
	logger.Init(cfg.Env)
	db.Init(cfg.DatabaseURL, cfg.AuditDatabaseURL)
	defer db.Close()

	ctx := context.Background()

	// 3. Hash Password using existing utility
	hash, err := service.HashPassword(*password)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	var phonePtr *string
	if *phone != "" {
		phonePtr = phone
	}

	// 4. Create User Object
	user := &models.User{
		Name:         *name,
		Email:        email,
		Phone:        phonePtr,
		PasswordHash: &hash,
		Role:         *role,
		IsVerified:   true, // Superusers are verified by default
		IsActive:     true,
	}

	// 5. Initialize Repository & Save
	authRepo := repository.NewAuthRepository()
	err = authRepo.WithTransaction(ctx, func(tr repository.AuthRepository) error {
		// Note: For superadmins, tenant_id can be null if the schema allows it.
		// If it fails due to NOT NULL constraint, you'll need to create/provide a system tenant ID.
		if err := tr.CreateUser(ctx, user); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to create superuser: %v", err)
	}

	fmt.Printf("✅ Superuser created successfully!\n")
	fmt.Printf("ID: %s\n", user.ID)
	fmt.Printf("Name: %s\n", user.Name)
	fmt.Printf("Email: %s\n", *user.Email)
	fmt.Printf("Role: %s\n", user.Role)
}
