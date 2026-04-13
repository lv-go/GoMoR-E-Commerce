import { getFavoritesFromLocalStorage } from "~/utils/favorites";

export default function FavoritesCount({ className }: { className?: string }) {
  const favorites = getFavoritesFromLocalStorage();
  const favoriteCount = favorites.length;

  return (
    <div className={className}>
      {favoriteCount > 0 && (
        <span className="px-1 py-0 text-sm text-white bg-pink-500 rounded-full">
          {favoriteCount}
        </span>
      )}
    </div>
  );
}
