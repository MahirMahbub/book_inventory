# Book Inventory System 
## Using GIN Framework and GORM
- Added dockerization with postgres
- Added Elasticsearch sync with Postgres using logstash
- Added Kibana for inspection
- Added swagger specification
- Added basic api for books and author
- Added bearer jwt auth from scratch
- Refresh Token Automatic Reuse Detection
- User Verification and Password Change Multiple Token-link Detection
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
    - Post Author
    - Get Author
    - Get Paginated Authors
    - Update Author
    - Delete Author
    - Create Admin
- user service
  - JWT Authentication
    - Register
    - Log In
    - User Verification with Resend and Old Token Detection on Resend and Expiry.
    - Password Change with Resend and Old Token Detection on Resend and Expiry.
    - Access Token and Refresh Token
- elastic search
   - Indexing book and author
   - Get paginated Authors by name search(fuzzy search in first_name, last_name) (only for user)
   - Get paginated Books by name search(fuzzy search in title, author's first_name, author's last_name) (only for user)

- Postgres to ELK Sync


                                  +---------------v------------------+
                                  |                                  |
                                  |              data                |
                                  |                                  |
                                  +---------------+------------------+
                                                  |
                                                  |
                                                  |
                                                  |
                                           +------v------+
                                           |             |
                                           |  PostgreSQL |
                                           |             |
                                           +------^------+
                                                  |
                                                  |
                                                  |
                                  +---------------v------------------+
                                  |                                  |
                                  |            Logstash              |
                                  |                                  |
                                  +---------------+------------------+
                                                  |
                                                  |
                                                  |
                                                  |
                                          +-------v--------+
                                          |                |
                                          | Elasticsearch  |
                                          |                |
                                          +----------------+
    
