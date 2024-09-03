job "example" {
  datacenters = ["dc1"]
  type        = "batch"

  group "example" {
    task "ping" {
      driver = "systemd-run"

      config {
        unitname = "nick-test.service"
        command = "ping"
        args = ["-c", "60", "www.google.com"]
      }
    }
  } 
}
