[![codecov](https://codecov.io/gh/packethost/pkg/branch/master/graph/badge.svg?token=ErCO6uOE4T)](https://codecov.io/gh/packethost/pkg)

# pkg

Common functionality/helpers for our go services

## Grafana Dashboard

In order to make easier the default go - gin API's dashboard creation, there is a base JSON exported one in
`./grafana/base_dashboard.json`.

For install it go to `Import dashboard`, copy the content of the json file and replace `"title": "DEFINE_TITLE"` at the bottom.

Also change all the references for your service `service=\"packet-invoice-production-internal-http\"` (leaved that as an example for testing).
_You service should provide a `/metrics` endpoint for prometheus._

Select the folder, and hit import.
