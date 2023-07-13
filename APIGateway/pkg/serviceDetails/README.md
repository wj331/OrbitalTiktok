## ServiceDetails handler

This package dynamically handles changes in `../configs/service_configs.json`.

Utilizes the import `github.com/fsnotify/fsnotify` that listens for file system events. If there is a change in `service_configs.json`, updates the pool of generic clients respectively.

Routes are handled via `/servicename/:version/:method`.