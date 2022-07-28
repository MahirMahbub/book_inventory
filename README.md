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
  - Get Book
  - Post Book
  - Get paginated Books
  - Update Book
  - Delete Book
  - Post Author by admin
  - Get Author
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
   - Get paginated Authors by name search(fuzzy search)
   - Get paginated Books by name search(fuzzy search)

    
