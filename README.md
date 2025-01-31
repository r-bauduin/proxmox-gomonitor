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

## Utilisation 
```bash
curl http://{IP}:59720/status
```

```json
[{"cpu_usage":0.106772891134822,"lxc_ratio_up":0,"lxc_running":0,"lxc_stopped":0,"lxc_total":0,"node":"lame-68","ram_usage":60.367929748335044,"vms_ratio_up":50,"vms_running":10,"vms_stopped":10,"vms_total":20,"system_disk_usage":70.34184645397838}]
```

Le service gère également les serveurs Proxmox avec plusieurs noeuds :
```json
[{"cpu_usage":0.190528034839412,"lxc_ratio_up":75,"lxc_running":3,"lxc_stopped":1,"lxc_total":4,"node":"sql-2","ram_usage":17.19418133700936,"vms_ratio_up":100,"vms_running":1,"vms_stopped":0,"vms_total":1,"system_disk_usage":70.34184645397838},{"cpu_usage":0.258928091458061,"lxc_ratio_up":80,"lxc_running":4,"lxc_stopped":1,"lxc_total":5,"node":"sql-1","ram_usage":48.70376346508337,"vms_ratio_up":80,"vms_running":8,"vms_stopped":2,"vms_total":10,"system_disk_usage":80}]
```
