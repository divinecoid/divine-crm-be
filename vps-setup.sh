#!/bin/bash

# ===========================================
# Divine CRM - VPS Deployment Script
# ===========================================
# VPS: 156.67.219.209
# User: divinecrm
# ===========================================

set -e

echo "🚀 Divine CRM Deployment Script"
echo "================================"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if running as root or with sudo
if [ "$EUID" -ne 0 ]; then 
    echo -e "${YELLOW}⚠️  Please run with sudo${NC}"
    exit 1
fi

# ===========================================
# STEP 1: Update System
# ===========================================
echo -e "\n${GREEN}[1/8] Updating system...${NC}"
apt update && apt upgrade -y

# ===========================================
# STEP 2: Install Dependencies
# ===========================================
echo -e "\n${GREEN}[2/8] Installing dependencies...${NC}"
apt install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    lsb-release \
    git \
    ufw \
    nginx \
    certbot \
    python3-certbot-nginx

# ===========================================
# STEP 3: Install Docker
# ===========================================
echo -e "\n${GREEN}[3/8] Installing Docker...${NC}"

# Remove old versions
apt remove -y docker docker-engine docker.io containerd runc 2>/dev/null || true

# Add Docker's official GPG key
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# Set up repository
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker Engine
apt update
apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# Add user to docker group
usermod -aG docker divinecrm

# Start Docker
systemctl enable docker
systemctl start docker

echo -e "${GREEN}✅ Docker installed: $(docker --version)${NC}"

# ===========================================
# STEP 4: Install Docker Compose
# ===========================================
echo -e "\n${GREEN}[4/8] Installing Docker Compose...${NC}"

DOCKER_COMPOSE_VERSION="v2.24.0"
curl -L "https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

echo -e "${GREEN}✅ Docker Compose installed: $(docker-compose --version)${NC}"

# ===========================================
# STEP 5: Setup Firewall
# ===========================================
echo -e "\n${GREEN}[5/8] Configuring firewall...${NC}"

ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow 80/tcp    # HTTP
ufw allow 443/tcp   # HTTPS
ufw allow 3000/tcp  # Frontend (dev)
ufw allow 3001/tcp  # Auth Service
ufw allow 3002/tcp  # CRM Service

echo "y" | ufw enable
ufw status

echo -e "${GREEN}✅ Firewall configured${NC}"

# ===========================================
# STEP 6: Create Project Directory
# ===========================================
echo -e "\n${GREEN}[6/8] Creating project directory...${NC}"

PROJECT_DIR="/home/divinecrm/divine-crm"
mkdir -p $PROJECT_DIR
chown -R divinecrm:divinecrm $PROJECT_DIR

echo -e "${GREEN}✅ Project directory created: $PROJECT_DIR${NC}"

# ===========================================
# STEP 7: Setup Nginx Reverse Proxy
# ===========================================
echo -e "\n${GREEN}[7/8] Configuring Nginx...${NC}"

# Create Nginx config for Divine CRM
cat > /etc/nginx/sites-available/divine-crm << 'NGINX_CONFIG'
# Divine CRM - Nginx Configuration

# Frontend
server {
    listen 80;
    server_name crm.yourdomain.com;  # Change this to your domain

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}

# Auth Service API
server {
    listen 80;
    server_name auth-api.yourdomain.com;  # Change this to your domain

    location / {
        proxy_pass http://localhost:3001;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# CRM Service API
server {
    listen 80;
    server_name api.yourdomain.com;  # Change this to your domain

    location / {
        proxy_pass http://localhost:3002;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # WebSocket support for webhooks
    location /api/v1/webhooks {
        proxy_pass http://localhost:3002;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
NGINX_CONFIG

# Enable site
ln -sf /etc/nginx/sites-available/divine-crm /etc/nginx/sites-enabled/
rm -f /etc/nginx/sites-enabled/default

# Test and reload Nginx
nginx -t && systemctl reload nginx

echo -e "${GREEN}✅ Nginx configured${NC}"

# ===========================================
# STEP 8: Create Helper Scripts
# ===========================================
echo -e "\n${GREEN}[8/8] Creating helper scripts...${NC}"

# Start script
cat > $PROJECT_DIR/start.sh << 'EOF'
#!/bin/bash
cd /home/divinecrm/divine-crm/divine-crm-be
docker-compose up -d
echo "✅ Divine CRM Backend started"

cd /home/divinecrm/divine-crm/divine-crm-fe
npm run build
pm2 start npm --name "divine-crm-fe" -- start
echo "✅ Divine CRM Frontend started"
EOF

# Stop script
cat > $PROJECT_DIR/stop.sh << 'EOF'
#!/bin/bash
cd /home/divinecrm/divine-crm/divine-crm-be
docker-compose down
echo "✅ Divine CRM Backend stopped"

pm2 stop divine-crm-fe
echo "✅ Divine CRM Frontend stopped"
EOF

# Logs script
cat > $PROJECT_DIR/logs.sh << 'EOF'
#!/bin/bash
echo "=== Auth Service Logs ==="
docker logs --tail 50 divine-auth-service

echo ""
echo "=== CRM Service Logs ==="
docker logs --tail 50 divine-crm-service

echo ""
echo "=== Frontend Logs ==="
pm2 logs divine-crm-fe --lines 50
EOF

# Update script
cat > $PROJECT_DIR/update.sh << 'EOF'
#!/bin/bash
cd /home/divinecrm/divine-crm

echo "📥 Pulling latest changes..."
git pull

echo "🔄 Rebuilding backend..."
cd divine-crm-be
docker-compose down
docker-compose up --build -d

echo "🔄 Rebuilding frontend..."
cd ../divine-crm-fe
npm install
npm run build
pm2 restart divine-crm-fe

echo "✅ Update complete!"
EOF

chmod +x $PROJECT_DIR/*.sh
chown -R divinecrm:divinecrm $PROJECT_DIR

echo -e "${GREEN}✅ Helper scripts created${NC}"

# ===========================================
# DONE
# ===========================================
echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}✅ VPS Setup Complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "Next steps:"
echo "1. Upload your project to: /home/divinecrm/divine-crm/"
echo "2. Configure .env files"
echo "3. Update Nginx with your domain"
echo "4. Run: cd /home/divinecrm/divine-crm && ./start.sh"
echo ""
echo -e "${YELLOW}For SSL (HTTPS):${NC}"
echo "certbot --nginx -d crm.yourdomain.com -d api.yourdomain.com -d auth-api.yourdomain.com"
echo ""