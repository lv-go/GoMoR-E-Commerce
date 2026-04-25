import { initializeApp } from "firebase/app";
import { connectAuthEmulator, getAuth, type User } from "firebase/auth";

const firebaseConfig = {
  apiKey: import.meta.env.VITE_FIREBASE_API_KEY,
  authDomain: import.meta.env.VITE_FIREBASE_AUTH_DOMAIN,
  projectId: import.meta.env.VITE_FIREBASE_PROJECT_ID,
}

const app = initializeApp(firebaseConfig)
export const auth = getAuth(app)

if (process.env.NODE_ENV === 'development') {
  connectAuthEmulator(auth, "http://localhost:9099");
}

export type UserWithRole = User & {
  role: string;
};

export async function getCurrentUser(): Promise<UserWithRole | null> {
  return new Promise<User | null>((resolve) => {
    if (auth.currentUser) {
      resolve(auth.currentUser);
    }
    const unsubscribe = auth.onAuthStateChanged((user) => {
      unsubscribe();
      resolve(user);
    });
  }).then((user) => {
    if (!user) {
      return null;
    }
    return user.getIdTokenResult().then((token) => {
      return {
        ...user,
        role: token.claims?.['role'] as string,
      };
    });
  });
}
