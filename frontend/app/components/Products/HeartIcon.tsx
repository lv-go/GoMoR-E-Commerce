import { FaHeart, FaRegHeart } from "react-icons/fa";

import { type Product } from "~/schemas/product";
import {
  addFavoriteToLocalStorage,
  getFavoritesFromLocalStorage,
  removeFavoriteFromLocalStorage,
} from "~/utils/favorites";

const HeartIcon = ({ product }: { product: Product }) => {
  const favorites = getFavoritesFromLocalStorage();
  const isFavorite = favorites.some((p: Product) => p._id === product._id);


  const toggleFavorites = () => {
    if (isFavorite) {
      removeFavoriteFromLocalStorage(product._id!);
    } else {
      addFavoriteToLocalStorage(product);
    }
  };

  return (
    <div
      className="absolute top-2 right-5 cursor-pointer"
      onClick={toggleFavorites}
    >
      {isFavorite ? (
        <FaHeart className="text-pink-500" />
      ) : (
        <FaRegHeart className="text-white" />
      )}
    </div>
  );
};

export default HeartIcon;
