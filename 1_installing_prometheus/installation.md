1. Install prometheus
docker run \
    --name prometheus \
    -p 9090:9090 \
    -v ./prometheus.yml:/etc/prometheus/prometheus.yml \
    prom/prometheus

2. convert to docker-compose