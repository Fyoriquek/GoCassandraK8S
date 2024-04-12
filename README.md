# Cassandra microservice on K8S (Golang)
</br>


![image](https://github.com/Fyoriquek/GoCassandraK8S/assets/166552285/de6f518d-091a-4ba1-ba38-00252d184b4a)



<b> I. Prepare environment </b>

Start by creating the namespace, in this case, cassandra-go
```
 kubectl create namespace cassandra-go
```
</br>

1. Create secrets (you can change the variable values under password folder)
   
CassandraK8S/CassandraSecrets# go run secrets.go

Output:
```
CassandraK8S/CassandraSecrets# kubectl get secret --namespace=cassandra-go
NAME         TYPE     DATA   AGE
dbasecret    Opaque   2      19s

```

</br>

2. Create Cassandra Statefulset:

```
CassandraK8S/CassandraStatefulSet# go run cassandra_statefulset.go
```

</br>

Output:
```
CassandraK8S/CassandraStatefulSet# kubectl get statefulset --namespace=cassandra-go
NAME        READY   AGE
cassandra   1/1     19s

CassandraK8S/CassandraStatefulSet# # kubectl get pods --namespace=cassandra-go
NAME                                     READY   STATUS             RESTARTS      AGE
cassandra-0                              1/1     Running            1             19s

```
</br>

3. Create Cassandra Service for StatefulSet:
```
CassandraK8S/CassandraServiceStatefulSet# go run cassandrastatefulsetservice.go
```
</br>

Output:
```
CassandraK8S/CassandraServiceStatefulSet# kubectl get svc --namespace=cassandra-go
NAME                    TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
cassandra               ClusterIP   None             <none>        9042/TCP         15s

```
</br>
4. Create C* index and table stats.books with the help of CassandraScript/cassandra.cql

```
CassandraScript# kubectl exec --namespace=cassandra-go -ti cassandra-0 -- cqlsh
Connected to Cassandra at 127.0.0.1:9042.
[cqlsh 5.0.1 | Cassandra 3.11.6 | CQL spec 3.4.4 | Native protocol v4]
Use HELP for help.
cqlsh>
cqlsh> CREATE KEYSPACE stats WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 1};
cqlsh> USE stats;
cqlsh:stats> CREATE TABLE stats.books (
         ...     id UUID,
         ...     authorfirstname varchar,
         ...     authorlastname  varchar,
         ...     price int,
         ...     bookname    varchar,
         ...     PRIMARY KEY(id)
         ... );
cqlsh:stats> CREATE INDEX ON stats.books(price);
cqlsh:stats>
```
</br>

5. Create the Docker image microimage for Cassandra Deployment:
```
CassandraMicroservices/CreateGetBooksMicroservice# docker build -t microimage .
``` 
5.b. Tag the image and push it to (in this case) local repository:

```
CassandraMicroservices/CreateGetBooksMicroservice#  docker tag microimage localhost:5000/microimage
CassandraMicroservices/CreateGetBooksMicroservice#  docker push localhost:5000/microimage
Using default tag: latest
The push refers to repository [localhost:5000/microimage]
<...etc...>
```
</br>

6. Create Cassandra Deployment cassandra-createbooks:

```
CassandraK8S/CassandraDeployment# go run deployment.go
```
Output:

```
CassandraK8S/CassandraDeployment# kubectl get deploy --namespace=cassandra-go
NAME                    READY   UP-TO-DATE   AVAILABLE   AGE
cassandra-createbooks   1/1     1            1           21s
CassandraK8S/CassandraDeployment#  kubectl get pods --namespace=cassandra-go
NAME                                     READY   STATUS    RESTARTS      AGE
cassandra-0                              1/1     Running   1             31s
cassandra-createbooks-6688c9334c-2449t   1/1     Running   1             10s
```

</br>

7. Create Cassandra Service cassandra-createbooks:

```
CassandraK8S/CassandraService# go run service.go
```
Output:
```
CassandraK8S/CassandraService#  kubectl get svc --namespace=cassandra-go -o wide
NAME                    TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE    SELECTOR
cassandra               ClusterIP   None             <none>        9042/TCP         45s    app=cassandra
cassandra-createbooks   NodePort    10.105.60.30     <none>        8083:30006/TCP   51s    app=cassandra-createbooks

```
</br>

<b> II. Perform checks </b>

8. From another terminal, run a curl command to create a new book:

```
curl -X POST localhost:30006/books/new \
-H 'Content-Type: application/json'  \
-d '{ "authorfirstname": "Haruki", "authorlastname": "Murakami", "bookname": "1Q84", "price": 55}'
```
You will get an index ID as response (in this case): {"id":"5ec9145d-f822-11ee-b41d-5e8662e7a658"}
</br>

9. Check the inserted rows:

```
curl -X GET localhost:30006/books/all
[
 {
  "id": "ab849e84-f81e-11ee-9377-92257ad4849b",
  "authorfirstname": "Haruki",
  "authorlastname": "Murakami",
  "bookname": "First Person Singular",
  "price": 42
 },
 {
  "id": "5ec9145d-f822-11ee-b41d-5e8662e7a658",
  "authorfirstname": "Haruki",
  "authorlastname": "Murakami",
  "bookname": "1Q84",
  "price": 55
 }
]
```
</br>

10. Check from the cqlsh side:
```
 kubectl exec  --namespace=cassandra-go -ti cassandra-0 -- cqlsh -e "select * from stats.books;"

 id                                   | authorfirstname | authorlastname | bookname              | price
--------------------------------------+-----------------+----------------+-----------------------+-------
 ab849e84-f81e-11ee-9377-92257ad4849b |          Haruki |       Murakami | First Person Singular |    42
 5ec9145d-f822-11ee-b41d-5e8662e7a658 |          Haruki |       Murakami |                  1Q84 |    55
```


