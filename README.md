# prem-gateway

`prem-gateway` is an API gateway designed to handle and manage various operations, directing requests from `prem-app` to either `prem-daemon` for Docker image management or directly to Docker images that provide `prem-services`.

## Table of Contents

- [Features](#features)
- [Services](#services)
- [Usage](#usage)
  - [Setup](#setup)
  - [Start/Stop](#startstop)
  - [Production](#production)
  - [Service Management](#service-management)
- [License](#license)

## Features

- API Gateway
- Domain Management
- TLS
- *(More features like Authentication/Authorization, Rate Limiting, Logging, and Metrics to be added soon.)*

## Services

- `dnsd` - [Read More](./dns/README.md)
- `authd` - [Read More](./auth/README.md)
- `controllerd` - [Read More](./controller/README.md)

## Usage

### Setup

Create a Docker network:

```bash
docker network create prem-gateway
```

Set file permission:

```bash
chmod 600 ./traefik/letsencrypt/acme.json
```

### Start/Stop

Start the `prem-gateway`:

```bash
make up
```

To stop the gateway:

```bash
make down
```

### Production

For production environments, use the Let's Encrypt production server:

```bash
make up LETSENCRYPT_PROD=true
```

### Service Management

To restart services and assign them with a subdomain/TLS certificate:

```bash
make up LETSENCRYPT_PROD=true SERVICES=premd,premapp
```

Run `prem-gateway` with `prem-app` and `prem-daemon`:

```bash
make runall PREMD_IMAGE={IMG} PREMAPP_IMAGE={IMG}
```

Stop `prem-gateway`, `prem-app`, and `prem-daemon`:

```bash
make stopall PREMD_IMAGE={IMG} PREMAPP_IMAGE={IMG}
```

## License

`prem-gateway` is licensed under [MIT License](./LICENSE).

---

Feel free to adapt this template to better fit the specifics and personality of your project!
