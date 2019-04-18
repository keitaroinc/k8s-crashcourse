const { events, Job } = require("brigadier");

events.on("push", function(e, project) {
  console.log("Starting build pipeline");

  var build = new Job("build-images");
  build.image = "docker:stable-dind";
  build.privileged = true;
  build.docker = {
    enabled: true
  }
  build.env = {
    DOCKER_DRIVER: "overlay"
  }

  build.tasks = [
    "docker build /src/sample-apps/api -f /src/sample-apps/api/Dockerfile -t api:latest",
    "docker build /src/sample-apps/backend -f /src/sample-apps/backend/Dockerfile -t backend:latest"
  ]

  build.streamLogs = true;
  build.run();
  events.emit("deploy", e, project)
})

events.on("deploy", function(e, project) {
  console.log("Starting deploy pipeline");

  var deploy = new Job("deploy-sample-apps");
  deploy.image = "bitnami/kubectl:1.13"

  deploy.tasks = [
    "kubectl scale --replicas=0 -f /src/sample-apps/api/deployment.yaml",
    "kubectl scale --replicas=1 -f /src/sample-apps/api/deployment.yaml",
    "kubectl scale --replicas=0 -f /src/sample-apps/backend/deployment.yaml",
    "kubectl scale --replicas=1 -f /src/sample-apps/backend/deployment.yaml"
  ]

  deploy.streamLogs = true;
  deploy.run();
})

