package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"crypto/tls"
)

const (
	baseURL  = "https://localhost:8006/api2/json"
	apiToken = "PVEAPIToken=root@pam!api=d8e5b982-4a93-4731-8780-44d8d54fc9c3"
)

type NodeStatus struct {
	CPU  float64 `json:"cpu"`
	RAM  struct {
		Used  int64 `json:"used"`
		Total int64 `json:"total"`
	} `json:"memory"`
}

type VM struct {
	VMID   string    `json:"vmid"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Node struct {
	Node string `json:"node"`
}

type APIResponse struct {
	Data json.RawMessage `json:"data"`
}

func fetchProxmoxData(endpoint string, target interface{}) error {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	url := baseURL + endpoint
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("erreur lors de la création de la requête : %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiToken)

	log.Printf("Envoi de la requête : %s", url)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erreur lors de l'exécution de la requête : %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Réponse brute : %s", string(body))
		return fmt.Errorf("erreur API : %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("erreur lors de la lecture de la réponse : %v", err)
	}

	log.Printf("Réponse brute : %s", string(body))

	var response APIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("erreur lors du parsing JSON : %v", err)
	}

	return json.Unmarshal(response.Data, target)
}

func getNodes() ([]Node, error) {
	var nodes []Node
	if err := fetchProxmoxData("/nodes", &nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}

func getNodeStatus(node string) (*NodeStatus, error) {
	var status NodeStatus
	if err := fetchProxmoxData("/nodes/"+node+"/status", &status); err != nil {
		return nil, err
	}
	return &status, nil
}

func getVMs(node string) ([]VM, error) {
	var vms []VM
	if err := fetchProxmoxData("/nodes/"+node+"/qemu", &vms); err != nil {
		return nil, err
	}
	return vms, nil
}

func getLXC(node string) ([]VM, error) {
    var lxc []VM
    if err := fetchProxmoxData("/nodes/"+node+"/lxc", &lxc); err != nil {
        return nil, err
    }
    return lxc, nil
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	nodes, err := getNodes()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des nœuds", http.StatusInternalServerError)
		log.Printf("Erreur lors de la récupération des nœuds : %v", err)
		return
	}

	if len(nodes) == 0 {
		// Aucun nœud trouvé
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	metrics := []map[string]interface{}{}

	for _, node := range nodes {
		status, err := getNodeStatus(node.Node)
		if err != nil {
			log.Printf("Erreur lors de la récupération du statut pour le nœud %s : %v", node.Node, err)
			continue
		}

		vms, err := getVMs(node.Node)
		if err != nil {
			log.Printf("Erreur lors de la récupération des VMs pour le nœud %s : %v", node.Node, err)
			continue
		}

		vmsRunning := 0
		vmsTotal := len(vms)
		for _, vm := range vms {
			if vm.Status == "running" {
				vmsRunning++
			}
		}

		lxc, err := getLXC(node.Node)
		if err != nil {
		    log.Printf("Erreur lors de la récupération des LXC pour le nœud %s : %v", node.Node, err)
		    continue
        }

        lxcRunning := 0
        lxcTotal := len(lxc)
        for _, lxc := range lxc {
            if lxc.Status == "running" {
                lxcRunning++
            }
        }

		metrics = append(metrics, map[string]interface{}{
			"node":        node.Node,
			"cpu_usage":   status.CPU * 100,
			"ram_usage":   float64(status.RAM.Used) / float64(status.RAM.Total) * 100,
			"vms_running": vmsRunning,
			"vms_stopped": vmsTotal - vmsRunning,
			"vms_total":   vmsTotal,
            "vms_ratio_up": func() float64 {
                if vmsTotal > 0 {
                    return (float64(vmsRunning) / float64(vmsTotal)) * 100
                }
                return 0
            }(),
    		"lxc_running": lxcRunning,
			"lxc_stopped": lxcTotal - lxcRunning,
			"lxc_total":   lxcTotal,
			"lxc_ratio_up": func() float64 {
			    if lxcTotal > 0 {
                    return (float64(lxcRunning) / float64(lxcTotal)) * 100
                }
                return 0
            }
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		http.Error(w, "Erreur lors de l'encodage JSON", http.StatusInternalServerError)
		log.Printf("Erreur lors de l'encodage JSON : %v", err)
	}
}

func main() {
	http.HandleFunc("/status", statusHandler)

	fmt.Println("Service de supervision Proxmox démarré sur le port 59720...")
	log.Fatal(http.ListenAndServe(":59720", nil))
}
