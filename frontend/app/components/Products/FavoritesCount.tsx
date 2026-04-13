import { useGetFavoritesCount } from "~/hooks/favorites";

export default function FavoritesCount({ className }: { className?: string }) {
  const { data: favoritesCount = 0 } = useGetFavoritesCount();

  return (
    <div className={className}>
      {favoritesCount > 0 && (
        <span className="px-1 py-0 text-sm text-white bg-pink-500 rounded-full">
          {favoritesCount}
        </span>
      )}
    </div>
  );
}
