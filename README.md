# KNICTL / OpenShift-Install Wrapper

As the name states, this is a small app for wrapping `knictl` and `openshift-install` operations for a specified site. 

## Dependencies

- Docker

## Usage

The app provides a familiar command structure to `knictl` and `openshift-install`.  Each command and subcommand will print the help menu on input error or `-h` or `--help`.  Refer to these for a detailed explanation of the commands since they're subject to change and thus likely to be stale if duplicated here.

Like `knictl`, the `create` subcommand expects a GitHub path to the site directory, e.g. `github.com/path/to/site`. 

`knictl` requires credential files to operate. For cloud platforms, these should be the default `$HOME/.<platform>` locations. To deploy the OpenShift cluster, a `pull-secret.json` file must exist in `$HOME/.kni`

## Build

- ` make` / `make build`
  - Pulls source code , compiles the binaries, and collects them into a runtime image. 
  - Optionally pass `BRANCH=<git ref>` to checkout and build a different branch of the `copejon/kni-openshift-installer` repository.

#### Container Execution

##### Required Bind-Mounts

When deploying to a cloud platform, it's necessary to inject credential files into the container.  The instructions below assume an AWS deployment, but you must substitute in the appropriate directory for your the cloud platform.

Openshift-install requires certain stateful data in order to tear down a cluster.  Losing this state means you'll have to manually, and painfully, perform the teardown yourself.  To protect users from finding themselves in this position, the app will check for the `*/.kni` directory and fail if it does not exist.  The directory will be created when the `--mount` command is supplied to the container.

Additionally, the `*/.kni` is expected to contain the `pull-secret.json`. 

> Note: `/root/.kni` is the default path expected by the app.  This path is configurable with the `--kni-dir` option.

See [deploy a blueprint](#deploy-a-blueprint) for an example of the bind-mount options.

##### Build the container image

```bash
make
```

***or if building a non-master branch***

```bash
make BRANCH=<git ref>
```

##### Deploy a blueprint

```bash
docker run --rm \
-v $HOME/.aws:/root/.aws:ro \
-v $HOME/.kni:/root/.kni \
localhost/kni-install create cluster --repo github.com/path/to/site
```

> **To only deploy a bare cluster (i.e. not apply workload manifests), use the `--bare-cluster` option**

##### *Or* Create Ignition Configs

```bash
docker run --rm \
-v $HOME/.aws:/root/.aws:ro \
-v $HOME/.kni:/root/.kni \
localhost/kni-install create ignition-configs --repo github.com/path/to/site
```

##### Tear down a cluster

```bash
docker run --rm \
-v $HOME/.aws:/root/.aws:ro \
-v $HOME/.kni:/root/.kni \
localhost/kni-install destroy cluster --repo github.com/path/to/site
```

# License

Copyright © 2020 Jonathan Cope jcope@redhat.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.