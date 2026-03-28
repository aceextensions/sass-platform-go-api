package templates

const BaseEmailLayout = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body { font-family: 'Inter', -apple-system, sans-serif; background-color: #f8fafc; margin: 0; padding: 0; }
        .wrapper { width: 100%; background-color: #f8fafc; padding-bottom: 40px; }
        .container { max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 12px; overflow: hidden; margin-top: 40px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1); }
        .header { background-color: #059669; padding: 32px; text-align: center; }
        .header h1 { color: #ffffff; margin: 0; font-size: 24px; font-weight: 700; }
        .content { padding: 40px 32px; color: #334155; line-height: 1.6; }
        .footer { padding: 32px; text-align: center; font-size: 12px; color: #94a3b8; background-color: #f1f5f9; }
        .button { display: inline-block; background-color: #059669; color: #ffffff !important; padding: 12px 24px; border-radius: 8px; text-decoration: none; font-weight: 600; margin-top: 24px; }
    </style>
</head>
<body>
    <div class="wrapper">
        <div class="container">
            <div class="header"><h1>Practixa</h1></div>
            <div class="content">
                {{ content | safe }}
            </div>
            <div class="footer">
                <p>&copy; 2026 Practixa. All rights reserved.</p>
                <p>Ensuring your business moves antigravity.</p>
            </div>
        </div>
    </div>
</body>
</html>
`

// Wrap template in base layout
// Deprecated: We now handle wrapping in the notification service for better rendering.
func Wrap(tmpl string) string {
	return tmpl
}
