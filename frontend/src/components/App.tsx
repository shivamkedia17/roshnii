// import { useEffect, useState } from "react";
import "../css/App.css";
// import { Login } from "./auth/Login";
import { MainLayout } from "./layout/MainLayout";
// import { useAuth } from "../context/AuthContext";

function App() {
  // const { isAuthenticated, isLoading } = useAuth();

  // if (isLoading) {
  //   return <div className="loading">Loading...</div>;
  // }

  // return <>{isAuthenticated ? <MainLayout /> : <Login />}</>;

  return (
    <>
      <MainLayout />;
    </>
  );
}

export default App;
