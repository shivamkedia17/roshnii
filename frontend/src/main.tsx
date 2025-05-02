import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "@/index.css";
import App from "@/components/App";
import { AuthProvider } from "@/context/AuthContext";
import { PhotoProvider } from "@/context/PhotoContext";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <AuthProvider>
      <PhotoProvider>
        <App />
      </PhotoProvider>
    </AuthProvider>
  </StrictMode>,
);
