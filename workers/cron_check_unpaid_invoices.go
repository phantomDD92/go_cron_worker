package workers

import (
	"go_proxy_worker/slack"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/invoice"
)

func CronCheckUnpaidInvoices() {
	s := gocron.NewScheduler(time.UTC)
	s.Every(24).Hours().Do(CheckUnpaidInvoices)
	s.StartBlocking()
}

func CheckUnpaidInvoices() {
	// fileName := "cron_check_unpaid_invoices.go"
	// emptyErrMap := make(map[string]interface{})

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	params := &stripe.InvoiceListParams{
		Status: stripe.String("open"),
	}

	i := invoice.List(params)

	var unpaidInvoices []*stripe.Invoice
	for i.Next() {
		inv := i.Invoice()

		// Check if the invoice is unpaid (open)
		if inv.Status == "open" {
			if inv.AmountDue > 24900 {
				unpaidInvoices = append(unpaidInvoices, inv)
			}
		}
	}

	slack.SlackUnpaidInvoiceAlert("#stripe-alerts", unpaidInvoices)
}
