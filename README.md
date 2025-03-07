# ncli

Simple netconf command line client

## Usage

```bash
ncli --help
Simple netconf command line client

Usage:
  ncli [command]

Available Commands:
  commit          Commit changes in candidate datastore
  completion      Generate the autocompletion script for the specified shell
  copy-config     Copy configuration from source to target datastore
  delete-config   Delete configuration from target datasource
  discard-changes Discard changes from candidate datastore
  edit-config     Send edit-config rpc to specified target datastore
  get             Send get rpc with specified filter or filter file
  get-config      Send get-config rpc with specified filter and source datasource
  get-schema      Get schema with specified identifier
  hello           Send hello request
  help            Help about any command
  kill-session    Kills session with the specified id
  rpc             Send rpc request
  validate        Validate changes in specified datastore

Flags:
  -h, --help                   help for ncli
      --host string            hostname or address of the device
      --lock string            wrap calls with lock/unlock - if applicable
      --logging-level string   set logging level - info,debug,critical
      --password string        password for authentication
      --port int               port of the device (default 830)
      --username string        username for authentication
  -v, --version                version for ncli

Use "ncli [command] --help" for more information about a command.
```

### Examples

```bash
ncli --host clab-single-leaf1 --username admin --password admin@123 hello
```

```bash
ncli --host clab-single-leaf1 --username admin --password admin@123 get --path /netconf-state/capabilities
```

```bash
ncli --host clab-single-leaf1 --username admin --password admin@123 get-config --source running --path /
ncli --host clab-single-leaf1 --username admin --password admin@123 get-config --source running --path /configuration/system/host-name
```

```bash
ncli --host clab-single-leaf1 --username admin --password admin@123 rpc --rpc /get-route-information
```

```bash
ncli --host clab-single-leaf1 --username admin --password admin@123 edit-config --target candidate --path /configuration/system/host-name --value leaf1
ncli --host clab-single-leaf1 --username admin --password admin@123 discard-changes

ncli --host clab-single-leaf1 --username admin --password admin@123 edit-config --target candidate --path /configuration/system/host-name --value leaf1
ncli --host clab-single-leaf1 --username admin --password admin@123 commit
```

