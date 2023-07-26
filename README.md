# diskexporter

A simple prometheus exporter that reports the disk usage of any directory on your machine.

## Usage:

```
docker run -d -p 1971:1971 -v /path/to/monitor:/mypath -e MONITORED_PATH=/mypath diskexporter:latest
```

Your metrics will be available at http://localhost:1971/
