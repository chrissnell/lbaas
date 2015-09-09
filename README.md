# lbaas
RESTful service to provide generic VIP management for a variety of load balancer backends


# Data Flow
Assumptions
-----------
1. One or more nodes are configured in Kubernetes
2. A service has been set up Kubernetes and there is a working NodePort fronting it

On ```lbaasd``` process startup
------------------------------
1. ```lbaasd``` launches the following engines
 * K8S nodes watcher
 * K8S services watcher
 * Load Balancer updater
 * IP registration renewal
2. ```lbaasd``` begins listening on the API service port

The ```lbaasd``` K8S nodes watcher engine
-------------------------------------
1. On startup, it does an initial poll of K8S for a list of nodes, then sends this list to the load balancer update engine via channel.
2. The engine then watches K8S for nodes changes and sends an updated list to the load balancer update engine when a change event is received.


The ```lbaasd``` K8S services watcher engine
-------------------------------------
1. On startup, the K8S services watcher engine checks etcd for the list of services that it should be watching and does an initial poll of K8S for NodePort mappings.  Once received, this list is sent to the load balancer update engine via channel
2. Subsequently, the engine then watches K8S for services changes and sends the load balancer update engine (via channel) any changes that it receives.

The ```lbaasd``` load balancer update engine
-------------------------------------
1. On startup, the load balancer update engine first polls etcd for the list of VIPs that it will be managing.
2. It then queries ```cidrd``` to register IPs for any VIPs lacking an IP UUID in etcd.  
3. Next, it queries ```cidrd``` with a ```GET``` to request an initial dump of all UUID->IP mappings for the existing VIPs.  
4. As soon as the engine has current IP registrations, NodePort mappings (from services watcher engine), and a nodes list (from nodes watcher engine), it performs an initial update of VIPs on the load balancer.   
5. Subsequently, the load balancer updaate engine listens on channels for update events from the nodes and services engines and updates VIPs as necessary.

The ```lbaasd``` IP renewal engine
1. On startup, the IP registration renewal engine reads a list of IP UUIDs from etcd and then issues a RENEW command to cidrd for each UUID.   
2. It repeats this process every N minutes


Once ```lbaasd``` has started
-----------------------------
1.  A VIP is created with ```lbaasctl```, a command line client to the ```lbaasd``` RESTful service.  A Kubernetes service name is supplied, along with a front-end port and protocol for the VIP to listen on.
  ```
   % lbaasctl <unique VIP id> <Kubernetes service name> <listen port> <listen protocol> <optional profiles to apply (ssl, etc)>
  ```
2. REST API controller records this new VIP in etcd under ```/vips/vip-id```.  This change is picked up by the nodes watcher, services watcher, load balancer updater, and IP renewal engines.
