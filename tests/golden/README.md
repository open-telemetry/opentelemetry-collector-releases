# Golden file tests

[Golden files](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/pkg/golden) are YAML representation
of OpenTelemetry signals.

Golden files are used in conjunction with the [cmd/golden](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/cmd/golden) utility
which can listen for OTLP traffic and compare it to incoming traffic. If the input matches the golden file, the program
exits with code 0, and exits with code 1 otherwise.

## Update Golden files

Golden files are typically not managed by hand. We use a feature of the golden command to capture input and write it to file.

To do so, in the docker-compose.yaml file, uncomment the parameter `--write-expected` and delete the data/expected.yaml file.

Run the Docker compose setup with (replacing VERSION with the right version to test for):
```shell
MY_UID="$(id -u)" MY_GID="$(id -g)" IMG=ghcr.io/open-telemetry/opentelemetry-collector-releases/opentelemetry-collector-contrib:VERSION docker-compose up -d --wait
```

Check in on your docker compose:

```shell
$> docker ps -a 
CONTAINER ID   IMAGE                                                                                             COMMAND                  CREATED              STATUS                          PORTS                      NAMES
4716ece9c21e   ghcr.io/open-telemetry/opentelemetry-collector-contrib/golden:latest                              "/golden --expected …"   About a minute ago   Exited (0) About a minute ago                              golden
13e72a5f997d   ghcr.io/open-telemetry/opentelemetry-collector-releases/opentelemetry-collector-contrib:0.139.0   "/otelcol-contrib --…"   About a minute ago   Up About a minute               4317-4318/tcp, 55679/tcp   collector
```

Note how golden exited with code 0, meaning it managed to perform its function.

The golden utility will write the file to disk, and you can add it via git to make a diff and understand the changes.

Make a PR with the changes ; CI runs the same steps and will check the match works.

