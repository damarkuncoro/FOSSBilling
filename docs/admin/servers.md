# Server Provisioning & Hosting Managers

FOSSBilling integrates with industry-standard control panels to provision web hosting, email accounts, and database quotas automatically upon invoice payment.

---

## 🖥️ Supported Control Panels

| Control Panel | Driver ID | Authentication | Supported Actions |
| :--- | :--- | :--- | :--- |
| **cPanel / WHM** | `cpanel` | WHM API Token / Access Hash | Create, Suspend, Unsuspend, Terminate, Change Password |
| **Plesk** | `plesk` | REST API Key / Secret Key | Create Subscription, Suspend, Unsuspend, Delete, Password |
| **DirectAdmin** | `directadmin` | DirectAdmin Login Key / API | CMD_API_ACCOUNT_USER, Suspend, Unsuspend, Delete |
| **HestiaCP** | `hestia` | Hestia API Access Key | `v-add-user`, `v-add-web-domain`, `v-suspend-user`, `v-delete-user` |
| **CentOS Web Panel (CWP)** | `cwp` | CWP API Key | `/v1/account`, `/v1/account/suspend`, `/v1/account/delete` |
| **Custom Webhooks** | `custom` | HMAC-SHA256 Secret Header | JSON Webhooks for custom Kubernetes / Docker / VPS nodes |

---

## ⚙️ Adding a New Server

1. Navigate to **Servers** (`/servers`) in the Administrator Portal.
2. Click **Add New Server**.
3. Enter Server Details:
   - **Name:** e.g., `cPanel US-East-1 Cluster`
   - **Hostname / IP:** e.g., `whm.provider.com` (or IP address)
   - **Driver:** Choose from dropdown (`cpanel`, `plesk`, `directadmin`, `hestia`, `cwp`, `custom`)
   - **Port:** (Default: `2087` for WHM, `8443` for Plesk, `2222` for DirectAdmin, `8083` for HestiaCP)
   - **Access Key / API Token:** Paste your secure API token.
4. Click **Test Connection** to verify live communication with the server daemon.
5. Save the server and assign it to your hosting products.
