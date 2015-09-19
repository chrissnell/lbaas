
# ```lbaasd``` operation

On ```lbaasd``` process startup
------------------------------
1. ```lbaasd``` launches the following engines
 * K8S nodes watcher - Responsible for polling K8S node objects
 * K8S services watcher - Responsible for polling K8S service objects
 * Load Balancer updater - Responsible for communicating with load balancer
 * IP registration renewal - Responible for maintaining leases with ```cidrd```
 * DNS registration - Updates etcd with latest DNS registrations
 
2. ```lbaasd``` begins listening on the API service port

The ```lbaasd``` K8S nodes watcher engine
-------------------------------------
1. On startup, it begins a watch on all nodes for a specified namespace/selector, sending node events into a buffered channel to be read by the load balancer updater engine.
2. Next, it does an initial poll of K8S for a list of nodes, then sends this list to the load balancer update engine via channel, which signals th beginning of a VIP-by-VIP inspection of LB pool members to ensure that they align with living Kubernetes nodes.


The ```lbaasd``` K8S services watcher engine
-------------------------------------
1. On startup, the K8S services watcher engine checks etcd (via kubernetes controller) for the list of services that it should be watching.
2. Next, it sets up a watcher (via kubernetes controller) on all services and notifies the LB update engine via channel if any changes are noted for any of the watched services.
3. Next, it does an initial poll of K8S for NodePort mappings (via kubernetes controller).  Once received, a call is made to the LB controller to check these mappings against existing VIPs.

The ```lbaasd``` load balancer update engine
-------------------------------------
1. On startup, the load balancer update engine first polls etcd (via LB controller -> etcd model) for the list of VIPs that it will be managing.
3. Next, it queries ```cidrd``` with a ```GET``` (via cidrd controller) to request an initial dump of all UUID->IP mappings for the existing VIPs.  
3. It then queries ```cidrd``` (via cidrd controller) to register IPs for any VIPs lacking an IP UUID in etcd and calls the DNS controller to make any required DNS updates (which propagate to etcd)  
4. As soon as the engine has current IP registrations, NodePort mappings (from services watcher engine), and a nodes list (from nodes watcher engine), it calls to the load balancer controller to do a VIP-by-VIP inspection of LB pool members to ensure that they align with living Kubernetes nodes.
5. sends corresponding update events to the DNS updater engine.  
6. Subsequently, the load balancer update engine listens on channels for update events from the nodes and services engines and updates VIPs as necessary.
7. The engine also gets a watcher (via lb controller -> etcd model) to watch for VIP deletion events.  If a VIP is deleted, the deletion is propagated to the load balancer (via lb controller), the services watcher (via kubernetes controller), and the DNS updater (via dns controller).   Subsequently, a DELETE event is sent to cidrd (via cidrd controller).

The ```lbaasd``` IP renewal engine
-------------------------------------
1. On startup, the IP registration renewal engine reads a list of IP UUIDs from etcd and then issues a RENEW command to cidrd for each UUID.   
2. It repeats this process every N minutes

The ```lbaasd``` DNS updater engine
-----------------------------------
1. On startup, the DNS updater engine starts up and listens on channels for events (consisting of a service name and an IP) from the load balancer update engine.  Whenever an incoming event is received, it updates [dns service] (via dns controller) with the latest mapping.

Once ```lbaasd``` has started
-----------------------------
1.  A VIP is created with ```lbaasctl```, a command line client to the ```lbaasd``` RESTful service.  A Kubernetes service name is supplied, along with a front-end port and protocol for the VIP to listen on.
  ```
   % lbaasctl <unique VIP id> <Kubernetes service name> <listen port> <listen protocol> <optional profiles to apply (ssl, etc)>
  ```
2. REST API controller records this new VIP in etcd under ```/vips/vip-id```.  This change is picked up by the nodes watcher, services watcher, load balancer updater, and IP renewal engines.

