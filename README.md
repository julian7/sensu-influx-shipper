# Sensu-go Influx handler server

This is a simple TCP server to take [sensu-go handlers](https://docs.sensu.io/sensu-go/latest/reference/handlers/#tcp-udp-handlers) data in TCP, and forward them to an InfluxDB server.

The rationale is simple: even if you have a medium-sized setting, metrics will come in almost every second. Firing up a new app to handle each one can take significant amount of resources, which is very easy to save.

This application can easily replace [sensu-influxdb-handler](https://github.com/sensu/sensu-influxdb-handler) by using a slightly different configuration.

## Legal

This project is licensed under [Blue Oak Model License v1.0.0](https://blueoakcouncil.org/license/1.0.0). It is not registered either at OSI or GNU, therefore GitHub is widely looking at the other direction. However, this is the license I'm most happy with: you can read and understand it with no legal degree, and there are no hidden or cryptic meanings in it.

The project is also governed with [Contributor Covenant](https://contributor-covenant.org/)'s [Code of Conduct](https://www.contributor-covenant.org/version/1/4/) in mind. I'm not copying it here, as a pledge for taking the verbatim version by the word, and we are not going to modify it in any way.

## Any issues?

Open a ticket, perhaps a pull request. We support [GitHub Flow](https://guides.github.com/introduction/flow/). You might want to [fork](https://guides.github.com/activities/forking/) this project first.