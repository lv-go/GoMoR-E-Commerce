import { useMutation, useQuery } from "@tanstack/react-query";
import type { Cart, CartItem } from "~/schemas/cart";
import type { Product } from "~/schemas/product";
import { newCart } from "~/schemas/cart";
import type { ShippingAddress } from "~/schemas/order";
import { queryClient } from "~/utils/query-client";

function getCart(): Cart {
  const cartJSON = sessionStorage.getItem("cart");
  if (!cartJSON) {
    return newCart();
  }
  return JSON.parse(cartJSON);
};

function setCart(cart: Cart) {
  sessionStorage.setItem("cart", JSON.stringify(cart));
  queryClient.invalidateQueries({ queryKey: ["cart"] });
};

export const useGetCart = () => {
  return useQuery({
    queryKey: ["cart"],
    queryFn: () => getCart(),
  });
}

export const useAddToCart = () => {
  return useMutation({
    mutationFn: async ({ product, quantity }: { product: Product, quantity: number }) => {
      const cart = getCart();
      if (!cart.cartItems.some((p) => p._id === product._id)) {
        cart.cartItems.push({
          _id: product._id || "",
          name: product.name,
          price: product.price,
          image: product.image,
          countInStock: product.countInStock,
          brand: product.brand,
          quantity: quantity,
        });
        setCart(cart);
      }
    }
  });
}

export const useRemoveFromCart = () => {
  return useMutation({
    mutationFn: async (productId: string) => {
      const cart = getCart();
      const updateCartItems = cart.cartItems.filter(p => p._id !== productId);
      cart.cartItems = updateCartItems;
      setCart(cart);
    }
  });
}

export const useUpdateCart = () => {
  return useMutation({
    mutationFn: async ({ _id, quantity }: { _id: string, quantity: number }) => {
      const cart = getCart();
      const updateCartItems = cart.cartItems
        .map(item => item._id === _id ? { ...item, quantity } : item);
      cart.cartItems = updateCartItems;
      setCart(cart);
    }
  });
}

export const useClearCart = () => {
  return useMutation({
    mutationFn: async () => {
      setCart(newCart());
    }
  });
}

export const useSaveShippingAddress = () => {
  return useMutation({
    mutationFn: async (shippingAddress: ShippingAddress) => {
      const cart = getCart();
      cart.shippingAddress = shippingAddress;
      setCart(cart);
    }
  });
}

export const useSavePaymentMethod = () => {
  return useMutation({
    mutationFn: async (paymentMethod: string) => {
      const cart = getCart();
      cart.paymentMethod = paymentMethod;
      setCart(cart);
    }
  });
}
