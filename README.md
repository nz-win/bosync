# bosync - Backorder Sync

## Summary

This is the first version of small cli program I created that fetches results from an internal API proxy and then updates a database so that the data can be used for reporting.

This repository demonstrates an early initial version of the program and only provides basic functionality.
The current version in prod is more flexible and allows you to configure different types of queires and push the data into various databases in different ways.

Authentication/Authorization between this program and the target proxy is handled by nginx using mTLS.

