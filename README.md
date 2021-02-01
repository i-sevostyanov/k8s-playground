[![Latest Release](https://img.shields.io/github/release/i-sevostyanov/k8s-playground.svg)](https://github.com/i-sevostyanov/k8s-playground/releases/latest)
![CI](https://github.com/i-sevostyanov/k8s-playground/workflows/CI/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/i-sevostyanov/k8s-playground)](https://goreportcard.com/report/github.com/i-sevostyanov/k8s-playground)
[![codecov](https://codecov.io/gh/i-sevostyanov/k8s-playground/branch/main/graph/badge.svg?token=JEFLNDSIY5)](https://codecov.io/gh/i-sevostyanov/k8s-playground)
[![GitHub license](https://img.shields.io/github/license/i-sevostyanov/k8s-playground)](https://github.com/i-sevostyanov/k8s-playground/blob/main/LICENSE)

## Kubernetes Playground
This project contains a ready-to-use Kubernetes playground. It uses Vagrant and Ansible to start k8s on 3 nodes and deploy
ZooKeeper, Kafka, Prometheus Operator, Kafka Exporter, and two golang applications (Consumer and Producer).

### Attention
Ansible playbook is not idempotent, please consider running it only once. The next attempts will fail.

### Prerequisites
You need the following installed to use this playground:
- `Vagrant`, tested with version 2.2.14
- `VirtualBox`, tested with Version 6.1

### Structure
```
k8s-playground
├─ cmd                
│  ├─ consumer         // Main package for consumer
│  ├─ producer         // Main package for producer
├─ deployment          // k8s manifests for consumer and producer 
├─ internal            // Internal packages for applications (by golang convention) 
│  ├─ bussiness        // Business logic of applications
│  ├─ infrastructure   // Application infrastructure code
├─ vagrant             // Contains a Vagrantfile
│  ├─ provision        // Contains a Ansible playbooks for k8s cluster and management machine
│  |  ├─ manifests     // k8s manifests for kafka, zookepeer, prometheus, etc
```

### Producer & Consumer

#### Business logic
`Producer` publishes a current timestamp into `input` topic every 5 seconds (value by default).

`Consumer` reads timestamp from `input` topic, converts it into date in RFC3339 format then publishes it into `output` topic.

#### CI/CD
Github Actions are invoked on every commit to handle code linting, tests, and code coverage. The release creation triggers build action. 
It builds docker image and pushes it to the Github Container Registry. 

### Bringing up the cluster
To bring up the cluster, clone this repository to a working directory.

```
git clone https://github.com/i-sevostyanov/k8s-playground.git
```

Change into the working directory and run `vagrant up` command:

```
cd k8s-playground/vagrant
vagrant up
```

After the command is complete, you can connect to the `management` machine and wait for a few minutes until all pods are started. 
You can check the status by running the following commands:

```
vagrant ssh management
kubectl get pods -o wide -A
```

### Everything is okay?
To check that the consumer and producer are working properly you need to connect to the `management` machine and execute the following command:
```shell
kubectl exec -i -t kafka-0 -- kafka-console-consumer --bootstrap-server localhost:9092 --topic output --from-beginning
```

If the screen displays dates in RFC3339 format with an interval of 5 seconds (by default), then everything works correctly:
```shell
2021-01-01T08:49:26Z
2021-01-01T08:49:31Z
2021-01-01T08:49:36Z
2021-01-01T08:49:41Z
2021-01-01T08:49:46Z
2021-01-01T08:49:51Z
```

### Monitoring

#### Grafana
To access the grafana dashboards, you need to set up port forwarding as follows:
```shell
vagrant ssh management
kubectl port-forward `kubectl get pods -l app=grafana -n monitoring -o name` --address 0.0.0.0 3000:3000 -n monitoring &
```
After that you can access grafana by going to the following address: [http://192.168.50.13:3000](http://192.168.50.13:3000).
In addition to Kubernetes dashboards, you will find a dashboard with consumer and producer metrics, as well as a Kafka dashboard.

Login: admin
Password: admin

#### Prometheus
Same for Prometheus:
```shell
vagrant ssh management
kubectl port-forward prometheus-k8s-0 --address 0.0.0.0 9090:9090 -n monitoring &
```
Prometheus web interface will be available at [http://192.168.50.13:9090](http://192.168.50.13:9090).

### What can be improved
- [ ] Idempotent Ansible Playbooks
- [ ] Linting for Ansible Playbooks ([ansible-lint](https://ansible-lint.readthedocs.io/en/latest/))
- [ ] Split Ansible Playbooks into roles and tasks
- [ ] Testing for Ansible Playbooks ([Molecule](https://molecule.readthedocs.io/en/latest/index.html))
- [ ] Parametrize kubernetes manifests with [helm](https://helm.sh) or [kustomize](https://kustomize.io)
- [ ] Collect metrics from JMX
- [ ] Add more metrics for consumer and producer 
