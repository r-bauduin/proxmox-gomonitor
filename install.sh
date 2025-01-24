#!/bin/bash

# Fonction pour vérifier si une commande existe
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Installer Go si nécessaire
install_go() {
    echo "Vérification de l'installation de Go..."
    if command_exists go; then
        echo "Go est déjà installé."
    else
        echo "Installation de Go..."
        wget https://go.dev/dl/go1.20.5.linux-amd64.tar.gz -O /tmp/go.tar.gz
        tar -C /usr/local -xzf /tmp/go.tar.gz
        export PATH=$PATH:/usr/local/go/bin
        echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
        echo "Go a été installé."
    fi
}

# Demander les informations Proxmox à l'utilisateur
configure_proxmox() {
    default_url="https://localhost:8006/api2/json"
    echo "Configuration des identifiants Proxmox..."
    read -rp "Entrez l'URL de l'API Proxmox (par défaut : $default_url): " proxmox_url
    proxmox_url=${proxmox_url:-$default_url}

    read -rp "Entrez le Token API Proxmox (ex: PVEAPIToken=root@pam!api=<token>): " proxmox_token

    sed -i "s|baseURL  = .*|baseURL  = \"$proxmox_url\"|" main.go
    sed -i "s|apiToken = .*|apiToken = \"$proxmox_token\"|" main.go

    echo "Configuration terminée."
}

# Compiler le projet
compile_project() {
    echo "Compilation du projet..."
    go build -o proxmox-monitor main.go
    if [ $? -eq 0 ]; then
        echo "Le projet a été compilé avec succès. L'exécutable est disponible sous le nom 'proxmox-monitor'."
        mv proxmox-monitor /usr/local/bin/
    else
        echo "Erreur lors de la compilation du projet."
        exit 1
    fi
}

# Configurer le service systemd
setup_service() {
    echo "Création du service systemd..."
    bash -c 'cat > /etc/systemd/system/proxmox-monitor.service <<EOF
[Unit]
Description=Proxmox Monitor Service
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/proxmox-monitor
Restart=on-failure
User=root

[Install]
WantedBy=multi-user.target
EOF'

    # Recharger systemd et activer le service
    systemctl daemon-reload
    systemctl enable proxmox-monitor
    systemctl start proxmox-monitor

    echo "Service Proxmox Monitor configuré et démarré."
}

install_go
configure_proxmox
compile_project
setup_service

echo "Installation terminée. Le service 'proxmox-monitor' est actif."
