# Book Inventory System 
## Using GIN Framework and GORM
- Added dockerization with postgres
- Added Elasticsearch sync with Postgres using logstash
- Added Kibana for inspection
- Added swagger specification
- Added basic api for books and author
- Added bearer jwt auth
## Services
- book service
  - user
    - Get Book
    - Post Book
    - Get Paginated Books
    - Update Book
    - Delete Book
    - Get Author
  - admin
    - Get Book
    - Post Book
    - Get Paginated Books
    - Update Book
    - Delete Book
    - Post Author by admin
    - Get Author by admin
    - Get Paginated Authors by admin
- user service
  - JWT Authentication
    - Register
    - Log In
    - User Verification with Resend
    - Password Change with Resend
    - Access Token and Refresh Token
    - Refresh Token Automatic Reuse Detection
    - User Verification and Password Change Multiple Token-link Detection
 - elastic search
   - Indexing book and author
   - Get paginated Authors by name search(fuzzy search in first_name, last_name) (only for user)
   - Get paginated Books by name search(fuzzy search in title, author's first_name, author's last_name) (only for user)

    
