import { useLogin } from "@/hooks/useAuth";

export function TempLogin() {
  const login = useLogin();

  return (
    <>
      <button onClick={() => login()}>Login With Google</button>;
    </>
  );
}
