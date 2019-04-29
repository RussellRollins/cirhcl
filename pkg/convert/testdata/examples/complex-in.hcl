job "build" {
  docker {
     image = "circleci/ruby:2.4.1"
  }

  checkout {}
  run {
    command = "mkdir -p my_workspace"
  }

  run {
   command = "echo \"Hello Build\" > my_workspace/echo-output"
  }

  run {
    command = "sleep 5"
  }

  persist_to_workspace {
    root = "my_workspace"
    paths = ["echo-output"]
  }
}

job "testa" {
  docker {
     image = "circleci/ruby:2.4.1"
  }

  attach_workspace {
    at = "my_workspace"
  }
  run {
    command = "echo \"Hello Test A\""
  }

  run {
    command = "cat my_workspace/echo-output"
  }

  run {
    command = "sleep 5"
  }
}

job "testb" {
  docker {
     image = "circleci/ruby:2.4.1"
  }

  checkout {}
  run {
    command = "echo \"Hello Test B\""
  }

  run {
    command = "sleep 5"
  }
}
job "deploy" {
  docker {
    image = "circleci/ruby:2.4.1"
  }
  checkout {}
  run {
    command = "echo \"Hello Deploy\""
  }

  run {
    command = "sleep 5"
  }
}

workflow "build_and_test" {
  workflow_job "build" {}
  workflow_job "testa" {
    requires = ["${workflow_job.build.name}"]
  }
  workflow_job "testb" {
    requires = ["${workflow_job.build.name}"]
  }
  workflow_job "deploy" {
    requires = ["${workflow_job.testa.name}", "${workflow_job.testb.name}"]
  }
}
