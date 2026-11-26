package customers

import (
	"github.com/jeremybumsted/plain-cli/internal/output"
	"github.com/jeremybumsted/plain-cli/internal/plain"
)

// printCustomerTable formats and prints customers as a table
func printCustomerTable(formatter *output.Formatter, customers []plain.Customer) error {
	if len(customers) == 0 {
		return formatter.Info("No customers found")
	}

	headers := []string{"ID", "Name", "Email", "External ID", "Updated"}
	rows := make([][]string, 0, len(customers))

	for _, c := range customers {
		name := c.FullName
		if len(name) > 30 {
			name = name[:27] + "..."
		}
		if name == "" {
			name = "(no name)"
		}

		email := c.Email.Email
		if email == "" {
			email = "-"
		}

		externalID := c.ExternalID
		if externalID == "" {
			externalID = "-"
		}

		updatedAt := c.UpdatedAt.ISO8601
		if len(updatedAt) > 10 {
			updatedAt = updatedAt[:10]
		}

		row := []string{
			c.ID,
			name,
			email,
			externalID,
			updatedAt,
		}
		rows = append(rows, row)
	}

	return formatter.PrintTable(headers, rows)
}

// displayCustomer displays a single customer's details
func displayCustomer(formatter *output.Formatter, customer *plain.Customer) error {
	if err := formatter.Print("Customer Details"); err != nil {
		return err
	}
	if err := formatter.Print("----------------"); err != nil {
		return err
	}
	if err := formatter.Printf("ID:          %s\n", customer.ID); err != nil {
		return err
	}
	if err := formatter.Printf("Name:        %s\n", customer.FullName); err != nil {
		return err
	}
	if customer.ShortName != "" {
		if err := formatter.Printf("Short Name:  %s\n", customer.ShortName); err != nil {
			return err
		}
	}
	if err := formatter.Printf("Email:       %s\n", customer.Email.Email); err != nil {
		return err
	}
	if customer.ExternalID != "" {
		if err := formatter.Printf("External ID: %s\n", customer.ExternalID); err != nil {
			return err
		}
	}

	createdAt, _ := customer.CreatedAt.Time()
	updatedAt, _ := customer.UpdatedAt.Time()
	if err := formatter.Printf("Created:     %s\n", createdAt.Format("2006-01-02 15:04:05")); err != nil {
		return err
	}
	if err := formatter.Printf("Updated:     %s\n", updatedAt.Format("2006-01-02 15:04:05")); err != nil {
		return err
	}

	return nil
}
