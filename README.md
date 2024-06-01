# prom-dirsize-exporter

> A [Prometheus](https://prometheus.io/) exporter that exports the size of directories.

![Go version](https://img.shields.io/github/go-mod/go-version/brpaz/prom-dirsize-exporter?style=for-the-badge)
[![Latest Release](https://img.shields.io/github/v/release/brpaz/prom-dirsize-exporter?style=for-the-badge](https://github.com/brpaz/prom-dirsize-exporter/releases/latest)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/brpaz/prom-dirsize-exporter/CI?style=for-the-badge)](https://github.com/brpaz/prom-dirsize-exporter/actions/CI)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](LICENSE)

## About

This Prometheus exporter, allows to export directory sizes as Prometheus metrics.

This exporter will export metrics in the following format:

```
directory_size_bytes{path="/path/to/your/directory",name="directory"} <size_in_bytes>
```

## Usage

The recommended way to use this exporter is with Docker.

```shell
docker run ghcr.io/brpaz/prom-dirsize-exporter:latest
```

> [!IMPORTANT]
> When using Docker, ensure that directories you want to measure are mounted in the container as a volume.

### Configuration

The exporter can be configured using both command line flags or envrionment variables. Any command line flag will take precedence over envrionment variables.

Below you can find a list of supported configurations:

| Name | Flag | Envrionment variable | Default value | Description |
| Port | --metrics-port   | METRICS_PORT | 8080 | The port that the exporter listen to |
| Directories to monitor | --directories | DIRECTORIES | [] | A list of directory paths to monitor, separated by ":" |
| Metrics Path | --metrics-path | METRICS_PATH | /metrics | The path where the metrics are exposed |


## Built With

- [Cobra](https://cobra.dev/)

## Contributing

All contributions are welcome. Please check [Contributing guide](CONTRIBUTING.md) for instructions howe to contribute to this project.


## 🫶 Support

If you find this project helpful and would like to support its development, there are a few ways you can contribute:

[![Sponsor me on GitHub](https://img.shields.io/badge/Sponsor-%E2%9D%A4-%23db61a2.svg?&logo=github&logoColor=red&&style=for-the-badge&labelColor=white)](https://github.com/sponsors/brpaz)

<a href="https://www.buymeacoffee.com/Z1Bu6asGV" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: auto !important;width: auto !important;" ></a>


## Author

👤 **Bruno Paz**

- Website: [brunopaz.dev](https://brunopaz.dev)
- Github: [@brpaz](https://github.com/brpaz)


## 📝 License

Copyright [Bruno Paz](https://github.com/brpaz).

This project is [MIT](https://opensource.org/licenses/MIT) licensed.


