import { useState } from "react";
import {
  AiOutlineHome,
  AiOutlineLogin,
  AiOutlineShopping,
  AiOutlineShoppingCart,
  AiOutlineUserAdd,
} from "react-icons/ai";
import { FaBox, FaBriefcase, FaChartBar, FaHeart, FaSignOutAlt, FaTags, FaUser, FaUsers } from "react-icons/fa";
import { Link, NavLink, useNavigate } from "react-router";
import { useFirebaseAuth } from "~/FirebaseAuthContext";
import { auth } from "~/firebase-config";
import { useClearCart, useGetCart } from "~/hooks/cart";
import { newCart } from "~/schemas/cart";
import FavoritesCount from "./Products/FavoritesCount";
import { useClearFavorites } from "~/hooks/favorites";

export default function Navigation() {
  const { user } = useFirebaseAuth();
  const { data: cart = newCart() } = useGetCart();
  const cartItems = cart.cartItems;
  const { mutate: clearCart } = useClearCart();
  const { mutate: clearFavorites } = useClearFavorites();

  const navigate = useNavigate();

  const logoutHandler = async () => {
    try {
      await auth.signOut();
      clearCart();
      clearFavorites();
      navigate("/login");
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <ul className="menu w-full grow">
      <li>
        <NavLink
          to="/"
          className="is-drawer-close:tooltip is-drawer-close:tooltip-right"
          data-tip="Home"
        >
          <AiOutlineHome className="my-1.5 inline-block size-4" />
          <span className="is-drawer-close:hidden">HOME</span>
        </NavLink>
      </li>

      <li>
        <Link
          to="/shop"
          className="is-drawer-close:tooltip is-drawer-close:tooltip-right"
          data-tip="Shop"
        >
          <AiOutlineShopping className="my-1.5 inline-block size-4" />
          <span className="is-drawer-close:hidden">SHOP</span>
        </Link>
      </li>


      <li>
        <Link to="/cart" data-tip="Cart" className="is-drawer-close:tooltip is-drawer-close:tooltip-right">
          <AiOutlineShoppingCart className="my-1.5 inline-block size-4" />
          <span className="is-drawer-close:hidden">Cart</span>

          <div className="absolute top-0 right-1">
            {cartItems.length > 0 && (
              <span>
                <span className="px-1 py-0 text-sm text-white bg-pink-500 rounded-full">
                  {cartItems.reduce((a: number, c) => a + c.quantity, 0)}
                </span>
              </span>
            )}
          </div>
        </Link>
      </li>

      <li>
        <Link to="/favorites" data-tip="Favorites" className="is-drawer-close:tooltip is-drawer-close:tooltip-right">
          <div className="my-1.5 inline-block size-4">
            <FaHeart />
          </div>
          <span className="is-drawer-close:hidden">Favorites</span>
          <FavoritesCount className="absolute top-0 right-1" />
        </Link>
      </li>
      {user && (
        <>
          {user.role === "admin" && (
            <>
              <li>
                <Link
                  to="/admin"
                  className="is-drawer-close:tooltip is-drawer-close:tooltip-right"
                  data-tip="Dashboard"
                >
                  <FaChartBar className="my-1.5 inline-block size-4" />
                  <span className="is-drawer-close:hidden">Dashboard</span>
                </Link>
              </li>
              <li>
                <Link
                  to="/admin/products"
                  className="is-drawer-close:tooltip is-drawer-close:tooltip-right"
                  data-tip="Products"
                >
                  <FaBox className="my-1.5 inline-block size-4" />
                  <span className="is-drawer-close:hidden">Products</span>
                </Link>
              </li>
              <li>
                <Link
                  to="/admin/categories"
                  className="is-drawer-close:tooltip is-drawer-close:tooltip-right"
                  data-tip="Categories"
                >
                  <FaTags className="my-1.5 inline-block size-4" />
                  <span className="is-drawer-close:hidden">Categories</span>
                </Link>
              </li>
              <li>
                <Link
                  to="/admin/orders"
                  className="is-drawer-close:tooltip is-drawer-close:tooltip-right"
                  data-tip="Orders"
                >
                  <FaBriefcase className="my-1.5 inline-block size-4" />
                  <span className="is-drawer-close:hidden">Orders</span>
                </Link>
              </li>
              <li>
                <Link
                  to="/admin/users"
                  className="is-drawer-close:tooltip is-drawer-close:tooltip-right"
                  data-tip="Users"
                >
                  <FaUsers className="my-1.5 inline-block size-4" />
                  <span className="is-drawer-close:hidden">Users</span>
                </Link>
              </li>
            </>
          )}

          <li>
            <Link to="/profile" className="is-drawer-close:tooltip is-drawer-close:tooltip-right" data-tip="Profile">
              <FaUser className="my-1.5 inline-block size-4" />
              <span className="is-drawer-close:hidden">Profile</span>
            </Link>
          </li>
          <li>
            <button
              onClick={logoutHandler}
              className="is-drawer-close:tooltip is-drawer-close:tooltip-right"
              data-tip="Logout"
            >
              <FaSignOutAlt className="my-1.5 inline-block size-4" />
              <span className="is-drawer-close:hidden">Logout</span>
            </button>
          </li>
        </>
      )
      }
      {
        !user && (
          <>
            <li>
              <Link
                to="/login"
                className="is-drawer-close:tooltip is-drawer-close:tooltip-right"
                data-tip="Login"
              >
                <AiOutlineLogin className="my-1.5 inline-block size-4" />
                <span className="is-drawer-close:hidden">LOGIN</span>
              </Link>
            </li>
            <li>
              <Link
                to="/register"
                className="is-drawer-close:tooltip is-drawer-close:tooltip-right"
                data-tip="Register"
              >
                <AiOutlineUserAdd className="my-1.5 inline-block size-4" />
                <span className="is-drawer-close:hidden">REGISTER</span>
              </Link>
            </li>
          </>
        )
      }
      {/* <li>
        <button
          onClick={toggleDropdown}
          className="flex items-center text-gray-800 focus:outline-none"
        >
          {user ? (
            <span className="text-white">{user.displayName}</span>
          ) : (
            <></>
          )}
          {user && (
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className={`h-4 w-4 ml-1 ${dropdownOpen ? "transform rotate-180" : ""
                }`}
              fill="none"
              viewBox="0 0 24 24"
              stroke="white"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d={dropdownOpen ? "M5 15l7-7 7 7" : "M19 9l-7 7-7-7"}
              />
            </svg>
          )}
        </button>

        {dropdownOpen && user && (
          <ul
            className={`absolute right-0 mt-2 mr-14 space-y-2 bg-white text-gray-600 ${!user.role ? "-top-20" : "-top-80"
              } `}
          >
             */}

    </ul >
  );
}
