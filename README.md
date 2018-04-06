# HK Agent

A simple console program that monitors HTTP traffic on a machine.

## Features

- [x] Consumes an actively written-to w3c-formatted HTTP access log (Common Log Format)
- [x] Every 10s, displays in the console the sections of the web site with the most hits as well as interesting summary statistics on the traffic as a whole.
- [x] Whenever the total traffic for the past 2 minutes exceeds a certain number on average, displays an alert
- [x] Whenever the total traffic drops again below that value on average for the past 2 minutes, displays a message saying that it recovered
- [x] All messages showing when alerting thresholds are crossed remain visible on the page for historical reasons

## Testing

TODO: Write a test for the alerting logic

## Potential future improvements

See the [issues](https://github.com/Ullaakut/hk-agent/issues?q=is%3Aopen+is%3Aissue+milestone%3A%22Potential+future+improvements%22) and [projects](https://github.com/Ullaakut/hk-agent/projects/2) pages for a list of possible improvements that could be done in the near future.