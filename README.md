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

It's necessary to inject credential files into the container.  The following assumes an AWS deployment, but you must substitute in the appropriate directory for your the cloud platform.

`openshift-install` writes cluster metadata back to the filesystem for later use during `destroy`.  It is **very** important that you preserve the data in  `*/.kni`  if you don't want to tear down the cluster by hand.  This can be done easily by bind-mounting a directory from the host machine (as done in the commands below), or by creating a docker volume.  To create a docker volume, change the `src=$HOME/.kni` to `src=NameYourVolume` and the `type=bind` arg in the same list to `type=volume`.  Do not change the AWS `--mount` command.

##### Required Bind-Mounts

Openshift-install requires certain stateful data in order to tear down a cluster.  Losing this state means you'll have to manually, and painfully, perform the teardown yourself.  To protect users from finding themselves in this position, the app will check for certain directories and fail if it cannot find them.  Users must provide stateful locations to the container to store openshift-install output.  The mount-points of these directories on the container are `/root/.kni`. If deploying on a cloud platform, the default credential location (e.g. `/root/.aws/`) must be mounted as well.

Additionally, the `/root/.kni` is expected to contain the `pull-secret.json` . 

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
--mount type=bind,src=$HOME/.aws,dst=/root/.aws,readonly \
--mount type=bind,src=$HOME/.kni,dst=/root/.kni 
localhost/kni-install create cluster --repo github.com/path/to/site
```

> **To only deploy a bare cluster (i.e. not apply workload manifests), use the `--bare-cluster` option**

##### *Or* Create Ignition Configs

```bash
docker run --rm \
--mount type=bind,src=$HOME/.aws,dst=/root/.aws,readonly \
--mount type=bind,src=$HOME/.kni,dst=/root/.kni \
localhost/kni-install create ignition-configs --repo github.com/path/to/site
```

##### Tear down a cluster

```bash
docker run --rm \
--mount type=bind,src=$HOME/.aws,dst=/root/.aws,readonly \
--mount type=bind,src=$HOME/.kni,dst=/root/.kni \
localhost/kni-install destroy cluster --repo github.com/path/to/site
```

##### Baremetal Teardown

// TODO

Thanks to @jtudelag whose container file I ~~stole~~ used as a base for the knictl build stage :)
