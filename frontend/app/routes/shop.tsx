import { useState } from "react";

import Loader from "~/components/Loader";
import ProductCard from "~/components/Products/ProductCard";
import { useGetPage } from "~/hooks/categories";
import { useGetBrands, useGetPage as useGetProductsPage } from "~/hooks/products";
import type { Route } from "./+types/shop";

export function meta({ }: Route.MetaArgs) {
  return [
    { title: "GoMoR-E-Commerce - Shop" },
    { name: "description", content: "Shop for the latest trends in fashion and accessories." },
  ];
}

export default function Shop() {
  const { data: categoriesPage, isLoading: isLoadingCategories } = useGetPage();
  const categories = categoriesPage?.items || [];

  const [categoriesChecked, setCategoriesChecked] = useState<string[]>([]);
  const [brandChecked, setBrandChecked] = useState<string>("");
  const [priceFilter, setPriceFilter] = useState("");

  const { data: productsPage, isLoading: isLoadingProducts } = useGetProductsPage({
    category: categoriesChecked,
    brand: brandChecked,
    price: priceFilter,
  });

  const products = productsPage?.items || [];

  const handleBrandClick = (brand: string) => {
    setBrandChecked(brand);
  };

  const handleCheck = (value: boolean, id: string) => {
    const updatedChecked = value
      ? [...categoriesChecked, id]
      : categoriesChecked.filter((c) => c !== id);
    setCategoriesChecked(updatedChecked);
  };

  const handleReset = () => {
    setCategoriesChecked([]);
    setBrandChecked("");
    setPriceFilter("");
  };

  const { data: brands, isLoading: isLoadingBrands } = useGetBrands();

  const handlePriceChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    // Update the price filter state when the user types in the input filed
    setPriceFilter(e.target.value);
  };

  return (
    <div className="container mx-auto">
      <div className="flex md:flex-row">
        <div className="bg-base-200 p-3 mt-2 mb-2">
          <h2 className="h4 text-center py-2 bg-base-300 rounded-full mb-2">
            Filter by Categories
          </h2>

          <div className="p-5 w-[15rem]">
            {categories?.map((c) => (
              <div key={c._id} className="mb-2">
                <div className="flex ietms-center mr-4">
                  <input
                    type="checkbox"
                    id={`category-checkbox-${c._id}`}
                    onChange={(e) => handleCheck(e.target.checked, c._id || "")}
                    checked={categoriesChecked.includes(c._id || "")}
                    className="checkbox"
                  />

                  <label
                    htmlFor={`category-checkbox-${c._id}`}
                    className="ml-2 text-sm font-medium"
                  >
                    {c.name}
                  </label>
                </div>
              </div>
            ))}
          </div>

          <h2 className="h4 text-center py-2 bg-base-300 rounded-full mb-2">
            Filter by Brands
          </h2>

          <div className="p-5">
            {brands?.map((brand, index) => (
              <div key={index} className="flex items-enter mr-4 mb-5">
                <input
                  type="radio"
                  id={`brand-radio-${index}`}
                  name="brand"
                  checked={brandChecked === brand}
                  onChange={() => handleBrandClick(brand)}
                  className="radio"
                />

                <label
                  htmlFor={`brand-radio-${index}`}
                  className="ml-2 text-sm font-medium"
                >
                  {brand}
                </label>
              </div>
            ))}
          </div>

          <h2 className="h4 text-center py-2 bg-base-300 rounded-full mb-2">
            Filer by Price
          </h2>

          <div className="p-5 w-[15rem]">
            <input
              type="text"
              placeholder="Enter Price"
              value={priceFilter}
              onChange={handlePriceChange}
              className="input input-bordered w-full"
            />
          </div>

          <div className="p-5 pt-0">
            <button
              className="btn btn-primary w-full my-4"
              onClick={handleReset}
            >
              Reset
            </button>
          </div>
        </div>

        <div className="p-3">
          <h2 className="h4 text-center mb-2">{products?.length} Products</h2>
          <div className="flex flex-wrap">
            {isLoadingProducts ? (
              <Loader />
            ) : (
              products?.map((p) => (
                <div className="p-3" key={p._id}>
                  <ProductCard p={p} />
                </div>
              ))
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
