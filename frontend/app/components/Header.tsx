import Loader from "./Loader";
import SmallProduct from "./Products/SmallProduct";
import ProductCarousel from "./Products/ProductCarousel";
import { useGetTopProducts } from "~/hooks/products";

const Header = () => {
  const { data, isLoading, error } = useGetTopProducts();
  const products = data?.items || [];

  if (isLoading) {
    return <Loader />;
  }

  if (error) {
    return <h1>ERROR</h1>;
  }

  return (
    <>
      <div className="flex justify-around">
        <div className="xl:block lg:hidden md:hidden:sm:hidden">
          <div className="grid grid-cols-2">
            {products.map((product) => (
              <div key={product._id}>
                <SmallProduct product={product} />
              </div>
            ))}
          </div>
        </div>
        <ProductCarousel />
      </div>
    </>
  );
};

export default Header;
