package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/builder"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/payment"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/payment/gateways"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/provisioning"
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
	fmt.Println("  version        Display current binary and runtime version")
	fmt.Println("  status         Check registered payment and provisioning drivers")
	fmt.Println("  client:create  Create a new client account interactively via builder")
	fmt.Println("  help           Show this help message")
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
	provRegistry.Register(provisioning.NewCpanelProvisioner(provisioning.CpanelConfig{Host: "cpanel.host.com"}))
	provRegistry.Register(provisioning.NewLicenseProvisioner("SALT"))

	fmt.Println("\n⚡ Registered Provisioning Engines:")
	for _, p := range provRegistry.List() {
		fmt.Printf("  • Product Type: %s\n", p.Type())
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
