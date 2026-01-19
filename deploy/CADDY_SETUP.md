# Quick Setup for IP-Based Access (No Domain)

## Setup Steps

### 1. Copy Caddyfile to EC2

On your EC2 instance:
```bash
sudo vi /etc/caddy/Caddyfile
```

Paste the content from `deploy/Caddyfile` (the simplified IP version).

### 2. Create log directory
```bash
sudo mkdir -p /var/log/caddy
sudo chown caddy:caddy /var/log/caddy
```

### 3. Restart Caddy
```bash
sudo systemctl restart caddy
sudo systemctl enable caddy
sudo systemctl status caddy
```

### 4. Update EC2 Security Group

Make sure these ports are open:
- ✅ **Port 80** (HTTP) - for web access
- ✅ **Port 22** (SSH) - for your management
- ❌ **Port 3000** - keep CLOSED (only accessible from localhost)

### 5. Access Your App

```bash
# From your browser or terminal
http://YOUR_EC2_PUBLIC_IP
```

## Finding Your EC2 Public IP

```bash
# On EC2 instance
curl http://169.254.169.254/latest/meta-data/public-ipv4

# Or check AWS Console
```

## Testing

```bash
# Test from your local machine
curl -I http://YOUR_EC2_PUBLIC_IP

# Should return HTTP 200 OK
```

## Important Notes

⚠️ **No HTTPS**: Without a domain, you can't get automatic SSL certificates from Let's Encrypt. Your traffic will be unencrypted (HTTP only).

For production, you should:
1. Get a domain name (cheap options: Namecheap, Google Domains, Cloudflare)
2. Point it to your EC2 IP
3. Update Caddyfile to use the domain
4. Caddy will automatically get HTTPS certificates

## When You Get a Domain

Simply replace the Caddyfile content with:

```caddy
yourdomain.com {
    reverse_proxy localhost:3000
}
```

That's it! Caddy handles everything else automatically (HTTPS, certificates, renewal).

## Troubleshooting

**Can't access via IP?**
```bash
# Check if Caddy is running
sudo systemctl status caddy

# Check if Go app is running
sudo systemctl status ct-finance

# Check if port 80 is listening
sudo netstat -tlnp | grep :80

# Check Go app directly
curl http://localhost:3000/api/health
```

**Caddy won't start?**
```bash
# Check Caddy logs
sudo journalctl -u caddy -n 50

# Validate Caddyfile syntax
caddy validate --config /etc/caddy/Caddyfile
```

**Still can't access?**
- Check AWS Security Group allows port 80 from 0.0.0.0/0
- Check if there's a firewall on EC2: `sudo iptables -L`
- Verify EC2 public IP hasn't changed
