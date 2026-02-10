# URL Shortener Backend (Go + PostgreSQL)

## Phase 1 – Go & HTTP Fundamentals
- Learned variable declaration, functions and struct
- Learned maps for the temporary data structure for storing URLs and short code
- Built the basic mind architecture for project.

## Phase 2 – URL Shortener Core Logic
- HTTP Server communication- GET and POST Request
- Learned different ways to connect to a network - for testing and real usage
- Decided on multiple routes, for POSTing a request and GETting a shorten url.
- Used JSON for sending the user request, learned coding and decoding JSON.
- Coded the redirection logic and added all the elementary error handling.

## Phase 3 – PostgreSQL Persistence with Docker
- Realized the limitation of maps, moved toward database for permanent storage.
- Created a docker supported machine for the database usage.
- Made docker PostgreSQL avaiable on the port and made the server listening to it.
- Modified the code to completely shift from maps to database service.


Current state:
- Fully working backend URL shortener
- Persistent storage via PostgreSQL
- Production-style architecture foundation
