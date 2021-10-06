
<details>
  <summary>Open for a hint</summary>

```shell
$ sudo touch /etc/ring-game/config.txt
$ sudo chmod 400 /etc/ring-game/config.txt
```
</details>
<details>
  <summary>Open for a hint</summary>

```shell
$ sudo touch /etc/ring-game/restraining-bolt.txt
$ sudo chmod 400 /etc/ring-game/restraining-bolt.txt

```
</details>

<details>
  <summary>Open for a hint</summary>

```yaml

apiVersion: v1
kind: Service
metadata:
  creationTimestamp: null
  labels:
    kots.io/backup: velero
    kots.io/kotsadm: "true"
  name: kotsadm
  namespace: default
spec:
  ports:
  - name: http
    port: 3000
    targetPort: http
  selector:
    app: kotsadm
  type: ClusterIP
status:
  loadBalancer: {}

```

</details>