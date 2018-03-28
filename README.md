# HK Agent

A simple console program that monitors HTTP traffic on a machine.

## Features

* Consumes an actively written-to w3c-formatted HTTP access log (https://en.wikipedia.org/wiki/Common_Log_Format)
* Every 10s, displays in the console the sections of the web site with the most hits as well as interesting summary statistics on the traffic as a whole.
* Whenever the total traffic for the past 2 minutes exceeds a certain number on average, displays an alert
* Whenever the total traffic drops again below that value on average for the past 2 minutes, displays a message saying that it recovered
* All messages showing when alerting thresholds are crossed remain visible on the page for historical reasons

## Testing

TODO: Write a test for the alerting logic

## Potential future improvements

TODO: Explain how I’d improve on this application design