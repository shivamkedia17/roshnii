import { useLogout } from "@/hooks/useAuth";

export function TempLogout() {
  const logout = useLogout();

  return (
    <>
      <div>
        <h1>Logged In!</h1>
        <button onClick={() => logout()}>Log Out</button>
      </div>
      ;
    </>
  );
}
