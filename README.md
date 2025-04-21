# Calendar 🗓️

[![License](https://img.shields.io/github/license/worldline-go/calendar?color=blue&style=flat-square)](https://raw.githubusercontent.com/worldline-go/calendar/main/LICENSE)
[![Coverage](https://img.shields.io/sonar/coverage/worldline-go_calendar?logo=sonarcloud&server=https%3A%2F%2Fsonarcloud.io&style=flat-square)](https://sonarcloud.io/summary/overall?id=worldline-go_calendar)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/worldline-go/calendar/test.yml?branch=main&logo=github&style=flat-square&label=ci)](https://github.com/worldline-go/calendar/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/worldline-go/calendar?style=flat-square)](https://goreportcard.com/report/github.com/worldline-go/calendar)
[![Web](https://img.shields.io/badge/web-document-blueviolet?style=flat-square)](https://worldline-go.github.io/calendar/)

This service that provides information about holidays in a given country or special code.

## Development

Use `make` to show help and create env, run tests, etc.

```sh
# Run compose-file to open postgresql in local
make env
# Start the service, it default reads `calendar.[yaml|yml|json|toml]` or use `CONFIG_FILE` env value for file path.
make run
```

## References

- https://datatracker.ietf.org/doc/html/rfc5545
- https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
- https://en.wikipedia.org/wiki/ISO_3166-1_alpha-3
- https://www.thunderbird.net/en-US/calendar/holidays/
- https://calendars.icloud.com/holidays/tr_tr.ics/
- https://github.com/ics-tools/viewer
