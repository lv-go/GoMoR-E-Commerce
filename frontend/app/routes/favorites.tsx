import Product from "../components/Products/Product";
import type { Route } from "./+types/favorites";

export async function clientLoader() {
  return JSON.parse(sessionStorage.getItem("favorite_products") || "[]")
}

export default function Favorites({ loaderData: favorites }: Route.ComponentProps) {
  return (
    <div className="ml-[10rem]">
      <h1 className="text-lg font-bold ml-[3rem] mt-[3rem]">
        FAVORITE PRODUCTS
      </h1>

      <div className="flex flex-wrap">
        {favorites.map((product) => (
          <Product key={product._id} product={product} />
        ))}
      </div>
    </div>
  );
}
