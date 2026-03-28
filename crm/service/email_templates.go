package service

// CRM specific templates
const (
	CustomerWelcomeTemplate = `
<h2 style="margin-top: 0; color: #1e293b; font-weight: 700;">Welcome to {{ tenant_name }}!</h2>
<p>Hello <strong>{{ name }}</strong>,</p>
<p>Thank you for choosing <strong>{{ tenant_name }}</strong>. We are excited to have you as our valued customer.</p>
<p>Your customer account has been successfully registered. You can now access our services and track your orders seamlessly.</p>
<p>If you have any questions or need assistance, please feel free to reach out to us.</p>
<p>Best regards,<br>The {{ tenant_name }} Team</p>
`
)
