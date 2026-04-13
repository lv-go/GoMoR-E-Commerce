import { useMutation, useQuery } from "@tanstack/react-query";
import type { Product } from "~/schemas/product";
import { queryClient } from "~/utils/query-client";

// Retrive favorites from a sessionStorage
const getFavoritesFromSessionStorage = (): Product[] => {
  const favoritesJSON = sessionStorage.getItem("favorites") || "[]";
  return JSON.parse(favoritesJSON);
};

export const useGetAllFavorites = () => {
  return useQuery({
    queryKey: ["favorites"],
    queryFn: getFavoritesFromSessionStorage,
  });
}

export const useGetIsFavoriteById = (productId: string) => {
  return useQuery({
    queryKey: ["favorites", productId],
    queryFn: () => getFavoritesFromSessionStorage().some((p: Product) => p._id === productId),
  });
}

export const useGetFavoritesCount = () => {
  return useQuery({
    queryKey: ["favorites-count"],
    queryFn: () => getFavoritesFromSessionStorage().length,
  });
}

export const useAddFavorite = () => {
  return useMutation({
    mutationFn: async (product: Product) => {
      const favorites = getFavoritesFromSessionStorage();
      if (!favorites.some((p: Product) => p._id === product._id)) {
        favorites.push(product);
        sessionStorage.setItem("favorites", JSON.stringify(favorites));
      }
      queryClient.invalidateQueries({ queryKey: ["favorites"] });
      queryClient.invalidateQueries({ queryKey: ["favorites-count"] });
    },
  });
}

export const useRemoveFavorite = () => {
  return useMutation({
    mutationFn: async (productId: string) => {
      const favorites = getFavoritesFromSessionStorage();
      const updateFavorites = favorites.filter(
        (product: Product) => product._id !== productId
      );
      sessionStorage.setItem("favorites", JSON.stringify(updateFavorites));
      queryClient.invalidateQueries({ queryKey: ["favorites"] });
      queryClient.invalidateQueries({ queryKey: ["favorites-count"] });
    },
  });
}

export const useClearFavorites = () => {
  return useMutation({
    mutationFn: async () => {
      sessionStorage.removeItem("favorites");
      queryClient.invalidateQueries({ queryKey: ["favorites"] });
      queryClient.invalidateQueries({ queryKey: ["favorites-count"] });
    },
  });
}
