# Localization & Internationalization (i18n)

FOSSBilling is designed for a global market, supporting multiple languages, currencies, and regional formatting out of the box.

---

## 🌍 Language Management

The system automatically detects the client's browser language and presents the appropriate translation if available.

### Supported Languages
FOSSBilling currently includes translations for:
- English (US/UK)
- Indonesian (Bahasa Indonesia)
- German (Deutsch)
- French (Français)
- Spanish (Español)
- Arabic (العربية) - Including full RTL support.

---

## 💰 Multi-Currency Billing

You can define multiple currencies in the Admin Portal under **Settings > Currencies**.

### Key Features
- **Default Currency:** The internal currency used for base calculations.
- **Conversion Rates:** Set manual rates or enable automatic updates via external providers.
- **Format Precision:** Configure how many decimal places to show for each currency (e.g., 0 for IDR, 2 for USD).

---

## 📅 Regional Formatting

FOSSBilling respects regional differences in data display:
- **Date Formats:** Choose between `YYYY-MM-DD`, `DD/MM/YYYY`, or `MM/DD/YYYY`.
- **Time Zones:** Set a global system timezone or allow clients to set their own.
- **Number Formats:** Correct handling of decimal and thousands separators based on locale.

---

## ✍️ Contributing Translations

Translation files are located in the `pkg/i18n/locales/` directory as JSON files. If you want to add a new language or improve an existing one:

1. Copy the `en_US.json` file.
2. Rename it to your locale code (e.g., `pt_BR.json`).
3. Translate the string values.
4. Submit a Pull Request to our repository.
