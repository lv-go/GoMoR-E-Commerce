import { FaHeart, FaRegHeart } from "react-icons/fa";
import { useAddFavorite, useGetIsFavoriteById, useRemoveFavorite } from "~/hooks/favorites";

import { type Product } from "~/schemas/product";

export default function HeartIcon({ product }: { product: Product }) {
  const { data: isFavorite = false } = useGetIsFavoriteById(product._id!);
  const { mutate: addFavorite } = useAddFavorite();
  const { mutate: removeFavorite } = useRemoveFavorite();



  const toggleFavorites = () => {
    if (isFavorite) {
      removeFavorite(product._id!);
    } else {
      addFavorite(product);
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
        <FaRegHeart className="text-black" />
      )}
    </div>
  );
}
