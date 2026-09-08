# Customizing Invoice PDFs

FOSSBilling generates high-quality PDF invoices using a built-in Go rendering engine. You can customize the look and feel of these documents to match your brand.

---

## 🎨 Branding the PDF

The simplest way to customize your invoices is through the **Company Information** settings in the Admin Portal:
- **Logo:** Upload your "Dark Logo" (optimized for white backgrounds) to be displayed in the top header.
- **Address:** Your official company address will be rendered in the seller details section.
- **Footer Text:** Add legal notes or payment instructions.

---

## 🏗️ Technical Architecture

The PDF generation logic resides in `pkg/pdf/`. It follows a structured template approach:

- **Layout Engine:** Uses a coordinate-based positioning system for maximum precision.
- **Font Support:** Supports embedded Unicode fonts for international character rendering (RTL, CJK, etc.).
- **Data Source:** Injects a `domain.Invoice` and `domain.Client` object into the generator.

---

## 🛠️ Advanced Customization

If you need to change the actual layout (e.g., move the tax column or add a custom field):

1. Navigate to `backend-go/pkg/pdf/invoice_generator.go`.
2. Locate the drawing functions for headers, line items, and footers.
3. Modify the coordinates (`X`, `Y`) or font styles.
4. Rebuild the backend binary: `make build`.

### Dynamic Data
You can inject any field from the invoice or client record. If you have custom client attributes, ensure they are fetched in the `InvoiceService` before being passed to the PDF generator.

---

## 📄 PDF Specifications
- **Format:** A4
- **DPI:** 72 (Standard PDF units)
- **Security:** Invoices are generated on-the-fly and never stored as files on the server for enhanced privacy.
