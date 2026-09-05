package pdf

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

// GenerateInvoicePDF creates a pixel-perfect, professionally styled PDF 1.4 document
func GenerateInvoicePDF(inv *domain.Invoice, client *domain.Client, compName, compAddr, compTax string) ([]byte, error) {
	if compName == "" {
		compName = "FOSSBilling Enterprise"
	}
	if compAddr == "" {
		compAddr = "Cyber 2 Tower Lt. 18, Jakarta Selatan"
	}

	var stream bytes.Buffer

	// 1. Top Brand Header Bar
	stream.WriteString("0.125 0.420 0.769 rg\n40 765 515 42 re f\n")
	stream.WriteString("1 1 1 rg\n")
	stream.WriteString("BT\n/F2 18 Tf\n55 780 Td\n(" + escapePDF(compName) + ") Tj\nET\n")
	stream.WriteString("BT\n/F2 18 Tf\n440 780 Td\n(INVOICE) Tj\nET\n")

	// 2. Info Cards (Company & Invoice Meta)
	stream.WriteString("0.96 0.97 0.99 rg\n40 670 250 82 re f\n")
	stream.WriteString("0.88 0.91 0.94 RG\n1 w\n40 670 250 82 re S\n")

	stream.WriteString("0.12 0.16 0.22 rg\n")
	stream.WriteString("BT\n/F2 10 Tf\n52 735 Td\n(ISSUER INFORMATION) Tj\nET\n")
	stream.WriteString("0.40 0.45 0.53 rg\n")
	stream.WriteString("BT\n/F1 9 Tf\n52 720 Td\n(" + escapePDF(compAddr) + ") Tj\nET\n")
	if compTax != "" {
		stream.WriteString("BT\n/F1 9 Tf\n52 706 Td\n(Tax ID / NPWP: " + escapePDF(compTax) + ") Tj\nET\n")
	}
	stream.WriteString("BT\n/F1 9 Tf\n52 692 Td\n(Official Automated Tax Statement) Tj\nET\n")

	// Right Meta Card
	stream.WriteString("0.96 0.97 0.99 rg\n305 670 250 82 re f\n")
	stream.WriteString("0.88 0.91 0.94 RG\n305 670 250 82 re S\n")

	stream.WriteString("0.12 0.16 0.22 rg\n")
	stream.WriteString("BT\n/F2 10 Tf\n317 735 Td\n(INVOICE DETAILS) Tj\nET\n")
	stream.WriteString("0.40 0.45 0.53 rg\n")
	stream.WriteString("BT\n/F1 9 Tf\n317 720 Td\n(Invoice No: #" + escapePDF(inv.Nr) + ") Tj\nET\n")
	stream.WriteString("BT\n/F1 9 Tf\n317 706 Td\n(Issue Date: " + inv.CreatedAt.Format("02 Jan 2006") + ") Tj\nET\n")
	stream.WriteString("BT\n/F1 9 Tf\n317 692 Td\n(Due Date: " + inv.DueAt.Format("02 Jan 2006") + ") Tj\nET\n")

	// Status Stamp Pill
	if inv.Status == domain.InvoiceStatusPaid {
		stream.WriteString("0.85 0.95 0.88 rg\n465 722 80 20 re f\n")
		stream.WriteString("0.08 0.45 0.22 rg\nBT\n/F2 9 Tf\n478 728 Td\n(PAID) Tj\nET\n")
	} else {
		stream.WriteString("0.99 0.88 0.88 rg\n465 722 80 20 re f\n")
		stream.WriteString("0.72 0.11 0.11 rg\nBT\n/F2 9 Tf\n473 728 Td\n(UNPAID) Tj\nET\n")
	}

	// 3. Customer Box
	stream.WriteString("0.98 0.98 0.99 rg\n40 595 515 62 re f\n")
	stream.WriteString("0.88 0.91 0.94 RG\n40 595 515 62 re S\n")
	stream.WriteString("0.12 0.16 0.22 rg\n")
	stream.WriteString("BT\n/F2 10 Tf\n52 640 Td\n(BILLED TO:) Tj\nET\n")
	stream.WriteString("0.20 0.25 0.33 rg\n")
	clientName := client.FirstName + " " + client.LastName
	if client.Company != "" {
		clientName += " (" + client.Company + ")"
	}
	stream.WriteString("BT\n/F2 9 Tf\n52 625 Td\n(" + escapePDF(clientName) + ") Tj\nET\n")
	stream.WriteString("0.40 0.45 0.53 rg\n")
	stream.WriteString("BT\n/F1 9 Tf\n52 610 Td\n(" + escapePDF(client.Email) + " - Country: " + escapePDF(client.Country) + ") Tj\nET\n")

	// 4. Items Table Header
	stream.WriteString("0.12 0.16 0.24 rg\n40 560 515 24 re f\n")
	stream.WriteString("1 1 1 rg\n")
	stream.WriteString("BT\n/F2 9 Tf\n52 568 Td\n(ITEM / SERVICE DESCRIPTION) Tj\nET\n")
	stream.WriteString("BT\n/F2 9 Tf\n355 568 Td\n(QTY) Tj\nET\n")
	stream.WriteString("BT\n/F2 9 Tf\n410 568 Td\n(UNIT PRICE) Tj\nET\n")
	stream.WriteString("BT\n/F2 9 Tf\n495 568 Td\n(TOTAL) Tj\nET\n")

	// Table Rows
	y := 538
	for i, item := range inv.Items {
		if i%2 == 1 {
			stream.WriteString(fmt.Sprintf("0.97 0.98 0.99 rg\n40 %d 515 22 re f\n", y-4))
		}
		stream.WriteString(fmt.Sprintf("0.88 0.91 0.94 RG\n40 %d 515 22 re S\n", y-4))
		stream.WriteString("0.15 0.20 0.28 rg\n")
		stream.WriteString(fmt.Sprintf("BT\n/F1 9 Tf\n52 %d Td\n(%s) Tj\nET\n", y+3, escapePDF(item.Title)))
		stream.WriteString(fmt.Sprintf("BT\n/F1 9 Tf\n362 %d Td\n(%d) Tj\nET\n", y+3, item.Quantity))

		priceStr := fmt.Sprintf("%s %s", inv.Currency, item.Price.String())
		totalStr := fmt.Sprintf("%s %s", inv.Currency, item.Price.String())
		stream.WriteString(fmt.Sprintf("BT\n/F1 9 Tf\n410 %d Td\n(%s) Tj\nET\n", y+3, escapePDF(priceStr)))
		stream.WriteString(fmt.Sprintf("BT\n/F2 9 Tf\n490 %d Td\n(%s) Tj\nET\n", y+3, escapePDF(totalStr)))
		y -= 22
	}

	// 5. Summary Card (Right-aligned)
	boxY := y - 75
	stream.WriteString(fmt.Sprintf("0.96 0.97 0.99 rg\n315 %d 240 75 re f\n", boxY))
	stream.WriteString(fmt.Sprintf("0.88 0.91 0.94 RG\n315 %d 240 75 re S\n", boxY))

	stream.WriteString("0.40 0.45 0.53 rg\n")
	stream.WriteString(fmt.Sprintf("BT\n/F1 9 Tf\n327 %d Td\n(Subtotal:) Tj\nET\n", boxY+55))
	stream.WriteString(fmt.Sprintf("BT\n/F1 9 Tf\n450 %d Td\n(%s %s) Tj\nET\n", boxY+55, inv.Currency, inv.Subtotal.String()))

	stream.WriteString(fmt.Sprintf("BT\n/F1 9 Tf\n327 %d Td\n(PPN / Tax (11%%):) Tj\nET\n", boxY+38))
	stream.WriteString(fmt.Sprintf("BT\n/F1 9 Tf\n450 %d Td\n(%s %s) Tj\nET\n", boxY+38, inv.Currency, inv.Tax.String()))

	// Grand Total Highlight
	stream.WriteString(fmt.Sprintf("0.125 0.420 0.769 rg\n315 %d 240 26 re f\n", boxY))
	stream.WriteString("1 1 1 rg\n")
	stream.WriteString(fmt.Sprintf("BT\n/F2 11 Tf\n327 %d Td\n(TOTAL AMOUNT:) Tj\nET\n", boxY+8))
	stream.WriteString(fmt.Sprintf("BT\n/F2 11 Tf\n450 %d Td\n(%s %s) Tj\nET\n", boxY+8, inv.Currency, inv.Total.String()))

	// 6. Footer
	stream.WriteString("0.55 0.60 0.68 rg\n")
	stream.WriteString("0.88 0.91 0.94 RG\n40 75 m 555 75 l S\n")
	stream.WriteString("BT\n/F1 8 Tf\n40 60 Td\n(Payment is due within 14 days of issue. For bank transfers, quote invoice number as transaction note.) Tj\nET\n")
	stream.WriteString("BT\n/F2 8 Tf\n40 46 Td\n(Thank you for choosing " + escapePDF(compName) + "!) Tj\nET\n")
	stream.WriteString("BT\n/F1 8 Tf\n440 46 Td\n(Generated by FOSSBilling) Tj\nET\n")

	streamBytes := stream.Bytes()

	var pdf bytes.Buffer
	offsets := make([]int, 7)
	pdf.WriteString("%PDF-1.4\n")

	offsets[1] = pdf.Len()
	pdf.WriteString("1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n")

	offsets[2] = pdf.Len()
	pdf.WriteString("2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n")

	offsets[3] = pdf.Len()
	pdf.WriteString("3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 595 842] /Contents 4 0 R /Resources << /Font << /F1 5 0 R /F2 6 0 R >> >> >>\nendobj\n")

	offsets[4] = pdf.Len()
	pdf.WriteString(fmt.Sprintf("4 0 obj\n<< /Length %d >>\nstream\n", len(streamBytes)))
	pdf.Write(streamBytes)
	pdf.WriteString("\nendstream\nendobj\n")

	offsets[5] = pdf.Len()
	pdf.WriteString("5 0 obj\n<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>\nendobj\n")

	offsets[6] = pdf.Len()
	pdf.WriteString("6 0 obj\n<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica-Bold >>\nendobj\n")

	xrefOffset := pdf.Len()
	pdf.WriteString("xref\n0 7\n0000000000 65535 f \n")
	for i := 1; i <= 6; i++ {
		pdf.WriteString(fmt.Sprintf("%010d 00000 n \n", offsets[i]))
	}

	pdf.WriteString(fmt.Sprintf("trailer\n<< /Size 7 /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", xrefOffset))

	return pdf.Bytes(), nil
}

func escapePDF(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "(", "\\(")
	s = strings.ReplaceAll(s, ")", "\\)")
	return s
}
