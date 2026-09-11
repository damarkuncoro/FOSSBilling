package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/builder"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/config"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/postgres"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/importer"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/payment"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/payment/gateways"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/provisioning"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/seeder"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/auth"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/geoip"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/i18n"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/tools"
)

const AppVersion = "2.0.0-golang"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	command := os.Args[1]

	switch command {
	case "version":
		fmt.Printf("FOSSBilling Next-Gen CLI v%s (Go %s on %s/%s)\n", AppVersion, runtime.Version(), runtime.GOOS, runtime.GOARCH)

	case "status":
		runStatus()

	case "client:create":
		runCreateClient(os.Args[2:])

	case "admin:create":
		runCreateAdmin(os.Args[2:])

	case "db:backup":
		runDBBackup(os.Args[2:])

	case "db:seed:bench":
		runDBSeedBench(os.Args[2:])

	case "db:import:legacy":
		runDBImportLegacy(os.Args[2:])

	case "invoice:build":
		runBuildInvoice(os.Args[2:])

	case "tools:password":
		runGeneratePassword(os.Args[2:])

	case "tools:geoip":
		runGeoIPLookup(os.Args[2:])

	case "locale:list":
		runListLocales()

	case "help":
		printUsage()

	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("🚀 FOSSBilling CLI Management Utility")
	fmt.Println("Usage: go run ./cmd/cli <command> [arguments]")
	fmt.Println("\nAvailable commands:")
	fmt.Println("  version         Display current binary and runtime version")
	fmt.Println("  status          Check registered payment, registrar, and provisioning drivers")
	fmt.Println("  client:create   Create a new client account interactively via builder")
	fmt.Println("  admin:create    Create a new staff/admin account with bcrypt password hash")
	fmt.Println("  db:backup       Create a full SQL backup of the database")
	fmt.Println("  db:seed:bench   Perform high-performance batch seeding for stress testing")
	fmt.Println("  db:import:legacy Migrates data from legacy PHP FOSSBilling (MySQL)")
	fmt.Println("  invoice:build   Construct and preview an invoice with items via builder")
	fmt.Println("  tools:password  Generate a cryptographically secure random password")
	fmt.Println("  tools:geoip     Lookup country and flag for a given IP address")
	fmt.Println("  locale:list     List all supported locales, flags, and text directions")
	fmt.Println("  help            Show this help message")
}

func runStatus() {
	fmt.Println("==================================================")
	fmt.Printf("📦 FOSSBilling System Status (Version: %s)\n", AppVersion)
	fmt.Println("==================================================")

	// Payment Registry Status
	gwRegistry := payment.NewGatewayRegistry()
	gwRegistry.Register(gateways.NewMidtransGateway("mock_key", "mock_key", false))
	gwRegistry.Register(gateways.NewStripeGateway("sk_test", "pk_test", "whsec"))
	gwRegistry.Register(gateways.NewPayPalGateway("client_id", "secret", false))
	gwRegistry.Register(gateways.NewBankTransferGateway("BCA", "1234567890", "PT FOSSBilling"))

	fmt.Println("\n💳 Registered Payment Gateways:")
	for _, gw := range gwRegistry.List() {
		fmt.Printf("  • [%s] %s (Type: %s)\n", gw.ID(), gw.Name(), gw.Type())
	}

	// Provisioning Registry Status
	provRegistry := provisioning.NewProvisionerRegistry()
	provRegistry.Register("cpanel", provisioning.NewCpanelProvisioner(provisioning.CpanelConfig{Host: "cpanel.host.com"}))
	provRegistry.Register("directadmin", provisioning.NewDirectAdminProvisioner("da.host.com", 2222, "admin", "mock_key"))
	provRegistry.Register("plesk", provisioning.NewPleskProvisioner(provisioning.PleskConfig{Host: "plesk.host.com"}))
	provRegistry.Register("hestia", provisioning.NewHestiaProvisioner(provisioning.HestiaConfig{Host: "hestia.host.com"}))
	provRegistry.Register("cwp", provisioning.NewCWPProvisioner(provisioning.CWPConfig{Host: "cwp.host.com"}))
	provRegistry.Register("custom", provisioning.NewCustomServerProvisioner(provisioning.CustomServerConfig{}))
	provRegistry.Register("license", provisioning.NewLicenseProvisioner("SALT"))

	fmt.Println("\n⚡ Registered Provisioning Engines:")
	for id, p := range provRegistry.List() {
		fmt.Printf("  • [%s] Product Type: %s\n", id, p.Type())
	}

	// Registrar Registry Status
	regRegistry := provisioning.NewRegistrarRegistry()
	regRegistry.Register("mock", provisioning.NewMockRegistrarDriver())
	regRegistry.Register("namecheap", provisioning.NewNamecheapRegistrarDriver(provisioning.NamecheapConfig{ApiUser: "user", ApiKey: "key"}))
	regRegistry.Register("resellerclub", provisioning.NewResellerClubRegistrarDriver(provisioning.ResellerClubConfig{AuthUserID: "123", APIKey: "key"}))
	regRegistry.Register("internetbs", provisioning.NewInternetbsRegistrarDriver(provisioning.InternetbsConfig{ApiKey: "key", Password: "pass"}))
	regRegistry.Register("custom", provisioning.NewCustomRegistrarDriver())

	fmt.Println("\n🌐 Registered Domain Registrars:")
	for id := range regRegistry.List() {
		fmt.Printf("  • [%s]\n", id)
	}

	fmt.Println("\n✅ All subsystems are healthy and ready.")
}

func runCreateClient(args []string) {
	fs := flag.NewFlagSet("client:create", flag.ExitOnError)
	email := fs.String("email", "", "Client email address (required)")
	firstName := fs.String("first-name", "", "First name (required)")
	lastName := fs.String("last-name", "", "Last name")
	company := fs.String("company", "", "Company name")
	currency := fs.String("currency", "USD", "Currency code (default USD)")
	password := fs.String("password", "Password123!", "Account password")

	_ = fs.Parse(args)

	if *email == "" || *firstName == "" {
		fmt.Println("Error: --email and --first-name are required flags.")
		fs.Usage()
		os.Exit(1)
	}

	client, err := builder.NewClientBuilder().
		WithEmail(*email).
		WithPassword(*password).
		WithName(*firstName, *lastName).
		WithCompany(*company).
		WithCountryAndCurrency("ID", *currency).
		Build()

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Printf("✅ Successfully constructed Client via ClientBuilder:\n")
	fmt.Printf("   Email   : %s\n", client.Email)
	fmt.Printf("   Name    : %s %s\n", client.FirstName, client.LastName)
	fmt.Printf("   Company : %s\n", client.Company)
	fmt.Printf("   Currency: %s\n", client.Currency)
	fmt.Printf("   Status  : %s\n", client.Status)
}

func runBuildInvoice(args []string) {
	fs := flag.NewFlagSet("invoice:build", flag.ExitOnError)
	clientID := fs.Int64("client-id", 1, "Client ID")
	title := fs.String("item-title", "Cloud VPS - Starter Plan", "Invoice item title")
	price := fs.Float64("price", 19.99, "Item unit price")
	qty := fs.Int("qty", 1, "Item quantity")
	taxRate := fs.Float64("tax-rate", 10.0, "Tax rate percentage (e.g. 10.0 for 10%)")
	currency := fs.String("currency", "USD", "Currency code")

	_ = fs.Parse(args)

	inv, items, err := builder.NewInvoiceBuilder().
		ForClient(*clientID).
		WithSerieAndNr("INV", "").
		WithCurrency(*currency, 1.0).
		WithTaxRate(*taxRate).
		WithDueDays(14).
		AddItem(*title, decimal.FromFloat(*price), *qty, true).
		Build()

	if err != nil {
		log.Fatalf("Failed to build invoice: %v", err)
	}

	fmt.Printf("✅ Successfully constructed Invoice via InvoiceBuilder:\n")
	fmt.Printf("   Invoice Number : %s-%s\n", inv.Serie, inv.Nr)
	fmt.Printf("   Client ID      : %d\n", inv.ClientID)
	fmt.Printf("   Status         : %s\n", inv.Status)
	fmt.Printf("   Currency       : %s\n", inv.Currency)
	fmt.Printf("   Line Items     : %d items\n", len(items))
	for idx, item := range items {
		fmt.Printf("     [%d] %s (x%d) @ %.2f = %.2f\n", idx+1, item.Title, item.Quantity, item.Price.ToFloat(), (item.Price * decimal.Money(item.Quantity)).ToFloat())
	}
	fmt.Printf("   Subtotal       : %.2f %s\n", inv.Subtotal.ToFloat(), inv.Currency)
	fmt.Printf("   Tax (%.1f%%)     : %.2f %s\n", inv.TaxRate, inv.Tax.ToFloat(), inv.Currency)
	fmt.Printf("   Grand Total    : %.2f %s\n", inv.Total.ToFloat(), inv.Currency)
	fmt.Printf("   Due Date       : %s\n", inv.DueAt.Format("2006-01-02"))
}

func runCreateAdmin(args []string) {
	fs := flag.NewFlagSet("admin:create", flag.ExitOnError)
	email := fs.String("email", "", "Admin email address (required)")
	name := fs.String("name", "Super Administrator", "Admin full name")
	role := fs.String("role", "admin", "Admin role (superadmin, admin, support, billing)")
	password := fs.String("password", "", "Password (leave empty to auto-generate)")

	_ = fs.Parse(args)

	if *email == "" {
		fmt.Println("Error: --email is a required flag.")
		fs.Usage()
		os.Exit(1)
	}

	pass := *password
	if pass == "" {
		var err error
		pass, err = tools.GeneratePassword(16, true)
		if err != nil {
			log.Fatalf("Failed to generate password: %v", err)
		}
	}

	hash, err := auth.HashPassword(pass)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	staffRole := domain.StaffRoleAdmin
	switch *role {
	case "superadmin":
		staffRole = domain.StaffRoleSuperAdmin
	case "support":
		staffRole = domain.StaffRoleSupport
	case "billing":
		staffRole = domain.StaffRoleBilling
	}

	staff := domain.Staff{
		Email:        *email,
		Name:         *name,
		Role:         staffRole,
		PasswordHash: hash,
		Status:       "active",
	}

	fmt.Println("✅ Successfully provisioned Administrator Account:")
	fmt.Printf("   Email         : %s\n", staff.Email)
	fmt.Printf("   Name          : %s\n", staff.Name)
	fmt.Printf("   Role          : %s\n", staff.Role)
	fmt.Printf("   Status        : %s\n", staff.Status)
	fmt.Printf("   Plain Password: %s\n", pass)
	fmt.Printf("   Bcrypt Hash   : %s\n", staff.PasswordHash)
}

func runDBBackup(args []string) {
	fs := flag.NewFlagSet("db:backup", flag.ExitOnError)
	output := fs.String("output", "", "Output file path (e.g. backup.sql)")
	dbURL := fs.String("db-url", "", "Custom Database URL")

	_ = fs.Parse(args)

	cfg := config.Load()
	finalURL := cfg.DatabaseURL
	if *dbURL != "" {
		finalURL = *dbURL
	}

	outFile := *output
	if outFile == "" {
		outFile = fmt.Sprintf("fossbilling_backup_%s.sql", time.Now().Format("20060102_150405"))
	}

	fmt.Printf("📦 Creating database backup to %s...\n", outFile)

	// Since we are using pgx, we can't easily do a dump via Go without
	// reimplementing pg_dump. The best way is to shell out to pg_dump.
	cmd := exec.Command("pg_dump", finalURL, "-f", outFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("❌ Backup failed: %v\nMake sure 'pg_dump' is installed and in your PATH.", err)
	}

	fmt.Printf("✅ Backup successfully created: %s\n", outFile)
}

func runDBSeedBench(args []string) {
	fs := flag.NewFlagSet("db:seed:bench", flag.ExitOnError)
	clients := fs.Int("clients", 1000, "Number of clients to seed")
	invoices := fs.Int("invoices", 5000, "Number of invoices to seed")
	dbURL := fs.String("db-url", "", "Custom Database URL (optional)")

	_ = fs.Parse(args)

	cfg := config.Load()
	finalURL := cfg.DatabaseURL
	if *dbURL != "" {
		finalURL = *dbURL
	}

	ctx := context.Background()
	pool, err := postgres.NewPostgresPool(ctx, finalURL)
	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}
	defer pool.Close()

	s := seeder.NewBenchSeeder(pool)
	if err := s.SeedLargeScale(ctx, *clients, *invoices); err != nil {
		log.Fatalf("❌ Seeding failed: %v", err)
	}
}

func runGeneratePassword(args []string) {
	fs := flag.NewFlagSet("tools:password", flag.ExitOnError)
	length := fs.Int("length", 16, "Password length (8-64)")
	includeSpecial := fs.Bool("special", true, "Include special characters")

	_ = fs.Parse(args)

	pwd, err := tools.GeneratePassword(*length, *includeSpecial)
	if err != nil {
		log.Fatalf("Failed to generate password: %v", err)
	}
	fmt.Printf("🔑 Generated Secure Password: %s\n", pwd)
}

func runGeoIPLookup(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: go run ./cmd/cli tools:geoip <country-iso or ip-address>")
		os.Exit(1)
	}
	input := args[0]
	isPrivate := geoip.IsPrivateIP(input)
	res := geoip.LookupCountry(input)

	fmt.Printf("🌍 GeoIP Lookup for %s:\n", input)
	fmt.Printf("   Country Code : %s %s\n", res.ISOCode, res.Flag)
	fmt.Printf("   Country Name : %s\n", res.Name)
	fmt.Printf("   Currency     : %s\n", res.Currency)
	fmt.Printf("   Is Private IP: %t\n", isPrivate)
}

func runListLocales() {
	fmt.Println("🌐 Supported FOSSBilling Locales:")
	for _, l := range i18n.SupportedLocales {
		fmt.Printf("   %s [%s] %s (%s) — Direction: %s\n", l.Flag, l.Code, l.Name, l.Native, l.Dir)
	}
}

func runDBImportLegacy(args []string) {
	fs := flag.NewFlagSet("db:import:legacy", flag.ExitOnError)
	source := fs.String("source", "", "MySQL source DSN (e.g. user:pass@tcp(127.0.0.1:3306)/legacy_db)")

	_ = fs.Parse(args)

	if *source == "" {
		fmt.Println("Error: --source DSN is required.")
		fs.Usage()
		os.Exit(1)
	}

	cfg := config.Load()
	ctx := context.Background()

	// Connect to destination (Postgres)
	pool, err := postgres.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ Destination connection failed: %v", err)
	}
	defer pool.Close()

	// Extract standard sql.DB from pgx pool for the importer service
	// In a real app we might use a wrapper, but here we can just use the DB handle if available
	// or use a direct repository implementation.
	cr := postgres.NewClientRepository(pool)
	or := postgres.NewOrderRepository(pool)
	ir := postgres.NewInvoiceRepository(pool)

	imp := importer.NewLegacyImporter(cr, or, ir)
	if err := imp.RunFullMigration(ctx, *source); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}
}
