# lbaas
RESTful service to provide generic VIP management for a variety of load balancer backends


# Data Flow
Assumptions
-----------
1. One or more nodes are configured in Kubernetes
2. A service has been set up Kubernetes and there is a working NodePort fronting it

```lbaas``` process
-------------------
1. A VIP is created with ```lbaasctl```, a command line client to the ```lbaasd``` RESTful service.  A Kubernetes service name is supplied, along with a front-end port and protocol for the VIP to listen on.
  ```
   % lbaasctl <unique VIP id> <Kubernetes service name> <listen port> <listen protocol> <optional profiles to apply (ssl, etc)>
  ```

2. ```lbaasd``` records this new VIP in etcd under ```/vips/vip-id```
3. The K8S nodes watcher engine watches Kubernetes for node changes and records a list of current nodes in a mutex-controlled shared data structure and sends an "update" signal via channel to the load balancer updater engine
3. The K8S services watcher engine is notified of this new VIP by channel and begins watching the VIP's backend K8S service for the NodePort.
4. The K8S services watcher engine records the NodePort in a mutex-controlled shared data struture and sends an "update" signal via channel to the load balancer updater engine
5. When the load balancer updater engine recieves an update signal from either the services watcher or node watcher, it reconfigures the VIP on the load balancer (but only if both the NodePort and backend nodes are present in the VIP configuration)
