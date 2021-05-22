# order service nosql
This handles order service using mongodb and only focus handling product stoct substraction which is the 'buggy' part on reported case (product stock showed minus).

# why product stock showed minus integer?

*This analyis assumes that order service already using database 'transaction' when product showed minus value because it involves payment*

The product on event 12.12 perfomed well means many customer ordering that product. It causes race condition on that product row (or document). After get data to validate stock - reqQty ad if pass then it continues using query to substract on database **without adding condition stock greater than requested item quantity (x) or when get data did not use either 'select ... for update' (for mysql/postgres)**. See below example for substract on db without adding 'stock > x'.
1. Example for above query(in mysql/postgres) : UPDATE product SET stock = stock - x WHERE id = ?; **(this query could cause product stock minus Integer)**

**To handle above case, we can do :**
1. If still using mysql on mysql or postgres, either add condition 'stock > x' on update query then abort when updated row is 0, or use 'select ... for update' to get item data and which blocks other transaction to do 'select ... for update' and write process. This probably causes db and order service performance decreasing because it even blocks reading.
- Example : UPDATE product SET stock = stock - x WHERE id = ? AND stock > ?; 

2. Use mongodb (the versions which have transaction feature). Using this, inside transaction, when trying to update document which during transaction there are other process modified the doc, the update inside transaction will return error write conflict. This will handle the race condition using because will return write conflict and then we can add retry logic on order service.

# How to use the service
1. Prepare mongodb server with it's replicas to be able to use transaction (docker-compose example for mongodb replicas is on folder docs)
2. Fill conf.json with respected values
3. Create database with name 'online_store' and collection with name 'product'
4. Add new document on 'product' (fill code and stock field)
5. run service using go run main.go on service directory or build it's executable file first.
6. hit : curl -X POST -d '{"id":"somegeneratedidfrommongodbdocument", "qty":1}' http://localhost:7789/order-product
7. To use hitter to test api : go to folder 'cmd' then 'go run main.go id loop'. For Example : 'go run main.go 60a8921abea0d20007330b9b 20'. Please change the args id using document id generated from your mongodb and loop with some positive integer.