package main

//import (
//	"context"
//	"github.com/SwanHtetAungPhyo/migration_service/providers/config"
//	"github.com/SwanHtetAungPhyo/migration_service/providers/database"
//	"go.uber.org/fx"
//	"gorm.io/gorm"
//	"log"
//)
//
//func main.go() {
//	app := fx.New(
//		fx.Provide(
//			config.LoadConfig,
//			func(cfg *config.Config) string {
//				return cfg.DB.DSN
//			},
//			database.NewDB,
//		),
//		fx.Invoke(func(db *gorm.DB, lc fx.Lifecycle) {
//			lc.Append(fx.Hook{
//				OnStart: func(ctx context.Context) error {
//					log.Println("Starting database migration...")
//					return database.Migrate(db)
//				},
//				OnStop: func(ctx context.Context) error {
//					log.Println("Closing database connection...")
//					sqlDB, err := db.DB()
//					if err != nil {
//						return err
//					}
//					return sqlDB.Close()
//				},
//			})
//		}),
//	)
//
//	app.Run()
//}

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/nedpals/supabase-go"
)

func main() {
	supabaseUrl := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")
	supabase := supabase.CreateClient(supabaseUrl, supabaseKey)

	user, err := supabase.Auth.SignUp(context.Background(), supabase.User{
		Email:    "user@example.com",
		Password: "securepassword",
	})
	if err != nil {
		log.Fatalf("Signup error: %v", err)
	}
	fmt.Printf("Signed up user: %+v\n", user)

	// 2. Send verification email
	err = supabase.Auth.SendMagicLink(context.Background(), "user@example.com")
	if err != nil {
		log.Fatalf("Verification email error: %v", err)
	}
	fmt.Println("Verification email sent")

	// 3. Forgot password flow
	err = supabase.Auth.ResetPasswordForEmail(context.Background(), "user@example.com")
	if err != nil {
		log.Fatalf("Password reset error: %v", err)
	}
	fmt.Println("Password reset email sent")

	// Note: In a real application, steps 4-5 would happen after user clicks the email link
	// 4. Verify email (when user clicks link)
	// verifyToken := "from-email-query-params"
	// verifiedUser, err := supabase.Auth.VerifyEmail(context.Background(), verifyToken)

	// 5. Reset password (when user submits new password)
	// resetToken := "from-email-query-params"
	// err = supabase.Auth.UpdateUser(context.Background(), resetToken, supabase.UserCredentials{
	//     Password: "new_secure_password",
	// })
}
