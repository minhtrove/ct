# CT Finance - Deployment Guide

## Prerequisites

### 1. EC2 Instance Setup
Your EC2 instance should have:
- **OS**: Amazon Linux 2023 (ARM64)
- **Security Group**: Allow inbound on port 3000 (or 80/443 if using nginx)
- **IAM Role**: Attached IAM role with SES permissions

### 2. IAM Role for EC2 (SES Access)
Create an IAM role with this policy:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ses:SendEmail",
        "ses:SendRawEmail"
      ],
      "Resource": "*"
    }
  ]
}
```

Attach this role to your EC2 instance.

### 3. GitHub Secrets Configuration

Add these secrets to your GitHub repository (Settings → Secrets and variables → Actions):

| Secret Name | Description | Example |
|------------|-------------|---------|
| `AWS_ACCESS_KEY_ID` | AWS access key for GitHub Actions | `AKIAIOSFODNN7EXAMPLE` |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key | `wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY` |
| `AWS_REGION` | AWS region | `us-east-1` |
| `AWS_SES_REGION` | SES region | `us-east-1` |
| `AWS_SES_FROM_EMAIL` | Verified SES email | `noreply@yourdomain.com` |
| `AWS_SES_FROM_NAME` | Email sender name | `CT Finance` |
| `EC2_HOST` | EC2 public IP or domain | `ec2-xx-xx-xx-xx.compute.amazonaws.com` |
| `EC2_USER` | SSH user | `ec2-user` |
| `SSH_PRIVATE_KEY` | EC2 SSH private key | Contents of your `.pem` file |
| `MONGODB_URI` | MongoDB Atlas connection string | `mongodb+srv://user:pass@cluster.mongodb.net/db` |
| `APP_NAME` | Application name | `CT Finance` |
| `APP_BASE_URL` | Application URL | `https://yourdomain.com` |

## Initial EC2 Setup

SSH into your EC2 instance and run:

```bash
# Update system
sudo yum update -y

# Create app directory
mkdir -p ~/app
cd ~/app

# Install required packages (if needed)
sudo yum install -y git

# Set up log directory
sudo mkdir -p /var/log/ct-finance
sudo chown ec2-user:ec2-user /var/log/ct-finance
```

## MongoDB Atlas Setup

1. Create a MongoDB Atlas cluster (free tier available)
2. Create a database user
3. Whitelist your EC2 IP or use 0.0.0.0/0 (less secure)
4. Get the connection string
5. Add to GitHub Secrets as `MONGODB_URI`

## AWS SES Setup

1. Verify your sender email address in AWS SES
2. If in sandbox mode, verify recipient emails too
3. Request production access if needed
4. Ensure your IAM role has SES permissions

## Deployment Process

### Manual First Deployment

1. **Push to main branch** or trigger manually:
   ```bash
   git push origin main
   ```

2. **Monitor GitHub Actions**:
   - Go to GitHub → Actions tab
   - Watch the deployment progress

3. **Verify deployment**:
   ```bash
   ssh ec2-user@your-ec2-host
   sudo systemctl status ct-finance
   sudo journalctl -u ct-finance -f
   ```

### Automatic Deployments

After initial setup, every push to `main` branch will automatically:
1. Build the Go binary for ARM64
2. Generate templ files
3. Package the application
4. Deploy to EC2
5. Restart the service
6. Run health check

## Nginx Setup (Optional)

For production, use Nginx as reverse proxy:

```nginx
server {
    listen 80;
    server_name yourdomain.com;

    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

For HTTPS, use Let's Encrypt:
```bash
sudo yum install -y certbot python3-certbot-nginx
sudo certbot --nginx -d yourdomain.com
```

## Troubleshooting

### Check service status
```bash
sudo systemctl status ct-finance
```

### View logs
```bash
sudo journalctl -u ct-finance -f
```

### Restart service
```bash
sudo systemctl restart ct-finance
```

### Check if app is listening
```bash
curl http://localhost:3000/api/health
```

### View environment variables
```bash
cat ~/app/.env
```

## Rollback

If deployment fails, the script automatically rolls back. To manually rollback:

```bash
cd ~/app
sudo systemctl stop ct-finance
mv ct-finance.old ct-finance
sudo systemctl start ct-finance
```

## Security Best Practices

1. ✅ Use IAM roles instead of access keys where possible
2. ✅ Keep SSH keys secure
3. ✅ Use HTTPS (Let's Encrypt)
4. ✅ Regularly update system packages
5. ✅ Use security groups to restrict access
6. ✅ Enable CloudWatch monitoring
7. ✅ Rotate credentials regularly
8. ✅ Use environment variables for secrets

## Monitoring

### CloudWatch Logs (Optional)
Install CloudWatch agent:

```bash
sudo yum install -y amazon-cloudwatch-agent
```

### Application Metrics
Access at: http://your-domain/api/health

## Support

For issues:
1. Check GitHub Actions logs
2. Check EC2 application logs: `sudo journalctl -u ct-finance`
3. Verify environment variables
4. Check network connectivity to MongoDB Atlas
5. Verify AWS SES permissions
