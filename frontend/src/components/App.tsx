import { useAuthContext } from "@/context/AuthContext";
import "@/css/App.css";
import { Loading } from "./common/Loading";

export default function App() {
  const { isAuthenticated, isLoading, login, logout } = useAuthContext();

  return (
    <>
      {isLoading ? (
        <Loading />
      ) : isAuthenticated ? (
        <div>
          <h1>Logged In!</h1>
          <button onClick={logout}>Log Out</button>
        </div>
      ) : (
        <button onClick={login}>Login With Google</button>
      )}
    </>
  );
}
