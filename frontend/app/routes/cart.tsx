import { FaTrash } from "react-icons/fa";
import { Link, useRevalidator } from "react-router";
import type { Cart, CartItem } from "~/schemas/cart";
import type { Product } from "~/schemas/product";
import { getCart, removeFromCart, updateCartItemQty } from "~/utils/cartUtils";
import type { Route } from "./+types/cart";

export function meta({ }: Route.MetaArgs) {
  return [
    { title: "GoMoR-E-Commerce - Cart" },
    { name: "description", content: "Your shopping cart." },
  ];
}

export async function clientLoader(): Promise<Cart> {
  return getCart();
}

export default function Cart({ loaderData: cart }: Route.ComponentProps) {
  const { cartItems = [] } = cart;
  const { revalidate } = useRevalidator();

  const addToCartHandler = (id: string, qty: number) => {
    updateCartItemQty(id, qty);
    revalidate();
  };

  const removeFromCartHandler = (id: string) => {
    removeFromCart(id);
    revalidate();
  };

  return (
    <>
      <div className="container flex justify-around items-start wrap mx-auto mt-8">
        {cartItems.length === 0 ? (
          <div>
            Your cart is empty <Link to="/shop">Go To Shop</Link>
          </div>
        ) : (
          <>
            <div className="flex flex-col w-[80%]">
              <h1 className="text-2xl font-semibold mb-4">Shopping Cart</h1>

              {cartItems.map((item: CartItem) => (
                <div key={item._id} className="flex items-center mb-[1rem] pb-2">
                  <div className="w-[5rem] h-[5rem]">
                    <img
                      src={item.image}
                      alt={item.name}
                      className="w-full h-full object-cover rounded"
                    />
                  </div>

                  <div className="flex-1 ml-4">
                    <Link to={`/product/${item._id}`} className="text-pink-500">
                      {item.name}
                    </Link>

                    <div className="mt-2 text-white">{item.brand}</div>
                    <div className="mt-2 text-white font-bold">
                      $ {item.price}
                    </div>
                  </div>

                  <div className="w-24">
                    <select
                      className="w-full p-1 border rounded text-black"
                      value={item.quantity}
                      onChange={(e) =>
                        addToCartHandler(item._id, Number(e.target.value))
                      }
                    >
                      {[...Array(item.countInStock).keys()].map((x) => (
                        <option key={x + 1} value={x + 1}>
                          {x + 1}
                        </option>
                      ))}
                    </select>
                  </div>

                  <div>
                    <button
                      className="text-red-500 mr-[5rem]"
                      onClick={() => removeFromCartHandler(item._id)}
                    >
                      <FaTrash className="ml-[1rem] mt-[.5rem]" />
                    </button>
                  </div>
                </div>
              ))}

              <div className="mt-8 w-[40rem]">
                <div className="p-4 rounded-lg">
                  <h2 className="text-xl font-semibold mb-2">
                    Items ({cartItems.reduce((acc: number, item: CartItem) => acc + item.quantity, 0)})
                  </h2>

                  <div className="text-2xl font-bold">
                    ${" "}
                    {cart.itemsPrice.toFixed(2)}
                  </div>

                  <Link
                    to="/shipping"
                    className="btn btn-primary mt-4 rounded-full text-lg w-full"
                  >
                    Proceed To Checkout
                  </Link>
                </div>
              </div>
            </div>
          </>
        )}
      </div>
    </>
  );
}
