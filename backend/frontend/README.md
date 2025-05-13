# Local Development

Start the frontend dev server:

```sh
bun run dev
```

Serving the webapp on the same domain as the backend:

- First, build the web app using `bun run dbuild`.
- This will write files to `frontend/dist`
- `frontend/dist/index.html` is the entry point of the web app.
- Make sure backend is configured to serve static assets (CSS, JS, favicons, etc.) from the `frontend/dist/assets` directory.

NOTE: To avoid browser's cookie issues please:

- Visit `http://127.0.0.1:8080`, NOT `http://localhost:8080` - the browser treats these as two separate origins and does not send the cookies set by `127.0.0.1` to `localhost`.
