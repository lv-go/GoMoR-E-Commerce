import type { Cart, CartItem } from "~/schemas/cart";
import type { ShippingAddress } from "~/schemas/order";

export const toFixedDecimals = (num: string) => {
  return (Math.round(Number(num) * 100) / 100).toFixed(2);
};

export const updateCart = (state: Cart) => {
  // Calculate the items price
  state.itemsPrice = state.cartItems.reduce((acc: number, item: CartItem) => acc + item.price * item.quantity, 0);

  // Calculate the shipping price
  state.shippingPrice = state.itemsPrice > 100 ? 0 : 10;

  // Calculate the tax price
  state.taxPrice = (0.15 * state.itemsPrice);

  // Calculate the total price
  state.totalPrice = (
    Number(state.itemsPrice) +
    Number(state.shippingPrice) +
    Number(state.taxPrice)
  );

  // Save the cart to localStorage
  localStorage.setItem("cart", JSON.stringify(state));

  return state;
};

export function getCart(): Cart {
  return JSON.parse(localStorage.getItem("cart") || "{}");
};

export function addCartItem(item: CartItem) {
  const cart = getCart();
  const existItem = cart.cartItems.find((x) => x._id === item._id);

  if (existItem) {
    cart.cartItems = cart.cartItems.map((x) =>
      x._id === existItem._id ? item : x
    );
  } else {
    cart.cartItems = [...cart.cartItems, item];
  }

  return updateCart(cart);
}

export function updateCartItemQty(id: string, quantity: number) {
  const cart = getCart();
  cart.cartItems = cart.cartItems.map((x) =>
    x._id === id ? { ...x, quantity } : x
  );
  return updateCart(cart);
}

export function removeFromCart(id: string) {
  const cart = getCart();
  cart.cartItems = cart.cartItems.filter((x) => x._id !== id);
  return updateCart(cart);
}

export function clearCartItems() {
  const cart = getCart();
  cart.cartItems = [];
  return updateCart(cart);
}

export function saveShippingAddress(data: ShippingAddress) {
  const cart = getCart();
  cart.shippingAddress = data;
  return updateCart(cart);
}

export function savePaymentMethod(paymentMethod: string) {
  const cart = getCart();
  cart.paymentMethod = paymentMethod;
  return updateCart(cart);
}
