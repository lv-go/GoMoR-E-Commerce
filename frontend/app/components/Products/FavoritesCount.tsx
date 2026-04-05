import { useSelector } from "react-redux";

export default function FavoritesCount({ className }: { className?: string }) {
  // const favorites = useSelector((state: any) => state.favorites);
  // const favoriteCount = favorites.length;
  const favoriteCount = 5;

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
