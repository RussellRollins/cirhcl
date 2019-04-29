version = "2"

job "build" {
  docker {
     image = "circleci/ruby:2.4.1"
  }

  checkout {}
  run {
    command = "echo \"A first hello\""
  }
}
