# Sensu-go Influx shipper

This is a simple TCP server to take [sensu-go handlers](https://docs.sensu.io/sensu-go/latest/reference/handlers/#tcp-udp-handlers) data in TCP, and forward them to an InfluxDB server.

The rationale is simple: even if you have a medium-sized setting, metrics will come in almost every second. Firing up a new app to handle each one can take significant amount of resources, which is very easy to save.

This application can easily replace [sensu-influxdb-handler](https://github.com/sensu/sensu-influxdb-handler) by using a slightly different configuration.

## Differences to sensu-influxdb-handler

The original handler is a simple, raw in to raw out data interface adapter. However, when the metric collector is outputting in graphite format, no tags will be saved. We try to set the following default tags to each data point, as defaults:

* host: entity's name
* ip: entity's first IP address of the first physical interface
* check: check's name

The application is also able to group metric data. By default, it splits metric names by the first dot, using the first part for metric name, and use the rest as metric key. If there are no dots in the name, metric key will use "value".

When grouping is enabled, metric key will take the value of the original name's last piece prefixed by an underscore, and the middle of the key will go to the "metric" tag. This setting can be very handy for graphite_plaintext input, displayed on Grafana.

Examples:

| original key  | influx name | influx key (grouped: false) | influx key (grouped: true) | extra tag (grouped: true) |
| :------------ | :---------- | :-------------------------- | :------------------------- | :------------------------ |
| disk          | disk        | value                       | value                      | (none)                    |
| disk.free     | disk        | free                        | value                      | metric=free               |
| disk.var.free | disk        | var.free                    | _free                      | metric=var                |

## Compile from source

Use [mage](https://magefile.org/) to build: `mage all`. This will do a full build into `dist/` subfolder. To do a build with publishing results, run `mage publish`.

## Usage

```text
Sensu-go event consuming TCP server for InfluxDB data shipping

Usage:
  sensu-influx-shipper serve [flags]

Flags:
  -a, --addr string       InfluxDB's TCP port (default "http://127.0.0.1:8086")
  -d, --database string   InfluxDB database (default "metrics")
  -g, --grouping          Group metric data
  -h, --help              help for serve
  -l, --listen string     TCP port to listen to (default "127.0.0.1:3333")
  -p, --pass string       InfluxDB password
  -u, --user string       InfluxDB username (default "metrics")

Global Flags:
      --config string      configuration file (default: /etc/sensu-influx-shipper[.yml]
  -L, --logfile string     log file. Possible values: none, stdout, stderr, or file name (default "stderr")
  -F, --logformat string   log format. Possible values: logfmt, or json (default "logfmt")
```

This command runs a TCP server on `listen` port (port number or on a specific interface in `ip:port` format), and accepts from [sensu go events](https://docs.sensu.io/sensu-go/5.15/reference/events/). Then it ships metric data found in events to an InfluxDB service (see `addr`, `database`, `user`, `pass` options).

## Configuration

The application accepts a configuration file too, of the following structure (the example is YAML-formatted, but other formats are supported by [cobra](https://github.com/spf13/cobra)):

```yaml
---
logfile: "log file"
logformat: "log format"
serve:
    addr: "InfluxDB server URL"
    database: "InfluxDB database name"
    grouping: true
    listen: "local TCP server listen port"
    pass: "InfluxDB password"
    user: "InfluxDB user"
```

Alternatively, configuration can come from environment variables, like `LOGFORMAT` or `SERVE_LISTEN`.

## Legal

This project is licensed under [Blue Oak Model License v1.0.0](https://blueoakcouncil.org/license/1.0.0). It is not registered either at OSI or GNU, therefore GitHub is widely looking at the other direction. However, this is the license I'm most happy with: you can read and understand it with no legal degree, and there are no hidden or cryptic meanings in it.

The project is also governed with [Contributor Covenant](https://contributor-covenant.org/)'s [Code of Conduct](https://www.contributor-covenant.org/version/1/4/) in mind. I'm not copying it here, as a pledge for taking the verbatim version by the word, and we are not going to modify it in any way.

## Any issues?

Open a ticket, perhaps a pull request. We support [GitHub Flow](https://guides.github.com/introduction/flow/). You might want to [fork](https://guides.github.com/activities/forking/) this project first.
