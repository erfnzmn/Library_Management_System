# Books Management

## must handle

- storing and managing book records.
- supporting searching and filtering.
- Tracking status of each copy.
- handling book lifecycle.
- Integrating with user, reservation and selling.

## key entities and attributes

1. id (unique identifier)
2. title
3. auther
4. ISBN
5. publisher
6. year_of_publication
7. edition
8. genre/category
9. language
10. description
11. reservation_status(available, reserved)
12. selling_status (available, sold out)
13. cover_ image
14. tags/keywords
15. price (for selling)

## Core Usecases

1. Add book records
    - Admin adds a new title.
    - Admin adds needed informations completing book details or system can fetch details from ISBN APIs.
2. Update book details
    - Admin can edit book details if needed.
3. Remove/Deactivate books
    - Admin can deactivate book records if needed.
4. Search/Browse for books
    - User can search for books using: name, auther, ISBN.
5. Filter records
    - User can filter list by: genre, tags, date of publication and status.
6. View books details
    - User clicks to see book details such as copies available, description, reservation/sell buttons and overall details.
7. Favorite - User wants ti add a specific book to his favorite books list. 8. Integrate with selling - copies marked as "for sale" can be purchased.

## Scenarios

Let’s walk through typical scenarios:

### Scenario A: A student wants to borrow a book

- Student searches for the title.
- System shows book details + copies.
- Student clicks “Reserve” (handled by reservation module).
- Status of that copy changes to “reserved”.

### Scenario B: A librarian adds new books

- Librarian enters book metadata.
- System generates a book record.
- the new book is added to the DB and book lists

### Scenario C: A user wants to buy a book

- User searches or browses books.
- System shows if the book is available for sale.
- User clicks “Buy”, handled by selling module.

### Scenario D: User wants to add a book to their favorites

- First, user needs to find the book he is looking for using search, filter or browsing through the list.
- Clicks on the book to see the details.
- On book detail page clicks on the "Add to favorites" button.
- System adds the book to users favorites and shows a success message.

## Integration points

1. Reservation module
    - tracks reservations and updates book status.
2. Selling module
    - marks copies as sold and generates receipts.

## Extra features (optional)

1. Digital book support.
2. Book rating/reviews from users.
