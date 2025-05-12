import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import App from "@/components/App";
import { Provider } from "@/context/Provider";

import "@/index.css";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <Provider>
      <App />
    </Provider>
  </StrictMode>,
);
