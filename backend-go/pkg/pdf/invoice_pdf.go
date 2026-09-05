package pdf

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/fossbilling/backend-go/internal/domain"
)

type InvoicePDFData struct {
	Invoice     *domain.Invoice
	Client      *domain.Client
	CompanyName string
	CompanyAddr string
	CompanyTax  string
}

const invoiceHTMLTemplate = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>Invoice #{{.Invoice.Nr}}</title>
  <style>
    body { font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; color: #333; margin: 0; padding: 40px; }
    .invoice-box { max-width: 800px; margin: auto; padding: 30px; border: 1px solid #eee; box-shadow: 0 0 10px rgba(0, 0, 0, 0.15); font-size: 14px; line-height: 24px; }
    .header { display: flex; justify-content: space-between; border-bottom: 2px solid #206bc4; padding-bottom: 20px; margin-bottom: 20px; }
    .company-title { font-size: 24px; font-weight: bold; color: #206bc4; }
    .invoice-title { font-size: 24px; font-weight: bold; text-align: right; color: #333; }
    .status-stamp { display: inline-block; padding: 4px 12px; border-radius: 4px; font-weight: bold; text-transform: uppercase; }
    .status-paid { background: #d4edda; color: #155724; }
    .status-unpaid { background: #f8d7da; color: #721c24; }
    .meta-table { width: 100%; margin-bottom: 20px; }
    .meta-table td { vertical-align: top; }
    .items-table { width: 100%; border-collapse: collapse; margin-top: 20px; }
    .items-table th { background: #f8fafc; border-bottom: 2px solid #dee2e6; padding: 10px; text-align: left; }
    .items-table td { padding: 10px; border-bottom: 1px solid #eee; }
    .total-section { margin-top: 20px; float: right; width: 300px; }
    .total-row { display: flex; justify-content: space-between; padding: 6px 0; }
    .grand-total { font-size: 18px; font-weight: bold; color: #206bc4; border-top: 2px solid #206bc4; }
  </style>
</head>
<body>
  <div class="invoice-box">
    <table style="width: 100%;">
      <tr>
        <td>
          <div class="company-title">{{.CompanyName}}</div>
          <div>{{.CompanyAddr}}</div>
          <div>NPWP/Tax ID: {{.CompanyTax}}</div>
        </td>
        <td style="text-align: right;">
          <div class="invoice-title">INVOICE</div>
          <div><strong>#{{.Invoice.Nr}}</strong></div>
          <div>Tanggal: {{.Invoice.CreatedAt.Format "02 Jan 2006"}}</div>
          <div>Jatuh Tempo: {{.Invoice.DueAt.Format "02 Jan 2006"}}</div>
          <div style="margin-top: 5px;">
            {{if eq .Invoice.Status "paid"}}
              <span class="status-stamp status-paid">LUNAS / PAID</span>
            {{else}}
              <span class="status-stamp status-unpaid">UNPAID</span>
            {{end}}
          </div>
        </td>
      </tr>
    </table>

    <hr style="border: 0; border-top: 1px solid #eee; margin: 20px 0;">

    <table class="meta-table">
      <tr>
        <td>
          <strong>Ditagihkan Kepada:</strong><br>
          {{.Client.FirstName}} {{.Client.LastName}}<br>
          {{if .Client.Company}}{{.Client.Company}}<br>{{end}}
          {{.Client.Email}}<br>
          {{.Client.Country}}
        </td>
      </tr>
    </table>

    <table class="items-table">
      <thead>
        <tr>
          <th>Deskripsi Layanan / Produk</th>
          <th style="text-align: center;">Periode</th>
          <th style="text-align: center;">Jumlah</th>
          <th style="text-align: right;">Harga Satuan</th>
          <th style="text-align: right;">Total</th>
        </tr>
      </thead>
      <tbody>
        {{range .Invoice.Items}}
        <tr>
          <td>{{.Title}}</td>
          <td style="text-align: center;">{{if .Period}}{{.Period}}{{else}}-{{end}}</td>
          <td style="text-align: center;">{{.Quantity}}</td>
          <td style="text-align: right;">{{$.Invoice.Currency}} {{.Price.String}}</td>
          <td style="text-align: right;">{{$.Invoice.Currency}} {{.Price.String}}</td>
        </tr>
        {{end}}
      </tbody>
    </table>

    <div class="total-section">
      <div class="total-row">
        <span>Subtotal:</span>
        <span>{{.Invoice.Currency}} {{.Invoice.Subtotal.String}}</span>
      </div>
      <div class="total-row">
        <span>Pajak (Tax):</span>
        <span>{{.Invoice.Currency}} {{.Invoice.Tax.String}}</span>
      </div>
      <div class="total-row grand-total">
        <span>Total:</span>
        <span>{{.Invoice.Currency}} {{.Invoice.Total.String}}</span>
      </div>
    </div>
    <div style="clear: both;"></div>

    <div style="margin-top: 40px; padding-top: 20px; border-top: 1px solid #eee; font-size: 12px; color: #888;">
      <p>Terima kasih atas kepercayaan Anda menggunakan layanan {{.CompanyName}}.</p>
    </div>
  </div>
</body>
</html>
`

// GenerateInvoiceHTML renders the complete styled HTML invoice
func GenerateInvoiceHTML(inv *domain.Invoice, client *domain.Client, companyName, companyAddr, companyTax string) ([]byte, error) {
	if companyName == "" {
		companyName = "FOSSBilling Enterprise"
	}
	if companyAddr == "" {
		companyAddr = "Cyber 2 Tower, Kuningan, Jakarta Selatan"
	}
	if companyTax == "" {
		companyTax = "01.234.567.8-901.000"
	}

	data := InvoicePDFData{
		Invoice:     inv,
		Client:      client,
		CompanyName: companyName,
		CompanyAddr: companyAddr,
		CompanyTax:  companyTax,
	}

	tmpl, err := template.New("invoicePDF").Parse(invoiceHTMLTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse invoice template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute invoice template: %w", err)
	}

	return buf.Bytes(), nil
}
