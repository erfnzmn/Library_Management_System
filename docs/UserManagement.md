# User Management

## Functional Requirements Extraction

- Ability to register a new user and assign a role  
- Ability to log in via email or username  
- Ability to recover a forgotten password  
- Ability to view and update user information  
- Ability to sign out or delete one’s own account  
- Ability for an admin to delete a user from the system  
- Display each user’s status in the user list (Active, Inactive, Suspended, etc.)  
- Ability for an admin to view the full list of users  
- Record each login’s timestamp and IP address for security auditing  

## Actors

- Admin  
- Normal User  

## Use Cases

**User:**  

- Sign up  
- Log in  
- Recover password  
- Update user  
- Sign out  
- Delete user  

**Manager:**  

- See user info  
- See users list  
- See user status  
- Delete or ban user  

## Usecase scenarios

### Sign Up

**Preconditions:**  

- The user must not already be logged in.

**Main Scenario:**  

1. The user clicks **Sign Up** and views the registration form.  
2. The user enters the required information.  
3. The system:  
   - Checks that no duplicate username or email exists.  
   - Verifies that the password meets policy requirements.  
   - Saves the submitted information.  
4. A success message is displayed.  
5. The system generates a token for the user and redirects them to the home page.

**Alternative Flows:**  

- **If the username or email is already in use:**  
  1. Display an error message.  
  2. Allow the user to either go to the login page or try a different username or email.  

- **If the password does not meet requirements:**
  1. Display an error message specifying the issue.  
  2. Allow the user to retry.  

- **If a database error occurs:**  
  1. Display a system error message.  
  2. Prompt the user to try again.

**Postconditions:**  

- A new user is created with the specified role.  
- The user is redirected to the home page under their username.

### Login

**Preconditions:**  

- The user must already be registered in the system.  
- The user must not be logged in during the current session.
**Main Scenario:**  

1. The user enters their email or username and password.  
2. The system validates that:  
   - The email or username exists in the database.  
   - The entered password matches the hashed password in the database.  
3. If validation succeeds:  
   - The system issues an authentication token.  
   - The login time and IP address are recorded in the Audit table.  
   - The user is redirected to their dashboard or the page they originally requested.  
   - A welcome/success message is displayed.

**Alternative Flows:**  

- **Email or username not found:**  
  - Display “No account found with that email or username.”  
  - Allow the user to retry or click Sign Up.

- **Incorrect password:**  
  - Display “Incorrect password.”  
  - Allow the user to retry or click Forgot Password.

- **Too many failed attempts (e.g., 5):**  
  - Temporarily lock the account for X minutes.  
  - Display “Account temporarily locked; please try again in 10 minutes.”

- **Server or database error when issuing token:**  
  - Display “Internal error; please try again.”

**Postconditions:**  

- The user is successfully authenticated and holds a valid token.  
- A record of the login event (timestamp and IP) exists in the Audit table.

### Password Recovery

**Preconditions:**  

- The user must already be registered in the system (an account with that email or username exists).  
- The user must have access to the email address associated with their account.  
- The email delivery service must be available and operational.

**Main Scenario:**  

1. The user clicks the “Recover Password” option on the login page.  
2. The user enters the email address for password recovery.  
3. The system:  
   - Checks whether that email exists in the database (i.e. is already registered).  
   - Generates a one-time token with a limited expiration.  
   - Sends the token along with a password reset link to the user’s email.  
4. The user opens their email and clicks the reset link.  
5. The system:  
   - Validates that the token is authentic, unexpired, and tied to that user.  
   - Displays the password reset form.  
6. The user enters a new password.  
7. The system:  
   - Verifies that the new password complies with the password policy (length, special characters, etc.).  
   - Hashes the new password and stores it in the database.  
   - Displays “Password successfully changed” and redirects the user to the login page.

**Alternative Flows:**

- **If the email is not found:**  
  - Display “No account found with that email.”  
  - Allow the user to retry or to register instead.

- **If an email delivery error occurs:**  
  - Enable the user to resend the recovery email.

- **If the token is invalid or expired:**  
  - Display “Reset link invalid or expired.”  
  - Return the user to the password recovery page to request a new link.

- **If the new password does not meet policy requirements:**  
  - Display the specific error (e.g. “Password must be at least 8 characters”).  
  - Allow the user to correct and resubmit.

**Postconditions:**

- The user’s password has been replaced with the new value.  
- The recovery token has been invalidated and can no longer be reused.

### Update User Info

**Preconditions:**

- The user must be successfully logged in.  
- The profile edit page or form must be accessible.  
- The user’s current data must be loaded into the form.

**Main Scenario:**

1. The user navigates to the profile edit page/form.  
2. The system displays the current information (name, email, phone, avatar, etc.) in the form.  
3. The user modifies the desired fields (for example, name or phone number).  
4. The user clicks “Save Changes.”  
5. The system:  
   - Validates that any new email (if changed) is unique.  
   - Verifies that the phone number is in the correct format.  
   - Ensures that any uploaded avatar meets allowed type and size.  
   - Confirms that all submitted data is valid.  
6. On successful validation:  
   - The system updates the user record in the database.  
   - Displays “Profile updated successfully.”  
   - The form reloads with the updated data.

**Alternative Flows:**

- **If any input is invalid:**  
  - Display the corresponding error message.  
  - Allow the user to correct and resubmit.

- **If a database error occurs:**  
  - Display a system error message.  
  - Prompt the user to try again.

**Postconditions:**

- The user’s updated information has been saved in the database.  
- If the email was changed and requires re-verification, the user is marked as unverified until they confirm again.

### Sign Out

**Preconditions:**

- The user must already be successfully logged in.

**Main scenario:**

1. The user clicks “Sign Out.”  
2. The system:  
   - Removes the token from the client.  
   - Deletes any HttpOnly cookies associated with the user’s session.  
   - Displays “You have successfully signed out” (or redirects directly to the login page).  
   - The user is taken to the public page (e.g., home or login).

**Postconditions:**

- The user’s session is fully terminated.  
- The user can no longer access protected areas without logging in again.

### Delete User

**Preconditions:**

- The user must already be successfully logged in.  
- All of the user’s active loans and reservations must be settled.

**Main scenario:**

1. The user navigates to the “Delete Account” page.  
2. The system displays a list of important notices (e.g., loss of reservation history, cancellation of future reservations).  
3. The user enters their password to confirm their identity.  
4. The system validates that the entered password is correct.  
5. The user clicks the “Confirm Delete” button.  
6. The system:  
   - Marks the user record as deleted (soft-delete) or removes it from the database.  
   - Invalidates all of the user’s active tokens.  
   - Sends an account-deletion confirmation email to the user.  
   - Displays “Your account has been successfully deleted” and redirects the user to the public page (home or sign-up).

**Alternative flow:**

- **If the user has any outstanding loans or reservations:**  
  - The system displays “Please return all borrowed books or cancel future reservations first.”  
  - The deletion process is aborted.

- **If the password is incorrect:**  
  - The system displays “The entered password is incorrect.”  
  - The user may try again.

- **If a database or server error occurs:**  
  - The system displays “Internal error; please try again.”

**Postconditions:**

- The user’s account is deleted or deactivated.  
- The user’s tokens are invalidated.

### Users List (Manager)

**Preconditions:**

- The user must be logged in as a manager or admin.

**Main scenario:**

1. The manager navigates to the “Users List” page.  
2. The system executes a query to retrieve users.  
3. The system displays the results along with controls:  
   - Search field (name, email)  
   - Status filter (Active, Suspended, Deleted)  
   - Pagination  

**Alternative flows:**

- **No users are found:** display “No users exist.”  
- **A system error occurs:** display “Internal error; please try again.”

**Postconditions:**

- Data is only read; no changes are made.

### User Detail (Manager)

**Preconditions:**

- The user must be logged in as a manager or admin.  
- The target user must exist in the system.

**Main scenario:**

1. The manager clicks on “User Details” for the target user.  
2. The system retrieves that user’s information from the database based on their ID.  
3. The system displays the information to the manager in a form.

**Alternative flows:**

- **If the user has been soft-deleted:** display “User not found.”  
- **If a database error occurs:** display “Internal error; please try again.”

**Postconditions:**

- No data is modified; it is only read.

### Ban User (Manager)

**Preconditions:**

- The user must be logged in as a manager.  
- The target user must not have been hard-deleted.

**Main scenario:**

1. The manager clicks “Delete/Ban.”  
2. The system shows a confirmation dialog: “Are you sure you want to delete or ban this user?”  
3. The manager chooses either “Soft-delete” or “Suspend.”  
4. The system executes the corresponding query.  
5. The system displays “Operation completed successfully.”  
6. The users list is refreshed to show the updated status.

**Alternative flows:**

- **If the manager cancels:** return to the list without changes.  
- **If the user is already soft-deleted:** disable the “Delete” option and show only “Unban” or “Restore.”  
- **If a query error occurs:** display “Internal error; please try again.”

**Postconditions:**

- The user’s status field in the users list is updated.  
