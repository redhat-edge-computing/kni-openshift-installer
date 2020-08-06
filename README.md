# KNICTL / OpenShift-Install Wrapper

As the name states, this is a small app for wrapping `knictl` and `openshift-install` operations for a specified site. 

## Dependencies

> NOTE: These are included in the container, only required when running the app non-containerized.

- `knictl` must be built from source and reachable in $PATH.  Code and build instructions here: https://github.com/akraino-edge-stack/kni-installer
- A bash/sh environment, doesn't work on windows.
- A container runtime and client, if running the app containerized (recommended)
- Note that `openshift-install` is **not** a prerequisite as knictl will download the binary during runtime.

## Build

The `hack/build.sh` script...builds the app.  Supply the `--image` flag to build the app container (recommended).  The container will have the necessary system dependencies installed.

## Usage

The app provides a familiar command structure to `knictl` and `openshift-install`.  Each command and subcommand will print the help menu on input error or `-h` or `--help`.  Refer to these for a detailed explanation of the commands since they're subject to change and thus likely to be stale if duplicated here.

Like `knictl`, the `create` subcommand expects a GitHub path to the site directory, e.g. `github.com/path/to/site`. 

`knictl` requires credential files to operate. For cloud platforms, these should be the default `$HOME/.<platform>` locations. To deploy the OpenShift cluster, a `pull-secret.json` file must exist in `$HOME/.kni`

#### Container Execution

It's necessary to inject credential files into the container.  The following assumes an AWS deployment, but you must substitute in the appropriate directory for your the cloud platform.

`openshift-install` writes cluster metadata back to the filesystem for later use during `destroy`.  It is **very** important that you preserve the data in  `*/.kni`  if you don't want to tear down the cluster by hand.  This can be done easily by bind-mounting a directory from the host machine (as done in the commands below), or by creating a docker volume.  To create a docker volume, change the `src=$HOME/.kni` to `src=NameYourVolume` and the `type=bind` arg in the same list to `type=volume`.  Do not change the AWS --mount command.

Bind-mounted directories and credential files must exist on the host prior to use.

```bash
mkdir $HOME/.kni
cp <pull-secret> $HOME/.kni
```

##### Build the container image

```bash
./hack/build.sh --image
```

##### Deploy a blueprint

```bash
docker run --rm \
--mount type=bind,src=$HOME/.aws,dst=/root/.aws,readonly \
--mount type=bind,src=$HOME/.kni,dst=/root/.kni 
quay.io/jcope/kni-install create cluster --repo github.com/path/to/site
```

##### *Or* Create Ignition Configs

```bash
docker run --rm \
--mount type=bind,src=$HOME/.aws,dst=/root/.aws,readonly \
--mount type=bind,src=$HOME/.kni,dst=/root/.kni \
quay.io/jcope/kni-install create ignition-configs --repo github.com/path/to/site
```

##### Tear down a cluster

```bash
docker run --rm \
--mount type=bind,src=$HOME/.aws,dst=/root/.aws,readonly \
--mount type=bind,src=$HOME/.kni,dst=/root/.kni \
quay.io/jcope/kni-install destroy cluster --repo github.com/path/to/site
```

##### Baremetal Teardown

// TODO

Thanks to @jtudelag whose container file I ~~stole~~ used as a base for the app container :)

#### Local Execution

To run the app locally, first ensure you have a kni directory, and place the pull-secret.json file in it.

```bash
make $HOME/.kni
cp <pull-secret> $HOME/.kni
```

##### Deploy a blueprint

```bash
kni-install create cluster --repo github.com/path/to/site
```

##### *Or* Create Ignition Configs

```bash
kni-install create ignition-configs --repo github.com/path/to/site
```

##### Tear down a cluster

```bash
kni-install destroy cluster --repo github.com/path/to/site
```



