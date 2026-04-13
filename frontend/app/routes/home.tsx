import { Link, useParams } from "react-router";
import Header from "~/components/Header";
import Loader from "~/components/Loader";
import Message from "~/components/Message";
import ProductInfo from "~/components/Products/ProductInfo";
import { useGetPage } from "~/hooks/products";
import type { Route } from "./+types/home";

export function meta({ }: Route.MetaArgs) {
  return [
    { title: "GoMoR-E-Commerce" },
    { name: "description", content: "Welcome to GoMoR-E-Commerce!" },
  ];
}

export default function Home() {
  const { keyword } = useParams();
  const { data, isLoading, isError, error } = useGetPage({ keyword });
  const products = data?.items || [];

  return <>
    {!keyword ? <Header /> : null}
    {isLoading ? (
      <Loader />
    ) : (
      <>
        <div className="flex justify-between items-center">
          <h1 className="ml-[20rem] mt-[10rem] text-[3rem]">
            Special Products
          </h1>

          <Link
            to="/shop"
            className="bg-pink-600 font-bold rounded-full py-2 px-10 mr-[18rem] mt-[10rem]"
          >
            Shop
          </Link>
        </div>

        <div>
          <div className="flex justify-center flex-wrap mt-[2rem]">
            {isError ? (
              <Message variant="error">
                {error?.message || "Something went wrong"}
              </Message>
            ) : (
              products.map((product) => (
                <div key={product._id}>
                  <ProductInfo product={product} />
                </div>
              ))
            )}
          </div>
        </div>
      </>
    )}
  </>;
}
