import type { Route } from "./+types/home";
import { Link, useParams } from "react-router";
import Header from "~/components/Header";
import Loader from "~/components/Loader";
import Message from "~/components/Message";
import Product from "~/components/Products/Product";
import { useGetProductsQuery, useGetTopProductsQuery } from "~/redux/api/productApiSlice";

export function meta({ }: Route.MetaArgs) {
  return [
    { title: "GoMoR-E-Commerce" },
    { name: "description", content: "Welcome to GoMoR-E-Commerce!" },
  ];
}

export default function Home() {
  const { keyword } = useParams();
  const { data, isLoading, isError, error } = useGetProductsQuery({ keyword });

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
                {error?.data?.message || error?.error || "Something went wrong"}
              </Message>
            ) : (
              data.products?.map((product) => (
                <div key={product._id}>
                  <Product product={product} />
                </div>
              ))
            )}
          </div>
        </div>
      </>
    )}
  </>;
}
