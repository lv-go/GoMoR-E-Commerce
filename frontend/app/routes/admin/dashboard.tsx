import { useEffect, useState } from "react";
import Chart, { type Props } from "react-apexcharts";
import { useGetTotalOrdersQuery, useGetTotalSalesByDateQuery, useGetTotalSalesQuery } from "~/hooks/orders";
import { useGetPage as useGetUsersPage } from "~/hooks/users";
import Loader from "../../components/Loader";
import OrderList from "./orders";
import { FaBriefcase, FaDollarSign, FaUsers } from "react-icons/fa";
import { Link } from "react-router";

export default function Dashboard() {
  const { data: sales, isLoading: isLoadingSales } = useGetTotalSalesQuery();
  const { data: usersPage, isLoading: isLoadingCustomers } = useGetUsersPage();
  const customers = usersPage?.items || [];
  const { data: orders, isLoading: isLoadingOrders } = useGetTotalOrdersQuery();
  const { data: salesDetail, isLoading: isLoadingSalesDetail } = useGetTotalSalesByDateQuery();

  const [state, setState] = useState<Props>({
    options: {
      chart: {
        type: "line",
      },
      tooltip: {
        theme: "dark",
      },
      colors: ["#00E396"],
      dataLabels: {
        enabled: true,
      },
      stroke: {
        curve: "smooth",
      },
      title: {
        text: "Sales Trend",
        align: "left",
      },
      grid: {
        borderColor: "#ccc",
      },
      markers: {
        size: 1,
      },
      xaxis: {
        categories: [],
        title: {
          text: "Date",
        },
      },
      yaxis: {
        title: {
          text: "Sales",
        },
        min: 0,
      },
      legend: {
        position: "top",
        horizontalAlign: "right",
        floating: true,
        offsetY: -25,
        offsetX: -5,
      },
    },
    series: [{ name: "Sales", data: [] }],
  });

  useEffect(() => {
    if (salesDetail) {
      const formattedSalesDate = salesDetail.map((item) => ({
        x: item._id,
        y: item.total,
      }));

      setState((prevState) => ({
        ...prevState,
        options: {
          ...prevState.options,
          xaxis: {
            categories: formattedSalesDate.map((item) => item.x),
          },
        },

        series: [
          { name: "Sales", data: formattedSalesDate.map((item) => item.y) },
        ],
      }));
    }
  }, [salesDetail]);

  return (
    <>
      <section className="xl:ml-[4rem] md:ml-[0rem]">
        <div className="container flex justify-around flex-wrap">
          <Link to="/admin/orders" className="rounded-lg bg-base-300 p-5 w-[20rem] mt-5">
            <div className="flex items-center justify-center font-bold rounded-full w-[3rem] bg-pink-500 text-center p-3">
              <FaDollarSign size={24} />
            </div>

            <p className="mt-5">Sales</p>
            <h1 className="text-xl font-bold">
              $ {isLoadingSales ? <Loader /> : sales?.toFixed(2)}
            </h1>
          </Link>
          <Link to="/admin/users" className="rounded-lg bg-base-300 p-5 w-[20rem] mt-5">
            <div className="font-bold rounded-full w-[3rem] bg-pink-500 text-center p-3">
              <FaUsers size={24} />
            </div>

            <p className="mt-5">Customers</p>
            <h1 className="text-xl font-bold">
              {isLoadingCustomers ? <Loader /> : customers?.length}
            </h1>
          </Link>
          <Link to="/admin/orders" className="rounded-lg bg-base-300 p-5 w-[20rem] mt-5">
            <div className="flex items-center justify-center font-bold rounded-full w-[3rem] bg-pink-500 text-center p-3">
              <FaBriefcase size={24} />
            </div>

            <p className="mt-5">All Orders</p>
            <h1 className="text-xl font-bold">
              {isLoadingOrders ? <Loader /> : orders}
            </h1>
          </Link>
        </div>

        <div className="ml-[10rem] mt-[4rem]">
          <Chart
            options={state.options}
            series={state.series}
            type="bar"
            width="70%"
          />
        </div>

        <div className="mt-[4rem]">
          <OrderList />
        </div>
      </section>
    </>
  );
}
