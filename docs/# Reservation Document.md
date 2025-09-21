# Reservation Document

## Usecases

- Place reservations
- View my reservations
- Cancel reservation
- Admin cancel
- Auto-Allocate Copy (System)
- Expire reservation (System)

## Scenarios

1.  

### Preconditions

- User is authenticated and not banned.
- User is under limit of active reservations.
- User has no active reservation for this title.

### Main success scenario

- User chooses the books and clicks on the title to open book details.
- Clicks on Reserve.
- User chooses the amount of time that it'll need the book.
- System checks if the book is available at the moment:
    1. if it is system will notify the user by a message that his reserved book is good to pickup.
    2. if not the system will notify the user that he is in the queue.

### Alternative Flows

- Already reserved same title → return ALREADY_RESERVED and show current reservation.
- A3: Over user reservation limit → return OVER_LIMIT.
- User blocked (overdue/fines/banned) → return USER_BLOCKED and show the reason.

2.  

### Preconditions

- User is authenticated.

### Main success scenario

- User opens MY RESERVATIONS.
- System returns a paginated list of users reservation history with filter (Date, status).

### Alternative Flows

- No reservations → return empty list.

3. 

### Preconditions

- User is authenticated and owns the reservation.
- Reservation status ∈ {QUEUED, READY_FOR_PICKUP}.

### Main success scenario

- Users opens his reservation list.
- User clicks Cancel on a reservation.
- System marks reservation CANCELLED.
- If the reservation was READY_FOR_PICKUP, the assigned copy is freed.
- System triggers reallocation (next user in queue) and sends notifications as needed.
- Return updated item or users reservation list.


