import { useEffect } from "react";
import { useAuthStore } from "~/lib/auth";

export function AuthBootstrap() {
  const refreshSession = useAuthStore((state) => state.refreshSession);

  useEffect(() => {
    refreshSession();
  }, [refreshSession]);

  return null;
}
