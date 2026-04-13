import type { Product } from "~/schemas/product";
import ProductComponent from "../components/Products/ProductInfo";
import type { Route } from "./+types/favorites";
import { useGetAllFavorites } from "~/hooks/favorites";

export function meta({ }: Route.MetaArgs) {
  return [
    { title: "Favorites | GoMoR-E-Commerce" },
    { name: "description", content: "Favorites | GoMoR-E-Commerce" },
  ];
}

export default function Favorites() {
  const { data: favorites = [] } = useGetAllFavorites();
  return (
    <div className="ml-[10rem]">
      <h1 className="text-lg font-bold ml-[3rem] mt-[3rem]">
        FAVORITE PRODUCTS
      </h1>

      <div className="flex flex-wrap">
        {favorites.map((product) => (
          <ProductComponent key={product._id} product={product} />
        ))}
      </div>
    </div>
  );
}
