# Helpdesk & Support Tickets

The integrated ticket helpdesk (`/support`) provides a unified channel for customer inquiries, technical troubleshooting, and billing requests.

---

## 🎧 Support Workflow

1. **Opening a Ticket:**
   - Select Department (General, Technical Support, Billing).
   - Set Priority (Low, Medium, High, Urgent).
   - Provide summary, message, and error logs.
2. **Conversation Threading:**
   - Reply directly to staff updates.
   - Attach logs or diagnostic screenshots.
3. **Automated Status Tracking:**
   - Tickets transition through `open`, `answered`, `waiting_on_client`, and `closed`.
   - Inactive resolved tickets are automatically closed by the background worker.
