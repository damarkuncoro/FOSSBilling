# Server Managers & Provisioning

FOSSBilling automates the lifecycle of hosting accounts by connecting directly to remote server managers.

---

## 🖥️ Supported Server Managers

FOSSBilling includes high-performance drivers for the following platforms:

### WHM / cPanel and FOSSBilling
cPanel is the industry standard for shared hosting.
- **Connection:** Uses the WHM JSON API v1.
- **Requirements:** A valid WHM API Token with `create-acct`, `manage-accounts`, and `modify-accounts` permissions.
- **Automated Actions:** Account creation, suspension for non-payment, and termination.

### HestiaCP and FOSSBilling
A popular open-source control panel for VPS and dedicated servers.
- **Connection:** Uses the Hestia CLI-over-API gateway.
- **Requirements:** API Access enabled in Hestia settings and an Access Key/Secret pair.
- **Features:** Supports automated user creation and web domain assignment.

### CWP (CentOS Web Panel) and FOSSBilling
A powerful free/pro panel for RPM-based distributions.
- **Connection:** Uses the CWP REST API v1.
- **Requirements:** API Key generated from the CWP admin dashboard.
- **Port:** Default communication occurs over port `2304`.

### DirectAdmin and FOSSBilling
A lightweight and fast control panel alternative.
- **Connection:** Uses the `CMD_API` interface.
- **Authentication:** Supports both Login Keys and standard administrator credentials.

---

## 🔧 Other Server Managers

If your control panel is not listed above, FOSSBilling offers two flexible alternatives:

### 1. Plesk REST v2
Modern driver for Windows and Linux Plesk nodes using the latest JSON REST API for subscription management.

### 2. Custom Webhook Provisioner
For advanced users with custom infrastructure (Kubernetes, Proxmox, or bespoke scripts).
- **How it works:** FOSSBilling sends a POST request with a JSON payload containing order details to your endpoint.
- **Security:** Supports custom headers and HMAC signature verification.

---

## 🧪 Testing Connectivity

Before assigning a server to a product, use the **Test Connection** button in the server management view. This perform a real-time handshake with the remote API to ensure:
1. The server is reachable over the network.
2. The API credentials have the correct permissions.
3. The server manager is responding with the expected protocol version.
