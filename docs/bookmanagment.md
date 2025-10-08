#bookmanagment

Redis Caching Scenarios – Library Management System
1. Books Module
🎯 Goal:
To improve performance and reduce database load for frequently accessed book data.
🧩 Scenarios:
Get All Books
Key: books:all
Cache Type: Read Cache (Full List)
Process:
When the user requests the list of all books, first check Redis for the key books:all.
If it exists → return cached data.
If not → fetch from the database, store it in Redis, and then return it.
Expiration: 5 minutes
Get Book by ID
Key: book:{id}
Cache Type: Read Cache (Single Record)
Process:
On request, check if book:{id} exists in Redis.
If found → return cached version.
If not found → get it from DB, cache it, and return it.
Expiration: 10 minutes
Search Books
Key: books:search:{query}
Cache Type: Read Cache (Query Result)
Process:
For each search query, store the result of the search in Redis.
When the same search term is used again, return data directly from cache.
Expiration: 3 minutes (shorter since search results can change often)
Add / Update / Delete Book
Cache Invalidation:
On any write operation (add, update, delete), remove:
books:all
book:{id}
Any related search keys (books:search:*)
This ensures consistency between Redis and the main database.