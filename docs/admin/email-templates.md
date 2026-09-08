# Email Templates

FOSSBilling automates customer communication through a series of customizable email templates. These templates use a simple placeholder system to inject dynamic data like client names, invoice IDs, and product titles.

---

## 📂 Template Categories

Emails are grouped by their functional area:
- **Client:** Registration welcome, password reset, and profile updates.
- **Invoice:** New invoice generated, payment reminder, and payment receipt.
- **Service:** Order activation, hosting account login details, and suspension notices.
- **Support:** Ticket received confirmation and staff reply notifications.

---

## ✏️ Customizing a Template

Navigate to **System > Email Templates** to see the list of active templates.

### Placeholder Variables
You can use double braces to inject data. Common variables include:
- `{{ client_name }}`: The first name of the recipient.
- `{{ invoice_nr }}`: The formatted invoice number.
- `{{ product_title }}`: The name of the service being discussed.
- `{{ company_name }}`: Your business name from Company Settings.

---

## 🧪 Testing Your Emails

Before going live, it is highly recommended to use the **Send Test Email** tool in the mail settings. This ensures your SMTP configuration is working and the HTML rendering of your templates is correct across different mail clients (Gmail, Outlook, etc.).

---

## 🌍 Multilingual Templates

If you have multiple languages enabled, you can provide translated versions of each email template. The system will automatically select the template version that matches the client's preferred language.
