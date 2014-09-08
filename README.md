# Deploy

Deploy orchestrates deployments on Marathon.

## API

Deploy has a resource-based API. In the docs below:

- `{pool}` is the name of the pool to deploy to
- `{app}` is the name of the app to deploy as
- `{version}` is the version of the app to deploy as

### To create a pool

- `PUT /pools/{pool}`

### To deploy an app

- `PUT /pools/{pool}/apps/{app}/versions/{version}`

You can set up rules for pools. E.g.:

- Set a specific env to deploy to
- Set specific env vars.
- Set resource constraints.

## To get data

Every sub-path of the deploy URI above is gettable. More docs later.

## Pools

Pools represent a broad configuration for a set of deployments. For example, they specify which Marathon instance to deploy to, and can set other env vars. One use for this might be to set a pool as a 'testing' pool, disabling discovery announcements, and perhaps alter logging rules.
