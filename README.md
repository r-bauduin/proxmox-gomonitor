# Proxmox GoMonitor

## Introduction
**Proxmox GoMonitor** est une solution simple de monitoring pour Proxmox, développée en Go.

## Useful Links
- [Documentation sur la génération des clés API Proxmox](https://pve.proxmox.com/wiki/Proxmox_VE_API#API_Tokens)

---

## Getting Started

### Installation

Clonez le repository en exécutant la commande suivante :
```bash
git clone https://github.com/r-bauduin/proxmox-gomonitor.git
cd proxmox-gomonitor && chmod +x install.sh && ./install.sh
```

Pendant l'installation, il vous sera demandé :

- **L'URL d'accès à Proxmox** (par défaut : `https://localhost:8006/api2/json`).
- **Le PVEAPIToken** pour l'accès API.

---

## Service Management

Une fois le service installé et configuré, vous pouvez utiliser les commandes systemd pour gérer le processus :

### Démarrer le service `proxmox-monitor` :
```bash
systemctl start proxmox-monitor
```


### Stopper le service `proxmox-monitor` :
```bash
systemctl stop proxmox-monitor
```


### Redémarrer le service `proxmox-monitor` :
```bash
systemctl restart proxmox-monitor
```

### 💡 Vous pouvez vérifier les logs du service avec la commande suivante :
```bash
journalctl -u proxmox-monitor -f
```
