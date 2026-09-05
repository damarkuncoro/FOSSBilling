package notification

const welcomeTemplate = `
<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
  <div style="max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e0e0e0; border-radius: 8px;">
    <h2 style="color: #206bc4;">Selamat Datang di {{.AppName}}, {{.FirstName}}!</h2>
    <p>Akun Anda telah berhasil dibuat. Anda sekarang dapat masuk ke Client Portal untuk mengelola layanan.</p>
    <p><strong>Email Login:</strong> {{.Email}}</p>
    <hr style="border: 0; border-top: 1px solid #eee; margin: 20px 0;">
    <p style="font-size: 12px; color: #888;">&copy; {{.AppName}}</p>
  </div>
</body>
</html>`

const invoiceCreatedTemplate = `
<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
  <div style="max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e0e0e0; border-radius: 8px;">
    <h2 style="color: #206bc4;">Tagihan Baru Terbit #{{.InvoiceNr}}</h2>
    <p>Halo {{.FirstName}},</p>
    <ul>
      <li><strong>Nomor Invoice:</strong> {{.InvoiceNr}}</li>
      <li><strong>Total Tagihan:</strong> {{.Currency}} {{.Total}}</li>
      <li><strong>Batas Pembayaran:</strong> {{.DueAt}}</li>
    </ul>
    <hr style="border: 0; border-top: 1px solid #eee; margin: 20px 0;">
    <p style="font-size: 12px; color: #888;">&copy; {{.AppName}}</p>
  </div>
</body>
</html>`

const paymentReceiptTemplate = `
<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
  <div style="max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e0e0e0; border-radius: 8px;">
    <h2 style="color: #2fb344;">Bukti Pembayaran Lunas #{{.InvoiceNr}}</h2>
    <p>Halo {{.FirstName}},</p>
    <ul>
      <li><strong>Nomor Invoice:</strong> {{.InvoiceNr}}</li>
      <li><strong>Jumlah Dibayar:</strong> {{.Currency}} {{.Total}}</li>
      <li><strong>ID Transaksi:</strong> {{.TxnID}}</li>
    </ul>
    <hr style="border: 0; border-top: 1px solid #eee; margin: 20px 0;">
    <p style="font-size: 12px; color: #888;">&copy; {{.AppName}}</p>
  </div>
</body>
</html>`

const ticketReplyTemplate = `
<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
  <div style="max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e0e0e0; border-radius: 8px;">
    <h2 style="color: #206bc4;">Balasan Baru pada Tiket [#{{.TicketID}}]</h2>
    <p>Halo {{.FirstName}}, staf kami telah membalas tiket Anda: <strong>{{.Subject}}</strong></p>
    <div style="background: #f8fafc; border-left: 4px solid #206bc4; padding: 12px; margin: 15px 0;">
      {{.Message}}
    </div>
    <hr style="border: 0; border-top: 1px solid #eee; margin: 20px 0;">
    <p style="font-size: 12px; color: #888;">&copy; {{.AppName}}</p>
  </div>
</body>
</html>`
