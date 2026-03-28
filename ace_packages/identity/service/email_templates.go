package service

const (
	InvitationTemplate = `
<h2 style="margin-top: 0; color: #1e293b; font-weight: 700;">You've been invited!</h2>
<p>Hello,</p>
<p><strong>{{ inviter_name }}</strong> has invited you to join the <strong>{{ tenant_name }}</strong> workspace on <strong>Practixa</strong> as a <strong>{{ role }}</strong>.</p>
<p>Practixa helps your team move with antigravity—managing products, categories, and accounting in one seamless interface.</p>
<p style="text-align: center;">
    <a href="{{ invite_link }}" class="button">Accept Invitation</a>
</p>
<p style="margin-top: 32px; font-size: 14px; color: #64748b;">
    If the button above doesn't work, copy and paste this link into your browser:<br>
    <span style="word-break: break-all; color: #059669;">{{ invite_link }}</span>
</p>
<p style="margin-top: 24px;">Welcome to the future of workspace management!</p>
<p>Best regards,<br>The Practixa Team</p>
`

	WelcomeTemplate = `
<h2 style="margin-top: 0; color: #1e293b; font-weight: 700;">Welcome to Practixa!</h2>
<p>Hello <strong>{{ name }}</strong>,</p>
<p>We're thrilled to have you on board. Your account has been successfully created and you're ready to start managing your workspace with antigravity.</p>
<p>Explore your dashboard, set up your products, and see how easy accounting can be when everything is in one place.</p>
<p style="text-align: center;">
    <a href="{{ dashboard_link }}" class="button">Go to Dashboard</a>
</p>
<p>If you have any questions, our support team is always here to help.</p>
<p>Best regards,<br>The Practixa Team</p>
`

	PasswordResetTemplate = `
<h2 style="margin-top: 0; color: #1e293b; font-weight: 700;">Password Reset Request</h2>
<p>Hello,</p>
<p>We received a request to reset the password for your Practixa account. If you didn't make this request, you can safely ignore this email.</p>
<p>To reset your password, click the button below:</p>
<p style="text-align: center;">
    <a href="{{ reset_link }}" class="button">Reset Password</a>
</p>
<p style="margin-top: 32px; font-size: 14px; color: #64748b;">
    This link will expire in 24 hours for security reasons.
</p>
<p>Best regards,<br>The Practixa Team</p>
`

	PasswordChangedTemplate = `
<h2 style="margin-top: 0; color: #1e293b; font-weight: 700;">Password Changed Successfully</h2>
<p>Hello,</p>
<p>This is a confirmation that the password for your Practixa account has been successfully changed.</p>
<p>If you did not make this change, please contact our support team immediately or reset your password to secure your account.</p>
<p style="text-align: center;">
    <a href="{{ dashboard_link }}" class="button">Go to Dashboard</a>
</p>
<p>Best regards,<br>The Practixa Team</p>
`
)
