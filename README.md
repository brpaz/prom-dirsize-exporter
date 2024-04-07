# prom-dirsize-exporter

> A [Prometheus](https://prometheus.io/) exporter to generate metrics regarding directory sizes

![Go version](https://img.shields.io/github/go-mod/go-version/brpaz/prom-dirsize-exporter?style=for-the-badge)
[![Latest Release](https://img.shields.io/github/v/release/brpaz/prom-dirsize-exporter?style=for-the-badge](https://github.com/brpaz/prom-dirsize-exporter/releases/latest)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/brpaz/prom-dirsize-exporter/CI?style=for-the-badge)](https://github.com/brpaz/prom-dirsize-exporter/actions/CI)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](LICENSE)

## About

Small prometheus exporter than generates metrics regarding directory sizes.

This exporter will export metrics like this:

```
directory_size_bytes{path="/path/to/your/directory",name="My Directory"} <size_in_bytes>
```

## Usage

The recommended way to use this exporter is with Docker.

```shell
docker run ghcr.io/brpaz/prom-dirsize-exporter-latest
```

> [!IMPORTANT]
> Please ensure the directories you want to measure are mounted in the container as a volume.

### Configuration

The exporter can be configured using both command line flags and envrionment variables. Any defined envrionment variable will take precedence
over command line flags.

Below you can find a list of supported configurations:

| Name | Flag | Envrionment variable | Default value | Description |
| Port | --port   | PORT | 8080 | The port that the exporter listen to |
| Directories to monitor | --directories | DIRECTORIES | [] | A list of directory paths to monitor, separated by ";" |
| Metrics Path | --metrics-path | METRICS_PATH | /metrics | The path where the metrics are exposed |


## Built With

- [Cobra](https://cobra.dev/)

## Contributing

All contributions are welcome. Please check [Contributing guide](CONTRIBUTING.md) for instructions howe to contribute to this project.

## Author

👤 **Bruno Paz**

- Website: [brunopaz.dev](https://brunopaz.dev)
- Github: [@brpaz](https://github.com/brpaz)


## 📝 License

Copyright [Bruno Paz](https://github.com/brpaz).

This project is [MIT](https://opensource.org/licenses/MIT) licensed.


