# Helm Chart Manifest Download Step 

helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm search bitnami
helm fetch bitnami/postgresql

# Prerequisites (Existing Clusters)

1. The velero CLI needs to be installed: 
    Option #1: For macOS, you can use Homebrew: 
    `brew install velero`
    Option #2A: Extract the github release tarball (linux): 
    ```s
    wget -O velero-cli.tar.gz https://github.com/vmware-tanzu/velero/releases/download/v1.3.1/velero-v1.3.1-linux-amd64.tar.gz
    tar -xvf velero-cli.tar.gz && rm velero-cli.tar.gz
    mv velero*linux-amd64/velero /usr/local/bin/velero && rm -r velero*linux-amd64
    ```
    Option #2B: Extract the github release tarball (macOS): 
    ```s
    wget -O velero-cli.tar.gz https://github.com/vmware-tanzu/velero/releases/download/v1.3.1/velero-v1.3.1-darwin-amd64.tar.gz
    tar -xvf velero-cli.tar.gz && rm velero-cli.tar.gz
    mv velero*darwin-amd64/velero /usr/local/bin/velero && rm -r velero*darwin-amd64
    ```

2. Configure your provider for velero, generating a secrets file. (Steps for GCP are below, but you can reference the following for other providers: https://github.com/vmware-tanzu/velero-plugin-for-gcp https://github.com/vmware-tanzu/velero-plugin-for-aws and https://github.com/vmware-tanzu/velero-plugin-for-microsoft-azure)
3. Ensure GCP is configured on your command line, and run the following script for convenience (step can be skipped if manual steps from above links were performed): 
```s
wget https://raw.githubusercontent.com/replicatedhq/replicated-automation/master/vendor/snapshot-feature/setup-gcp.sh && chmod +x setup-gcp.sh
./setup-gcp.sh | tee gcp-setup.log
BUCKET_NAME=$(cat gcp-setup.log | grep "BUCKET_NAME" | awk '{print $2}')
```
4. Add the velero plugin, along with the secret created in step 3. This plugin is included automatically for embedded cluster installs. Change the provider to GCP, AWS, or AZURE

AWS:   `velero install --use-restic --plugins velero/plugin-for-aws:v1.0.0 --no-default-backup-location --no-secret`
GCP:   `velero install --use-restic --plugins velero/plugin-for-gcp:v1.0.0 --no-default-backup-location --no-secret`
Azure: `velero install --use-restic --plugins velero/plugin-for-microsoft-azure:v1.0.0 --no-default-backup-location --no-secret`



# Example installation how-to
1. Ensure velero is installed for existing cluster installs (see above). It is automatically included in embedded installs. 
2. Install the kots application: 
    `curl -sSL https://k8s.kurl.sh/postgres-snapshots-austin | sudo bash` (embedded)
    `curl https://kots.io/install | bash && kubectl kots install postgres-snapshots/austin` (existing)
3. Enter the password and upload the license provided in this repository
4. Click "download postgres snapshots from the internet" if not in an airgapped environment. Otherwise, provide your airgap local registry credentials. 
5. Click the "Version History" tab and deploy the application
6. Click the "Snapshots" tab and select GCP for your provider 
8. Copy the generated key to the clipboard and paste in "Service Account" section in the admin console. 
9. Copy the BUCKET_NAME value from the log and paste in "Bucket" section in the admin console. 
10. Enter a desired path (/ is fine)
11. Click "Update Settings"
12. Go to the "Snapshots" tab and click "Start a snapshot" 
