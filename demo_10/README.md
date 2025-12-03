### Panel for observability

![alt text](image-1.png)

![alt text](image-2.png)

![alt text](image-3.png)

![alt text](image-4.png)

### Alert to Discord
If latency (P99) of book requests exceeds 2 seconds, send an alert to Discord.

![alt text](image.png)

### Load test
Using vegeta to load testing the API and watch how the metrics change in Grafana.

```bash
vegeta attack -duration=30s -rate=100 -targets=targets.txt
```
