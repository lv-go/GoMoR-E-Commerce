
import React from 'react'
import AdminMenu from '../../components/Admin/AdminMenu'
import { Outlet, redirect } from 'react-router'

export async function clientLoader() {
  const userInfo = localStorage.getItem("userInfo");
  if (!userInfo) {
    return redirect("/login");
  }
  return null;
}

export default function Layout() {
  return (
    <>
      <AdminMenu />
      <Outlet />
    </>
  )
}
