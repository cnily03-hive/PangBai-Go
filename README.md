# PangBai-HTTP

This is the challenge for NewStarCTF 2024 in the category of Web, Week 4.

The participants need to leak JWT key via SSTI and take exploit of SSRF to read any file, and then read environment file under `proc/` directory to get the flag.

The challenge provide the whole source code including [hind.md](hind.md) to participants.

## Deployment

> [!NOTE]
> If the development is at ichunqiu platform, please modify [docker-compose.yml](docker-compose.yml) to change `Dockerfile` into `Dockerfile.icq` and the environment variable `FLAG` to `ICQ_FLAG`.

Docker is provided. You can run the following command to start the environment quickly:

```bash
docker compose build # Build the image
docker compose up -d # Start the container
```

## Exploit

The exploit script is provided in [exploit/exp.py](exploit/exp.py). Just run the script with Python.

```bash
python3 exp.py '172.18.0.2:8000' # Please replace the origin
```

## License

Copyright (c) Cnily03. All rights reserved.

Licensed under the [MIT](LICENSE) License.
