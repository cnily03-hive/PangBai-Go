# PangBai-HTTP

This is the challenge for NewStarCTF 2024 in the category of Web, Week 4.

The participants need to leak JWT key via SSTI and take exploit of SSRF to read any file, and then read `/proc/self/environ` to get the flag.

The challenge provides the whole source code including [hind.md](hind.md) to participants.

## Deployment

> [!NOTE]
> If the development is at ichunqiu platform, please modify [docker-compose.yml](docker-compose.yml) to change `Dockerfile` into `Dockerfile.icq` and the environment variable `FLAG` to `ICQ_FLAG`.

Docker is provided. You can run the following command to start the environment quickly:

```bash
docker compose build # Build the image
docker compose up -d # Start the container
```

## Exploit

The `/eye` route has vulnerability of SSTI. By passing `input` GET parameter, the server will render the template with the input.

Therefore, visit `/eye?input={{.Config.JwtKey}}` to leak the JWT key.

The `/favorite` route needs us to send a PUT request with cookie and body. We can easily sign the Cookie whose `user` is `Papa` with the leaked key, and use `{{ .Curl "url" }}` to take exploit of SSRF. Note that the Gopher protocol is useful to send raw TCP request of HTTP.

After changing the value of `SignaturePath` into `/proc/self/environ`, we can read the environment variables of the server by simply visiting `/favorite`.

The exploit script is provided on [exploit/exp.py](exploit/exp.py). Just run it with Python.

```bash
python3 exp.py '172.18.0.2:8000' # Please replace the origin
```

## License

Copyright (c) Cnily03. All rights reserved.

Licensed under the [MIT](LICENSE) License.
