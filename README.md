# cirhcl

Pronounced like "circle" cirhcl is an HCL2 based configuration library that can parse HCL to a valid CircleCI configuration YAML. This is not a general purpose HCL to YAML converter. Instead, it is opinionated in favor of HCL that closely matches the CircleCI Configuration Spec.

## How to Use

_Warning:_ cirhcl is a prototype/toy/learning tool. Many circle YAML features aren't supported, the documentation is limited, and I wouldn't trust this for day-to-day or "production" type work.

Install cirhcl:

```
cd cmd/cirhcl
go install
```

Create a simple config.hcl in your .circleci directory:

```hcl
job "build" {
  docker {
     image = "circleci/ruby:2.4.1"
  }

  checkout {}
  run {
    command = "echo \"A first hello\""
  }
}
```

Use cirhcl to convert that to a YAML version:

```
cirhcl config.hcl
cat config.yml
```

## Supported Features

To get a more comprehensive view of the supported features, check out the example inputs in `pkg/convert/testdata/examples`.

### Jobs

Jobs can be specified as HCL blocks.

### Executors

Only the Docker executor is supported and only image can be configured.

### Steps

The following steps can be added in order to a job block.

* checkout
* run
* persist_to_workspace
* attach_workspace

### Workflow

A workflow can be specified as an HCL block.

### Workflow Jobs

Individual Workflow steps can be configured, they can also be used as variable inputs to other Workflow Jobs, in order to specify the `requires` block with HCL checking to make sure you only specify valid dependencies.
