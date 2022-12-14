# OwnTracks Recorder

Docker based deployment

```bash
# persistance data directory for sqlite
mkdir data
sudo chown -R 1000:1000 data

# run container in background
docker-compose up -d

# run migrations
docker-compose exec owntracks-recorder /app/ot-recorder migrate up
```
