# Example installation how-to (Embedding Cluster with kURL) 

1. SSH into an Ubuntu 18 environment
2. Install the kots application, embedding a kubernetes cluster with Velero: 
```
curl -sSL https://k8s.kurl.sh/postgres-snapshots-austin | sudo bash
```
3. Follow the instructions on the command line to access the application with the provided password. 
4. Upload the license file provided in this repository
5. Click "Download postgres snapshots from the internet" if not in an airgapped environment. Otherwise, provide your airgap local registry credentials. 
6. Go to the "Snapshots" tab and click "Start a snapshot." This will create a snapshot and store in the existing cluster. See other guides for connecting to other backup locations.   
