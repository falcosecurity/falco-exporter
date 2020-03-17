# Release Process

Currently, our release process is partially automated. We use: 
- the [goreleaser](https://github.com/goreleaser/goreleaser) tool for artifacts
- the provided `Makefile` for the container image

When we release we do the following process:

1. We decide together (usually in the #falco channel in [slack](https://sysdig.slack.com)) what's the next version to tag
2. A person with repository rights does the tag
3. The same person runs commands in their machine following the "Release commands" section below
4. The tag is live on [Github](https://github.com/falcosecurity/falco-exporter/releases) with the artifacts, and the container image is live on [DockerHub](https://hub.docker.com/r/falcosecurity/falco-exporter) with proper tags

## Release commands

Tag the version

```bash
git tag -a v0.1.0-rc.0 -m "v0.1.0-rc.0"
git push origin v0.1.0-rc.0
```

Run goreleaser, make sure to export your GitHub token first

```
export GITHUB_TOKEN=<YOUR_GH_TOKEN>
goreleaser --rm-dist
```

Finally, build and publish the container image

```
make image/build
make image/push
make image/latest
```

## TODO

- [ ] Setup goreleaser [on a CI system](https://goreleaser.com/ci/), ie., CircleCI
- [ ] Build and publish images using a CI system, ie., CircleCI