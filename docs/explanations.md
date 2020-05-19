# Controller
The controller is responsible for accepting routing directives.  
The controller is injected with a db object and a transaction engine(te) object  
which are used to query data from the database or execute some transaction.  
The data produced by queries or the te are marhsalled by the controller in a byte slice  
and returned back to the calling router method.  
The controller is also responsible for locking critical areas of the code such as read and writes to maps  
to prevent data races.

# Database 
The database package includes a Storage Interface and  Database struct.  
The db struct contains four maps: Users, Orders, Currencies, Wallets.
Every wallet maintains a list of currency ids which it owns.
The Users map maintains a list of active users using the API.
The storage interface is used accros the application to make DI easier and allow mocking for testing.  

# Router
The Router package implements 3 mehtods:  
1. HandleRequest which is responsible for unmarshalling the incoming request.  
      * This method takes an additional argument which is a function and acts as a middleware.  
2. IdentifyUser is responsible for checking if the calling entitity has provided a valid userID to the API  
      * This method rejects calls without userID excluding the registration call where an uuid is generated  
      * IdentifyUser is used as the middleware here to guard the routes for the HandleRequest method  
3. RouteRequest implements a switch which multiplexes the API functionality and calls the correct controller method.  
      * This method sends the byte slice returned from the controller or any errors that might have occured  

# Transaction Engine (te)
This package implements the logic for the Limit Order functionaility. 
1. Check the incoming transaction for the type of Order either buyer or seller.  
2. Try place the oreder i.e. try to find matching buyer/seller.  
3. If there are any mathcing offers a list of them is returned and passed to Execute Transaction method.  
4. Execute Transaction sets the exchange rate between the offer to the selling price(Price field) of the new order that has arrived.  
5. The te goes through the list trying to fill the new order until the sum to invest is reduces to 0 or the entity no longer has.  
   either money or tokens to transact with 
6. Once the order is filled it is then marked as Deleted 