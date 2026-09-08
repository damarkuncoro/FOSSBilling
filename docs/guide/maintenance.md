# Maintenance & Troubleshooting

Keeping your FOSSBilling instance healthy requires regular maintenance and the ability to resolve issues quickly.

---

## 🛠️ Routine Maintenance

- **Log Rotation:** Periodically clear the `activity_logs` and `audit_logs` if they become too large.
- **Database Backups:** Use `pg_dump` to create regular backups of your PostgreSQL database.
- **Cache Clearing:** If you notice configuration changes are not reflecting in the portals, use the **Clear Cache** button in the Admin Dashboard.

---

## 🔐 Manual Admin Password Reset

If you have lost access to your primary administrator account and cannot use the "Forgot Password" feature, you can reset it manually via the CLI or Database.

### Option 1: Via CLI (Recommended)
Use the built-in management utility from your server terminal:
```bash
cd backend-go
./bin/cli staff reset-password --email admin@fossbilling.org --password NewSecurePass123!
```

### Option 2: Via SQL
If the CLI is unavailable, run this query against your PostgreSQL database:
```sql
-- Replace with your new bcrypt hash (cost 12)
UPDATE staff 
SET password_hash = '$2a$12$...' 
WHERE email = 'admin@fossbilling.org';
```

---

## 🔍 Troubleshooting Issues

### 1. Portals not connecting to API
- Ensure the `APP_URL` in your configuration matches the URL you are using to access the API.
- Check browser console for CORS errors. Ensure your API server has the correct allowed origins.

### 2. Provisioning Failures
- Verify the server manager credentials (API Token/Key).
- Check the **Worker Logs** using Docker:
  ```bash
  docker logs fossbilling-worker
  ```
- Use the **Test Connection** button in the Admin Server Manager to verify network reachability.

### 3. Invoices not generating
- Ensure the system **Cron Job** is running.
- FOSSBilling generates invoices 14 days before the due date by default.

---

## 🚩 Error Reporting

If you encounter a bug or a system crash, please gather the following information before opening an issue:

1. **Environment Details:** OS, Docker version, and FOSSBilling version.
2. **API Logs:** Run `docker logs fossbilling-api` to see the stack trace.
3. **Frontend Console:** Open Developer Tools (F12) in your browser and check the "Console" and "Network" tabs for failed requests.
4. **Steps to Reproduce:** A clear list of actions that lead to the error.

Report bugs at: [https://github.com/damarkuncoro/FOSSBilling/issues](https://github.com/damarkuncoro/FOSSBilling/issues)
