# Dev mode
- run `air` for live reloads and automatic templ generation
- on sql changes, run `make sqlc`

# Deployment
- run `make build` to build and zip a linux version
- scp zip onto server
- unzip
- cd senatus
- `nohup ./main &` and ctrl+c
- logs will be in nohup.out

# Todos
- handle date changes / dates in general
- date based pagination
- hide event if the time has passed and the following events is in already progress
- double check on delete