# opentelemetry-collector-releases Release Procedure

This document describes the steps to follow when doing a new release of the
opentelemetry-collector-releases repository. This release depends on the [opentelemetry-collector release][1] and the [opentelemetry-collector-contrib release][2] repositories.

For general information about all Collector repositories release procedures, see the
[opentelemetry-collector release process documentation][1].

## Releasing opentelemetry-collector-releases

0. Ensure that the [opentelemetry-collector][1] and [opentelemetry-collector-contrib][2] release procedures have been followed and a new
   opentelemetry-collector and opentelemetry-collector-contrib version have been released. The opentelemetry-collector-releases release
   should be done after both of these releases.
1. Review and merge the 3 Renovate PRs for core and contrib components, as well as OCB.
2. Run the GitHub Action workflow "[Update Version in Distributions and Prepare PR](https://github.com/open-telemetry/opentelemetry-collector-releases/actions/workflows/update-version.yaml)" which will update the minor version automatically (e.g. v0.116.0 -> v0.117.0) or manually provide a new version if releasing a bugfix or skipping a version. Select "create pr" option.
The PR needs to be manually closed and re-opened once to trigger pipelines.
   -  🛑 **Do not move forward until this PR is merged.** 🛑
3. Check out main and ensure it has the "Update version from ..." commit in your local
   copy by pulling in the latest from
   `open-telemetry/opentelemetry-collector-releases`. Assuming your upstream
   remote is named `upstream`, you can try running:
   - `git checkout main && git fetch upstream && git rebase upstream/main`
4. Create a tag for the new release version by running:
   
   ⚠️ If you set your remote using `https` you need to include `REMOTE=https://github.com/open-telemetry/opentelemetry-collector-releases.git` in each command. ⚠️
   
   - `make push-tags TAG=v0.85.0`
5. Wait for the new tag build to pass successfully.
6. Ensure the "Release Core", "Release Contrib", "Release k8s", "Release OTLP", "Release Builder" and "Release OpAMP Suporvisor" actions pass, this will
    1. push new container images to `https://hub.docker.com/repository/docker/otel/opentelemetry-collector`, `https://hub.docker.com/repository/docker/otel/opentelemetry-collector-contrib` and `https://hub.docker.com/repository/docker/otel/opentelemetry-collector-k8s` as well as their respective counterparts on GHCR
    2. create a Github release for the tag and push all the build artifacts to the Github release. See [example](https://github.com/open-telemetry/opentelemetry-collector-releases/actions/workflows/release-core.yaml).
    3. build and release ocb and opampsupervisor binaries under a separate tagged Github release, e.g. `cmd/{builder,opampsupervisor}/v0.85.0`
    4. build and push ocb and opampsupervisor Docker images to `https://hub.docker.com/r/otel/opentelemetry-collector-builder` and the GitHub Container Registry within the releases repository (and opampsupervisor respectively)
7. Update the release notes with the CHANGELOG.md updates.

## Post-release steps

After the release is complete, the release manager should do the following steps:

1. Create an issue or update existing issues for each problem encountered throughout the release in
the opentelemetry-collector-releases repository and label them with the `release:retro` label.
Communicate the list of issues to the core release manager.

## Bugfix releases

See the [opentelemetry-collector release procedure][1] document for the bugfix release criteria and
process.

## Release schedule and release manager rotation

See the [opentelemetry-collector release procedure][1] document for the release schedule and release
manager rotation.

[1]: https://github.com/open-telemetry/opentelemetry-collector/blob/main/docs/release.md
[2]: https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/docs/release.md
