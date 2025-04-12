
─────────────────────────────────────────────
Diagram Overview
─────────────────────────────────────────────

           User Browser (SPA)
                  │
       (1) Click “Sign in with Google”
                  │
                  ▼
         GET /api/auth/google/login
                  │
           [Backend Handler]
                  │  ──(2) Generate OAuth “state” and set as a cookie
                  │
          Redirect to Google’s OAuth URL
          ──────────────────────────────►
                  │                    Google’s Authorization Endpoint
                  │
                  │       (3) User logs in and approves access
                  │◄───────────────────────────────
                  │
         Google redirects with code and state
                  │
         GET /api/auth/google/callback?code=XYZ&state=ABC
                  │
           [Backend Handler]
                  │  ──(4) Verify state (compare with cookie)
                  │  ──(5) Exchange the code for an access token
                  │  ──(6) Call Google’s API to fetch user info
                  │  ──(7) Lookup/Create user in your database
                  │  ──(8) Generate internal JWT token
                  │  ──(9) Set JWT as an HTTP-only cookie and redirect
                  │
                  ▼
Redirect back to the SPA (e.g., GET /)
─────────────────────────────────────────────

─────────────────────────────────────────────
Mapping Instructions to the Flow
─────────────────────────────────────────────

1. Browser Interaction & Login Trigger
   • Frontend (components/auth/Login.tsx):
     – When the user clicks "Sign in with Google," the click handler redirects the browser to the backend endpoint (GET /api/auth/google/login).

2. Initiate OAuth Flow (Backend – OAuthLogin handler)
   • In internal/auth/google_oauth.go → OAuthLogin()
     – Generates a random “state” string (using generateState()).
     – Sets an HTTP-only cookie for “oauthstate” to save the generated state.
     – Uses the oauth2.Config to build the auth URL, and then redirects the request to Google’s Authorization Endpoint.

3. User Login & Authorization on Google
   • Google’s Authorization Endpoint:
     – The user logs into their Google account and grants permission.

4. Callback Handling (Backend – OAuthCallback handler)
   • In internal/auth/google_oauth.go → OAuthCallback()
     – Reads the “state” value returned in the query string and compares it with the value stored in the cookie for security.
     – Retrieves the authorization code from the query string.
     – Exchanges the code for an access token, using the oauth2.Config.Exchange() method.

5. Retrieving User Information from Google
   • In OAuthCallback() after the token exchange:
     – Uses the new access token (via oauth2.Config.Client()) to call Google’s user info endpoint (https://www.googleapis.com/oauth2/v2/userinfo).
     – Parses the JSON response into a GoogleUser struct.

6. Mapping User Information & Session Creation
   • Also in OAuthCallback():
     – Checks if the returned GoogleUser has a verified email.
     – Uses the configured UserStore interface to either find an existing user by email or create a new user record in your database (user creation logic as specified).
     – Generates a JWT token (via the JWTService) that maps the user's user_id and email.

7. Completing the Flow
   • The backend then sets the JWT as an HTTP-only cookie (or alternatively passes the token to the frontend via a URL fragment).
   • Finally, the backend redirects the user back to the SPA’s home route, completing the login process.

─────────────────────────────────────────────
Summary
─────────────────────────────────────────────

• The Authorization Code Flow is initiated by the SPA redirecting to the backend’s OAuthLogin endpoint.
• The backend handles security by generating a “state” cookie, redirecting to Google, and later verifying this state on return.
• On receiving the callback, the backend exchanges the received code for an access token, retrieves user info from Google, and ensures the login state is secure.
• Finally, the backend maps user information to your database, creates a JWT, sets it in an HTTP-only cookie, and redirects the user back to the SPA.
