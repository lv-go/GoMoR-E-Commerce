
import { Outlet, redirect } from 'react-router';
import { getCurrentUser } from '~/firebase-config';
import AdminMenu from '../../components/Admin/AdminMenu';

export async function clientLoader() {
  const user = await getCurrentUser();
  if (!user) {
    return redirect("/login");
  }
  if (user.role !== "user") {
    return redirect("/");
  }
  return null;
}

export default function Layout() {
  return (
    <>
      <Outlet />
    </>
  )
}
