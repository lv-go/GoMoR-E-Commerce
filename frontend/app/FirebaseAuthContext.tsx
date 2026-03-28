import * as React from "react";
import { auth, type UserWithRole } from "./firebase-config";
import { createContext, useContext, useState, type ReactNode } from "react";

type ContextState = { user: UserWithRole | null };

const FirebaseAuthContext =
  createContext<ContextState | undefined>(undefined);

export function FirebaseAuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<UserWithRole | null>(null);
  const value = { user };

  React.useEffect(() => auth.onAuthStateChanged((user) => {
    if (!user) {
      setUser(null);
      return;
    }
    user.getIdTokenResult().then((token) => {
      setUser({
        ...user,
        role: token.claims?.['role'] as string,
      });
    });
  }), []);

  return (
    <FirebaseAuthContext.Provider value={value}>
      {children}
    </FirebaseAuthContext.Provider>
  );
};

export function useFirebaseAuth() {
  const context = useContext(FirebaseAuthContext);
  if (context === undefined) {
    throw new Error("useFirebaseAuth must be used within FirebaseAuthProvider");
  }
  return context;
}
