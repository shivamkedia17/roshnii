import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import App from "@/components/App";
import { Provider } from "@/context/Provider";
import { AuthProvider } from "@/context/AuthContext";
import "@/index.css";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <Provider>
      <AuthProvider>
        <App />
      </AuthProvider>
    </Provider>
  </StrictMode>,
);
