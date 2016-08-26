# berkeley

Distributed clock synchronization with the Berkeley algorithm.

```
Usage: ./berkeley (-m or -s) -addr=0.0.0.0:0 [-slaves=slavesJsonFile.json]
  -m      Run program as master node that will compute the synchronization algorithm
  -s      Run program as slave node that will listen for requests from the master node 
          for its current time, and receives a synchronization value
  -addr   IP:Port string for the program to run under eg. "-addr=127.0.0.1:1337"
  -slaves Name of json file containing the list of slaves nodes addresses

```

The Slaves file will look something like this

```
[
  "127.0.0.1:1234",
  "127.0.0.1:1235",
  "127.0.0.1:1236",
  ...
]
